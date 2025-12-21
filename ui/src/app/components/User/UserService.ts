import {SESSION_URL, USER_BOARD_ADMIN_LIST, USER_BOARDS_URL, USER_PROFILE_URL, USER_URL} from "../Shared/Api.ts";
import {Subject} from "rxjs";
import Cookies from "js-cookie";
import {IUser} from "./IUser.ts";
import {IUserDetails} from "./IUserDetails.ts";
import {IUserRole} from "./IUserRole.ts";
import {IBoard} from "../Boards/types/IBoard.ts";

export enum UserRole {
    USER_ADMIN = "User Admin",
    REVIEWER = "Reviewer",
    SUPER_MESSAGE_BOARD_ADMIN = "Super Message Board Admin",
}

export interface ISignInResponse {
    user: IUser;
    sessionId: string;
    error?: string;
}

function setSessionCookie(sessionId: string) {
    Cookies.set("sessionId", sessionId, {
        expires: 30
    });
}

export function removeSessionCookie() {
    Cookies.remove("sessionId");
}

export function userHasRole(role: UserRole, roles: IUserRole[]): boolean {
    for (let j = 0; j < roles.length; j++) {
        if (roles[j].name === role) {
            return true;
        }
    }
    return false;
}

export function isReviewer(roles: IUserRole[]): boolean {
    return userHasRole(UserRole.REVIEWER, roles);
}

export function isUserAdmin(roles: IUserRole[]): boolean {
    return userHasRole(UserRole.USER_ADMIN, roles);
}

export function isSuperMessageBoardAdmin(roles: IUserRole[]): boolean {
    return userHasRole(UserRole.SUPER_MESSAGE_BOARD_ADMIN, roles);
}

export function getUsers(): Subject<IUser[]> {
    const users$ = new Subject<IUser[]>();
    fetch(`${USER_URL}`, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                users$.next(resp?.results?.users || null);
            });
        } else {
            users$.error(res.statusText);
        }
    }).catch((err) => {
        users$.error(err);
    });
    return users$;
}

export function getUserDetailsBySessionId(): Subject<IUserDetails> {
    const user$ = new Subject<IUserDetails>();
    fetch(`${SESSION_URL}`, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    user$.next(resp.results);
                } else {
                    user$.error(resp?.message || "Unknown error");
                }
            }).catch((err) => {
                user$.error(err);
            });
        } else {
            user$.error(res.statusText);
        }
    }).catch((err) => {
        user$.error(err);
    });
    return user$;
}

export function ValidateUsernameAndPassword(username: string, password: string): Subject<ISignInResponse> {
    const validate$ = new Subject<ISignInResponse>();
    fetch(`${USER_URL}/sign-in`, {
        method: 'POST',
        body: JSON.stringify({
            username: username,
            password: password
        }),
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    if (resp?.results?.sessionId && resp?.results?.user) {
                        setSessionCookie(resp?.results?.sessionId);
                        validate$.next(resp.results);
                    } else {
                        console.error("No session id returned from server");
                        validate$.error("Error signing in");
                    }
                } else {
                    validate$.error({
                        error: resp?.message || "Unknown error"
                    });
                }
            });
        } else {
            validate$.error({
                error: res.statusText
            });
        }
    }).catch((err) => {
        validate$.error({
            error: err
        });
    });
    return validate$;
}

export function getUserProfileBySlug(slug: string) {
    const user$ = new Subject<IUserDetails>();
    fetch(`${USER_PROFILE_URL}/${slug}`, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    user$.next(resp.results);
                } else {
                    res.json().then(resp => {
                        user$.error(resp?.message || "Unknown error");
                    });
                }
            })
        } else {
            user$.error(res.statusText);
        }
    });
    return user$;
}

export function getJoinedBoards(): Subject<IBoard[]> {
    const boards$ = new Subject<IBoard[]>();
    fetch(USER_BOARDS_URL, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    boards$.next(resp.results.boards);
                } else {
                    res.json().then(resp => {
                        boards$.error(resp?.message || "Unknown error");
                    });
                }
            })
        } else {
            boards$.error(res.statusText);
        }
    });
    return boards$;
}

export function userJoinBoard(boardId: number): Subject<boolean> {
    const joinBoard$ = new Subject<boolean>();
    fetch(`${USER_BOARDS_URL}/${boardId}`, {
        method: 'POST',
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    joinBoard$.next(true);
                } else {
                    joinBoard$.error({
                        error: resp?.message || "Unknown error"
                    });
                }
            });
        } else {
            joinBoard$.error({
                error: res.statusText
            });
        }
    }).catch((err) => {
        joinBoard$.error({
            error: err
        });
    });
    return joinBoard$;
}

export function getUserAdminBoards(): Subject<IBoard[]> {
    const boards$ = new Subject<IBoard[]>();
    fetch(USER_BOARD_ADMIN_LIST, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    boards$.next(resp.results.boards);
                } else {
                    boards$.error(resp?.message || "Unknown error");
                }
            }).catch((err) => {
                boards$.error(err);
            });
        } else {
            boards$.error(res.statusText);
        }
    }).catch((err) => {
        boards$.error(err);
    });
    return boards$;
}