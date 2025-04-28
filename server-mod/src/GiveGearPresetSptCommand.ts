import {inject, injectable} from "tsyringe";
import {ISptCommand} from "@spt/helpers/Dialogue/Commando/SptCommands/ISptCommand";
import {ItemHelper} from "@spt/helpers/ItemHelper";
import {ISendMessageRequest} from "@spt/models/eft/dialog/ISendMessageRequest";
import {ISptProfile, IUserDialogInfo} from "@spt/models/eft/profile/ISptProfile";
import {MailSendService} from "@spt/services/MailSendService";
import {ICloner} from "@spt/utils/cloners/ICloner";
import {SaveServer} from '@spt/servers/SaveServer';


@injectable()
export class GiveGearPresetSptCommand implements ISptCommand {
    private static commandRegex = /^spt give-gear-preset ([a-z]{2,5} ?".+"|\w+) ([a-z]{2,5} ?".+"|\w+)$/;

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
        return 'TODO';
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

        if (!result[1] || !result[2]) {
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

        const gearPresetItems = this.itemHelper.findAndReturnChildrenAsItems(equipmentBuild.Items, result[2])
        // TODO filter out: Pockets, SecuredContainer, ArmBand
        let itemsToSend = this.cloner.clone(gearPresetItems);
        if (itemsToSend.length > 0) {
            // clear slotId from main item so we can accept it in UI
            itemsToSend[0].slotId = undefined;
        }
        itemsToSend = this.itemHelper.replaceIDs(itemsToSend);
        this.itemHelper.setFoundInRaid(itemsToSend);
        this.mailSendService.sendSystemMessageToPlayer(sessionId, "SPT GIVE", itemsToSend);
        return request.dialogId;
    }
}
