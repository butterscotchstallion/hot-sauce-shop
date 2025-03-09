import {Rating} from "primereact/rating";

interface ISpiceRatingProps {
    rating: number;
}

export default function SpiceRating(props: ISpiceRatingProps) {
    return (
        <Rating value={props.rating}
                readOnly={true}
                cancel={false}
                title={`Spice Rating: ${props.rating}/5`}
                onIcon={<img src="/images/rating/hot-pepper-rating-on.png" alt="custom-image-active" width="25"
                             height="25"/>}
                offIcon={<img src="/images/rating/hot-pepper-rating-off.png" alt="custom-image" width="25"
                              height="25"/>}
        />
    )
}