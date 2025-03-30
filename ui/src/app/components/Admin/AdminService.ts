import {Subject} from "rxjs";
import {IUserRole} from "../User/IUserRole.ts";
import {ADMIN_USER_DETAIL_URL, ADMIN_USER_ROLE_LIST} from "../Shared/Api.ts";
import {IUser} from "../User/IUser.ts";

export function updateUser(user: IUser, roles: IUserRole[]): Subject<boolean> {
    const updateUser$ = new Subject<boolean>();
    fetch(`${ADMIN_USER_DETAIL_URL}/${user.slug}`, {
        method: 'PUT',
        credentials: 'include',
        body: JSON.stringify({
            user: user,
            roles: roles
        }),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then((response: Response) => {
        if (response.ok) {
            response.json().then((results) => {
                if (results.status === "OK") {
                    updateUser$.next(true);
                } else {
                    updateUser$.error(results.message);
                }
            });
        } else {
            updateUser$.error(response.statusText);
        }
    });
    return updateUser$;
}

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