import path from 'node:path';
import fs from "node:fs";
import {DependencyContainer} from 'tsyringe';
import {DatabaseServer} from '@spt/servers/DatabaseServer';
import {SaveServer} from '@spt/servers/SaveServer';
import {LogTextColor} from '@spt/models/spt/logging/LogTextColor';
import {Watermark} from '@spt/utils/Watermark';
import {PreSptModLoader} from '@spt/loaders/PreSptModLoader';
import {CommandoDialogueChatBot} from "@spt/helpers/Dialogue/CommandoDialogueChatBot";
import type {IPreSptLoadMod} from '@spt/models/external/IPreSptLoadMod';
import type {ILogger} from '@spt/models/spt/utils/ILogger';
import type {StaticRouterModService} from '@spt/services/mod/staticRouter/StaticRouterModService';
import type {DynamicRouterModService} from '@spt/services/mod/dynamicRouter/DynamicRouterModService';
import {MessageType} from "@spt/models/enums/MessageType";
import {ISendMessageRequest} from "@spt/models/eft/dialog/ISendMessageRequest";
import {SptCommandoCommands} from "@spt/helpers/Dialogue/Commando/SptCommandoCommands";
import {ProfileHelper} from "@spt/helpers/ProfileHelper";
import {GiveUserPresetSptCommand} from './GiveUserPresetSptCommand';

class GiveUI implements IPreSptLoadMod {
    public preSptLoad(container: DependencyContainer): void {

        container.register<GiveUserPresetSptCommand>("GiveUserPresetSptCommand", GiveUserPresetSptCommand);
        container.resolve<SptCommandoCommands>("SptCommandoCommands").registerSptCommandoCommand(container.resolve<GiveUserPresetSptCommand>("GiveUserPresetSptCommand"));

        const logger = container.resolve<ILogger>('WinstonLogger');
        const databaseServer = container.resolve<DatabaseServer>('DatabaseServer');
        const saveServer = container.resolve<SaveServer>('SaveServer');
        const watermark = container.resolve<Watermark>('Watermark');
        const preAkiModLoader = container.resolve<PreSptModLoader>('PreSptModLoader');
        const commando = container.resolve<CommandoDialogueChatBot>('CommandoDialogueChatBot');
        const profileHelper = container.resolve<ProfileHelper>('ProfileHelper');

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
                        return Promise.resolve(JSON.stringify({version, path: serverPath, modVersion, maxLevel}));
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
                    url: '/give-ui/give',
                    action: (_url, request, sessionId, _output) => {
                        const command = `spt give ${request.itemId} ${request.amount}`;
                        logger.log(`[give-ui] Running command: [${command}]`, LogTextColor.GREEN);
                        const message: ISendMessageRequest = {
                            dialogId: sessionId,
                            type: MessageType.SYSTEM_MESSAGE,
                            text: command,
                            replyTo: undefined,
                        };
                        const response = commando.handleMessage(sessionId, message);
                        return Promise.resolve(JSON.stringify({response}));
                    },
                },
                {
                    url: '/give-ui/give-user-preset',
                    action: (_url, request, sessionId, _output) => {
                        const command = `spt give-user-preset ${request.itemId}`;
                        logger.log(`[give-ui] Running command: [${command}]`, LogTextColor.GREEN);
                        const message: ISendMessageRequest = {
                            dialogId: sessionId,
                            type: MessageType.SYSTEM_MESSAGE,
                            text: command,
                            replyTo: undefined,
                        };
                        const response = commando.handleMessage(sessionId, message);
                        return Promise.resolve(JSON.stringify({response}));
                    },
                },
                {
                    url: '/give-ui/update-trader-rep',
                    action: (_url, request, sessionId, _output) => {
                        const repCommand = `spt trader ${request.nickname} rep ${request.rep}`;
                        logger.log(`[give-ui] Running command: [${repCommand}]`, LogTextColor.GREEN);
                        const repMessage: ISendMessageRequest = {
                            dialogId: sessionId,
                            type: MessageType.SYSTEM_MESSAGE,
                            text: repCommand,
                            replyTo: undefined,
                        };
                        const response = commando.handleMessage(sessionId, repMessage);
                        return Promise.resolve(JSON.stringify({response}));
                    },
                },
                {
                    url: '/give-ui/update-trader-spend',
                    action: (_url, request, sessionId, _output) => {
                        const spendCommand = `spt trader ${request.nickname} spend ${request.spend}`;
                        logger.log(`[give-ui] Running command: [${spendCommand}]`, LogTextColor.GREEN);
                        const spendMessage: ISendMessageRequest = {
                            dialogId: sessionId,
                            type: MessageType.SYSTEM_MESSAGE,
                            text: spendCommand,
                            replyTo: undefined,
                        };
                        const response = commando.handleMessage(sessionId, spendMessage);
                        return Promise.resolve(JSON.stringify({response}));
                    },
                },
                {
                    url: '/give-ui/update-level',
                    action: (_url, request, sessionId, _output) => {
                        const command = `spt profile level ${request.level}`;
                        logger.log(`[give-ui] Running command: [${command}]`, LogTextColor.GREEN);
                        const message: ISendMessageRequest = {
                            dialogId: sessionId,
                            type: MessageType.SYSTEM_MESSAGE,
                            text: command,
                            replyTo: undefined,
                        };
                        const response = commando.handleMessage(sessionId, message);
                        return Promise.resolve(JSON.stringify({response}));
                    },
                },
                {
                    url: '/give-ui/update-skill',
                    action: (_url, request, sessionId, _output) => {
                        const command = `spt profile skill ${request.skill} ${request.progress}`;
                        logger.log(`[give-ui] Running command: [${command}]`, LogTextColor.GREEN);
                        const message: ISendMessageRequest = {
                            dialogId: sessionId,
                            type: MessageType.SYSTEM_MESSAGE,
                            text: command,
                            replyTo: undefined,
                        };
                        const response = commando.handleMessage(sessionId, message);
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
                            console.log(`${cacheID} cache not found in index.json`)
                            return Promise.resolve(JSON.stringify({error: 404}));
                        }
                    } catch (e) {
                        console.log('sptappdata not found')
                        return Promise.resolve(JSON.stringify({error: 404}));
                    }


                },
            }],
            'give-ui-top-level-dynamic-route',
        )
    }
}

module.exports = {mod: new GiveUI()};
