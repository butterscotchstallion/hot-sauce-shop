import {IUserRole} from "./IUserRole.ts";
import {RefObject, useEffect, useRef} from "react";
import {Tag} from "primereact/tag";

interface IUserRoleListProps {
    roles: IUserRole[];
}

export function UserRoleList(props: IUserRoleListProps) {
    const roleNames: RefObject<string[]> = useRef<string[]>([]);

    useEffect(() => {
        roleNames.current = props.roles.map((role: IUserRole) => role.name);
    }, [props.roles]);

    return (
        <>
            <ul>
                {roleNames.current.map((roleName: string) => (
                    <li className="mb-2"><Tag key={roleName} value={roleName}></Tag></li>
                ))}
            </ul>
        </>
    )
}