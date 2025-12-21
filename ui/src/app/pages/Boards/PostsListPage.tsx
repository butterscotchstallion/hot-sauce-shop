import {ReactElement, RefObject, SetStateAction, useEffect, useRef, useState} from "react";
import {IBoardPost} from "../../components/Boards/types/IBoardPost.ts";
import {getBoardByBoardSlug, getBoards, getPosts, getTotalPostReplyMap} from "../../components/Boards/BoardsService.ts";
import PostList from "../../components/Boards/PostList.tsx";
import {getUserVoteMap} from "../../components/Boards/VoteService.ts";
import {Subject} from "rxjs";
import {NavigateFunction, Params, useNavigate, useParams} from "react-router";
import {Button} from "primereact/button";
import {IUser} from "../../components/User/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {IBoard} from "../../components/Boards/types/IBoard.ts";
import {getJoinedBoards, userJoinBoard} from "../../components/User/UserService.ts";
import {Toast} from "primereact/toast";
import {BoardDetailsSidebar} from "../../components/Boards/BoardDetailsSidebar.tsx";
import {BoardListSidebar} from "../../components/Boards/BoardListSidebar.tsx";
import {setPageTitle} from "../../components/Shared/PageTitle.ts";
import {IBoardDetails} from "../../components/Boards/types/IBoardDetails.ts";
import {Paginator} from "primereact/paginator";
import {IBoardPostsResponse} from "../../components/Boards/types/IBoardPostsResponse.ts";
import {Skeleton} from "primereact/skeleton";

/**
 * Handles multiple scenarios where post(s) are displayed:
 * - All posts, unfiltered
 * - Filtered posts
 * - A single post
 */
