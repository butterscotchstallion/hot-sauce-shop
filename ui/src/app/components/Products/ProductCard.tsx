import {IProduct} from "./IProduct.ts";
import {Card} from "primereact/card";
import {Button} from "primereact/button";
import {NavLink} from "react-router";
import ProductImage from "./ProductImage.tsx";
import SpiceRating from "./SpiceRating.tsx";
import {Toast} from "primereact/toast";
import {RefObject, useEffect, useMemo, useState} from "react";
import {ICartMap} from "../Cart/ICartMap.ts";
import {RootState} from "../../store.ts";
import {useDispatch, useSelector} from "react-redux";
import {addCartItem} from "../Cart/CartService.ts";
import {cartItemAdded} from "../Cart/Cart.slice.ts";

interface IProductCardProps {
    product: IProduct,
    toast: RefObject<Toast | null>
}

export default function ProductCard(props: IProductCardProps) {
    const cartMap: ICartMap = useSelector((state: RootState) => {
        return state.cart.nameQuantityMap;
    });
    const dispatch = useDispatch();
    const [isAddingToCart, setIsAddingToCart] = useState<boolean>(false);
    const [productQuantityMap, setProductQuantityMap] = useState<Map<string, number>>(
        useMemo(() => new Map<string, number>(), [])
    );

    useEffect(() => {
        let qty: number = 0;
        if (cartMap && typeof cartMap[props.product.name] === "number") {
            qty = cartMap[props.product.name];
        }
        const productQuantityMapCopy = productQuantityMap;
        productQuantityMapCopy.set(props.product.name, qty);
        setProductQuantityMap(productQuantityMapCopy);
    }, [cartMap, productQuantityMap, props.product.name]);

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
                dispatch(cartItemAdded(product));
            },
            error: (err: string) => {
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
                        badge={props.product.name in cartMap ? cartMap[props.product.name].toString() : '0'}
                        disabled={isAddingToCart}
                        onClick={() => addToCart(props.product)}/>
                </div>
            </Card>
        </div>
    )
}