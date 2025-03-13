import {IProduct} from "./IProduct.ts";
import {Card} from "primereact/card";
import {Button} from "primereact/button";
import {NavLink} from "react-router";
import ProductImage from "./ProductImage.tsx";
import SpiceRating from "./SpiceRating.tsx";
import {addCartItem} from "../Cart/CartService.ts";
import {useToastContext} from "../Shared/ToastContext.tsx";

interface IProductCardProps {
    product: IProduct
}

export default function ProductCard(props: IProductCardProps) {
    const {showToast} = useToastContext();

    function addToCart(inventoryItemId: number) {
        addCartItem({
            inventoryItemId,
            userId: 1,
            overrideQuantity: false,
            quantity: 1,
        }).subscribe({
            next: () => {
                showToast({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'Item added to cart',
                    life: 3000,
                });
            },
            error: (err) => {
                showToast({
                    severity: 'error',
                    summary: 'Error',
                    detail: 'Error adding item added to cart: ' + err,
                    life: 3000,
                });
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
                        onClick={() => addToCart(props.product.id)}/>
                </div>
            </Card>
        </div>
    )
}