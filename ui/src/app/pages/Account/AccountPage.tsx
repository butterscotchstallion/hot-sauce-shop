import {useEffect, useState} from "react";
import {IUser} from "../../components/User/IUser.ts";
import {getUserBySlug} from "../../components/User/AdminService.ts";

export function AccountPage() {
    const [account] = useState<IUser | null>(null);

    useEffect(() => {
        const account$: Subscription = getUserBySlug()
    }, []);

    return (
        <>
            <h1 className="text-3xl font-bold mb-4">Account Settings</h1>

            <section className="mt-4">
                <table className="w-1/3">
                    <tbody>
                    <tr>
                        <td className="w-1/2 mb-4">
                            <strong>Username</strong>
                        </td>
                        <td>
                            SauceBoss
                        </td>
                    </tr>
                    <tr>
                        <td className="w-1/2  mb-4">
                            <strong>Created</strong>
                        </td>
                        <td>
                            May 1, 2021
                        </td>
                    </tr>
                    <tr>
                        <td className="w-1/2  mb-4">
                            <strong>Roles</strong>
                        </td>
                        <td>
                            Role List Here
                        </td>
                    </tr>
                    </tbody>
                </table>
            </section>
        </>
    )
}