import {Suspense, useEffect, useState} from "react";
import Throbber from "../components/Shared/Throbber.tsx";
import {Params, useNavigate, useParams} from "react-router";
import {getProductDetail} from "../components/Products/ProductService.ts";
import {IProduct} from "../components/Products/IProduct.ts";
import {Subscription} from "rxjs";
import ProductImage from "../components/Products/ProductImage.tsx";
import {Button} from "primereact/button";
import SpiceRating from "../components/Products/SpiceRating.tsx";
import ReviewCard from "../components/Reviews/ReviewCard.tsx";
import {Card} from "primereact/card";
import {IProductDetail} from "../components/Products/IProductDetail.ts";
import {ITag} from "../components/Tag/ITag.ts";
import {Tag} from "primereact/tag";

export default function ProductDetailPage() {
    const params: Readonly<Params<string>> = useParams();
    const productSlug: string | undefined = params.slug;
    const [product, setProduct] = useState<IProduct>();
    const [productTags, setProductTags] = useState<ITag[]>([])
    const navigate = useNavigate();
    const review = {
        id: 1,
        comment: "This is a review",
        rating: 5,
        name: "Jane Doe",
        title: "Great product",
    }
    const productTagList = productTags.map((tag: ITag) => {
        return <Tag key={tag.id} severity="info" value={tag.name} className="mr-2"></Tag>
    });

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

                        <Button label="Edit" icon="pi pi-pencil" onClick={() => {
                            navigate("/admin/products/edit/" + product.slug);
                        }}/>
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
                                <Card>
                                    <p>{product.description}</p>

                                    <section className="mt-4">
                                        {productTagList}
                                    </section>
                                </Card>
                            </section>

                            <section className="mb-2">
                                <h2 className="font-bold text-lg mb-4">Ingredients</h2>
                                <Card>
                                    <p>Aged Cayenne Red Peppers, Distilled Vinegar, Water, Salt and Garlic Powder.</p>
                                </Card>
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