import {IProduct} from "./IProduct.ts";
import {ReactElement, useEffect, useState} from "react";
import {IReview} from "../Reviews/IReview.ts";
import {getProductReviews} from "./ProductService.ts";
import {Subscription} from "rxjs";
import ReviewCard from "../Reviews/ReviewCard.tsx";

interface IProductReviewListProps {
    product: IProduct;
}

export function ProductReviewList(props: IProductReviewListProps): ReactElement {
    const [reviews, setReviews] = useState<IReview[]>([]);
    const reviewList: ReactElement[] = reviews.map((review: IReview) => {
        return <ReviewCard key={review.id} review={review}/>
    });
    useEffect(() => {
        const reviews$: Subscription = getProductReviews(props.product.slug).subscribe({
            next: (reviews: IReview[]) => setReviews(reviews),
            error: (err) => console.error(err),
        });
        return () => {
            reviews$.unsubscribe();
        }
    }, [props.product.slug]);
    return (
        <>
            {reviewList}
        </>
    )
}