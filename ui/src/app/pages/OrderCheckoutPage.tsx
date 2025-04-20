import {CartItemsDataTable} from "../components/Cart/CartItemsDataTable.tsx";
import {DataTable} from "primereact/datatable";
import {IDeliveryOption} from "../components/Orders/IDeliveryOption.ts";
import {Column} from "primereact/column";
import {useEffect, useState} from "react";
import {Button} from "primereact/button";
import {useSelector} from "react-redux";
import {RootState} from "../store.ts";

interface IOrderTotalItems {
    name: string;
    price: number;
}

export function OrderCheckoutPage() {
    const cartSubtotal: number = useSelector((state: RootState) => state.cart.cartSubtotal);
    const deliveryOptions: IDeliveryOption[] = [
        {name: "Two Day", price: 9.99, deliveryDate: new Date("2025-04-21").toLocaleDateString()},
        {name: "Three Day", price: 4.99, deliveryDate: new Date("2025-04-22").toLocaleDateString()},
    ];
    const [selectedDeliveryOption, setSelectedDeliveryOption] = useState<IDeliveryOption>(deliveryOptions[0]);
    const [orderTotal, setOrderTotal] = useState<string>(cartSubtotal.toFixed(2));
    const orderTotalItems: IOrderTotalItems[] = [
        {name: "Items", price: cartSubtotal},
        {name: "Shipping & Handling", price: selectedDeliveryOption.price},
        {name: "Estimated Taxes", price: 0},
    ];
    const orderTotalPriceFormatted = (row) => {
        return `$${row.price}`;
    }
    useEffect(() => {
        const newOrderTotal: string = (parseFloat(orderTotal) + selectedDeliveryOption.price).toFixed(2);
        setOrderTotal(newOrderTotal);
    }, [cartSubtotal, selectedDeliveryOption]);

    return (
        <>
            <h1 className="text-2xl font-bold mb-4">Checkout</h1>

            <section className="flex justify-between gap-4">
                <CartItemsDataTable hideSubtotal={true}/>

                <aside className="w-2/3">
                    <section className="flex flex-col gap-4">
                        <section>
                            <h4 className="text-xl font-bold mb-2">Delivery Options</h4>
                            <DataTable
                                value={deliveryOptions}
                                selection={selectedDeliveryOption}
                                onSelectionChange={(e) => setSelectedDeliveryOption(e.value)}>
                                <Column selectionMode="single" headerStyle={{width: '3rem'}}></Column>
                                <Column field="name" header="Name"/>
                                <Column field="price" header="Price"/>
                                <Column field="deliveryDate" header="Delivery Date"/>
                            </DataTable>
                        </section>

                        <section>
                            <DataTable value={orderTotalItems}>
                                <Column field="name"/>
                                <Column field="price" body={orderTotalPriceFormatted}/>
                            </DataTable>
                        </section>

                        <section>
                            <section className="flex justify-between gap-4">
                                <h4 className="text-xl font-bold mb-2">
                                    Order total: ${orderTotal}
                                </h4>
                                <Button label="Place Order" icon="pi pi-cart-arrow-down"/>
                            </section>
                        </section>
                    </section>
                </aside>
            </section>
        </>
    )
}