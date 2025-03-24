import {AdminUserForm} from "../../components/Admin/AdminUserForm.tsx";

export interface IAdminUserPageProps {
    isNewUser: boolean;
}

export function AdminUserPage(props: IAdminUserPageProps) {
    return (
        <>
            <section className="flex">
                <AdminUserForm isNewUser={props.isNewUser} user={user}/>
            </section>
        </>
    )
}