import {IBoard} from "./IBoard.ts";
import {Subject} from "rxjs";
import {
    BOARD_DETAILS_URL,
    BOARD_PIN_POST_URL,
    BOARD_POSTS_URL,
    BOARD_SETTINGS_URL,
    BOARD_TOTAL_POSTS_URL,
    BOARD_TOTAL_REPLIES,
    BOARDS_URL,
    POST_DETAILS_URL,
    POSTS_URL
} from "../Shared/Api.ts";
import {IBoardPost} from "./IBoardPost.ts";
import {INewBoardPost} from "./INewBoardPost.ts";
import {IBoardDetails} from "./IBoardDetails.ts";
import {IBoardSettings} from "./IBoardSettings.ts";
import {getUserAdminBoards} from "../User/UserService.ts";
import {IBoardPostsResponse} from "./IBoardPostsResponse.ts";

export function getBoards(): Subject<IBoard[]> {
    const boards$ = new Subject<IBoard[]>();
    fetch(`${BOARDS_URL}`, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                boards$.next(resp?.results?.boards || null);
            });
        } else {
            boards$.error(res.statusText);
        }
    }).catch((err) => {
        boards$.error(err);
    });
    return boards$;
}

interface IGetPostParameters {
    parentId?: number;
    boardSlug?: string;
    postSlug?: string;
    offset?: number;
    perPage?: number;
}

export function getPosts({
                             parentId,
                             boardSlug,
                             postSlug,
                             offset = 0,
                             perPage = 10
                         }: IGetPostParameters): Subject<IBoardPostsResponse> {
    const posts$ = new Subject<IBoardPostsResponse>();
    let url: string = POSTS_URL;

    /**
     * - If there's a parentId, we're getting replies to a post.
     * - If there's a boardSlug, we're getting posts from a board.
     * - If there's a postSlug, we're getting a single post.
     */
    if (parentId && parentId > 0) {
        url += `?parentId=${parentId}`;
    } else {
        // All posts on board
        if (boardSlug && !postSlug) {
            url += `?boardSlug=${boardSlug}`;
        }
        // Post detail
        if (boardSlug && postSlug) {
            url += `?boardSlug=${boardSlug}&postSlug=${postSlug}`;
        }
    }

    if (url.indexOf("?") !== -1) {
        url += "&";
    } else {
        url += "?";
    }
    url += `offset=${offset}&perPage=${perPage}`;

    fetch(url, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                posts$.next(resp.results);
            });
        } else {
            posts$.error(res.statusText);
        }
    }).catch((err) => {
        posts$.error(err);
    });
    return posts$;
}

export function getPostDetail(boardSlug: string, postSlug: string): Subject<IBoardPost> {
    const post$ = new Subject<IBoardPost>();
    let url: string = POST_DETAILS_URL.replace(':boardSlug', boardSlug);
    url = url.replace(':postSlug', postSlug);
    fetch(url, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                post$.next(resp.results.post);
            });
        } else {
            res.json().then(resp => {
                post$.error(resp?.message || "Unknown error");
            });
        }
    }).catch((err) => {
        post$.error(err);
    });
    return post$;
}

export function getBoardByBoardSlug(boardSlug: string): Subject<IBoardDetails> {
    const board$ = new Subject<IBoardDetails>();
    fetch(BOARD_DETAILS_URL.replace(':slug', boardSlug), {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                board$.next(resp.results);
            });
        } else {
            board$.error(res.statusText);
        }
    }).catch((err) => {
        board$.error(err);
    });
    return board$;
}

export function addPost(post: INewBoardPost, boardSlug: string, postImages: File[]): Subject<IBoardPost> {
    const addPost$ = new Subject<IBoardPost>();
    const formData: FormData = new FormData();
    formData.append("title", post.title);
    formData.append("postText", post.postText);
    postImages.forEach(image => {
        formData.append("postImages", image);
    });
    if (post?.parentId) {
        formData.append("parentId", post.parentId.toString());
    }
    fetch(BOARD_POSTS_URL.replace(':slug', boardSlug), {
        credentials: 'include',
        method: 'POST',
        body: formData
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                addPost$.next(resp.results.post);
            });
        } else {
            res.json().then(resp => {
                addPost$.error(resp?.message || "Unknown error");
            });
        }
    }).catch((err) => {
        addPost$.error(err);
    });
    return addPost$;
}

export function getTotalPostsByBoardSlug(boardSlug: string): Subject<number> {
    const totalPosts$ = new Subject<number>();
    fetch(`${BOARD_TOTAL_POSTS_URL}/${boardSlug}`, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                totalPosts$.next(resp.results.totalPosts);
            });
        } else {
            totalPosts$.error(res.statusText);
        }
    }).catch((err) => {
        totalPosts$.error(err);
    });
    return totalPosts$;
}

export function getTotalPostReplyMap(boardSlug: string): Subject<Map<number, number>> {
    const replyMap$ = new Subject<Map<number, number>>();
    let url = BOARD_TOTAL_REPLIES;
    if (boardSlug) {
        url += `?boardSlug=${boardSlug}`;
    }
    fetch(url, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                replyMap$.next(resp.results.totalPostReplyMap);
            });
        } else {
            replyMap$.error(res.statusText);
        }
    }).catch((err) => {
        replyMap$.error(err);
    });
    return replyMap$;
}

export function pinPost(post: IBoardPost, boardSlug: string): Subject<boolean> {
    const pinPost$ = new Subject<boolean>();
    fetch(`${BOARD_PIN_POST_URL}/${boardSlug}/${post.slug}`, {
        credentials: 'include',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        }
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(_ => {
                pinPost$.next(true);
            });
        } else {
            res.json().then(resp => {
                pinPost$.error(resp?.message || "Unknown error");
            });
        }
    }).catch((err) => {
        pinPost$.error(err);
    });
    return pinPost$;
}

export function getBoardSettings(boardSlug: string): Subject<IBoardSettings> {
    const settings$ = new Subject<IBoardSettings>();
    fetch(`${BOARD_SETTINGS_URL}/${boardSlug}`, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                settings$.next(resp.results);
            });
        } else {
            settings$.error(res.statusText);
        }
    }).catch((err) => {
        settings$.error(err);
    });
    return settings$;
}

export function isSettingsAreaAvailable(boardSlug: string): Subject<boolean> {
    const settingsAvailable$ = new Subject<boolean>();
    getUserAdminBoards().subscribe({
        next: (boards: IBoard[]) => {
            let isSettingsAreaAvailable = false;
            for (let j = 0; j < boards.length; j++) {
                if (boards[j].slug === boardSlug) {
                    isSettingsAreaAvailable = true;
                    break;
                }
            }
            settingsAvailable$.next(isSettingsAreaAvailable);
        },
        error: (err) => {
            settingsAvailable$.error(err);
        }
    });
    return settingsAvailable$;
}