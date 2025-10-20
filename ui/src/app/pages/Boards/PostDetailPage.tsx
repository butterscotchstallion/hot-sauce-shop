import {Suspense, useEffect, useRef, useState} from "react";
import {IBoardPost} from "../../components/Boards/IBoardPost.ts";
import {Subject} from "rxjs";
import {getPostDetail, getPosts} from "../../components/Boards/BoardsService.ts";
import {Params, useParams} from "react-router";
import BoardPost from "../../components/Boards/BoardPost.tsx";
import PostList from "../../components/Boards/PostList.tsx";
import Throbber from "../../components/Shared/Throbber.tsx";
import {IUser} from "../../components/User/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import AddEditPostForm from "../../components/Boards/AddEditPostForm.tsx";
import {Card} from "primereact/card";

export default function PostDetailPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string | undefined = params?.boardSlug;
    const postSlug: string | undefined = params?.postSlug;
    const [post, setPost] = useState<IBoardPost>();
    const [replies, setReplies] = useState<IBoardPost[]>([]);
    const getPost$ = useRef<Subject<IBoardPost> | null>(null);
    const getReplies$ = useRef<Subject<IBoardPost[]> | null>(null);
    const user: IUser | null = useSelector((state: RootState) => state.user.user);

    const getReplies = () => {
        if (post) {
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
        }
    };

    useEffect(() => {
        if (boardSlug && postSlug) {
            // Post detail
            getPost$.current = getPostDetail(boardSlug, postSlug);
            getPost$.current.subscribe({
                next: (post: IBoardPost) => {
                    setPost(post);
                    getReplies();
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
            {replies.length === 0 && (
                <section>
                    <h1 className="text-2xl font-bold mb-4">No replies yet.</h1>
                    <p>Be the first to reply!</p>
                </section>
            )}
            {user && boardSlug && post && (
                <>
                    <section className="mt-4">
                        <Card>
                            <h1 className="text-2xl font-bold mb-4">Comment</h1>
                            <AddEditPostForm boardSlug={boardSlug}
                                             parentId={post.id}
                                             addPostCallback={() => getReplies()}/>
                        </Card>
                    </section>
                </>
            )}
        </>
    )
};