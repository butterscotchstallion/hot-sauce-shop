import AddEditPostForm from "../../components/Boards/AddEditPostForm.tsx";
import {Params, useParams, useSearchParams} from "react-router";
import {IBoard} from "../../components/Boards/types/IBoard.ts";
import {RefObject, useEffect, useRef, useState} from "react";
import {getBoardByBoardSlug} from "../../components/Boards/BoardsService.ts";
import {Subject} from "rxjs";
import {IBoardDetails} from "../../components/Boards/types/IBoardDetails.ts";

export default function NewPostPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string | undefined = params?.slug;
    const [searchParams] = useSearchParams();
    const [board, setBoard] = useState<IBoard>();
    const board$: RefObject<Subject<IBoardDetails> | null> = useRef<Subject<IBoardDetails>>(null);
    const parentSlug = useRef<string>('');

    useEffect(() => {
        const parentSlugParam: string | null = searchParams.get("parentSlug");
        if (parentSlugParam && parentSlugParam.length > 0) {
            parentSlug.current = parentSlugParam;
        }
        if (boardSlug) {
            board$.current = getBoardByBoardSlug(boardSlug);
            board$.current.subscribe({
                next: (details: IBoardDetails) => {
                    setBoard(details.board);
                    console.info("Board set to " + details.board.displayName);
                },
                error: (error: Error) => console.error(error),
            });
        }
        return () => {
            board$?.current?.unsubscribe();
        }
    }, [boardSlug]);

    return (
        <>
            <h1 className="text-3xl font-bold mb-4">
                {board && `${board.displayName} -`} New Post
            </h1>
            <section className="mt-4 w-1/2">
                {board && (
                    <AddEditPostForm
                        boardSlug={board.slug}
                        parentSlug={parentSlug.current}
                    />
                )}
            </section>
        </>
    )
}