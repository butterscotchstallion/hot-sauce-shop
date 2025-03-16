import {IUser} from "../User/IUser.ts";

export interface ICart {
    id: number;
    inventoryItemId: number;
    name: string;
    price: number;
    quantity: number;
    user: IUser;
    createdAt: Date;
    updatedAt: Date;
}