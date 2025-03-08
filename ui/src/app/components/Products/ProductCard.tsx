import {IProduct} from "./IProduct.ts";
import {Card} from "primereact/card";

interface IProductCardProps {
    product: IProduct
}

export default function ProductCard(props: IProductCardProps) {
    return (
        <div className="w-[200px]">
            <div className="bg-orange-400 h-[200px]"></div>
            <Card title={props.product.name}>
                {props.product.shortDescription}
            </Card>
        </div>
    )
}