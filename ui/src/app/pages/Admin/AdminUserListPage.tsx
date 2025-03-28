import {ReactElement, useEffect, useState} from "react";
import {IUser} from "../../components/User/IUser.ts";
import {getUsers} from "../../components/User/UserService.ts";
import Throbber from "../../components/Shared/Throbber.tsx";
import {DataTable} from "primereact/datatable";
import {Column} from "primereact/column";
import {NavLink} from "react-router";

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
    const createdDateTemplate = (rowData: IUser) => {
        return new Date(rowData.createdAt).toLocaleDateString();
    }
    const updatedDateTemplate = (rowData: IUser) => {
        return new Date(rowData.updatedAt).toLocaleDateString();
    }
    const userTemplate = (rowData: IUser) => {
        return <NavLink to={`/admin/users/edit/${rowData.slug}`}>{rowData.username}</NavLink>
    }
    const usersTable: ReactElement = (
        <DataTable value={users} className="w-full" stripedRows>
            <Column field="username" header="Name" sortable body={userTemplate}></Column>
            <Column field="createdAt" header="Created" body={createdDateTemplate} sortable></Column>
            <Column field="updatedAt" header="Updated" body={updatedDateTemplate} sortable></Column>
        </DataTable>
    )

    return (
        <>
            {isLoading ? <Throbber/> : usersTable}
        </>
    )
}