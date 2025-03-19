import {USER_URL} from "../Shared/Api.ts";
import {Subject} from "rxjs";

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
                validate$.next(resp?.status === "OK");
            });
        } else {
            validate$.error(res.statusText);
        }
    }).catch((err) => {
        validate$.error(err);
    });
    return validate$;
}