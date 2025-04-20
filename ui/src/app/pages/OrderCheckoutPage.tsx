import {CartItemsDataTable} from "../components/Cart/CartItemsDataTable.tsx";

export function OrderCheckoutPage() {
    return (
        <>
            <h1 className="text-2xl font-bold mb-4">Checkout</h1>

            <CartItemsDataTable/>
        </>
    )
}