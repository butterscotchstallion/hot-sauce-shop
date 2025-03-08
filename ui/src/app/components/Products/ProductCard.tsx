import {IProduct} from "./IProduct.ts";
import {Card} from "primereact/card";

interface IProductCardProps {
    product: IProduct
}

export default function ProductCard(props: IProductCardProps) {
    return (
        <div className="w-[200px]">
            <Card><img
                className="text-center mx-auto"
                src="/images/hot-pepper.png"
                alt={props.product.shortDescription}
            /></Card>
            <Card title={props.product.name}>
                <p>{props.product.shortDescription}</p>
            </Card>
        </div>
    )
}