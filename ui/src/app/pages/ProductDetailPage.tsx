import {ReactElement, Suspense, useEffect, useState} from "react";
import Throbber from "../components/Shared/Throbber.tsx";
import {Params, useNavigate, useParams} from "react-router";
import {getProductDetail, getProductReviews} from "../components/Products/ProductService.ts";
import {IProduct} from "../components/Products/IProduct.ts";
import {Subscription} from "rxjs";
import ProductImage from "../components/Products/ProductImage.tsx";
import {Button} from "primereact/button";
import SpiceRating from "../components/Products/SpiceRating.tsx";
import {Card} from "primereact/card";
import {IProductDetail} from "../components/Products/IProductDetail.ts";
import {ITag} from "../components/Tag/ITag.ts";
import {Tag} from "primereact/tag";
import {IUser} from "../components/User/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../store.ts";
import {ProductReviewForm} from "../components/Products/ProductReviewForm.tsx";
import {ProductReviewList} from "../components/Products/ProductReviewList.tsx";
import {IReview} from "../components/Reviews/IReview.ts";
import {IUserRole} from "../components/User/IUserRole.ts";
import {userHasRole, UserRole} from "../components/User/UserService.ts";

export default function ProductDetailPage() {
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const userRoles: IUserRole[] | [] = useSelector((state: RootState) => state.user.roles);
    const params: Readonly<Params<string>> = useParams();
    const productSlug: string | undefined = params.slug;
    const [product, setProduct] = useState<IProduct>();
    const [productTags, setProductTags] = useState<ITag[]>([])
    const navigate = useNavigate();
    const [reviews, setReviews] = useState<IReview[]>([]);
    const productTagList = productTags.map((tag: ITag) => {
        return <Tag key={tag.id} severity="info" value={tag.name} className="mr-2"></Tag>
    });

    const reviewSubmittedCallback = () => {
        loadReviews();
    }
    const addReviewFormArea: ReactElement = (
        user && userHasRole(UserRole.REVIEWER, userRoles) && product ?
            <ProductReviewForm reviewSubmittedCallback={reviewSubmittedCallback} product={product}/>
            : <p className="mb-4">Sign in to add a review</p>
    )
    const loadReviews: () => (Subscription | undefined) = (): Subscription | undefined => {
        if (productSlug) {
            return getProductReviews(productSlug).subscribe({
                next: (reviews: IReview[]) => setReviews(reviews),
                error: (err) => console.error(err),
            });
        }
    };

    useEffect(() => {
        if (productSlug) {
            const reviews$: Subscription | undefined = loadReviews();
            const product$: Subscription = getProductDetail(productSlug).subscribe({
                next: (productDetail: IProductDetail) => {
                    setProduct(productDetail.product);
                    setProductTags(productDetail.tags);
                },
                error: (err) => {
                    console.error(err);
                }
            });
            return () => {
                product$.unsubscribe();
                reviews$?.unsubscribe();
            }
        }
    }, []);

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
                                {addReviewFormArea}

                                <Suspense fallback={<Throbber/>}>
                                    <ProductReviewList reviews={reviews} product={product}/>
                                </Suspense>
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