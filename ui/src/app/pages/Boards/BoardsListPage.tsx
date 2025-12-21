import {BoardsList} from "../../components/Boards/BoardsList.tsx";
import {Suspense, useEffect, useState} from "react";
import {getBoards} from "../../components/Boards/BoardsService.ts";
import {IBoard} from "../../components/Boards/types/IBoard.ts";
import Throbber from "../../components/Shared/Throbber.tsx";

export default function BoardsListPage() {
    const [boards, setBoards] = useState<IBoard[]>([]);

    useEffect(() => {
        getBoards().subscribe({
            next: (boards: IBoard[]) => {
                setBoards(boards);
            },
            error: (err) => {
                console.error(err);
            }
        });
    }, []);

    return (
        <>
            <h1 className="text-3xl font-bold mb-4">Message Boards</h1>
            <section className="mt-4">
                <Suspense fallback={<Throbber/>}>
                    <BoardsList boards={boards}/>
                </Suspense>
            </section>
        </>
    )
}