import {Suspense, useEffect, useState} from "react";
import Throbber from "../components/Shared/Throbber.tsx";
import {Params, useParams} from "react-router";
import {getProductDetail} from "../components/Products/ProductService.ts";
import {IProduct} from "../components/Products/IProduct.ts";
import {Subscription} from "rxjs";
import ProductImage from "../components/Products/ProductImage.tsx";
import {Button} from "primereact/button";
import SpiceRating from "../components/Products/SpiceRating.tsx";
import ReviewCard from "../components/Reviews/ReviewCard.tsx";

export default function ProductDetailPage() {
    const params: Readonly<Params<string>> = useParams();
    const productSlug: string | undefined = params.slug;
    const [product, setProduct] = useState<IProduct>();
    const review = {
        id: 1,
        comment: "This is a review",
        rating: 5,
        name: "John Doe",
        title: "Great product",
    }
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

                    <section className="flex gap-10">
                        <aside>
                            <ul>
                                <li className="mb-4"><ProductImage product={product}/></li>
                                <li className="mb-4"><SpiceRating rating={product.spiceRating}/></li>
                                <li className="mb-4"><span
                                    className="text-green-200 font-bold pt-4">${product.price.toFixed(2)}</span>
                                </li>
                                <li><Button className="w-full" label="Add to Cart" icon="pi pi-cart-plus"/></li>
                            </ul>
                        </aside>

                        <div>
                            <section className="mb-2">
                                <h2 className="font-bold text-lg mb-4">Description</h2>
                                <p>{product.description}</p>
                            </section>

                            <section className="mb-2">
                                <h2 className="font-bold text-lg mb-4">Ingredients</h2>
                                <p>Aged Cayenne Red Peppers, Distilled Vinegar, Water, Salt and Garlic Powder.</p>
                            </section>

                            <section className="mt-10">
                                <h2 className="font-bold text-lg mb-4">Reviews</h2>
                                <ReviewCard review={review}/>
                            </section>
                        </div>
                    </section>
                </Suspense>
            ) : (
                <>Product not found.</>
            )}
        </>
    )
}