import {ManipulateType} from "dayjs";

export interface IShippingOption {
    id: string;
    name: string;
    description: string;
    timeToShipUnit: ManipulateType | undefined;
    timeToShipUnitQuantity: number;
    price: number;
    deliveryDate: string;
}