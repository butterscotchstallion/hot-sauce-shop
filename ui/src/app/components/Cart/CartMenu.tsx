import * as React from "react";
import {useEffect, useRef} from "react";
import {Button} from "primereact/button";
import {Menu} from "primereact/menu";
import {getCartItems} from "./CartService.ts";
import {ICart} from "./ICart.ts";
import {Toast} from "primereact/toast";

export default function CartMenu() {
    const toast = useRef<Toast>(null);
    const cartMenu = React.useRef<Menu>(null);
    const [numCartItems, setNumCartItems] = React.useState<number>(0);
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
        const cartItems$ = getCartItems().subscribe({
            next: (cartItems: ICart[]) => {
                setNumCartItems(cartItems.length);
            },
            error: () => {
                toast.current?.show({severity: 'error', summary: 'Error', detail: 'Error loading cart items'});
            }
        })

        return () => {
            cartItems$.unsubscribe();
        }
    }, []);

    return (
        <>
            <Menu model={items} ref={cartMenu} popup id="popup_menu_right" popupAlignment="right"/>
            <Button
                label="Cart"
                icon="pi pi-shopping-cart"
                className="mr-2"
                badge={numCartItems.toString()}
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