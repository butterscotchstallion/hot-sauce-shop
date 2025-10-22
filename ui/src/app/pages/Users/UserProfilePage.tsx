import {Params, useParams} from "react-router";
import {ReactElement, useEffect, useRef, useState} from "react";
import {IUser} from "../../components/User/IUser.ts";
import {Subject} from "rxjs";
import {getUserProfileBySlug} from "../../components/User/UserService.ts";
import {IUserDetails} from "../../components/User/IUserDetails.ts";
import {IUserRole} from "../../components/User/IUserRole.ts";
import TimeAgo from "react-timeago";
import {Card} from "primereact/card";
import {Button} from "primereact/button";
import {UserRoleList} from "../../components/User/UserRoleList.tsx";

export default function UserProfilePage() {
    const params: Readonly<Params<string>> = useParams();
    const userSlug: string | undefined = params?.slug;
    const [user, setUser] = useState<IUser | null>(null);
    const [userRoles, setUserRoles] = useState<IUserRole[]>([]);
    const createdAtFormatted = useRef<string | null>(null);
    const user$ = useRef<Subject<IUserDetails> | null>(null);
    const [userPostCount, setUserPostCount] = useState<number>(0);
    const userAvatar: ReactElement = (
        user && user.avatarFilename ? <>
            <aside className={"w-[250px]"}>
                <img
                    width={'250px'}
                    src={`/images/avatars/${user.avatarFilename}`}
                    alt={user.username}/>
            </aside>
        </> : <></>
    )

    useEffect(() => {
        if (userSlug) {
            user$.current = getUserProfileBySlug(userSlug);
            user$.current.subscribe({
                next: (details: IUserDetails) => {
                    setUser(details.user);
                    setUserRoles(details.roles);
                    setUserPostCount(details.userPostCount);
                    createdAtFormatted.current = new Date(details.user.createdAt).toLocaleDateString();
                },
                error: (error: Error) => console.error(error),
            })
            return () => {
                user$.current?.unsubscribe();
            }
        }
    }, [userSlug]);

    return (
        <>
            <h1 className="text-3xl font-bold mb-4">User Profile</h1>

            {user && (
                <section className="flex gap-4 w-full">
                    {userAvatar}
                    <div className={"w-2/3"}>
                        <Card>
                            <section className="flex justify-between">
                                <ul>
                                    <li className="mb-2">
                                        <strong
                                            className="pr-2 mb-1 block">Created</strong>
                                        {<TimeAgo date={user.createdAt}
                                                  title={createdAtFormatted.current}/>}
                                    </li>
                                    {/*<li className="mb-2">*/}
                                    {/*    <strong*/}
                                    {/*        className="pr-2 mb-1 block">Last*/}
                                    {/*        Updated</strong> {user?.updatedAt ? createdAtFormatted.current : 'never'}*/}
                                    {/*</li>*/}
                                    <li className="mb-2">
                                        <strong
                                            className="pr-2 mb-1 block">Karma</strong> 570,260
                                    </li>
                                    <li className="mb-2">
                                        <strong
                                            className="pr-2 mb-1 block">Posts</strong> {userPostCount}
                                    </li>
                                    <li className="mb-2">
                                        <strong
                                            className="pr-2 mb-1 block">Roles</strong> <UserRoleList roles={userRoles}/>
                                    </li>
                                </ul>

                                <div>
                                    <Button label="Follow" icon="pi pi-user-plus" className="mr-2"></Button>
                                </div>
                            </section>
                        </Card>
                    </div>
                </section>
            )}
        </>
    )
}