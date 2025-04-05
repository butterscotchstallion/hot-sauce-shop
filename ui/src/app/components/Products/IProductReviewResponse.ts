import {IReview} from "../Reviews/IReview.ts";
import {IReviewRatingDistribution} from "./IReviewRatingDistribution.ts";

export interface IProductReviewResponse {
    reviews: IReview[];
    reviewRatingDistributions: IReviewRatingDistribution[]
}