import {IUser} from "../User/types/IUser.ts";

export interface ICart {
    id: number;
    inventoryItemId: number;
    inventoryItemSlug: string;
    inventoryItemShortDescription: string;
    name: string;
    price: number;
    quantity: number;
    user: IUser;
    createdAt: Date;
    updatedAt: Date;
}