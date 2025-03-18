import {useEffect, useState} from "react";
import {IProduct} from "../components/Products/IProduct.ts";
import AdminInventoryItemForm from "../components/Admin/AdminInventoryItemForm.tsx";
import {Params, useParams} from "react-router";
import {Subscription} from "rxjs";
import {getProductDetail} from "../components/Products/ProductService.ts";

export default function AdminInventoryPage() {
    const [product, setProduct] = useState<IProduct | undefined>();
    const params: Readonly<Params<string>> = useParams();
    const productSlug: string | undefined = params.slug;

    useEffect(() => {
        if (productSlug) {
            const product$: Subscription = getProductDetail(productSlug).subscribe({
                next: (productDetail: IProduct) => {
                    setProduct(productDetail);
                },
                error: (err) => {
                    console.error(err);
                }
            })
            return () => {
                product$.unsubscribe();
            }
        }
    }, [productSlug])

    return (
        <>
            <h1 className="font-bold text-2xl mb-4">Admin - Edit Product</h1>

            <section className="flex">
                <AdminInventoryItemForm product={product}/>
            </section>
        </>
    )
}