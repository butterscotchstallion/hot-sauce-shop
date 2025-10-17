import {useEffect, useState} from "react";
import {IBoardPost} from "../../components/Boards/IBoardPost.ts";
import {getPosts} from "../../components/Boards/BoardsService.ts";
import PostList from "../../components/Boards/PostList.tsx";

/**
 * All posts, unfiltered
 */
export default function PostsListPage() {
    const [posts, setPosts] = useState<IBoardPost[]>([]);

    useEffect(() => {
        const $posts = getPosts().subscribe({
            next: (posts: IBoardPost[]) => {
                setPosts(posts);
            },
            error: (err) => {
                console.error(err);
            }
        });
        return () => {
            $posts.unsubscribe();
        }
    }, []);

    return (
        <>
            <PostList posts={posts}/>
        </>
    )
}