import {inject, injectable} from "tsyringe";
import {ISptCommand} from "@spt/helpers/Dialogue/Commando/SptCommands/ISptCommand";
import {ItemHelper} from "@spt/helpers/ItemHelper";
import {ISendMessageRequest} from "@spt/models/eft/dialog/ISendMessageRequest";
import {ISptProfile} from "@spt/models/eft/profile/ISptProfile";
import {IUserDialogInfo} from "@spt/models/eft/profile/IUserDialogInfo";
import {MailSendService} from "@spt/services/MailSendService";
import type {ICloner} from "@spt/utils/cloners/ICloner";
import {SaveServer} from '@spt/servers/SaveServer';


@injectable()
export class GiveGearPresetSptCommand implements ISptCommand {
    private static commandRegex = /^spt give-gear-preset ([a-z]{2,5} ?".+"|\w+)$/;

    public constructor(
        @inject("ItemHelper") protected itemHelper: ItemHelper,
        @inject("MailSendService") protected mailSendService: MailSendService,
        @inject("PrimaryCloner") protected cloner: ICloner,
        @inject("SaveServer") protected saveServer: SaveServer,
    ) {
    }

    public getCommand(): string {
        return "give-gear-preset";
    }

    public getCommandHelp(): string {
        return 'spt give-gear-preset\n========\nSends items to the player through the message system.\n\n\tspt give-user-preset [equipmentBuilds.Id]';
    }

    public performAction(commandHandler: IUserDialogInfo, sessionId: string, request: ISendMessageRequest): string {
        if (!GiveGearPresetSptCommand.commandRegex.test(request.text)) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                'Invalid use of give command. Use "help" for more information.',
            );
            return request.dialogId;
        }

        const result = GiveGearPresetSptCommand.commandRegex.exec(request.text);

        if (!result[1]) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                `Invalid use of give command. Use "help" for more information.`,
            );
            return request.dialogId;
        }

        const profile: ISptProfile = this.saveServer.getProfiles()[sessionId];
        const equipmentBuilds = profile.userbuilds.equipmentBuilds;
        const equipmentBuild = equipmentBuilds.find((eb) => eb.Id === result[1]);
        if (!equipmentBuild) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                `Couldn't find equipment build for Id: ${result[1]}`,
            );
            return request.dialogId;
        }

        let itemsToSend = this.cloner.clone(equipmentBuild.Items);
        itemsToSend.shift(); // remove default inventory item
        itemsToSend = itemsToSend.filter(item => item.slotId !== "Pockets" && item.slotId !== "SecuredContainer" && item.slotId !== "ArmBand" && item.slotId !== "Dogtag");
        itemsToSend = this.itemHelper.replaceIDs(itemsToSend);
        this.itemHelper.setFoundInRaid(itemsToSend);

        this.mailSendService.sendSystemMessageToPlayer(sessionId, "SPT GIVE", itemsToSend);
        return request.dialogId;
    }
}
