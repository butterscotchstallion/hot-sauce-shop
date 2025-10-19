import {NavigateFunction, Params, useNavigate, useParams} from "react-router";
import {useEffect, useState} from "react";
import {IBoardPost} from "../../components/Boards/IBoardPost.ts";
import {getPostsByBoardSlug} from "../../components/Boards/BoardsService.ts";
import PostList from "../../components/Boards/PostList.tsx";
import {RootState} from "../../store.ts";
import {IUser} from "../../components/User/IUser.ts";
import {useSelector} from "react-redux";
import {Button} from "primereact/button";

/**
 * Posts for a specific board
 * @constructor
 */
export default function BoardPostListPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string | undefined = params?.slug;
    const [boardPosts, setBoardPosts] = useState<IBoardPost[]>([]);
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const navigate: NavigateFunction = useNavigate();
    const navigateToNewPostPage = () => {
        navigate(`/boards/${boardSlug}/posts/new`)
    }

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
            {user && (
                <section>
                    <div className="flex w-full justify-end">
                        <Button onClick={() => navigateToNewPostPage()}>
                            <i className="pi pi-envelope mr-2"></i> Create Post
                        </Button>
                    </div>
                </section>
            )}
            <PostList posts={boardPosts}/>
        </>
    )
}