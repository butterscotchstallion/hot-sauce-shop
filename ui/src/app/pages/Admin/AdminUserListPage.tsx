import {ReactElement, useEffect, useState} from "react";
import {IUser} from "../../components/User/IUser.ts";
import {getUsers} from "../../components/User/UserService.ts";
import Throbber from "../../components/Shared/Throbber.tsx";
import {DataTable} from "primereact/datatable";
import {Column} from "primereact/column";

export function AdminUserListPage() {
    const [users, setUsers] = useState<IUser[]>([]);
    const [isLoading, setIsLoading] = useState<boolean>(true);

    useEffect(() => {
        setIsLoading(true);
        getUsers().subscribe({
            next: (users: IUser[]) => {
                setUsers(users);
                setIsLoading(false);
            },
            error: (err) => {
                console.error(err);
                setIsLoading(false);
            }
        });
    }, []);
    const createdDateTemplate = (rowData: IUser, _) => {
        return new Date(rowData.createdAt).toLocaleDateString();
    }
    const updatedDateTemplate = (rowData: IUser, _) => {
        return new Date(rowData.updatedAt).toLocaleDateString();
    }

    const usersTable: ReactElement = (
        <DataTable value={users} className="w-full">
            <Column field="username" header="Name"></Column>
            <Column field="createdAt" header="Created" body={createdDateTemplate}></Column>
            <Column field="updatedAt" header="Updated" body={updatedDateTemplate}></Column>
        </DataTable>
    )

    return (
        <>
            {isLoading ? <Throbber/> : usersTable}
        </>
    )
}