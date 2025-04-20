import {InputText} from "primereact/inputtext";
import {DataTable, DataTableFilterMeta} from "primereact/datatable";
import {Column} from "primereact/column";
import {ICart} from "./ICart.ts";
import * as React from "react";
import {ReactElement, RefObject, useEffect, useRef, useState} from "react";
import {Dropdown, DropdownChangeEvent} from "primereact/dropdown";
import {FilterMatchMode} from "primereact/api";
import {addCartItem, deleteCartItem} from "./CartService.ts";
import {cartItemRemoved, IIDQuantityMap, IInitialCartState, setCartItemQuantity} from "./Cart.slice.ts";
import {Toast} from "primereact/toast";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {Button} from "primereact/button";
import {confirmDialog} from "primereact/confirmdialog";

export function CartItemsDataTable() {
    const dispatch = useDispatch();
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const [cartItemsQuantity, setCartItemsQuantity] = React.useState<number>(0);
    const [cartSubtotal, setCartSubtotal] = React.useState<number>(0);
    const cartState: IInitialCartState = useSelector((state: RootState) => state.cart);
    const idQuantityMap: IIDQuantityMap = useSelector((state: RootState) => state.cart.idQuantityMap);
    const [filters, setFilters] = useState<DataTableFilterMeta>({
        global: {value: null, matchMode: FilterMatchMode.CONTAINS},
        name: {value: null, matchMode: FilterMatchMode.STARTS_WITH},
    });
    const [globalFilterValue, setGlobalFilterValue] = useState<string>('');
    const quantityOptions = Array.from(
        {length: 50},
        (_, i) => ({
            label: String(i + 1),
            value: i + 1,
        }),
    );
    let deleteCartItemInventoryId: number = 0;

    const quantityColTpl = (cartItem: ICart): ReactElement => {
        return <Dropdown value={cartItem.quantity}
                         onChange={(e: DropdownChangeEvent) => setCartItemQuantityFromMenu(cartItem, e.value)}
                         options={quantityOptions}
                         optionLabel="label"
                         optionValue="value"
                         className="w-full md:w-14rem"/>
    }

    const onGlobalFilterChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        const value: string = e.target.value;
        const _filters = {...filters};

        // @ts-expect-error value is known to exist
        _filters['global'].value = value;

        setFilters(_filters);
        setGlobalFilterValue(value);
    };
    const sidebarCartNameTpl = (cartItem: ICart): ReactElement => {
        return <div className='flex items-center gap-2'>{cartItem.name}</div>;
    }
    const priceColumnTpl = (cartItem: ICart): ReactElement => {
        return <>{cartItem.price.toFixed(2)}</>
    }

    function setCartItemQuantityFromMenu(cartItem: ICart, quantity: number) {
        console.log("Setting quantity for cart item " + cartItem.id + " to " + quantity);
        addCartItem({
            inventoryItemId: cartItem.inventoryItemId,
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
                setCartSubtotal(recalculateSubtotal(cartState.items));
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

    const onRemoveCartConfirmed = () => {
        if (deleteCartItemInventoryId > 0) {
            deleteCartItem({
                inventoryItemId: deleteCartItemInventoryId,
            }).subscribe({
                next: () => {
                    toast.current?.show({
                        severity: 'success',
                        summary: 'Success',
                        detail: 'Cart item removed',
                        life: 3000,
                    });
                    dispatch(cartItemRemoved({
                        id: deleteCartItemInventoryId,
                    }));
                },
                error: (message: string) => {
                    toast.current?.show({
                        severity: 'error',
                        summary: 'Error',
                        detail: 'Error removing cart item: ' + message,
                        life: 3000,
                    })
                }
            });
        } else {
            console.error("No cart item to remove");
        }
    }

    const onRemoveCartRejected = () => {
        deleteCartItemInventoryId = 0;
    }

    const openRemoveCartConfirmation = (cartItem: ICart) => {
        deleteCartItemInventoryId = cartItem.inventoryItemId;
        confirmDialog({
            header: "Remove Cart Item",
            message: 'Are you sure you want to remove ' + cartItem.name + ' (' + cartItem.quantity + ') from your cart?',
            icon: 'pi pi-exclamation-triangle',
            defaultFocus: 'accept',
            accept: onRemoveCartConfirmed,
            reject: onRemoveCartRejected
        });
    };

    const removeCartItemTpl = (cartItem: ICart) => {
        return <Button
            onClick={() => openRemoveCartConfirmation(cartItem)}
            title="Remove cart item"
            severity={"danger"}
            icon="pi pi-trash"
            className="p-button-rounded p-button-text"/>
    };

    useEffect(() => {
        const newTotal: number = calculateCartItemsTotal(cartState.items);
        setCartItemsQuantity(newTotal)
        const newSubtotal: number = recalculateSubtotal(cartState.items);
        setCartSubtotal(newSubtotal);
    }, [cartState, idQuantityMap, cartItemsQuantity]);

    return (
        <>
            <section className="flex flex-col gap-4 w-full">
                <section className="mb-4 w-full">
                    <InputText value={globalFilterValue}
                               onChange={onGlobalFilterChange}
                               placeholder="Filter cart items"/>

                    <DataTable className="mt-4 w-full"
                               value={cartState.items}
                               filters={filters}
                               globalFilterFields={['name']}>
                        <Column
                            sortable
                            filterField="name"
                            filterMatchMode="contains"
                            className="w-[40%] max-w-[80px] whitespace-nowrap overflow-hidden text-ellipsis"
                            field="name"
                            header="Name"
                            body={sidebarCartNameTpl}></Column>
                        <Column sortable
                                className="w-[20%]"
                                field="price"
                                header="Price"
                                body={priceColumnTpl}></Column>
                        <Column sortable
                                className="w-[5%]"
                                body={quantityColTpl}
                                field="quantity"
                                header="Quantity"></Column>
                        <Column className="w-[5%]"
                                header="Remove"
                                body={removeCartItemTpl}/>
                    </DataTable>
                </section>
                <h3 className="text-xl font-bold">Total: ${cartSubtotal.toFixed(2)}</h3>
            </section>

            <Toast ref={toast}/>
        </>
    )
}