import {Subject} from "rxjs";
import {IUser} from "./IUser.ts";
import {ADMIN_USER_DETAIL_URL} from "../Shared/Api.ts";

export function getUserBySlug(slug: string): Subject<IUser> {
    const user$ = new Subject<IUser>();
    fetch(`${ADMIN_USER_DETAIL_URL}/${slug}`).then((response: Response) => {
        if (response.ok) {
            response.json().then((results) => {
                if (results.status !== "OK") {
                    user$.next(results.user);
                } else {
                    user$.error(results.message);
                }
            });
        } else {
            user$.error(response.statusText);
        }
    }).catch((err) => {
        user$.error(err);
    });
    return user$;
}