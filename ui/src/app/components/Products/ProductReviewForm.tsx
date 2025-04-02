import {InputText} from "primereact/inputtext";
import {InputTextarea} from "primereact/inputtextarea";
import {Rating, RatingChangeEvent} from "primereact/rating";
import {useState} from "react";
import {Card} from "primereact/card";
import SpiceRating from "./SpiceRating.tsx";
import {Button} from "primereact/button";
import {IProduct} from "./IProduct.ts";
import {addReview} from "./ProductService.ts";
import {IReview} from "../Reviews/IReview.ts";

interface IProductReviewFormProps {
    product: IProduct;
}

export function ProductReviewForm(props: IProductReviewFormProps) {
    const [productRating, setProductRating] = useState<number>(0);
    const [spiceRating, setSpiceRating] = useState<number>(0);
    const [reviewTitle, setReviewTitle] = useState<string>("");
    const [reviewComment, setReviewComment] = useState<string>("");

    const onSubmit = () => {
        const review: IReview = {
            title: reviewTitle,
            rating: productRating,
            spiceRating: spiceRating,
            comment: reviewComment,
        }
        addReview(review, props.product.slug).subscribe({
            next: (success: boolean) => {

            }
        })
    };

    return (
        <>
            <Card title="Add Review" className="mb-4">
                <form onSubmit={onSubmit}>
                    <div className="flex flex-col gap-4">
                        <section>
                            <label className="block mb-2" htmlFor="title">Title</label>
                            <InputText id="title" type="text" required={true} minLength={10} maxLength={255}/>
                        </section>

                        <section>
                            <label className="block mb-2" htmlFor="rating">Overall Rating</label>
                            <Rating value={productRating}
                                    required={true}
                                    onChange={(e: RatingChangeEvent) => setProductRating(e.value || 1)}
                                    cancel={false}/>
                        </section>

                        <section>
                            <label className="block mb-2" htmlFor="rating">Spice Rating</label>
                            <SpiceRating readOnly={false}
                                         rating={spiceRating}
                                         onChange={(e: number) => setSpiceRating(e)}/>
                        </section>

                        <section>
                            <label className="block mb-2" htmlFor="reviewComment">Review</label>
                            <InputTextarea
                                required={true}
                                id="reviewComment"
                                rows={5}
                                cols={40}
                                minLength={10}
                                maxLength={1000}/>
                        </section>

                        <section className="mt-4 flex justify-right">
                            <Button type="submit" icon="pi pi-save" label="Submit Review"/>
                        </section>
                    </div>
                </form>
            </Card>
        </>
    )
}