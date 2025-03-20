import {SESSION_URL, USER_URL} from "../Shared/Api.ts";
import {Subject} from "rxjs";
import Cookies from "js-cookie";
import {IUser} from "./IUser.ts";

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

/*
function removeSessionCookie() {
    Cookies.remove("sessionId");
}
*/

export function getUserBySessionId(): Subject<IUser> {
    const user$ = new Subject<IUser>();
    fetch(`${SESSION_URL}`, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    user$.next(resp?.results?.user || null);
                } else {
                    user$.error(resp?.message || "Unknown error");
                }
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