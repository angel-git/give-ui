using System.Collections.Frozen;
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
using SPTarkov.Server.Core.Utils.Cloners;

namespace GiveUI.Command;

[Injectable]
public class GiveBundleSptCommand(
    ISptLogger<GiveBundleSptCommand> logger,
    MailSendService mailSendService,
    PresetHelper presetHelper,
    ICloner cloner,
    ItemHelper itemHelper) : ISptCommand
{
    private static readonly Regex _commandRegex = new(@"^spt give-bundle\s+((\{[^}]+\}\s*)+)$");

    // Exception for flares
    protected static readonly FrozenSet<MongoId> _excludedPresetItems =
    [
        ItemTpl.FLARE_RSP30_REACTIVE_SIGNAL_CARTRIDGE_RED,
        ItemTpl.FLARE_RSP30_REACTIVE_SIGNAL_CARTRIDGE_GREEN,
        ItemTpl.FLARE_RSP30_REACTIVE_SIGNAL_CARTRIDGE_YELLOW,
    ];


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
        var itemsFromBundle = JsonSerializer.Deserialize<Dictionary<string, int>>(json);
        if (itemsFromBundle != null)
            foreach (var (itemId, quantity) in itemsFromBundle)
            {
                var checkedItem = itemHelper.GetItem(itemId);
                if (checkedItem.Key)
                {
                    var preset = presetHelper.GetDefaultPreset(checkedItem.Value.Id);
                    if (preset is not null && !_excludedPresetItems.Contains(checkedItem.Value.Id))
                    {
                        for (var i = 0; i < quantity; i++)
                        {
                            var items = cloner.Clone(preset.Items);
                            items = items.ReplaceIDs().ToList();
                            itemsToSend.AddRange(items);
                        }
                    }
                    else if (itemHelper.IsOfBaseclass(checkedItem.Value.Id, BaseClasses.AMMO_BOX))
                    {
                        for (var i = 0; i < quantity; i++)
                        {
                            List<Item> ammoBoxArray =
                            [
                                new() { Id = new MongoId(), Template = checkedItem.Value.Id },
                                // DO NOT generate the ammo box cartridges, the mail service does it for us! :)
                                // _itemHelper.addCartridgesToAmmoBox(ammoBoxArray, checkedItem[1]);
                            ];
                            // DO NOT generate the ammo box cartridges, the mail service does it for us! :)
                            // _itemHelper.addCartridgesToAmmoBox(ammoBoxArray, checkedItem[1]);
                            itemsToSend.AddRange(ammoBoxArray);
                        }
                    }
                    else
                    {
                        if (checkedItem.Value.Properties.StackMaxSize == 1)
                        {
                            for (var i = 0; i < quantity; i++)
                            {
                                itemsToSend.Add(
                                    new Item
                                    {
                                        Id = new MongoId(),
                                        Template = checkedItem.Value.Id,
                                        Upd = itemHelper.GenerateUpdForItem(checkedItem.Value),
                                    }
                                );
                            }
                        }
                        else
                        {
                            var itemToSend = new Item
                            {
                                Id = new MongoId(),
                                Template = checkedItem.Value.Id,
                                Upd = itemHelper.GenerateUpdForItem(checkedItem.Value),
                            };
                            itemToSend.Upd.StackObjectsCount = quantity;
                            try
                            {
                                itemsToSend.AddRange(itemHelper.SplitStack(itemToSend));
                            }
                            catch
                            {
                                mailSendService.SendUserMessageToPlayer(
                                    sessionId,
                                    commandHandler,
                                    "Too many items requested. Please lower the amount and try again."
                                );

                                return new ValueTask<string>(request.DialogId);
                            }
                        }
                    }
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
            mailSendService.SendSystemMessageToPlayer(sessionId, "BUNDLE SENT", itemsToSend);
        }

        return new ValueTask<string>(request.DialogId);
    }
}