import {IProduct} from "./IProduct.ts";
import {Card} from "primereact/card";
import {Button} from "primereact/button";
import {NavLink} from "react-router";
import ProductImage from "./ProductImage.tsx";
import SpiceRating from "./SpiceRating.tsx";
import {addCartItem} from "../Cart/CartService.ts";
import {Toast} from "primereact/toast";
import {RefObject, useState} from "react";

interface IProductCardProps {
    product: IProduct,
    toast: RefObject<Toast | null>
}

export default function ProductCard(props: IProductCardProps) {
    const [isAddingToCart, setIsAddingToCart] = useState<boolean>(false);

    function addToCart(product: IProduct) {
        setIsAddingToCart(true);
        addCartItem({
            inventoryItemId: product.id,
            userId: 1,
            overrideQuantity: false,
            quantity: 1,
        }).subscribe({
            next: () => {
                props.toast.current?.show({
                    severity: 'success',
                    summary: 'Success',
                    detail: product.name + ' added to cart',
                    life: 3000,
                });
                setIsAddingToCart(false);
            },
            error: (err) => {
                props.toast.current?.show({
                    severity: 'error',
                    summary: 'Error',
                    detail: 'Error adding ' + product.name + ' to cart: ' + err,
                    life: 3000,
                });
                setIsAddingToCart(false);
            }
        })
    }

    return (
        <div className="w-[260px]">
            <NavLink to={props.product.slug}>
                <ProductImage product={props.product}/>
            </NavLink>
            <Card>
                <NavLink to={props.product.slug}>
                    <h2 className="max-w-60 whitespace-nowrap text-2xl font-bold overflow-hidden text-ellipsis"
                        title={props.product.name}>
                        {props.product.name}
                    </h2>
                    <p className="pb-4">{props.product.shortDescription}</p>

                    <section className="mb-4">
                        <SpiceRating rating={props.product.spiceRating}/>
                    </section>
                </NavLink>

                <div className="flex justify-between">
                    <span className="text-green-200 font-bold pt-4">${props.product.price.toFixed(2)}</span>
                    <Button
                        label="Add"
                        icon="pi pi-shopping-cart"
                        disabled={isAddingToCart}
                        onClick={() => addToCart(props.product)}/>
                </div>
            </Card>
        </div>
    )
}