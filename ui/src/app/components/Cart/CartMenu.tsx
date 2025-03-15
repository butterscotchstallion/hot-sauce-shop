import * as React from "react";
import {RefObject, useEffect, useRef} from "react";
import {Button} from "primereact/button";
import {Menu} from "primereact/menu";
import {Toast} from "primereact/toast";
import {RootState} from "../../store.ts";
import {useSelector} from "react-redux";
import {ICart} from "./ICart.ts";
import {IIDQuantityMap, IInitialCartState} from "./Cart.slice.ts";

export default function CartMenu() {
    const cartState: IInitialCartState = useSelector((state: RootState) => state.cart);
    const idQuantityMap: IIDQuantityMap = useSelector((state: RootState) => state.cart.idQuantityMap);
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const cartMenu: RefObject<Menu | null> = React.useRef<Menu>(null);
    const [cartItemsQuantity, setCartItemsQuantity] = React.useState<number>(0);
    const [cartSubtotal, setCartSubtotal] = React.useState<number>(0);
    const items = [
        {
            label: 'Cart',
            items: [
                {
                    label: 'Refresh',
                    icon: 'pi pi-refresh'
                },
                {
                    label: 'Export',
                    icon: 'pi pi-upload'
                }
            ]
        }
    ];

    useEffect(() => {
        const newTotal: number = calculateCartItemsTotal(cartState.items);
        setCartItemsQuantity(newTotal);
        console.log("Total cart items set to: " + newTotal);
    }, [cartState, idQuantityMap, cartItemsQuantity]);

    function calculateCartItemsTotal(cartItems: ICart[]): number {
        return cartItems.reduce(
            (sum: number, item: ICart) => sum + item.quantity,
            0,
        );
    }

    function recalculateSubtotal(cartItems: ICart[]) {
        setCartSubtotal(cartItems.reduce(
            (acc: number, item: ICart) =>
                acc + item.price * item.quantity,
            0,
        ));
    }

    return (
        <>
            <Menu model={items} ref={cartMenu} popup id="popup_menu_right" popupAlignment="right"/>
            <Button
                label="Cart"
                icon="pi pi-shopping-cart"
                className="mr-2"
                badge={cartItemsQuantity.toString()}
                onClick={(event) => {
                    if (cartMenu?.current) {
                        return cartMenu.current.toggle(event);
                    }
                }}
                aria-controls="popup_menu_right"
                aria-haspopup/>
            <Toast ref={toast}/>
        </>
    )
}