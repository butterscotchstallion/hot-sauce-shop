export interface IAddCartItemRequest {
    inventoryItemId: number;
    userId: number;
    quantity: number;
    overrideQuantity: boolean;
}