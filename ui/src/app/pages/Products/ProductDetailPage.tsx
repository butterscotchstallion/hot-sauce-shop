import {ReactElement, Suspense, useEffect, useState} from "react";
import Throbber from "../../components/Shared/Throbber.tsx";
import {Params, useNavigate, useParams} from "react-router";
import {getProductDetail, getProductReviews} from "../../components/Products/ProductService.ts";
import {IProduct} from "../../components/Products/IProduct.ts";
import {Subscription} from "rxjs";
import ProductImage from "../../components/Products/ProductImage.tsx";
import {Button} from "primereact/button";
import SpiceRating from "../../components/Products/SpiceRating.tsx";
import {Card} from "primereact/card";
import {IProductDetail} from "../../components/Products/IProductDetail.ts";
import {ITag} from "../../components/Tag/ITag.ts";
import {Tag} from "primereact/tag";
import {IUser} from "../../components/User/types/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {ProductReviewForm} from "../../components/Products/ProductReviewForm.tsx";
import {ProductReviewList} from "../../components/Products/ProductReviewList.tsx";
import {IReview} from "../../components/Reviews/IReview.ts";
import {IUserRole} from "../../components/User/types/IUserRole.ts";
import {userHasRole, UserRole} from "../../components/User/UserService.ts";
import {Chart} from "primereact/chart";
import {ChartData, ChartOptions} from "chart.js";
import {IProductReviewResponse} from "../../components/Products/IProductReviewResponse.ts";
import {IReviewRatingDistribution} from "../../components/Products/IReviewRatingDistribution.ts";

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
    const [chartOptions] = useState<ChartOptions>({});
    const [chartData, setChartData] = useState<ChartData>();
    const [isReviewInsightsVisible, setIsReviewInsightsVisible] = useState<boolean>(false);

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
                next: (resp: IProductReviewResponse) => {
                    setReviews(resp.reviews);
                    setReviewDataAndOptions(resp.reviewRatingDistributions);
                },
                error: (err) => console.error(err),
            });
        }
    };
    const setReviewDataAndOptions = (ratingDistribution: IReviewRatingDistribution[]) => {
        const documentStyle: CSSStyleDeclaration = getComputedStyle(document.documentElement);
        const labels: string[] = [];
        const chartData: number[] = [];
        ratingDistribution.forEach((ratingDistribution: IReviewRatingDistribution) => {
            labels.push(ratingDistribution.rating.toString());
            chartData.push(ratingDistribution.count);
        });
        const data = {
            labels: labels,
            datasets: [
                {
                    data: chartData,
                    backgroundColor: [
                        documentStyle.getPropertyValue('--blue-500'),
                        documentStyle.getPropertyValue('--yellow-500'),
                        documentStyle.getPropertyValue('--green-500')
                    ],
                    hoverBackgroundColor: [
                        documentStyle.getPropertyValue('--blue-400'),
                        documentStyle.getPropertyValue('--yellow-400'),
                        documentStyle.getPropertyValue('--green-400')
                    ]
                }
            ]
        }
        setChartData(data);
    }
    const toggleReviewInsightsButton = () => {
        setIsReviewInsightsVisible(!isReviewInsightsVisible);
    }

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
                        <aside className="w-[250px]">
                            <ul>
                                <li className="mb-4"><ProductImage product={product}/></li>
                                <li className="mb-4"><SpiceRating rating={product.spiceRating}/></li>
                                <li className="mb-4"><span
                                    className="text-green-200 font-bold pt-4">${product.price.toFixed(2)}</span>
                                </li>
                                <li><Button className="w-full" label="Add to Cart" icon="pi pi-cart-plus"/></li>
                            </ul>
                        </aside>

                        <div className="w-full">
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
                                <div className="flex justify-between mb-2">
                                    <h2 className="font-bold text-lg mb-4">Reviews ({reviews.length})</h2>

                                    <section className="w-1/2 flex gap-4 justify-end">
                                        <Button
                                            size="small"
                                            className="h-[35px]"
                                            severity="info"
                                            label="Toggle Review Insights"
                                            icon="pi pi-chart-bar"
                                            onClick={() => toggleReviewInsightsButton()}
                                        />
                                        {user && userHasRole(UserRole.REVIEWER, userRoles) && (
                                            <Button
                                                size="small"
                                                className="h-[35px]"
                                                label="Add Review"
                                                icon="pi pi-pencil"
                                                onClick={() => {
                                                    document.getElementById("add-review-area")?.scrollIntoView({
                                                        behavior: "smooth",
                                                        block: "start",
                                                        inline: "nearest"
                                                    });
                                                }}/>
                                        )}
                                    </section>
                                </div>

                                {isReviewInsightsVisible && (
                                    <section id="review-insights" className="mb-4">
                                        <Card title="Review Insights">
                                            <section>
                                                <h3 className="text-1xl font-bold mb-2">Review Ratings</h3>
                                                <Chart type="doughnut"
                                                       data={chartData}
                                                       options={chartOptions}
                                                       className="w-1/2 md:w-15rem"/>
                                            </section>
                                        </Card>
                                    </section>
                                )}

                                <div id="reviews">
                                    <Suspense fallback={<Throbber/>}>
                                        <ProductReviewList reviews={reviews} product={product}/>
                                    </Suspense>
                                </div>

                                <div id="add-review-area">
                                    {addReviewFormArea}
                                </div>
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