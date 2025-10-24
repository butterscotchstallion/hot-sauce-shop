import {IBoardPost} from "./IBoardPost.ts";
import {ReactElement} from "react";
import BoardPost from "./BoardPost.tsx";

interface IPostListProps {
    posts: IBoardPost[];
    voteMap: Map<number, number>;
    replyMap: Map<number, number>;
}

export default function PostList({posts, voteMap, replyMap}: IPostListProps) {
    const postList = posts?.map((post: IBoardPost): ReactElement => {
        return (
            <BoardPost
                boardPost={post}
                key={`post-${post.id}`}
                voteMap={voteMap}
                replyMap={replyMap}
            />
        )
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