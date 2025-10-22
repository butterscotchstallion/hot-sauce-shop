import {useEffect, useRef, useState} from "react";
import {IBoardPost} from "../../components/Boards/IBoardPost.ts";
import {getPosts} from "../../components/Boards/BoardsService.ts";
import PostList from "../../components/Boards/PostList.tsx";
import {getUserVoteMap} from "../../components/Boards/VoteService.ts";
import {Subject} from "rxjs";
import {NavigateFunction, Params, useNavigate, useParams} from "react-router";
import {Button} from "primereact/button";
import {IUser} from "../../components/User/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";

/**
 * Handles multiple scenarios where post(s) are displayed:
 * - All posts, unfiltered
 * - Filtered posts
 * - A single post
 */
export default function PostsListPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string | undefined = params?.boardSlug;
    const postSlug: string | undefined = params?.postSlug;
    const [posts, setPosts] = useState<IBoardPost[]>([]);
    const userVoteMap$ = useRef<Subject<Map<number, number>> | null>(null);
    const [userVoteMap, setuserVoteMap] = useState<Map<number, number>>(new Map());
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const navigate: NavigateFunction = useNavigate();
    const navigateToNewPostPage = () => {
        navigate(`/boards/${boardSlug}/posts/new`)
    }

    useEffect(() => {
        console.log("Fetching posts for boardSlug: " + boardSlug + " and postSlug: " + postSlug);
        const $posts = getPosts({
            postSlug,
            boardSlug,
        }).subscribe({
            next: (posts: IBoardPost[]) => {
                setPosts(posts);
            },
            error: (err) => {
                console.error(err);
            }
        });
        return () => {
            $posts.unsubscribe();
            if (userVoteMap$.current) (
                userVoteMap$.current.unsubscribe()
            )
        }
    }, [boardSlug, postSlug]);

    useEffect(() => {
        userVoteMap$.current = getUserVoteMap();
        userVoteMap$.current.subscribe({
            next: (voteMap: Map<number, number>) => {
                setuserVoteMap(voteMap);
            },
            error: (err) => {
                console.error(err);
            }
        })
    }, [posts])

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
            <PostList posts={posts} voteMap={userVoteMap}/>
        </>
    )
}