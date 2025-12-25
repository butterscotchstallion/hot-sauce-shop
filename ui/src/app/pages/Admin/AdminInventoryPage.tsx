import {useEffect, useState} from "react";
import {IProduct} from "../../components/Products/types/IProduct.ts";
import AdminInventoryItemForm from "../../components/Admin/AdminInventoryItemForm.tsx";
import {Params, useParams} from "react-router";
import {Subscription} from "rxjs";
import {getProductDetail} from "../../components/Products/ProductService.ts";
import {ITag} from "../../components/Tag/ITag.ts";
import {IProductDetail} from "../../components/Products/types/IProductDetail.ts";
import {getTags} from "../../components/Tag/TagService.ts";

export interface IAdminInventoryPageProps {
    isNewProduct: boolean;
}

export default function AdminInventoryPage(props: IAdminInventoryPageProps) {
    const [product, setProduct] = useState<IProduct | undefined>();
    const [productTags, setProductTags] = useState<ITag[]>([]);
    const [availableTags, setAvailableTags] = useState<ITag[]>([]);
    const params: Readonly<Params<string>> = useParams();
    const productSlug: string | undefined = params.slug;

    useEffect(() => {
        if (productSlug) {
            const product$: Subscription = getProductDetail(productSlug).subscribe({
                next: (productDetail: IProductDetail) => {
                    setProduct(productDetail.product);
                    setProductTags(productDetail.tags);
                },
                error: (err) => {
                    console.error(err);
                }
            });
            const tags$: Subscription = getTags().subscribe({
                next: (tags: ITag[]) => setAvailableTags(tags),
                error: (err) => {
                    console.error(err);
                }
            });
            return () => {
                product$.unsubscribe();
                tags$.unsubscribe();
            }
        }
    }, [productSlug])

    return (
        <>
            <section className="flex">
                <AdminInventoryItemForm
                    isNewProduct={props.isNewProduct}
                    product={product}
                    productTags={productTags}
                    availableTags={availableTags}
                />
            </section>
        </>
    )
}