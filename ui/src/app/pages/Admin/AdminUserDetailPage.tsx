import {AdminUserForm} from "../../components/Admin/AdminUserForm.tsx";
import {useEffect, useState} from "react";
import {IUser} from "../../components/User/IUser.ts";
import {Params, useParams} from "react-router";

export interface IAdminUserPageProps {
    isNewUser: boolean;
}

export function AdminUserDetailPage(props: IAdminUserPageProps) {
    const [user, setUser] = useState<IUser | null>(null);
    const params: Readonly<Params<string>> = useParams();
    const userSlug: string | undefined = params?.slug;
    
    useEffect(() => {
        getUserBySlug()
    }, []);

    return (
        <>
            <section className="flex">
                <AdminUserForm isNewUser={props.isNewUser} user={user}/>
            </section>
        </>
    )
}