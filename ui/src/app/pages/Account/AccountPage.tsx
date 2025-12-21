import {IUser} from "../../components/User/types/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {UserRoleList} from "../../components/User/UserRoleList.tsx";
import {IUserRole} from "../../components/User/types/IUserRole.ts";

export function AccountPage() {
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const roles: IUserRole[] = useSelector((state: RootState) => state.user.roles);
    return (
        <>
            <h1 className="text-3xl font-bold mb-4">Account Settings</h1>

            <section className="mt-4">
                {user ? (
                    <table className="w-1/4">
                        <tbody>
                        <tr>
                            <td className="w-1/2 pb-2">
                                <strong>Username</strong>
                            </td>
                            <td className="pb-2">
                                {user.username}
                            </td>
                        </tr>
                        <tr>
                            <td className="w-1/2 pb-2">
                                <strong>Created</strong>
                            </td>
                            <td className="pb-2">
                                {new Date(user.createdAt).toLocaleDateString()}
                            </td>
                        </tr>
                        <tr>
                            <td className="w-1/2 pb-2 align-top">
                                <strong>Roles</strong>
                            </td>
                            <td className="pb-2">
                                {roles ? (
                                    <UserRoleList roles={roles}/>
                                ) : 'No user roles available'}
                            </td>
                        </tr>
                        </tbody>
                    </table>
                ) : 'Error getting user details.'}
            </section>
        </>
    )
}