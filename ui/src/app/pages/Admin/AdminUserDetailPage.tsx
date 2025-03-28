import {AdminUserForm} from "../../components/Admin/AdminUserForm.tsx";
import {RefObject, useEffect, useRef, useState} from "react";
import {IUser} from "../../components/User/IUser.ts";
import {Params, useParams} from "react-router";
import {getUserBySlug} from "../../components/User/AdminService.ts";
import {Subscription} from "rxjs";
import {Toast} from "primereact/toast";

export interface IAdminUserPageProps {
    isNewUser: boolean;
}

export function AdminUserDetailPage(props: IAdminUserPageProps) {
    const [user, setUser] = useState<IUser | null>(null);
    const params: Readonly<Params<string>> = useParams();
    const userSlug: string | undefined = params?.slug;
    const toast: RefObject<Toast | null> = useRef<Toast | null>(null);

    useEffect(() => {
        let user$: Subscription;
        if (userSlug) {
            user$ = getUserBySlug(userSlug).subscribe((user: IUser) => setUser(user));
        } else {
            if (toast.current) {
                toast.current.show({severity: 'error', summary: 'Error', detail: 'User not found'});
            }
        }
        return () => {
            user$.unsubscribe();
        }
    }, [userSlug]);

    return (
        <>
            <section className="flex">
                <AdminUserForm isNewUser={props.isNewUser} user={user}/>
            </section>

            <Toast ref={toast}/>
        </>
    )
}