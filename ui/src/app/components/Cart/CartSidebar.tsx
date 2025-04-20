import * as React from "react";
import {RefObject, useEffect, useRef} from "react";
import {Button} from "primereact/button";
import {Toast} from "primereact/toast";
import {RootState} from "../../store.ts";
import {useSelector} from "react-redux";
import {ICart} from "./ICart.ts";
import {IIDQuantityMap, IInitialCartState} from "./Cart.slice.ts";
import {Sidebar} from "primereact/sidebar";
import {ConfirmDialog} from "primereact/confirmdialog";
import {NavigateFunction, useNavigate} from "react-router";
import {CartItemsDataTable} from "./CartItemsDataTable.tsx";

export default function CartSidebar() {
    const [sidebarVisible, setSidebarVisible] = React.useState<boolean>(false);
    const cartState: IInitialCartState = useSelector((state: RootState) => state.cart);
    const idQuantityMap: IIDQuantityMap = useSelector((state: RootState) => state.cart.idQuantityMap);
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const [cartItemsQuantity, setCartItemsQuantity] = React.useState<number>(0);
    const navigate: NavigateFunction = useNavigate();

    useEffect(() => {
        const newTotal: number = calculateCartItemsTotal(cartState.items);
        setCartItemsQuantity(newTotal)
    }, [cartState, idQuantityMap, cartItemsQuantity]);

    function calculateCartItemsTotal(cartItems: ICart[]): number {
        return cartItems.reduce(
            (sum: number, item: ICart) => sum + item.quantity,
            0,
        );
    }

    const goToCheckOut = () => {
        navigate('/orders/checkout');
        setSidebarVisible(false);
    }
    return (
        <>
            <Button
                label="Cart"
                icon="pi pi-shopping-cart"
                className="mr-2"
                badge={cartItemsQuantity.toString()}
                onClick={() => setSidebarVisible(true)}
                aria-controls="popup_menu_right"
                aria-haspopup/>

            <Sidebar
                style={{width: '33rem'}}
                position={"right"}
                visible={sidebarVisible}
                onHide={() => setSidebarVisible(false)}
            >
                <h2 className="text-2xl font-bold">Cart</h2>
                <section className="mt-4 cart-table-area">
                    <CartItemsDataTable/>
                </section>

                <section className="mt-4 mb-4 flex justify-between">
                    <Button
                        onClick={() => goToCheckOut()}
                        label="Checkout"
                        icon="pi pi-shopping-cart"
                        className="p-button-rounded"/>
                </section>
            </Sidebar>

            <Toast ref={toast}/>
            <ConfirmDialog/>
        </>
    )
}