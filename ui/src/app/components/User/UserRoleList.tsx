import {IUserRole} from "./IUserRole.ts";
import {useEffect, useState} from "react";
import {Tag} from "primereact/tag";

interface IUserRoleListProps {
    roles: IUserRole[];
}

export function UserRoleList({roles}: IUserRoleListProps) {
    const [allRoles, setAllRoles] = useState<IUserRole[]>([]);

    useEffect(() => {
        setAllRoles(roles);
    }, [roles]);

    return (
        <>
            <ul>
                {allRoles.map((role: IUserRole) => (
                    <li className="mb-2 inline-block mr-2" key={role.name}>
                        <Tag value={role.name} severity={role.colorClass}></Tag>
                    </li>
                ))}
            </ul>
        </>
    )
}