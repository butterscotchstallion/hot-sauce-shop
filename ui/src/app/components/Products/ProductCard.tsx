import {IProduct} from "./IProduct.ts";
import {Card} from "primereact/card";
import {Button} from "primereact/button";

interface IProductCardProps {
    product: IProduct
}

export default function ProductCard(props: IProductCardProps) {
    return (
        <div className="w-[260px]">
            <Card><img
                className="text-center mx-auto"
                src="/images/hot-pepper.png"
                alt={props.product.shortDescription}
            /></Card>
            <Card title={props.product.name}>
                <p className="pb-4">{props.product.shortDescription}</p>

                <div className="flex justify-between">
                    <span className="text-green-200 font-bold">${props.product.price.toFixed(2)}</span>
                    <Button label="Add" icon="pi pi-shopping-cart" className="p-button-success"/>
                </div>
            </Card>
        </div>
    )
}