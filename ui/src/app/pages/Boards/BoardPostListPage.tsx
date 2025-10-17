import {Params, useParams} from "react-router";
import {useEffect, useState} from "react";
import {IBoardPost} from "../../components/Boards/IBoardPost.ts";
import {getPostsByBoardSlug} from "../../components/Boards/BoardsService.ts";
import PostList from "../../components/Boards/PostList.tsx";

export default function BoardPostListPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string | undefined = params?.slug;
    const [boardPosts, setBoardPosts] = useState<IBoardPost[]>([]);

    useEffect(() => {
        if (boardSlug) {
            getPostsByBoardSlug(boardSlug).subscribe({
                next: (posts: IBoardPost[]) => {
                    setBoardPosts(posts);
                },
                error: (err) => {
                    console.error(err);
                }
            });
        }
    }, [boardSlug])

    return (
        <>
            <PostList posts={boardPosts} boardSlug={boardSlug}/>
        </>
    )
}