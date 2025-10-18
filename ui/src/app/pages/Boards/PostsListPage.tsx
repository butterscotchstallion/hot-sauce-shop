import {useEffect, useState} from "react";
import {IBoardPost} from "../../components/Boards/IBoardPost.ts";
import {getPosts} from "../../components/Boards/BoardsService.ts";
import PostList from "../../components/Boards/PostList.tsx";
import {IUser} from "../../components/User/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {Button} from "primereact/button";
import {NavigateFunction, useNavigate} from "react-router";

/**
 * All posts, unfiltered
 */
export default function PostsListPage() {
    const [posts, setPosts] = useState<IBoardPost[]>([]);
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const navigate: NavigateFunction = useNavigate();

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

    const navigateToNewPostPage = () => {
        navigate('/posts/new')
    }

    return (
        <>
            {user && (
                <section>
                    <div className="flex w-full justify-end">
                        <Button onClick={() => navigateToNewPostPage()}>
                            <i className="pi pi-envelope mr-2"></i> Create Post
                        </Button>
                    </div>
                </section>
            )}
            <PostList posts={posts}/>
        </>
    )
}