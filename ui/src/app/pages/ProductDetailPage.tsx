import {Suspense, useEffect, useState} from "react";
import Throbber from "../components/Shared/Throbber.tsx";
import {Params, useParams} from "react-router";
import {getProductDetail} from "../components/Products/ProductService.ts";
import {IProduct} from "../components/Products/IProduct.ts";
import {Subscription} from "rxjs";
import ProductImage from "../components/Products/ProductImage.tsx";

export default function ProductDetailPage() {
    const params: Readonly<Params<string>> = useParams();
    const productSlug: string | undefined = params.slug;
    const [product, setProduct] = useState<IProduct>();

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
            {product ? (
                <Suspense fallback={<Throbber/>}>
                    <section className="flex justify-between mb-4">
                        <h1 className="text-xl font-bold">{product.name}</h1>
                    </section>

                    <section className="flex gap-4">
                        <aside>
                            <ProductImage product={product}/>
                        </aside>

                        <div>
                            <h2 className="font-bold text-lg mb-4">Description</h2>
                            <p>{product.description}</p>
                        </div>
                    </section>
                </Suspense>
            ) : (
                <>Product not found.</>
            )}
        </>
    )
}