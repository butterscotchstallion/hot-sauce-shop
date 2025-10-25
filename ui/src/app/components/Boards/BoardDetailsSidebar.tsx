import {IBoard} from "./IBoard.ts";
import {Card} from "primereact/card";
import TimeAgo from "react-timeago";
import {useEffect, useState} from "react";
import dayjs from "dayjs";
import {NavLink} from "react-router";
import {IUser} from "../User/IUser.ts";

interface IBoardSidebarProps {
    board: IBoard;
    totalPosts: number;
    moderators: IUser[];
}

export function BoardDetailsSidebar({board, totalPosts, moderators}: IBoardSidebarProps) {
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
                    <li className="mb-2">
                        <strong
                            className="pr-2 mb-1 block">Moderators</strong>
                        {moderators.length > 0 ? (
                            <ul>
                                {moderators.map((moderator: IUser) => (
                                    <li key={moderator.id}>
                                        <NavLink to={`/users/${moderator.slug}`}>{moderator.username}</NavLink>
                                    </li>
                                ))}
                            </ul>
                        ) : "No moderators assigned."}
                    </li>
                </ul>
            </Card>
        </>
    )
}