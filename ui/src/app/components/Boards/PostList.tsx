import {IBoardPost} from "./types/IBoardPost.ts";
import {ReactElement} from "react";
import BoardPost from "./BoardPost.tsx";

interface IPostListProps {
    posts: IBoardPost[];
    voteMap: Map<number, number>;
    replyMap: Map<number, number>;
    isCurrentUserBoardMod: boolean;
}

export default function PostList({posts, voteMap, replyMap, isCurrentUserBoardMod}: IPostListProps) {
    return (
        <>
            <section className="mt-4">
                {posts?.map((post: IBoardPost): ReactElement => {
                    return (
                        <BoardPost
                            boardPost={post}
                            key={`post-${post.id}`}
                            voteMap={voteMap}
                            replyMap={replyMap}
                            isCurrentUserBoardMod={isCurrentUserBoardMod}
                        />
                    )
                })}
            </section>
        </>
    )
}