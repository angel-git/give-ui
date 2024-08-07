import { TradeCallbacks } from '@spt/callbacks/TradeCallbacks';
import { HandledRoute, ItemEventRouterDefinition } from '@spt/di/Router';
import { IPmcData } from '@spt/models/eft/common/IPmcData';
import { IItemEventRouterResponse } from '@spt/models/eft/itemEvent/IItemEventRouterResponse';
export declare class TradeItemEventRouter extends ItemEventRouterDefinition {
  protected tradeCallbacks: TradeCallbacks;
  constructor(tradeCallbacks: TradeCallbacks);
  getHandledRoutes(): HandledRoute[];
  handleItemEvent(
    url: string,
    pmcData: IPmcData,
    body: any,
    sessionID: string,
  ): Promise<IItemEventRouterResponse>;
}
