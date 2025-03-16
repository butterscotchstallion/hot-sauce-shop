import * as React from "react";
import {RefObject, useEffect, useRef} from "react";
import {Button} from "primereact/button";
import {Toast} from "primereact/toast";
import {RootState} from "../../store.ts";
import {useDispatch, useSelector} from "react-redux";
import {ICart} from "./ICart.ts";
import {IIDQuantityMap, IInitialCartState, setCartItemQuantity} from "./Cart.slice.ts";
import {Sidebar} from "primereact/sidebar";
import {DataTable} from "primereact/datatable";
import {Column} from "primereact/column";
import {confirmPopup, ConfirmPopup} from "primereact/confirmpopup";
import {Dropdown} from "primereact/dropdown";
import {addCartItem} from "./CartService.ts";

export default function CartMenu() {
    const dispatch = useDispatch();
    const [sidebarVisible, setSidebarVisible] = React.useState<boolean>(false);
    const cartState: IInitialCartState = useSelector((state: RootState) => state.cart);
    const idQuantityMap: IIDQuantityMap = useSelector((state: RootState) => state.cart.idQuantityMap);
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const [cartItemsQuantity, setCartItemsQuantity] = React.useState<number>(0);
    const [cartSubtotal, setCartSubtotal] = React.useState<number>(0);

    useEffect(() => {
        const newTotal: number = calculateCartItemsTotal(cartState.items);
        setCartItemsQuantity(newTotal)
        const newSubtotal: number = recalculateSubtotal(cartState.items);
        setCartSubtotal(newSubtotal);
        console.log("Total cart items set to: " + newTotal);
    }, [cartState, idQuantityMap, cartItemsQuantity]);

    function calculateCartItemsTotal(cartItems: ICart[]): number {
        return cartItems.reduce(
            (sum: number, item: ICart) => sum + item.quantity,
            0,
        );
    }

    function recalculateSubtotal(cartItems: ICart[]): number {
        return cartItems.reduce(
            (acc: number, item: ICart) =>
                acc + item.price * item.quantity,
            0,
        );
    }

    const removeCartItemTpl = (cartItem: ICart) => {
        return <Button
            onClick={() => openRemoveCartConfirmation(cartItem)}
            title="Remove cart item"
            severity={"danger"}
            icon="pi pi-trash"
            className="p-button-rounded p-button-text"/>
    };

    const onRemoveCartConfirmed = () => {

    }

    const onRemoveCartRejected = () => {

    }

    const quantityOptions = Array.from(
        {length: 50},
        (_, i) => ({
            label: String(i + 1),
            value: i + 1,
        }),
    );

    const openRemoveCartConfirmation = (event) => {
        confirmPopup({
            target: event.currentTarget,
            message: 'Are you sure you want to remove this cart item?',
            icon: 'pi pi-exclamation-triangle',
            defaultFocus: 'accept',
            accept: onRemoveCartConfirmed,
            reject: onRemoveCartRejected
        });
    };

    function setCartItemQuantityFromMenu(cartItem: ICart, quantity: number) {
        console.log("Setting quantity for cart item " + cartItem.id + " to " + quantity);
        addCartItem({
            inventoryItemId: cartItem.id,
            userId: 1,
            overrideQuantity: true,
            quantity: quantity,
        }).subscribe({
            next: () => {
                toast.current?.show({
                    severity: 'success',
                    summary: 'Success',
                    detail: 'Cart item quantity updated',
                    life: 3000,
                });
                dispatch(setCartItemQuantity({
                    id: cartItem.inventoryItemId,
                    quantity
                }));
            },
            error: (err: string) => {
                toast.current?.show({
                    severity: 'error',
                    summary: 'Error',
                    detail: 'Error updating cart item quantity: ' + err,
                })
            }
        });
    }

    const quantityColTpl = (cartItem: ICart) => {
        return <Dropdown value={cartItem.quantity}
                         onChange={(e) => setCartItemQuantityFromMenu(cartItem, e.value)}
                         options={quantityOptions}
                         optionLabel="label"
                         optionValue="value"
                         placeholder="Select a City" className="w-full md:w-14rem"/>
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
                style={{width: '28rem'}}
                position={"right"}
                visible={sidebarVisible}
                onHide={() => setSidebarVisible(false)}
            >
                <h2 className="text-2xl font-bold">Cart</h2>
                <section className="mt-4 cart-table-area">
                    <DataTable className="w-full" value={cartState.items}>
                        <Column
                            className="w-[40%] max-w-[80px] whitespace-nowrap overflow-hidden text-ellipsis"
                            field="name"
                            header="Name"></Column>
                        <Column className="w-[20%]" field="price" header="Price"></Column>
                        <Column className="w-[5%]" body={quantityColTpl} header="Quantity"></Column>
                        <Column className="w-[5%]" header="Remove" body={removeCartItemTpl}/>
                    </DataTable>
                </section>

                <section className="mt-4 mb-4 flex justify-between">
                    <h3 className="text-xl font-bold">Total: ${cartSubtotal.toFixed(2)}</h3>
                    <Button label="Checkout" icon="pi pi-shopping-cart" className="p-button-rounded"/>
                </section>
            </Sidebar>

            <Toast ref={toast}/>
            <ConfirmPopup/>
        </>
    )
}