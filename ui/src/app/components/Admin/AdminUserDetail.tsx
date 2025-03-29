import {IUser} from "../User/IUser.ts";
import {useEffect} from "react";

export interface IAdminUserFormProps {
    isNewUser: boolean;
    user: IUser;
}

export function AdminUserDetail(props: IAdminUserFormProps) {
    const roleNameList: string[] = [];

    useEffect(() => {
        
    }, [])

    const roleList = (

    )
    return (
        <>
            <h1 className="text-2xl font-bold w-full mb-4">{props.user.username}</h1>

            <section className="flex gap-4 w-full">
                <aside className={"w-[250px]"}>
                    <img
                        width={'250px'}
                        src={`/images/avatars/${props.user.avatarFilename}`}
                        alt={props.user.username}/>
                </aside>
                <div className={"w-2/3"}>
                    <ul>
                        <li><strong>Created</strong> {new Date(props.user.createdAt).toLocaleDateString()}</li>
                        <li><strong>Roles</strong> {roleList}</li>
                    </ul>
                </div>
            </section>
        </>
    )
}