import {IProduct} from "./IProduct.ts";

export interface IProductsResults {
    inventory: IProduct[];
    total: number;
}