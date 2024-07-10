import path from 'node:path';
import { DependencyContainer } from 'tsyringe';
import { DatabaseServer } from '@spt/servers/DatabaseServer';
import { SaveServer } from '@spt/servers/SaveServer';
import { LogTextColor } from '@spt/models/spt/logging/LogTextColor';
import { Watermark } from '@spt/utils/Watermark';
import { PreSptModLoader } from '@spt/loaders/PreSptModLoader';
import type { IPreSptLoadMod } from '@spt/models/external/IPreSptLoadMod';
import type { ILogger } from '@spt/models/spt/utils/ILogger';
import type { StaticRouterModService } from '@spt/services/mod/staticRouter/StaticRouterModService';

class GiveUI implements IPreSptLoadMod {
  public preSptLoad(container: DependencyContainer): void {
    const logger = container.resolve<ILogger>('WinstonLogger');
    const databaseServer = container.resolve<DatabaseServer>('DatabaseServer');
    const saveServer = container.resolve<SaveServer>('SaveServer');
    const watermark = container.resolve<Watermark>('Watermark');
    const preAkiModLoader = container.resolve<PreSptModLoader>('PreSptModLoader');

    const staticRouterModService =
      container.resolve<StaticRouterModService>('StaticRouterModService');

    // Hook up a new static route
    staticRouterModService.registerStaticRouter(
      'GiveUIModRouter',
      [
        {
          url: '/give-ui/server',
          action: (url, info, sessionId, output) => {
            logger.log(`[give-ui] Loading server info`, LogTextColor.GREEN);
            const version = watermark.getVersionTag();
            const serverPath = path.resolve();
            const modsInstalled = Object.values(preAkiModLoader.getImportedModDetails());
            const tarkovStashMod = modsInstalled.find((m) => m.name === 'give-ui');
            const modVersion = tarkovStashMod?.version;
            return Promise.resolve(JSON.stringify({ version, path: serverPath, modVersion }));
          },
        },
        {
          url: '/give-ui/profiles',
          action: (url, info, sessionId, output) => {
            logger.log(`[give-ui] Loading profiles`, LogTextColor.GREEN);
            return Promise.resolve(JSON.stringify(saveServer.getProfiles()));
          },
        },
        {
          url: '/give-ui/items',
          action: (url, info, sessionId, output) => {
            logger.log(`[give-ui] Loading items`, LogTextColor.GREEN);
            return Promise.resolve(JSON.stringify(databaseServer.getTables().templates.items));
          },
        },
        {
          url: '/give-ui/globals-presets',
          action: (url, info, sessionId, output) => {
            logger.log(`[give-ui] Loading global presets`, LogTextColor.GREEN);
            return Promise.resolve(JSON.stringify(databaseServer.getTables().globals.ItemPresets));
          },
        },
      ],
      'give-ui-top-level-route',
    );
  }
}

module.exports = { mod: new GiveUI() };
