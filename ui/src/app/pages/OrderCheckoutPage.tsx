import {CartItemsDataTable} from "../components/Cart/CartItemsDataTable.tsx";
import {DataTable} from "primereact/datatable";
import {IDeliveryOption} from "../components/Orders/IDeliveryOption.ts";
import {Column} from "primereact/column";
import {useState} from "react";

export function OrderCheckoutPage() {
    const deliveryOptions: IDeliveryOption[] = [
        {name: "Two Day", price: 9.99, deliveryDate: new Date("2025-04-21").toLocaleDateString()},
        {name: "Three Day", price: 4.99, deliveryDate: new Date("2025-04-22").toLocaleDateString()},
    ];
    const [selectedDeliveryOption, setSelectedDeliveryOption] = useState<IDeliveryOption>(deliveryOptions[0]);
    return (
        <>
            <h1 className="text-2xl font-bold mb-4">Checkout</h1>

            <section className="flex justify-between gap-4">
                <CartItemsDataTable/>

                <aside className="w-2/3">
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
                </aside>
            </section>
        </>
    )
}