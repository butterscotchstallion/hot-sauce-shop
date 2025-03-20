import {Rating, RatingChangeEvent} from "primereact/rating";
import {useState} from "react";
import {Nullable} from "primereact/ts-helpers";

interface ISpiceRatingProps {
    rating: number;
    readOnly?: boolean;
}

export default function SpiceRating(props: ISpiceRatingProps) {
    const [rating, setRating] = useState<number>(props.rating);

    const updateRating = (value: Nullable<number>) => {
        if (value && !props.readOnly) {
            setRating(value);
        }
    }

    return (
        <Rating value={rating}
                readOnly={props.readOnly}
                cancel={false}
                title={`Spice Rating: ${props.rating}/5`}
                onChange={(e: RatingChangeEvent) => updateRating(e.value)}
                onIcon={<img src="/images/rating/hot-pepper-rating-on.png" alt="custom-image-active" width="25"
                             height="25"/>}
                offIcon={<img src="/images/rating/hot-pepper-rating-off.png" alt="custom-image" width="25"
                              height="25"/>}
        />
    )
}