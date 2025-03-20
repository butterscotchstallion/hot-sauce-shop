import {SESSION_URL, USER_URL} from "../Shared/Api.ts";
import {Subject, Subscription} from "rxjs";
import Cookies from "js-cookie";
import {IUser} from "./IUser.ts";
import {setSignedIn, setUser} from "./User.slice.ts";
import {Dispatch, UnknownAction} from "@reduxjs/toolkit";
import {useDispatch} from "react-redux";

function setSessionCookie(sessionId: string) {
    Cookies.set("sessionId", sessionId, {
        expires: 30
    });
}

function removeSessionCookie() {
    Cookies.remove("sessionId");
}

export function getUserBySessionIdAndStore(dispatch: Dispatch<UnknownAction>): Subscription {
    return getUserBySessionId().subscribe({
        next: (user: IUser) => {
            dispatch(setUser(user));
            dispatch(setSignedIn(true));
            console.info("Set user to " + user.username)
        },
        error: (err) => {
            console.error("Error getting session from DB: " + err);
        }
    });
}

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

export function ValidateUsernameAndPassword(username: string, password: string) {
    const validate$ = new Subject<boolean>();
    const dispatch = useDispatch();
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
                        dispatch(setUser(resp?.results?.user));
                        dispatch(setSignedIn(true));
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