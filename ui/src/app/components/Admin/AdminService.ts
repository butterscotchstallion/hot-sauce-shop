import {Subject} from "rxjs";
import {IUserRole} from "../User/IUserRole.ts";
import {ADMIN_USER_ROLE_LIST} from "../Shared/Api.ts";

export function getRoleList(): Subject<IUserRole[]> {
    const roleList$ = new Subject<IUserRole[]>();
    fetch(ADMIN_USER_ROLE_LIST, {
        credentials: 'include'
    }).then((response: Response) => {
        if (response.ok) {
            response.json().then((results) => {
                if (results.status === "OK") {
                    roleList$.next(results.roles);
                } else {
                    roleList$.error(results.message);
                }
            });
        } else {
            roleList$.error(response.statusText);
        }
    });
    return roleList$;
}