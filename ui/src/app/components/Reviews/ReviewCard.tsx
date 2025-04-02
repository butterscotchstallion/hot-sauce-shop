import {IReview} from "./IReview.ts";
import {Card} from "primereact/card";
import {Avatar} from "primereact/avatar";
import {NavLink} from "react-router";

interface IReviewCardProps {
    review: IReview
}

export default function ReviewCard(props: IReviewCardProps) {
    return (
        <Card title={props.review.title}>
            <div className="flex gap-4">
                <aside className="flex flex-row text-center items-center gap-2 w-[100px]">
                    <NavLink to={`/users/${encodeURI(props.review.usernameSlug)}`}>
                        <Avatar
                            className="text-center mx-auto cursor-pointer block w-[50px] h-[50px]"
                            image={`/images/avatars/${props.review.userAvatarFilename}`}
                            shape="circle"/>
                        <div>{props.review.username}</div>
                    </NavLink>
                </aside>
                <div>
                    {props.review.comment}
                </div>
            </div>
        </Card>
    )
}