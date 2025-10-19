import {IBoardPost} from "./IBoardPost.ts";
import {ReactElement} from "react";
import BoardPost from "./BoardPost.tsx";

interface IPostListProps {
    posts: IBoardPost[];
}

export default function PostList({posts}: IPostListProps) {
    const postList = posts?.map((post: IBoardPost): ReactElement => {
        return <BoardPost post={post}/>
    });
    return (
        <>
            {posts?.length > 0 && (
                <section className="mt-4">
                    {postList}
                </section>
            )}
        </>
    )
}