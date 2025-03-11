import {IProduct} from "../Products/IProduct.ts";
import {IUser} from "../User/IUser.ts";

export interface ICart {
    inventory: IProduct;
    user: IUser;
    createdAt: Date;
    updatedAt: Date;
}