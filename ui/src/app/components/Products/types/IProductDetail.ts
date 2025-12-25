import {IProduct} from "./IProduct.ts";
import {ITag} from "../../Tag/ITag.ts";

export interface IProductDetail {
    product: IProduct;
    tags: ITag[];
}