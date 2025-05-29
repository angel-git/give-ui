using SPTarkov.DI.Annotations;
using SPTarkov.Server.Core.Context;
using SPTarkov.Server.Core.DI;
using SPTarkov.Server.Core.Helpers;
using SPTarkov.Server.Core.Helpers.Dialogue;
using SPTarkov.Server.Core.Models.Eft.Dialog;
using SPTarkov.Server.Core.Models.Enums;
using SPTarkov.Server.Core.Models.Spt.Mod;
using SPTarkov.Server.Core.Servers;
using SPTarkov.Server.Core.Services;
using SPTarkov.Server.Core.Utils;

namespace SPTarkov.Server.Core.GiveUI;

[Injectable]
public class GiveUIStaticRouter : StaticRouter
{
    public GiveUIStaticRouter(
        JsonUtil jsonUtil,
        Watermark watermark,
        ApplicationContext applicationContext,
        ProfileHelper profileHelper,
        GiftService giftService,
        SaveServer saveServer,
        DatabaseServer databaseServer,
        CommandoDialogChatBot commandoDialogChatBot,
        SptDialogueChatBot sptDialogueChatBot
    ) : base(
        jsonUtil, [
            new RouteAction(
                "/give-ui/server",
                (
                    url,
                    info,
                    sessionID,
                    output
                ) =>
                {
                    var version = watermark.GetVersionTag();
                    var mods = applicationContext?.GetLatestValue(ContextVariableType.LOADED_MOD_ASSEMBLIES)
                        ?.GetValue<List<SptMod>>() ?? [];
                    var modVersion = mods.Find(m => m.ModMetadata?.Name == "give-ui")?.ModMetadata?.Version ?? "0";
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
                    sessionID,
                    output
                ) => jsonUtil.Serialize(saveServer.GetProfiles()) ?? "{}"),
            new RouteAction(
                "/give-ui/items",
                (
                    url,
                    info,
                    sessionID,
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
                    sessionID,
                    output
                ) =>
                {
                    var command = (info as GiveUIMessageRequest)?.Message ?? "";
                    var message = new SendMessageRequest
                    {
                        DialogId = sessionID,
                        Type = MessageType.SystemMessage,
                        Text = command
                    };
                    var response = commandoDialogChatBot.HandleMessage(sessionID ?? "", message);
                    return jsonUtil.Serialize(response) ?? "{}";
                },
                typeof(GiveUIMessageRequest)
                ),
            new RouteAction(
                "/give-ui/spt",
                (
                    url,
                    info,
                    sessionID,
                    output
                ) =>
                {
                    var command = (info as GiveUIMessageRequest)?.Message ?? "";
                    var message = new SendMessageRequest
                    {
                        DialogId = sessionID,
                        Type = MessageType.SystemMessage,
                        Text = command
                    };
                    var response = sptDialogueChatBot.HandleMessage(sessionID ?? "", message);
                    return jsonUtil.Serialize(response) ?? "{}";
                },
                typeof(GiveUIMessageRequest)
                )
        ])
    {
    }
}
