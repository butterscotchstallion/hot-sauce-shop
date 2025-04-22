export interface IShippingOption {
    id: string;
    name: string;
    description: string;
    timeToShipUnit: string;
    timeToShipUnitQuantity: number;
    price: number;
    deliveryDate: string;
}