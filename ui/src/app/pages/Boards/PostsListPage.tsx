import {RefObject, useEffect, useRef, useState} from "react";
import {IBoardPost} from "../../components/Boards/IBoardPost.ts";
import {
    getBoardByBoardSlug,
    getBoards,
    getPosts,
    getTotalPostReplyMap,
    getTotalPostsByBoardSlug
} from "../../components/Boards/BoardsService.ts";
import PostList from "../../components/Boards/PostList.tsx";
import {getUserVoteMap} from "../../components/Boards/VoteService.ts";
import {Subject} from "rxjs";
import {NavigateFunction, Params, useNavigate, useParams} from "react-router";
import {Button} from "primereact/button";
import {IUser} from "../../components/User/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {IBoard} from "../../components/Boards/IBoard.ts";
import {getJoinedBoards, userJoinBoard} from "../../components/User/UserService.ts";
import {Toast} from "primereact/toast";
import {BoardDetailsSidebar} from "../../components/Boards/BoardDetailsSidebar.tsx";
import {BoardListSidebar} from "../../components/Boards/BoardListSidebar.tsx";

/**
 * Handles multiple scenarios where post(s) are displayed:
 * - All posts, unfiltered
 * - Filtered posts
 * - A single post
 */
export default function PostsListPage() {
    const toast: RefObject<Toast | null> = useRef(null);
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string = params?.boardSlug || '';
    const postSlug: string = params?.postSlug || '';
    const [posts, setPosts] = useState<IBoardPost[]>([]);
    const [board, setBoard] = useState<IBoard>();
    const [boardTotalPosts, setBoardTotalPosts] = useState<number>(0);
    const [userIsBoardMember, setUserIsBoardMember] = useState<boolean>(false);
    const userVoteMap$ = useRef<Subject<Map<number, number>> | null>(null);
    const [userVoteMap, setUserVoteMap] = useState<Map<number, number>>(new Map());
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const navigate: NavigateFunction = useNavigate();
    const navigateToNewPostPage = () => {
        navigate(`/boards/${boardSlug}/posts/new`)
    }
    const joinBoard$ = useRef<Subject<boolean>>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [boards, setBoards] = useState<IBoard[]>([]);
    const [totalPostReplyMap, setTotalPostReplyMap] = useState<Map<number, number>>(new Map());

    const joinBoard = () => {
        if (board) {
            setIsLoading(true);
            joinBoard$.current = userJoinBoard(board.id);
            joinBoard$.current.subscribe({
                next: () => {
                    setUserIsBoardMember(true);
                    setIsLoading(false);
                    if (toast.current) {
                        toast?.current.show({
                            severity: 'success',
                            summary: 'Success',
                            detail: `You are now a member of ${board.displayName}!`,
                            life: 3000,
                        })
                    }
                },
                error: (err) => {
                    console.log(err);
                    setIsLoading(false);
                    if (toast.current) {
                        toast?.current.show({
                            severity: 'error',
                            summary: 'Error',
                            detail: `There was a problem joining ${board.displayName}.`,
                            life: 3000,
                        })
                    }
                }
            });
        }
    }

    // When the board/post changes, get filtered posts based on the scenario
    useEffect(() => {
        console.log("Fetching posts for boardSlug: " + boardSlug + " and postSlug: " + postSlug);
        const posts$: Subject<IBoardPost[]> = getPosts({
            postSlug,
            boardSlug,
        });
        posts$.subscribe({
            next: (posts: IBoardPost[]) => {
                setPosts(posts);
            },
            error: (err) => {
                console.error(err);
            }
        });
        // When viewing a specific board...
        let board$: Subject<IBoard>;
        let boardTotalPosts$: Subject<number>;
        let getBoards$: Subject<IBoard[]>;
        const replyMap$: Subject<Map<number, number>> = getTotalPostReplyMap(boardSlug);
        replyMap$.subscribe({
            next: (totalPostReplyMap: Map<number, number>) => {
                setTotalPostReplyMap(totalPostReplyMap)
            },
            error: (err) => console.error(err)
        })
        if (boardSlug) {
            board$ = getBoardByBoardSlug(boardSlug);
            board$.subscribe({
                next: (board: IBoard) => setBoard(board),
                error: (err) => console.error(err),
            });
            boardTotalPosts$ = getTotalPostsByBoardSlug(boardSlug);
            boardTotalPosts$.subscribe({
                next: (totalPosts: number) => setBoardTotalPosts(totalPosts),
                error: (err) => console.error(err),
            });
        } else {
            // on "All posts" page
            getBoards$ = getBoards();
            getBoards$.subscribe({
                next: (boardList: IBoard[]) => setBoards(boardList),
                error: (err) => console.error(err),
            })
        }
        return () => {
            posts$.unsubscribe();
            board$?.unsubscribe();
            boardTotalPosts$?.unsubscribe();
            getBoards$?.unsubscribe();
            //replyMap$.unsubscribe();
            if (joinBoard$.current) {
                joinBoard$.current.unsubscribe();
            }
            if (userVoteMap$.current) (
                userVoteMap$.current.unsubscribe()
            )
        }
    }, [boardSlug, postSlug]);

    // When posts/user changes, get the posts the user has voted on
    useEffect(() => {
        if (user) {
            userVoteMap$.current = getUserVoteMap();
            userVoteMap$.current.subscribe({
                next: (voteMap: Map<number, number>) => setUserVoteMap(voteMap),
                error: (err) => console.error(err)
            })
        }
    }, [posts, user])

    // When board slug/user changes, get the boards of which the user is a member
    useEffect(() => {
        let joinedBoards$: Subject<IBoard[]>;
        if (user) {
            joinedBoards$ = getJoinedBoards();
            joinedBoards$.subscribe({
                next: (boards: IBoard[]) => {
                    boards.forEach((board: IBoard) => {
                        if (board.slug == boardSlug) {
                            setUserIsBoardMember(true);
                        }
                    })
                },
                error: (err) => console.error(err),
            })
        }
        return () => {
            joinedBoards$?.unsubscribe();
        }
    }, [boardSlug, user]);

    return (
        <>
            {user && boardSlug && (
                <section>
                    <div className="flex w-full justify-end">
                        {userIsBoardMember ? (
                            <Button onClick={() => navigateToNewPostPage()}>
                                <i className="pi pi-envelope mr-2"></i> Create Post
                            </Button>
                        ) : (
                            <Button
                                onClick={() => joinBoard()}
                                disabled={isLoading}><i className="pi pi-users mr-2"></i>
                                Join
                            </Button>
                        )}
                    </div>
                </section>
            )}
            <section className="flex justify-space-between gap-2 w-full">
                <section className="w-3/4">
                    <PostList posts={posts} voteMap={userVoteMap} replyMap={totalPostReplyMap}/>
                </section>
                <section className="w-1/4 mt-4">
                    {boardSlug ? (
                        <>
                            {board ?
                                <BoardDetailsSidebar board={board}
                                                     totalPosts={boardTotalPosts}/> : 'Loading board info...'}
                        </>
                    ) : (
                        <BoardListSidebar boards={boards}/>
                    )}
                </section>
            </section>
            <Toast ref={toast}/>
        </>
    )
}