using System.Text.Json;
using System.Text.RegularExpressions;
using SPTarkov.DI.Annotations;
using SPTarkov.Server.Core.Extensions;
using SPTarkov.Server.Core.Helpers;
using SPTarkov.Server.Core.Helpers.Dialog.Commando.SptCommands;
using SPTarkov.Server.Core.Models.Common;
using SPTarkov.Server.Core.Models.Eft.Common.Tables;
using SPTarkov.Server.Core.Models.Eft.Dialog;
using SPTarkov.Server.Core.Models.Eft.Profile;
using SPTarkov.Server.Core.Models.Utils;
using SPTarkov.Server.Core.Services;

namespace GiveUI.Command;

[Injectable]
public class GiveBundleSptCommand(
    ISptLogger<GiveBundleSptCommand> logger,
    MailSendService mailSendService,
    ItemHelper itemHelper) : ISptCommand
{
    private static readonly Regex _commandRegex = new(@"^spt give-bundle\s+((\{[^}]+\}\s*)+)$");


    string ISptCommand.Command => "give-bundle";

    string ISptCommand.CommandHelp =>
        "spt give-bundle\n========\nSends bundle items. This message is not meant to be written by user";

    public ValueTask<string> PerformAction(UserDialogInfo commandHandler, MongoId sessionId, SendMessageRequest request)
    {
        var match = _commandRegex.Match(request.Text);
        if (!match.Success)
        {
            mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                "Invalid use of give command. Use 'help' for more information."
            );
            return new ValueTask<string>(request.DialogId);
        }

        var itemsToSend = new List<Item>();
        // Extract the items string
        var json = match.Groups[1].Value;
        var items = JsonSerializer.Deserialize<Dictionary<string, int>>(json);
        if (items != null)
            foreach (var (itemId, quantity) in items)
            {
                var item = itemHelper.GetItem(itemId);
                if (item.Key)
                {
                    var itemToSend = new Item
                    {
                        Id = new MongoId(),
                        Template = item.Value.Id,
                        Upd = itemHelper.GenerateUpdForItem(item.Value),
                    };
                    itemToSend.Upd.StackObjectsCount = quantity;
                    itemsToSend.Add(itemToSend);
                }
                else
                {
                    logger.Warning($"Could not find item: {itemId} for bundle, skipping it");
                }
            }

        if (itemsToSend.Count == 0)
        {
            logger.Warning($"No items in the bundle, nothing to send");
        }
        else
        {
            itemsToSend = itemsToSend.ReplaceIDs().ToList();
            itemHelper.SetFoundInRaid(itemsToSend);
            mailSendService.SendSystemMessageToPlayer(sessionId, "SPT GIVE", itemsToSend);
        }

        return new ValueTask<string>(request.DialogId);
    }
}