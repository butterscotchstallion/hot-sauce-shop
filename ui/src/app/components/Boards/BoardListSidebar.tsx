import {Card} from "primereact/card";
import {IBoard} from "./types/IBoard.ts";
import {NavLink} from "react-router";

interface BoardListSidebarProps {
    boards: IBoard[];
}

export function BoardListSidebar({boards}: BoardListSidebarProps) {
    return (
        <Card title="Boards">
            {boards.length > 0 && (
                <ul>
                    {boards.map((board: IBoard) => (
                        <li className="mb-2" key={board.id}>
                            <NavLink to={`/boards/${board.slug}`}>{board.displayName}</NavLink>
                        </li>
                    ))}
                </ul>
            )}
            {boards.length === 0 && "No boards available."}
        </Card>
    )
}