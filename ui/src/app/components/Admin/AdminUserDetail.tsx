import {IUser} from "../User/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {IUserRole} from "../User/IUserRole.ts";
import {ReactElement} from "react";
import {Tag} from "primereact/tag";

export interface IAdminUserFormProps {
    isNewUser: boolean;
    user: IUser;
}

export function AdminUserDetail(props: IAdminUserFormProps) {
    const userRoles: IUserRole[] = useSelector((state: RootState) => state.user.roles);
    const userRoleNameList: string[] = userRoles.map((role: IUserRole) => role.name);
    const userRoleList: ReactElement[] = (
        userRoleNameList.map((roleName: string) => <Tag severity="info" value={roleName}></Tag>)
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
                        <li className="mb-2">
                            <strong
                                className="pr-2">Created</strong> {new Date(props.user.createdAt).toLocaleDateString()}
                        </li>
                        <li>
                            <strong className="pr-2">Roles</strong> {userRoleList || 'No roles set'}
                        </li>
                    </ul>
                </div>
            </section>
        </>
    )
}