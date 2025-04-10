import {InputText} from "primereact/inputtext";
import {DataTable} from "primereact/datatable";
import {Column} from "primereact/column";
import {ICart} from "./ICart.ts";

interface ICartItemsDataTableProps {
    cartItems: ICart[];
}

export function CartItemsDataTable(props: ICartItemsDataTableProps) {

    return (
        <>
            <InputText value={globalFilterValue}
                       onChange={onGlobalFilterChange}
                       placeholder="Filter cart items"/>
            <DataTable className="w-full"
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
        </>
    )
}