import AddEditPostForm from "../../components/Boards/AddEditPostForm.tsx";
import {Params, useParams} from "react-router";
import {IBoard} from "../../components/Boards/IBoard.ts";
import {RefObject, useEffect, useRef, useState} from "react";
import {getBoardByBoardSlug} from "../../components/Boards/BoardsService.ts";
import {Subject} from "rxjs";

export default function NewPostPage() {
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string | undefined = params?.slug;
    const [board, setBoard] = useState<IBoard>();
    const board$: RefObject<Subject<IBoard> | null> = useRef<Subject<IBoard>>(null);

    useEffect(() => {
        if (boardSlug) {
            board$.current = getBoardByBoardSlug(boardSlug);
            board$.current.subscribe({
                next: (board: IBoard) => setBoard(board),
                error: (error: Error) => console.error(error),
            });
        }
        return () => {
            board$?.current?.unsubscribe();
        }
    }, []);

    return (
        <>
            <h1 className="text-3xl font-bold mb-4">New Post</h1>
            <section className="mt-4 w-1/2">
                {board && (
                    <AddEditPostForm boardId={board.id}/>
                )}
            </section>
        </>
    )
}