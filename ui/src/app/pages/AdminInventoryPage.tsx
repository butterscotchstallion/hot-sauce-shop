import {useEffect, useState} from "react";
import {IProduct} from "../components/Products/IProduct.ts";
import AdminInventoryItemForm from "../components/Admin/AdminInventoryItemForm.tsx";
import {Params, useParams} from "react-router";
import {Subscription} from "rxjs";
import {getProductDetail} from "../components/Products/ProductService.ts";

export interface IAdminInventoryPageProps {
    isNewProduct: boolean;
}

export default function AdminInventoryPage(props: IAdminInventoryPageProps) {
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
            });
            return () => {
                product$.unsubscribe();
            }
        }
    }, [productSlug])

    return (
        <>
            <section className="flex">
                <AdminInventoryItemForm isNewProduct={props.isNewProduct} product={product}/>
            </section>
        </>
    )
}