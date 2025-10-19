import {IBoard} from "./IBoard.ts";
import {Subject} from "rxjs";
import {BOARD_DETAILS_URL, BOARD_POSTS_URL, BOARDS_URL, POST_DETAILS_URL, POSTS_URL} from "../Shared/Api.ts";
import {IBoardPost} from "./IBoardPost.ts";
import {INewBoardPost} from "./INewBoardPost.ts";

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

export function getPostsByBoardSlug(boardSlug: string): Subject<IBoardPost[]> {
    const posts$ = new Subject<IBoardPost[]>();
    fetch(BOARD_POSTS_URL.replace(':slug', boardSlug), {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                posts$.next(resp.results.posts);
            });
        } else {
            posts$.error(res.statusText);
        }
    }).catch((err) => {
        posts$.error(err);
    });
    return posts$;
}

export function getPosts(): Subject<IBoardPost[]> {
    const posts$ = new Subject<IBoardPost[]>();
    fetch(POSTS_URL, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                posts$.next(resp.results.posts);
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

export function getBoardByBoardSlug(boardSlug: string): Subject<IBoard> {
    const board$ = new Subject<IBoard>();
    fetch(BOARD_DETAILS_URL.replace(':slug', boardSlug), {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                board$.next(resp.results.board);
            });
        } else {
            board$.error(res.statusText);
        }
    }).catch((err) => {
        board$.error(err);
    });
    return board$;
}

export function addPost(post: INewBoardPost, boardSlug: string): Subject<IBoardPost> {
    const addPost$ = new Subject<IBoardPost>();
    fetch(BOARD_POSTS_URL.replace(':slug', boardSlug), {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify(post),
        headers: {
            'Content-Type': 'application/json'
        }
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
