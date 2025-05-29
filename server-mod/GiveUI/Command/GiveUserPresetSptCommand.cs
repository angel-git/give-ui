using System.Text.RegularExpressions;
using SPTarkov.DI.Annotations;
using SPTarkov.Server.Core.Helpers;
using SPTarkov.Server.Core.Helpers.Dialog.Commando.SptCommands;
using SPTarkov.Server.Core.Models.Eft.Dialog;
using SPTarkov.Server.Core.Models.Eft.Profile;
using SPTarkov.Server.Core.Servers;
using SPTarkov.Server.Core.Services;

namespace SPTarkov.Server.Core.GiveUI;

[Injectable]
public class GiveUserPresetSptCommand(
    MailSendService _mailSendService,
    SaveServer _saveServer,
    Utils.Cloners.FastCloner _cloner,
    ItemHelper _itemHelper) : ISptCommand
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
        var userPresetId = result.Groups[1].Value;
        if (string.IsNullOrEmpty(userPresetId))
        {
            _mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                "Invalid use of give command. Use 'help' for more information."
            );
            return request.DialogId;
        }

        var profile = _saveServer.GetProfiles()[sessionId];
        var weaponBuilds = profile.UserBuildData?.WeaponBuilds;
        var weaponBuild = weaponBuilds.Find((wb) => wb.Id == userPresetId);
        if (weaponBuild == null)
        {
            _mailSendService.SendUserMessageToPlayer(
                sessionId,
                commandHandler,
                $"Couldn't find weapon build for Id: {userPresetId}"
            );
            return request.DialogId;
        }

        var itemsToSend = _cloner.Clone(weaponBuild.Items);
        itemsToSend = _itemHelper.ReplaceIDs(itemsToSend);
        _itemHelper.SetFoundInRaid(itemsToSend);
        _mailSendService.SendSystemMessageToPlayer(sessionId, "SPT GIVE", itemsToSend);

        return request.DialogId;
    }
}
