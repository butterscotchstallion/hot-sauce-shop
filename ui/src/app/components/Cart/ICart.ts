import {IProduct} from "../Products/IProduct.ts";
import {IUser} from "../User/IUser.ts";

export interface ICart {
    id: number;
    inventory: IProduct;
    user: IUser;
    createdAt: Date;
    updatedAt: Date;
}