export default function PostsListPage() {
    const [postsLoading, setPostsLoading] = useState<boolean>(true);
    const [totalPosts, setTotalPosts] = useState<number>(0);
    const [offset, setOffset] = useState<number>(0);
    const [perPage, setPerPage] = useState<number>(10);
    const toast: RefObject<Toast | null> = useRef(null);
    const params: Readonly<Params<string>> = useParams();
    const boardSlug: string = params?.boardSlug || '';
    const postSlug: string = params?.postSlug || '';
    const [posts, setPosts] = useState<IBoardPost[]>([]);
    const [postReplies, setPostReplies] = useState<IBoardPost[]>([]);
    const [board, setBoard] = useState<IBoard>();
    const [userIsBoardMember, setUserIsBoardMember] = useState<boolean>(false);
    const userVoteMap$ = useRef<Subject<Map<number, number>> | null>(null);
    const [userVoteMap, setUserVoteMap] = useState<Map<number, number>>(new Map());
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const navigate: NavigateFunction = useNavigate();
    const navigateToNewPostPage = () => {
        let replyParam: string = "";
        // If there is a post slug, we're viewing a specific post and there should
        // only be one post in the list
        if (postSlug) {
            replyParam += "?parentId=" + posts[0].id;
        }
        const url = `/boards/${boardSlug}/posts/new${replyParam}`;
        navigate(url);
    }
    const joinBoard$ = useRef<Subject<boolean>>(null);
    const [isLoading, setIsLoading] = useState(false);
    const [boards, setBoards] = useState<IBoard[]>([]);
    const [totalPostReplyMap, setTotalPostReplyMap] = useState<Map<number, number>>(new Map());
    const [isCurrentUserBoardMod, setIsCurrentUserBoardMod] = useState<boolean>(false);
    const [boardDetails, setBoardDetails] = useState<IBoardDetails | undefined>();
    const onPageChange = (event: { first: SetStateAction<number>; rows: SetStateAction<number>; }) => {
        setOffset(event.first);
        setPerPage(event.rows);
    };

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

    const isCurrentUserInBoardMods = (mods: IUser[]): boolean => {
        for (const mod of mods) {
            if (mod.id === user?.id) {
                return true;
            }
        }
        return false;
    }
    const skeletonList: ReactElement[] = new Array(10).fill(0).map((): ReactElement => (
        <Skeleton size="260px"></Skeleton>
    ));

    // When the board/post changes, get filtered posts based on the scenario
    useEffect(() => {
        console.log("Fetching posts for boardSlug: " + boardSlug + " and postSlug: " + postSlug);
        const postsResponse$: Subject<IBoardPostsResponse> = getPosts({
            postSlug,
            boardSlug,
            offset,
            perPage,
        });
        postsResponse$.subscribe({
            next: (postsResponse: IBoardPostsResponse) => {
                setPostsLoading(false);
                setPosts(postsResponse.posts);
                setTotalPosts(postsResponse.totalPosts);
            },
            error: (err) => {
                console.error(err);
                setPostsLoading(false);
            }
        });
        // When viewing a specific board...
        let board$: Subject<IBoardDetails>;
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
                next: (boardDetails: IBoardDetails) => {
                    setPageTitle(boardDetails.board.displayName);
                    setBoard(boardDetails.board);
                    setIsCurrentUserBoardMod(isCurrentUserInBoardMods(boardDetails.moderators));
                    setBoardDetails(boardDetails);
                },
                error: (err) => console.error(err),
            });
        } else {
            // on "All posts" page
            setPageTitle("All Posts");
            getBoards$ = getBoards();
            getBoards$.subscribe({
                next: (boardList: IBoard[]) => setBoards(boardList),
                error: (err) => console.error(err),
            })
        }

        return () => {
            postsResponse$.unsubscribe();
            board$?.unsubscribe();
            getBoards$?.unsubscribe();
            replyMap$.unsubscribe();

            if (joinBoard$.current) {
                joinBoard$.current.unsubscribe();
            }
            if (userVoteMap$.current) (
                userVoteMap$.current.unsubscribe()
            )
        }
    }, [boardSlug, postSlug, offset]);

    useEffect(() => {
        let getPostReplies$: Subject<IBoardPostsResponse>;
        // Viewing a specific post
        if (postSlug && posts.length === 1) {
            setPageTitle(posts[0].title);
            const postID: number = posts[0].id;
            console.info(`Fetching replies for ${postID}`)
            getPostReplies$ = getPosts({
                parentId: postID
            })
            getPostReplies$.subscribe({
                next: (repliesResponse: IBoardPostsResponse) => setPostReplies(repliesResponse.posts),
                error: (err) => console.error(err),
            })
        }
        return () => {
            getPostReplies$?.unsubscribe();
        }
    }, [posts]);

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
                <section className="flex justify-between gap-2">
                    {board && <h1 className="text-3xl font-bold mb-4">{board.displayName}</h1>}
                    <div className="justify-end flex gap-2 items-center">
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
                    {isLoading ? skeletonList : <></>}
                    {posts.length > 0 ? (
                        <>
                            <PostList
                                posts={posts}
                                voteMap={userVoteMap}
                                replyMap={totalPostReplyMap}
                                isCurrentUserBoardMod={isCurrentUserBoardMod}
                            />
                            {posts.length >= perPage ? (
                                <div className="card mt-4 mb-4">
                                    <Paginator first={offset}
                                               rows={perPage}
                                               totalRecords={totalPosts}
                                               rowsPerPageOptions={[10, 20, 30]}
                                               onPageChange={onPageChange}/>
                                </div>
                            ) : ""}
                        </>
                    ) : "No posts available."}
                </section>
                <section className="w-1/4 mt-4">
                    {boardSlug ? (
                        <>
                            {board ?
                                <BoardDetailsSidebar boardDetails={boardDetails}/> : 'Loading board info...'}
                        </>
                    ) : (
                        <BoardListSidebar boards={boards}/>
                    )}
                </section>
            </section>

            {postReplies.length > 0 && (
                <section className="mt-4">
                    <h1 className="text-3xl font-bold mb-4">Comments</h1>
                    <section className="w-3/4">
                        <PostList
                            posts={postReplies}
                            voteMap={userVoteMap}
                            replyMap={totalPostReplyMap}
                            isCurrentUserBoardMod={isCurrentUserBoardMod}
                        />
                    </section>
                </section>
            )}

            {postSlug && postReplies.length === 0 && (
                <p>No comments on this post yet. Be the first!</p>
            )}

            <Toast ref={toast}/>
        </>
    )
}