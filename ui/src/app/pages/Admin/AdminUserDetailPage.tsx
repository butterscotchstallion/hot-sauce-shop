import {AdminUserDetail} from "../../components/Admin/AdminUserDetail.tsx";
import {RefObject, useEffect, useRef, useState} from "react";
import {IUser} from "../../components/User/types/IUser.ts";
import {Params, useParams} from "react-router";
import {getUserBySlug} from "../../components/User/AdminService.ts";
import {Subscription} from "rxjs";
import {Messages} from "primereact/messages";
import {AdminBreadcrumbs} from "../../components/Admin/AdminBreadcrumbs.tsx";

export interface IAdminUserPageProps {
    isNewUser: boolean;
}

export function AdminUserDetailPage(props: IAdminUserPageProps) {
    const [user, setUser] = useState<IUser | null>(null);
    const params: Readonly<Params<string>> = useParams();
    const userSlug: string | undefined = params?.slug;
    const msgs: RefObject<Messages | null> = useRef(null);

    const showErrorUserNotFound = () => {
        if (msgs.current) {
            msgs.current.clear();
            msgs.current.show([
                {
                    severity: 'error',
                    summary: 'Error',
                    detail: 'There was a problem loading the user.',
                    sticky: true,
                    closable: false
                }
            ]);
        }
    };

    useEffect(() => {
        let user$: Subscription;
        if (userSlug) {
            user$ = getUserBySlug(userSlug).subscribe({
                next: (user: IUser) => setUser(user),
                error: () => showErrorUserNotFound()
            });
        } else {
            showErrorUserNotFound();
        }

        return () => {
            user$.unsubscribe();
        }
    }, [userSlug]);

    const adminUserDetail = (
        user ? <AdminUserDetail user={user} isNewUser={props.isNewUser}/> : ''
    )

    return (
        <>
            <section className="mb-2">
                <AdminBreadcrumbs icon={"pi pi-user"} label="User Detail"/>
            </section>

            {adminUserDetail}
            <Messages ref={msgs}/>
        </>
    )
}