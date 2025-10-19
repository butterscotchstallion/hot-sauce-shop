import {useEffect, useRef, useState} from "react";
import {IBoardPost} from "../../components/Boards/IBoardPost.ts";
import {Subject} from "rxjs";
import {getPostDetail} from "../../components/Boards/BoardsService.ts";
import {Params, useParams} from "react-router";
import BoardPost from "../../components/Boards/BoardPost.tsx";

export default function PostDetailPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string | undefined = params?.boardSlug;
    const postSlug: string | undefined = params?.postSlug;
    const [post, setPost] = useState<IBoardPost>();
    const getPost$ = useRef<Subject<IBoardPost> | null>(null);

    useEffect(() => {
        if (boardSlug && postSlug) {
            getPost$.current = getPostDetail(boardSlug, postSlug);
            getPost$.current.subscribe({
                next: (post: IBoardPost) => {
                    setPost(post);
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
        </>
    )
};