import {IReview} from "./IReview.ts";
import {Card} from "primereact/card";
import {Avatar} from "primereact/avatar";
import {NavLink} from "react-router";
import {Rating} from "primereact/rating";
import SpiceRating from "../Products/SpiceRating.tsx";

interface IReviewCardProps {
    review: IReview
}

export default function ReviewCard(props: IReviewCardProps) {
    return (
        <Card title={props.review.title} className="w-full mb-4">
            <div className="flex gap-4">
                <aside className="flex flex-row text-center items-start gap-2 w-[100px]">
                    <NavLink to={`/users/${encodeURI(props.review.usernameSlug)}`}>
                        <Avatar
                            size={"large"}
                            className="text-center mx-auto cursor-pointer block w-[75px] h-[75px]"
                            image={`/images/avatars/${props.review.userAvatarFilename}`}
                            shape="circle"/>
                        <div>
                            {props.review.username}
                        </div>
                        <div>{new Date(props.review.createdAt).toLocaleDateString()}</div>
                    </NavLink>
                </aside>
                <section className="flex flex-col w-full gap-4">
                    <div className="min-h-[60px]">
                        {props.review.comment}
                    </div>
                    <section className="flex w-1/3 mt-6 justify-between">
                        <section className="w-1/2">
                            <p className="mb-2">Rating</p>
                            <Rating value={props.review.rating} readOnly={true} cancel={false}/>
                        </section>
                        <section className="w-1/2">
                            <p className="mb-2">Spice Rating</p>
                            <SpiceRating rating={props.review.spiceRating} readOnly={true}/>
                        </section>
                    </section>
                </section>
            </div>
        </Card>
    )
}