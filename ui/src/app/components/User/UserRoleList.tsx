import {IUserRole} from "./IUserRole.ts";
import {useEffect, useState} from "react";
import {Tag} from "primereact/tag";

interface IUserRoleListProps {
    roles: IUserRole[];
}

export function UserRoleList({roles}: IUserRoleListProps) {
    const [roleNames, setRoleNames] = useState<string[]>([]);

    useEffect(() => {
        setRoleNames(roles.map((role: IUserRole) => role.name));
    }, [roles]);

    return (
        <>
            <ul>
                {roleNames.map((roleName: string) => (
                    <li className="mb-2 inline-block mr-2" key={roleName}>
                        <Tag value={roleName} severity="info"></Tag>
                    </li>
                ))}
            </ul>
        </>
    )
}