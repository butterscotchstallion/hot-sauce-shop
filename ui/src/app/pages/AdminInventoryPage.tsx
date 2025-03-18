import ProductAutocomplete from "../components/Products/ProductAutocomplete.tsx";
import {useState} from "react";
import {IProduct} from "../components/Products/IProduct.ts";
import AdminInventoryItemForm from "../components/Admin/AdminInventoryItemForm.tsx";

export default function AdminInventoryPage() {
    const newProduct: IProduct = {};
    const [product, setProduct] = useState<IProduct | null>(newProduct);

    function setValue(key: string, value: string) {
        if (product) {
            product[key] = value;
        }
    }

    return (
        <>
            <section className="flex mb-4 gap-4">
                <ProductAutocomplete/>
            </section>

            <section className="flex max-w-1/2">
                <AdminInventoryItemForm/>
            </section>
        </>
    )
}