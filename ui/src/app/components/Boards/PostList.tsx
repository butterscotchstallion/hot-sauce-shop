import {IBoardPost} from "./IBoardPost.ts";
import {ReactElement} from "react";
import BoardPost from "./BoardPost.tsx";

interface IPostListProps {
    posts: IBoardPost[];
    voteMap: Map<number, number>;
    replyMap: Map<number, number>;
    isCurrentUserBoardMod: boolean;
}

export default function PostList({posts, voteMap, replyMap, isCurrentUserBoardMod}: IPostListProps) {
    const postList = posts?.map((post: IBoardPost): ReactElement => {
        return (
            <BoardPost
                boardPost={post}
                key={`post-${post.id}`}
                voteMap={voteMap}
                replyMap={replyMap}
                isCurrentUserBoardMod={isCurrentUserBoardMod}
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