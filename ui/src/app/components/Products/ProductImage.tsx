import {IProduct} from "./IProduct.ts";
import {Card} from "primereact/card";
import {ReactElement} from "react";

interface IProductImageProps {
    product: IProduct;
}

export default function ProductImage(props: IProductImageProps): ReactElement {
    return (
        <Card><img
            width="128"
            height="128"
            className="text-center mx-auto"
            src="/images/hot-pepper.png"
            alt={props.product.shortDescription}
        /></Card>
    )
}