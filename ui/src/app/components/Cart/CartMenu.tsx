import * as React from "react";
import {useEffect, useRef} from "react";
import {Button} from "primereact/button";
import {Menu} from "primereact/menu";
import {Toast} from "primereact/toast";
import {RootState} from "../../store.ts";
import {useSelector} from "react-redux";
import {ICart} from "./ICart.ts";

export default function CartMenu() {
    const cart = useSelector((state: RootState) => state.cart);
    const toast = useRef<Toast>(null);
    const cartMenu = React.useRef<Menu>(null);
    const [cartItemsQuantity, setCartItemsQuantity] = React.useState<number>(0);
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
        calculateCartItemsTotal(cart.items);
    }, [cart]);

    function calculateCartItemsTotal(cartItems: ICart[]) {
        setCartItemsQuantity(cartItems.reduce(
            (sum, item) => sum + item.quantity,
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