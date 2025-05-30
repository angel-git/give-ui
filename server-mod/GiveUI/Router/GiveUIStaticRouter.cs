using SPTarkov.DI.Annotations;
using SPTarkov.Server.Core.Controllers;
using SPTarkov.Server.Core.DI;
using SPTarkov.Server.Core.Helpers;
using SPTarkov.Server.Core.Helpers.Dialogue;
using SPTarkov.Server.Core.Models.Eft.Dialog;
using SPTarkov.Server.Core.Models.Enums;
using SPTarkov.Server.Core.Models.Spt.Mod;
using SPTarkov.Server.Core.Servers;
using SPTarkov.Server.Core.Services;
using SPTarkov.Server.Core.Utils;

namespace GiveUI.Router;

[Injectable]
public class GiveUIStaticRouter : StaticRouter
{
    public GiveUIStaticRouter(
        JsonUtil jsonUtil,
        Watermark watermark,
        ProfileHelper profileHelper,
        GiftService giftService,
        SaveServer saveServer,
        DatabaseServer databaseServer,
        LauncherController launcherController,
        CommandoDialogChatBot commandoDialogChatBot,
        SptDialogueChatBot sptDialogueChatBot
    ) : base(
        jsonUtil, [
            new RouteAction(
                "/give-ui/server",
                (
                    url,
                    info,
                    sessionId,
                    output
                ) =>
                {
                    var version = watermark.GetVersionTag();
                    var loadedMods = launcherController.GetLoadedServerMods();
                    var modVersion = "-1";
                    if (loadedMods.ContainsKey("give-ui"))
                    {
                        modVersion = loadedMods["give-ui"].Version;
                    }

                    var maxLevel = profileHelper.GetMaxLevel();
                    var gifts = giftService.GetGifts();
                    return jsonUtil.Serialize(new
                    {
                        version,
                        modVersion,
                        maxLevel,
                        gifts
                    }) ?? "{}";
                }
            ),
            new RouteAction(
                "/give-ui/profiles",
                (
                    url,
                    info,
                    sessionId,
                    output
                ) => jsonUtil.Serialize(saveServer.GetProfiles()) ?? "{}"),
            new RouteAction(
                "/give-ui/items",
                (
                    url,
                    info,
                    sessionId,
                    output
                ) =>
                {
                    var items = databaseServer.GetTables().Templates?.Items;
                    var globalPresets = databaseServer.GetTables().Globals?.ItemPresets;
                    return jsonUtil.Serialize(new
                    {
                        items,
                        globalPresets
                    }) ?? "{}";
                }),
            new RouteAction(
                "/give-ui/commando",
                (
                    url,
                    info,
                    sessionId,
                    output
                ) =>
                {
                    var command = (info as GiveUIMessageRequest)?.Message ?? "";
                    var message = new SendMessageRequest
                    {
                        DialogId = sessionId,
                        Type = MessageType.SystemMessage,
                        Text = command
                    };
                    var response = commandoDialogChatBot.HandleMessage(sessionId ?? "", message);
                    return jsonUtil.Serialize(response) ?? "{}";
                },
                typeof(GiveUIMessageRequest)
            ),
            new RouteAction(
                "/give-ui/spt",
                (
                    url,
                    info,
                    sessionId,
                    output
                ) =>
                {
                    var command = (info as GiveUIMessageRequest)?.Message ?? "";
                    var message = new SendMessageRequest
                    {
                        DialogId = sessionId,
                        Type = MessageType.SystemMessage,
                        Text = command
                    };
                    var response = sptDialogueChatBot.HandleMessage(sessionId ?? "", message);
                    return jsonUtil.Serialize(response) ?? "{}";
                },
                typeof(GiveUIMessageRequest)
            )
        ])
    {
    }
}