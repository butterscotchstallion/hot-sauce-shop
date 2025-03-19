import {SESSION_URL, USER_URL} from "../Shared/Api.ts";
import {Subject} from "rxjs";
import Cookies from "js-cookie";
import {IUser} from "./IUser.ts";

function setSessionCookie(sessionId: string) {
    Cookies.set("sessionId", sessionId, {
        expires: 30
    });
}

export function getUserBySessionId(): Subject<IUser> {
    const user$ = new Subject<IUser>();
    fetch(`${SESSION_URL}`).then((res: Response) => {
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

export function ValidateUsernameAndPassword(username: string, password: string) {
    const validate$ = new Subject<boolean>();
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
                    if (resp?.results?.sessionId) {
                        setSessionCookie(resp?.results?.sessionId);
                        validate$.next(true);
                    } else {
                        console.error("No session id returned from server");
                        validate$.error("Error signing in");
                    }
                } else {
                    validate$.error(resp?.message || "Unknown error");
                }
            });
        } else {
            validate$.error(res.statusText);
        }
    }).catch((err) => {
        validate$.error(err);
    });
    return validate$;
}