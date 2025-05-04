import path from 'node:path';
import fs from "node:fs";
import {DependencyContainer} from 'tsyringe';
import {DatabaseServer} from '@spt/servers/DatabaseServer';
import {SaveServer} from '@spt/servers/SaveServer';
import {LogTextColor} from '@spt/models/spt/logging/LogTextColor';
import {Watermark} from '@spt/utils/Watermark';
import {PreSptModLoader} from '@spt/loaders/PreSptModLoader';
import {CommandoDialogueChatBot} from "@spt/helpers/Dialogue/CommandoDialogueChatBot";
import {SptDialogueChatBot} from "@spt/helpers/Dialogue/SptDialogueChatBot";
import type {IPreSptLoadMod} from '@spt/models/external/IPreSptLoadMod';
import type {ILogger} from '@spt/models/spt/utils/ILogger';
import type {StaticRouterModService} from '@spt/services/mod/staticRouter/StaticRouterModService';
import type {DynamicRouterModService} from '@spt/services/mod/dynamicRouter/DynamicRouterModService';
import {MessageType} from "@spt/models/enums/MessageType";
import {ISendMessageRequest} from "@spt/models/eft/dialog/ISendMessageRequest";
import {SptCommandoCommands} from "@spt/helpers/Dialogue/Commando/SptCommandoCommands";
import {ProfileHelper} from "@spt/helpers/ProfileHelper";
import { GiftService } from "@spt/services/GiftService";
import {GiveUserPresetSptCommand} from './GiveUserPresetSptCommand';
import {GiveStashItemSptCommand} from "./GiveStashItemSptCommand";

class GiveUI implements IPreSptLoadMod {
    public preSptLoad(container: DependencyContainer): void {

        container.register<GiveUserPresetSptCommand>("GiveUserPresetSptCommand", GiveUserPresetSptCommand);
        container.register<GiveStashItemSptCommand>("GiveStashItemSptCommand", GiveStashItemSptCommand);
        container.resolve<SptCommandoCommands>("SptCommandoCommands").registerSptCommandoCommand(container.resolve<GiveUserPresetSptCommand>("GiveUserPresetSptCommand"));
        container.resolve<SptCommandoCommands>("SptCommandoCommands").registerSptCommandoCommand(container.resolve<GiveStashItemSptCommand>("GiveStashItemSptCommand"));

        const logger = container.resolve<ILogger>('WinstonLogger');
        const databaseServer = container.resolve<DatabaseServer>('DatabaseServer');
        const saveServer = container.resolve<SaveServer>('SaveServer');
        const watermark = container.resolve<Watermark>('Watermark');
        const preAkiModLoader = container.resolve<PreSptModLoader>('PreSptModLoader');
        const commandoDialog = container.resolve<CommandoDialogueChatBot>('CommandoDialogueChatBot');
        const sptDialog = container.resolve<SptDialogueChatBot>('SptDialogueChatBot');
        const profileHelper = container.resolve<ProfileHelper>('ProfileHelper');
        const giftService = container.resolve<GiftService>('GiftService');

        const staticRouterModService =
            container.resolve<StaticRouterModService>('StaticRouterModService');

        const dynamicRouterModService =
            container.resolve<DynamicRouterModService>('DynamicRouterModService');

        // Hook up a new static route
        staticRouterModService.registerStaticRouter(
            'GiveUIModRouter',
            [
                {
                    url: '/give-ui/server',
                    action: (_url, _info, _sessionId, _output) => {
                        logger.log(`[give-ui] Loading server info`, LogTextColor.GREEN);
                        const version = watermark.getVersionTag();
                        const serverPath = path.resolve();
                        const modsInstalled = Object.values(preAkiModLoader.getImportedModDetails());
                        const giveUiMod = modsInstalled.find((m) => m.name === 'give-ui');
                        const modVersion = giveUiMod?.version;
                        const maxLevel = profileHelper.getMaxLevel();
                        const gifts = giftService.getGifts();
                        return Promise.resolve(JSON.stringify({version, path: serverPath, modVersion, maxLevel, gifts}));
                    },
                },
                {
                    url: '/give-ui/profiles',
                    action: (_url, _info, _sessionId, _output) => {
                        logger.log(`[give-ui] Loading profiles`, LogTextColor.GREEN);
                        return Promise.resolve(JSON.stringify(saveServer.getProfiles()));
                    },
                },
                {
                    url: '/give-ui/items',
                    action: (_url, _info, _sessionId, _output) => {
                        logger.log(`[give-ui] Loading items`, LogTextColor.GREEN);
                        return Promise.resolve(JSON.stringify({
                            items: databaseServer.getTables().templates.items,
                            globalPresets: databaseServer.getTables().globals.ItemPresets
                        }));
                    },
                },
                {
                    url: '/give-ui/commando',
                    action: (_url, request, sessionId, _output) => {
                        const command = request.message;
                        logger.log(`[give-ui] Sending to commando: [${command}]`, LogTextColor.GREEN);
                        const message: ISendMessageRequest = {
                            dialogId: sessionId,
                            type: MessageType.SYSTEM_MESSAGE,
                            text: command,
                            replyTo: undefined,
                        };
                        const response = commandoDialog.handleMessage(sessionId, message);
                        return Promise.resolve(JSON.stringify({response}));
                    },
                },
                {
                    url: '/give-ui/spt',
                    action: (_url, request, sessionId, _output) => {
                        const command = request.message;
                        logger.log(`[give-ui] Sending to spt: [${command}]`, LogTextColor.GREEN);
                        const message: ISendMessageRequest = {
                            dialogId: sessionId,
                            type: MessageType.SYSTEM_MESSAGE,
                            text: command,
                            replyTo: undefined,
                        };
                        const response = sptDialog.handleMessage(sessionId, message);
                        return Promise.resolve(JSON.stringify({response}));
                    },
                },
            ],
            'give-ui-top-level-route',
        );

        dynamicRouterModService.registerDynamicRouter(
            'GiveUIDynamicModRouter',
            [{
                url: '/give-ui/cache',
                action: (url, _request, _sessionId, _output) => {
                    const cacheID = url.replace('/give-ui/cache/', '');
                    const serverPath = path.resolve();
                    const cachePath = path.join(serverPath, 'user', 'sptappdata', 'live');
                    try {
                        const indexJson = fs.readFileSync(path.join(cachePath, 'index.json'), 'utf8');
                        const index = JSON.parse(indexJson);
                        const image= index[cacheID]
                        try {
                            const file = fs.readFileSync(path.join(cachePath, `${image}.png`),  {encoding: 'base64'});
                            return Promise.resolve(JSON.stringify({imageBase64: file}));
                        } catch (e) {
                            return Promise.resolve(JSON.stringify({error: 404}));
                        }
                    } catch (e) {
                        return Promise.resolve(JSON.stringify({error: 404}));
                    }
                },
            }],
            'give-ui-top-level-dynamic-route',
        )
    }
}

module.exports = {mod: new GiveUI()};
