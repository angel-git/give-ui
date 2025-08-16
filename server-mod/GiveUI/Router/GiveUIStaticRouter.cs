using SPTarkov.DI.Annotations;
using SPTarkov.Server.Core.Controllers;
using SPTarkov.Server.Core.DI;
using SPTarkov.Server.Core.Helpers;
using SPTarkov.Server.Core.Helpers.Dialogue;
using SPTarkov.Server.Core.Models.Eft.Dialog;
using SPTarkov.Server.Core.Models.Eft.Quests;
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
        DatabaseService databaseService,
        LauncherController launcherController,
        CommandoDialogChatBot commandoDialogChatBot,
        SptDialogueChatBot sptDialogueChatBot,
        QuestHelper questHelper
    ) : base(
        jsonUtil, [
            new RouteAction(
                "/give-ui/server",
                async (
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
                    return await new ValueTask<string>(jsonUtil.Serialize(new
                    {
                        version,
                        modVersion,
                        maxLevel,
                        gifts
                    }) ?? "{}");
                }
            ),
            new RouteAction(
                "/give-ui/profiles",
                async (
                    url,
                    info,
                    sessionId,
                    output
                ) => await new ValueTask<string>(jsonUtil.Serialize(saveServer.GetProfiles()) ?? "{}")
            ),
            new RouteAction(
                "/give-ui/items",
                async (
                    url,
                    info,
                    sessionId,
                    output
                ) =>
                {
                    var items = databaseService.GetTemplates().Items;
                    var globalPresets = databaseService.GetGlobals().ItemPresets;
                    return await new ValueTask<string>(jsonUtil.Serialize(new
                    {
                        items,
                        globalPresets
                    }) ?? "{}");
                }),
            new RouteAction(
                "/give-ui/commando",
                async (
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
                        Text = command,
                        ReplyTo = "",
                    };
                    return await commandoDialogChatBot.HandleMessage(sessionId, message);
                },
                typeof(GiveUIMessageRequest)
            ),
            new RouteAction(
                "/give-ui/spt",
                async (
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
                        Text = command,
                        ReplyTo = "",
                    };
                    return await sptDialogueChatBot.HandleMessage(sessionId, message);
                },
                typeof(GiveUIMessageRequest)
            )
            ,
            new RouteAction(
                "/give-ui/quest",
                async (
                    url,
                    info,
                    sessionId,
                    output
                ) =>
                {
                    var questId = (info as GiveUIQuestRequest)?.QuestId ?? "";
                    var completeQuestRequestData = new CompleteQuestRequestData
                    {
                        QuestId = questId
                    };
                    questHelper.CompleteQuest(saveServer.GetProfiles()[sessionId].CharacterData!.PmcData!, completeQuestRequestData, sessionId);
                    return await new ValueTask<string>("{\"ok\": true}");
                },
                typeof(GiveUIQuestRequest)
            )
        ])
    {
    }
}