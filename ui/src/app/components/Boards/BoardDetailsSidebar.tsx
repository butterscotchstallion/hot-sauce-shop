import {Card} from "primereact/card";
import TimeAgo from "react-timeago";
import {useEffect, useState} from "react";
import dayjs from "dayjs";
import {NavLink} from "react-router";
import {IUser} from "../User/IUser.ts";
import {IBoardDetails} from "./IBoardDetails.ts";
import {BoardSettingsButton} from "./BoardSettingsButton.tsx";

interface IBoardSidebarProps {
    boardDetails: IBoardDetails | undefined;
}

export function BoardDetailsSidebar({boardDetails}: IBoardSidebarProps) {
    const [createdAtFormatted, setCreatedAtFormatted] = useState<string>();

    useEffect(() => {
        if (!boardDetails?.board) return;
        setCreatedAtFormatted(dayjs(boardDetails.board.createdAt).format('MMMM D, YYYY'))
    }, [boardDetails?.board]);

    const header = () => (
        <>
            {boardDetails && (
                <section className="flex justify-between items-center gap-2 p-4">
                    <h2 className="text-xl">{boardDetails.board.displayName}</h2>
                    <BoardSettingsButton/>
                </section>
            )}
        </>
    )
    return (
        <>
            {boardDetails && (
                <Card header={header}>
                    <ul>
                        <li className="mb-2">
                            {boardDetails.board.description}
                        </li>
                        <li className="mb-2">
                            <strong
                                className="pr-2 mb-1 block">Created By</strong>
                            <NavLink
                                to={`/users/${boardDetails.board.createdByUserSlug}`}>{boardDetails.board.createdByUsername}</NavLink>
                        </li>
                        <li className="mb-2">
                            <strong
                                className="pr-2 mb-1 block">Creation Date</strong>
                            <span className="cursor-help">
                            {<TimeAgo
                                date={boardDetails.board.createdAt}
                                title={createdAtFormatted}/>}
                        </span>
                        </li>
                        <li className="mb-2">
                            <strong
                                className="pr-2 mb-1 block">Total Posts</strong>
                            {boardDetails.totalPosts}
                        </li>
                        <li className="mb-2">
                            <strong
                                className="pr-2 mb-1 block">Members</strong>
                            {boardDetails.numBoardMembers}
                        </li>
                        <li className="mb-2">
                            <strong
                                className="pr-2 mb-1 block">Moderators</strong>
                            {boardDetails.moderators.length > 0 ? (
                                <ul>
                                    {boardDetails.moderators.map((moderator: IUser) => (
                                        <li key={moderator.id}>
                                            <NavLink to={`/users/${moderator.slug}`}>{moderator.username}</NavLink>
                                        </li>
                                    ))}
                                </ul>
                            ) : "No moderators assigned."}
                        </li>
                    </ul>
                </Card>
            )}
        </>
    )
}