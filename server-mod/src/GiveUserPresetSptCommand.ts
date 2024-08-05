import {inject, injectable} from "tsyringe";
import {ISptCommand} from "@spt/helpers/Dialogue/Commando/SptCommands/ISptCommand";
import {ItemHelper} from "@spt/helpers/ItemHelper";
import {ISendMessageRequest} from "@spt/models/eft/dialog/ISendMessageRequest";
import {ISptProfile, IUserDialogInfo} from "@spt/models/eft/profile/ISptProfile";
import {MailSendService} from "@spt/services/MailSendService";
import {ICloner} from "@spt/utils/cloners/ICloner";
import {SaveServer} from '@spt/servers/SaveServer';


@injectable()
export class GiveUserPresetSptCommand implements ISptCommand {
    private static commandRegex = /^spt give-user-preset ((([a-z]{2,5}) )?"(.+)"|\w+)$/;

    public constructor(
        @inject("ItemHelper") protected itemHelper: ItemHelper,
        @inject("MailSendService") protected mailSendService: MailSendService,
        @inject("PrimaryCloner") protected cloner: ICloner,
        @inject("SaveServer") protected saveServer: SaveServer,
    ) {
    }

    public getCommand(): string {
        return "give-user-preset";
    }

    public getCommandHelp(): string {
        return 'spt give-user-preset\n========\nSends items to the player through the message system.\n\n\tspt give-user-preset [weaponBuilds.Id]';
    }

    public performAction(commandHandler: IUserDialogInfo, sessionId: string, request: ISendMessageRequest): string {
        if (!GiveUserPresetSptCommand.commandRegex.test(request.text)) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                'Invalid use of give command. Use "help" for more information.',
            );
            return request.dialogId;
        }

        const result = GiveUserPresetSptCommand.commandRegex.exec(request.text);

        if (!result[1]) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                `Invalid use of give command. Use "help" for more information.`,
            );
            return request.dialogId;
        }

        const profile: ISptProfile = this.saveServer.getProfiles()[sessionId];
        const weaponBuilds = profile.userbuilds.weaponBuilds;
        const weaponBuild = weaponBuilds.find((wb) => wb.Id === result[1]);
        if (!weaponBuild) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                `Couldn't find weapon build for Id: ${result[1]}`,
            );
            return request.dialogId;
        }

        let itemsToSend = this.cloner.clone(weaponBuild.Items);
        itemsToSend = this.itemHelper.replaceIDs(itemsToSend);
        this.itemHelper.setFoundInRaid(itemsToSend);
        this.mailSendService.sendSystemMessageToPlayer(sessionId, "SPT GIVE", itemsToSend);

        return request.dialogId;
    }
}
