import {IProduct} from "./types/IProduct.ts";
import {ReactElement} from "react";
import {IReview} from "../Reviews/IReview.ts";
import ReviewCard from "../Reviews/ReviewCard.tsx";

interface IProductReviewListProps {
    product: IProduct;
    reviews: IReview[];
}

export function ProductReviewList(props: IProductReviewListProps): ReactElement {
    const reviewList: ReactElement[] = props.reviews.map((review: IReview) => {
        return <ReviewCard key={review.id} review={review}/>
    });
    return (
        <>
            {reviewList}
        </>
    )
}