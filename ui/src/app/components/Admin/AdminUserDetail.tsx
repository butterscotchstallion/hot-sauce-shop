import {IUser} from "../User/IUser.ts";

export interface IAdminUserFormProps {
    isNewUser: boolean;
    user: IUser | null;
}

export function AdminUserDetail(props: IAdminUserFormProps) {

    return (
        <>
            <h1 className="text-2xl font-bold">{props.user.username}</h1>
        </>
    )
}