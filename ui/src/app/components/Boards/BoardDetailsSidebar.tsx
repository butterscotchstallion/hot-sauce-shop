import {Card} from "primereact/card";
import TimeAgo from "react-timeago";
import {useEffect, useRef, useState} from "react";
import dayjs from "dayjs";
import {NavLink, Params, useParams} from "react-router";
import {IUser} from "../User/types/IUser.ts";
import {IBoardDetails} from "./types/IBoardDetails.ts";
import {BoardSettingsButton} from "./BoardSettingsButton.tsx";
import {isSettingsAreaAvailable} from "./BoardsService.ts";
import {BehaviorSubject} from "rxjs";
import {IUserRole} from "../User/types/IUserRole.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";

interface IBoardSidebarProps {
    boardDetails: IBoardDetails | undefined;
}

export function BoardDetailsSidebar({boardDetails}: IBoardSidebarProps) {
    const [createdAtFormatted, setCreatedAtFormatted] = useState<string>();
    const settingsAvailable = useRef<BehaviorSubject<boolean>>(null);
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string = params?.boardSlug || '';
    const [settingsAreaAvailable, setSettingsAreaAvailable] = useState<boolean>(false);
    const roles: IUserRole[] = useSelector((state: RootState) => state.user.roles);

    useEffect(() => {
        if (!boardDetails?.board) return;
        setCreatedAtFormatted(dayjs(boardDetails.board.createdAt).format('MMMM D, YYYY'));

        settingsAvailable.current = isSettingsAreaAvailable(boardSlug, roles)
        settingsAvailable.current.subscribe({
            next: (available: boolean) => {
                console.log(`Board settings available: ${available}`);
                setSettingsAreaAvailable(available)
            },
            error: (err: string) => {
                console.error(err);
                setSettingsAreaAvailable(false)
            }
        });
        return () => {
            settingsAvailable.current?.unsubscribe();
        }
    }, [boardDetails?.board, boardSlug, roles]);

    const header = () => (
        <>
            {boardDetails && (
                <section className="flex justify-between items-center gap-2 p-4 pb-0">
                    <h2 className="text-xl">{boardDetails.board.displayName}</h2>
                    <BoardSettingsButton
                        settingsAreaAvailable={settingsAreaAvailable}
                        boardSlug={boardSlug}
                    />
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