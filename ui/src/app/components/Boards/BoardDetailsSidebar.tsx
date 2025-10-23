import {IBoard} from "./IBoard.ts";
import {Card} from "primereact/card";
import TimeAgo from "react-timeago";
import {useEffect, useState} from "react";
import dayjs from "dayjs";
import {NavLink} from "react-router";

interface IBoardSidebarProps {
    board: IBoard;
    totalPosts: number;
}

export function BoardDetailsSidebar({board, totalPosts}: IBoardSidebarProps) {
    const [createdAtFormatted, setCreatedAtFormatted] = useState<string>();

    useEffect(() => {
        setCreatedAtFormatted(dayjs(board.createdAt).format('MMMM D, YYYY'))
    }, [board]);

    return (
        <>
            <Card title={board.displayName}>
                <ul>
                    <li className="mb-2">
                        {board.description}
                    </li>
                    <li className="mb-2">
                        <strong
                            className="pr-2 mb-1 block">Created By</strong>
                        <NavLink to={`/users/${board.createdByUserSlug}`}>{board.createdByUsername}</NavLink>
                    </li>
                    <li className="mb-2">
                        <strong
                            className="pr-2 mb-1 block">Creation Date</strong>
                        <span className="cursor-help">
                            {<TimeAgo
                                date={board.createdAt}
                                title={createdAtFormatted}/>}
                        </span>
                    </li>
                    <li className="mb-2">
                        <strong
                            className="pr-2 mb-1 block">Total Posts</strong>
                        {totalPosts}
                    </li>
                </ul>
            </Card>
        </>
    )
}