import {inject, injectable} from "tsyringe";
import {ISptCommand} from "@spt/helpers/Dialogue/Commando/SptCommands/ISptCommand";
import {ItemHelper} from "@spt/helpers/ItemHelper";
import {ISendMessageRequest} from "@spt/models/eft/dialog/ISendMessageRequest";
import {ISptProfile} from "@spt/models/eft/profile/ISptProfile";
import {IUserDialogInfo} from "@spt/models/eft/profile/IUserDialogInfo";
import {MailSendService} from "@spt/services/MailSendService";
import type {ICloner} from "@spt/utils/cloners/ICloner";
import {SaveServer} from '@spt/servers/SaveServer';
import {IItem} from "@spt/models/eft/common/tables/IItem";


interface InventoryItemHash {
    byItemId: Record<string, IItem>;
    byParentId: Record<string, IItem[]>;
}

@injectable()
export class GiveStashItemSptCommand implements ISptCommand {
    private static commandRegex = /^spt give-user-stash-item ((([a-z]{2,5}) )?"(.+)"|\w+)$/;

    public constructor(
        @inject("ItemHelper") protected itemHelper: ItemHelper,
        @inject("MailSendService") protected mailSendService: MailSendService,
        @inject("PrimaryCloner") protected cloner: ICloner,
        @inject("SaveServer") protected saveServer: SaveServer,
    ) {
    }

    public getCommand(): string {
        return "give-user-stash-item";
    }

    public getCommandHelp(): string {
        return 'spt give-user-stash-item\n========\nSends items to the player through the message system.\n\n\tspt give-user-stash-item [stashItem.Id]';
    }

    public performAction(commandHandler: IUserDialogInfo, sessionId: string, request: ISendMessageRequest): string {
        if (!GiveStashItemSptCommand.commandRegex.test(request.text)) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                'Invalid use of give command. Use "help" for more information.',
            );
            return request.dialogId;
        }

        const result = GiveStashItemSptCommand.commandRegex.exec(request.text);

        if (!result[1]) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                `Invalid use of give command. Use "help" for more information.`,
            );
            return request.dialogId;
        }

        const profile: ISptProfile = this.saveServer.getProfiles()[sessionId];

        const itemId = result[1];
        const inventoryItemHash = this.getInventoryItemHash(profile.characters.pmc.Inventory.items)

        const itemToAdd = inventoryItemHash.byItemId[itemId]
        if (!itemToAdd) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                `Couldn't find item with id: ${itemId}`,
            );
            return request.dialogId;
        }


        const checkedItem = this.itemHelper.getItem(itemToAdd._tpl);
        if (!checkedItem[0]) {
            this.mailSendService.sendUserMessageToPlayer(
                sessionId,
                commandHandler,
                "That item could not be found. Please refine your request and try again.",
            );
            return request.dialogId;
        }
        // TODO
        // const quantity = checkedItem[1]._props.StackMaxSize

        const allChild = this.getAllDescendantsIncludingSelf(itemId, inventoryItemHash)
        let itemsToSend = this.cloner.clone(allChild);
        // console.log('original allChild', allChild);

        itemsToSend = this.itemHelper.replaceIDs(itemsToSend);
        this.itemHelper.setFoundInRaid(itemsToSend);
        this.mailSendService.sendSystemMessageToPlayer(sessionId, "SPT GIVE", itemsToSend);

        // console.log('itemsToSend final', itemsToSend);

        return request.dialogId;
    }

    private getInventoryItemHash(inventoryItem: IItem[]): InventoryItemHash {
        const inventoryItemHash: InventoryItemHash = { byItemId: {}, byParentId: {} };
        for (const item of inventoryItem) {
            inventoryItemHash.byItemId[item._id] = item;

            if (!("parentId" in item)) {
                continue;
            }

            if (!(item.parentId in inventoryItemHash.byParentId)) {
                inventoryItemHash.byParentId[item.parentId] = [];
            }
            inventoryItemHash.byParentId[item.parentId].push(item);
        }
        return inventoryItemHash;
    }

    private getAllDescendantsIncludingSelf(parentId: string, hash: InventoryItemHash): IItem[] {
        const result: IItem[] = [];

        const self = hash.byItemId[parentId];
        if (self) {
            result.push(self); // Add the original item first
        }

        const directChildren = hash.byParentId[parentId] || [];

        for (const child of directChildren) {
            result.push(...this.getAllDescendantsIncludingSelf(child._id, hash));
        }

        return result;
    }
}
