import {Params, useParams} from "react-router";
import {ReactElement, useEffect, useRef, useState} from "react";
import {Subject} from "rxjs";
import {getUserProfileBySlug} from "../../components/User/UserService.ts";
import {IUserDetails} from "../../components/User/IUserDetails.ts";
import TimeAgo from "react-timeago";
import {Card} from "primereact/card";
import {Button} from "primereact/button";
import {UserRoleList} from "../../components/User/UserRoleList.tsx";

export default function UserProfilePage() {
    const params: Readonly<Params<string>> = useParams();
    const userSlug: string | undefined = params?.slug;
    const createdAtFormatted = useRef<string>(undefined);
    const user$ = useRef<Subject<IUserDetails> | null>(null);
    const [details, setDetails] = useState<IUserDetails>();

    const userAvatar: ReactElement = (
        details?.user && details.user.avatarFilename ? <>
            <aside className={"w-[250px]"}>
                <img
                    width={'250px'}
                    src={`/images/avatars/${details.user.avatarFilename}`}
                    alt={details.user.username}/>
            </aside>
        </> : <></>
    )

    useEffect(() => {
        if (userSlug) {
            user$.current = getUserProfileBySlug(userSlug);
            user$.current.subscribe({
                next: (details: IUserDetails) => {
                    setDetails(details);
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

            {details && (
                <section className="flex gap-4 w-full">
                    {userAvatar}
                    <div className={"w-2/3"}>
                        <Card>
                            <section className="flex justify-between">
                                <ul>
                                    <li className="mb-2">
                                        <strong
                                            className="pr-2 mb-1 block">Created</strong>
                                        {<TimeAgo date={details.user.createdAt}
                                                  title={createdAtFormatted.current}/>}
                                    </li>
                                    <li className="mb-2">
                                        <strong
                                            className="pr-2 mb-1 block">Karma</strong> {details.userPostVoteSum}
                                    </li>
                                    <li className="mb-2">
                                        <strong
                                            className="pr-2 mb-1 block">Posts</strong> {details.userPostCount}
                                    </li>
                                    <li className="mb-2">
                                        <strong className="pr-2 mb-1 block">Roles</strong>
                                        {details.roles ? <UserRoleList roles={details.roles}/> : 'No rules assigned'}
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