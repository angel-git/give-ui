using System.Text.RegularExpressions;
using SPTarkov.DI.Annotations;
using SPTarkov.Server.Core.Helpers;
using SPTarkov.Server.Core.Helpers.Dialog.Commando.SptCommands;
using SPTarkov.Server.Core.Models.Eft.Common.Tables;
using SPTarkov.Server.Core.Models.Eft.Dialog;
using SPTarkov.Server.Core.Models.Eft.Profile;
using SPTarkov.Server.Core.Servers;
using SPTarkov.Server.Core.Services;

namespace SPTarkov.Server.Core.GiveUI;

[Injectable]
public class GiveStashItemSptCommand(
    MailSendService _mailSendService,
    SaveServer _saveServer,
    Utils.Cloners.FastCloner _cloner,
    ItemHelper _itemHelper): ISptCommand
{

    private static readonly Regex _commandRegex = new(@"^spt give-user-stash-item ((([a-z]{2,5}) )?""(.+)""|\w+)$");

    public string GetCommand()
    {
        return "give-user-stash-item";
    }

    public string GetCommandHelp()
    {
        return
            "spt give-user-stash-item\n========\nSends items to the player through the message system.\n\n\tspt give-user-stash-item [stashItem.Id]";
    }

    public string PerformAction(UserDialogInfo commandHandler, string sessionId, SendMessageRequest request)
    {
        if (!_commandRegex.IsMatch(request.Text))
        {
            _mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                "Invalid use of give command. Use 'help' for more information."
            );
            return request.DialogId;
        }

        var result = _commandRegex.Match(request.Text);
        var itemId = result.Groups[1].Value;
        if (string.IsNullOrEmpty(itemId))
        {
            _mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                "Invalid use of give command. Use 'help' for more information."
            );
            return request.DialogId;
        }

        var profile = _saveServer.GetProfiles()[sessionId];

        var inventoryItemHash = GetInventoryItemHash(profile.CharacterData.PmcData.Inventory.Items);
        var itemToAdd = inventoryItemHash.ByItemId[itemId];
        if (itemToAdd == null)
        {
            _mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                $"Couldn't find item with Id: {itemId}"
            );
            return request.DialogId;
        }

        var checkedItem = _itemHelper.GetItem(itemToAdd.Template);
        if (!checkedItem.Key) {
            _mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                $"Couldn't find template with id: ${itemToAdd.Template}"
            );
            return request.DialogId;
        }

        var allChild = GetAllDescendantsIncludingSelf(itemId, inventoryItemHash);
        var itemsToSend = _cloner.Clone(allChild);

        itemsToSend = _itemHelper.ReplaceIDs(itemsToSend);
        _itemHelper.SetFoundInRaid(itemsToSend);
        _mailSendService.SendSystemMessageToPlayer(sessionId, "SPT GIVE", itemsToSend);

        return request.DialogId;

    }

    private InventoryItemHash GetInventoryItemHash(List<Item> inventoryItems)
    {
        var inventoryItemHash = new InventoryItemHash();

        foreach (var item in inventoryItems)
        {
            inventoryItemHash.ByItemId[item.Id] = item;

            if (item.ParentId == null)
            {
                continue;
            }

            if (!inventoryItemHash.ByParentId.ContainsKey(item.ParentId))
            {
                inventoryItemHash.ByParentId[item.ParentId] = [];
            }

            inventoryItemHash.ByParentId[item.ParentId].Add(item);
        }

        return inventoryItemHash;
    }

    private List<Item> GetAllDescendantsIncludingSelf(string parentId, InventoryItemHash hash)
    {
        var result = new List<Item>();

        if (hash.ByItemId.TryGetValue(parentId, out var self))
        {
            result.Add(self);
        }

        if (hash.ByParentId.TryGetValue(parentId, out var directChildren))
        {
            foreach (var child in directChildren)
            {
                result.AddRange(GetAllDescendantsIncludingSelf(child.Id, hash));
            }
        }

        return result;
    }

}
