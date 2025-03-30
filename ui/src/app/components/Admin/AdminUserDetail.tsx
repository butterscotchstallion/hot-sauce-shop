import {IUser} from "../User/IUser.ts";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {IUserRole} from "../User/IUserRole.ts";
import {ReactElement, RefObject, useEffect, useRef, useState} from "react";
import {PickList} from "primereact/picklist";
import {Card} from "primereact/card";
import {getRoleList, updateUser} from "./AdminService.ts";
import {Subscription} from "rxjs";
import {Button} from "primereact/button";
import {Toast} from "primereact/toast";

export interface IAdminUserFormProps {
    isNewUser: boolean;
    user: IUser;
}

export function AdminUserDetail(props: IAdminUserFormProps) {
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const userRoles: IUserRole[] = useSelector((state: RootState) => state.user.roles);
    const [sourceRoles, setSourceRoles] = useState<IUserRole[]>([]);
    const [targetRoles, setTargetRoles] = useState<IUserRole[]>([]);
    const userAvatar: ReactElement = (
        props.user.avatarFilename ? <>
            <aside className={"w-[250px]"}>
                <img
                    width={'250px'}
                    src={`/images/avatars/${props.user.avatarFilename}`}
                    alt={props.user.username}/>
            </aside>
        </> : <></>
    )
    const onChange = (event) => {
        setSourceRoles(event.source);
        setTargetRoles(event.target);
    };

    const itemTemplate = (role: IUserRole) => {
        return (
            <div className="flex flex-wrap p-2 align-items-center gap-3">
                <div className="flex-1 flex flex-column gap-2">
                    <span className="font-bold">{role.name}</span>
                </div>
            </div>
        );
    };

    useEffect(() => {
        setTargetRoles(userRoles);
        const roles$: Subscription = getRoleList().subscribe({
            next: (roles: IUserRole[]) => setSourceRoles(roles),
            error: (err) => {
                console.error(err);
            }
        });
        return () => {
            roles$.unsubscribe();
        }
    }, []);

    const save = () => {
        updateUser(props.user, targetRoles).subscribe({
            next: () => {
                if (toast.current) {
                    toast.current.show({
                        severity: 'success',
                        summary: 'Success',
                        detail: 'User saved successfully',
                        life: 3000,
                    })
                }
            },
            error: (err) => {
                if (toast.current) {
                    toast.current.show({
                        severity: 'error',
                        summary: 'Error',
                        detail: 'Error saving user: ' + err,
                        life: 3000,
                    })
                }
            }
        });
    }

    return (
        <>
            <section className="flex justify-between mb-4">
                <h1 className="text-2xl font-bold w-full mb-4">{props.user.username}</h1>

                <Button onClick={() => save()} label="Save" icon="pi pi-save"></Button>
            </section>

            <section className="flex gap-4 w-full">
                {userAvatar}
                <div className={"w-2/3"}>
                    <ul>
                        <li className="mb-2">
                            <strong
                                className="pr-2">Created</strong> {new Date(props.user.createdAt).toLocaleDateString()}
                        </li>
                        <li>
                            <Card title="User Roles">
                                <PickList dataKey="id"
                                          source={sourceRoles}
                                          target={targetRoles}
                                          onChange={onChange}
                                          itemTemplate={itemTemplate}
                                          breakpoint="1280px"
                                          sourceHeader="Available"
                                          targetHeader="Selected"/>
                            </Card>
                        </li>
                    </ul>
                </div>
            </section>
            <Toast ref={toast}/>
        </>
    )
}