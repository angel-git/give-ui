using System.Text.RegularExpressions;
using SPTarkov.DI.Annotations;
using SPTarkov.Server.Core.Helpers;
using SPTarkov.Server.Core.Helpers.Dialog.Commando.SptCommands;
using SPTarkov.Server.Core.Models.Eft.Dialog;
using SPTarkov.Server.Core.Models.Eft.Profile;
using SPTarkov.Server.Core.Servers;
using SPTarkov.Server.Core.Services;
using SPTarkov.Server.Core.Utils.Cloners;

namespace GiveUI.Command;

[Injectable]
public class GiveGearPresetSptCommand(
    MailSendService mailSendService,
    SaveServer saveServer,
    ICloner cloner,
    ItemHelper itemHelper) : ISptCommand
{
    private static readonly Regex _commandRegex = new(@"^spt give-gear-preset ([a-z]{2,5} ?"".+""|\w+)$");


    public string GetCommand()
    {
        return "give-gear-preset";
    }

    public string GetCommandHelp()
    {
        return
            "spt give-gear-preset\n========\nSends items to the player through the message system.\n\n\tspt give-user-preset [equipmentBuilds.Id]";
    }

    public string PerformAction(UserDialogInfo commandHandler, string sessionId, SendMessageRequest request)
    {
        if (request.Text == null || !_commandRegex.IsMatch(request.Text))
        {
            mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                "Invalid use of give command. Use 'help' for more information."
            );
            return request.DialogId ?? "";
        }

        var result = _commandRegex.Match(request.Text);

        var equipmentBuildId = result.Groups[1].Value;
        if (string.IsNullOrEmpty(equipmentBuildId))
        {
            mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                "Invalid use of give command. Use 'help' for more information."
            );
            return request.DialogId ?? "";
        }


        var profile = saveServer.GetProfiles()[sessionId];
        var equipmentBuilds = profile.UserBuildData?.EquipmentBuilds;
        var equipmentBuild = equipmentBuilds?.Find((eb) => eb.Id == equipmentBuildId);
        if (equipmentBuild == null)
        {
            mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                $"Couldn't find equipment build for Id: {equipmentBuildId}"
            );
            return request.DialogId ?? "";
        }

        var itemsToSend = cloner.Clone(equipmentBuild.Items) ?? [];
        itemsToSend.RemoveAt(0); // remove default inventory item
        itemsToSend = itemsToSend.Where(item =>
            item.SlotId != "Pockets" && item.SlotId != "SecuredContainer" && item.SlotId != "ArmBand" &&
            item.SlotId != "Dogtag").ToList();
        itemsToSend = itemHelper.ReplaceIDs(itemsToSend);
        itemHelper.SetFoundInRaid(itemsToSend);

        mailSendService.SendSystemMessageToPlayer(sessionId, "SPT GIVE", itemsToSend);
        return request.DialogId ?? "";
    }
}