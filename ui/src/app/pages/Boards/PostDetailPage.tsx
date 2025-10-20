import {Suspense, useEffect, useRef, useState} from "react";
import {IBoardPost} from "../../components/Boards/IBoardPost.ts";
import {Subject} from "rxjs";
import {getPostDetail, getPosts} from "../../components/Boards/BoardsService.ts";
import {Params, useParams} from "react-router";
import BoardPost from "../../components/Boards/BoardPost.tsx";
import PostList from "../../components/Boards/PostList.tsx";
import Throbber from "../../components/Shared/Throbber.tsx";

export default function PostDetailPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string | undefined = params?.boardSlug;
    const postSlug: string | undefined = params?.postSlug;
    const [post, setPost] = useState<IBoardPost>();
    const [replies, setReplies] = useState<IBoardPost[]>([]);
    const getPost$ = useRef<Subject<IBoardPost> | null>(null);
    const getReplies$ = useRef<Subject<IBoardPost[]> | null>(null);

    useEffect(() => {
        if (boardSlug && postSlug) {
            // Post detail
            getPost$.current = getPostDetail(boardSlug, postSlug);
            getPost$.current.subscribe({
                next: (post: IBoardPost) => {
                    setPost(post);
                    // Replies
                    getReplies$.current = getPosts(post.id);
                    getReplies$.current.subscribe({
                        next: (replies: IBoardPost[]) => {
                            setReplies(replies);
                        },
                        error: (err) => {
                            console.error(err);
                        }
                    });
                },
                error: (err) => {
                    console.error(err);
                }
            });
        } else {
            // TODO: redirect to 404 page
        }
        return () => {
            getPost$?.current?.unsubscribe();
        }
    }, []);

    return (
        <>
            {post && <BoardPost post={post}/>}
            {replies && replies.length > 0 && (
                <Suspense fallback={<Throbber/>}>
                    <PostList posts={replies}/>
                </Suspense>
            )}
        </>
    )
};