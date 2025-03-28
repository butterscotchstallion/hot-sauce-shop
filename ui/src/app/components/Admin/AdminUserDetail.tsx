import {IUser} from "../User/IUser.ts";

export interface IAdminUserFormProps {
    isNewUser: boolean;
    user: IUser | null;
}

export function AdminUserDetail(props: IAdminUserFormProps) {

    return (
        <>
            <div className={"flex flex-row gap-4"}>
                <h1 className="text-2xl font-bold">{props.user.username}</h1>

                <section className="flex gap-4">
                    <aside className={"w-1/3"}>
                        user avatar
                    </aside>
                    <div className={"w-2/3"}>
                        <ul>
                            <li>Created: {new Date(props.user.createdAt).toLocaleDateString()}</li>
                        </ul>
                    </div>
                </section>
            </div>
        </>
    )
}