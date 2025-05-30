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
public class GiveUserPresetSptCommand(
    MailSendService mailSendService,
    SaveServer saveServer,
    ICloner cloner,
    ItemHelper itemHelper) : ISptCommand
{
    private static readonly Regex _commandRegex = new(@"^spt give-user-preset ((([a-z]{2,5}) )?""(.+)""|\w+)$");

    public string GetCommand()
    {
        return "give-user-preset";
    }

    public string GetCommandHelp()
    {
        return
            "spt give-user-preset\n========\nSends items to the player through the message system.\n\n\tspt give-user-preset [weaponBuilds.Id]";
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
        var userPresetId = result.Groups[1].Value;
        if (string.IsNullOrEmpty(userPresetId))
        {
            mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                "Invalid use of give command. Use 'help' for more information."
            );
            return request.DialogId ?? "";
        }

        var profile = saveServer.GetProfiles()[sessionId];
        var weaponBuilds = profile.UserBuildData?.WeaponBuilds ?? [];
        var weaponBuild = weaponBuilds.Find((wb) => wb.Id == userPresetId);
        if (weaponBuild == null)
        {
            mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                $"Couldn't find weapon build for Id: {userPresetId}"
            );
            return request.DialogId ?? "";
        }

        var itemsToSend = cloner.Clone(weaponBuild.Items) ?? [];
        itemsToSend = itemHelper.ReplaceIDs(itemsToSend);
        itemHelper.SetFoundInRaid(itemsToSend);
        mailSendService.SendSystemMessageToPlayer(sessionId, "SPT GIVE", itemsToSend);

        return request.DialogId ?? "";
    }
}