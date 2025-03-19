import {Menu} from "primereact/menu";
import {Avatar} from "primereact/avatar";
import {MenuItem} from "primereact/menuitem";
import {RefObject, useRef, useState} from "react";
import {Button} from "primereact/button";
import {UserSignInModal} from "./UserSignInModal.tsx";

export default function UserAvatarMenu() {
    const isAuthenticated = false;
    const menu: RefObject<Menu | null> = useRef<Menu>(null);
    const [signInModalVisible, setSignInModalVisible] = useState<boolean>(false);
    const items: MenuItem[] = [
        {
            label: "User Menu",
            items: [
                {
                    label: 'Admin',
                    icon: 'pi pi-lock',
                    url: '/admin'
                },
                {
                    label: 'Account Settings',
                    icon: 'pi pi-user',
                    url: '/account'
                },
                {
                    label: 'Sign Out',
                    icon: 'pi pi-sign-out',
                    url: '/sign-out'
                }
            ]
        }
    ];

    const signInButton = () => {
        return (
            <Button
                onClick={() => setSignInModalVisible(true)}
                label="Sign In"
                icon="pi pi-lock"/>
        )
    }

    const avatarWithMenu = () => {
        return (
            <>
                <Avatar
                    onClick={(event) => menu?.current?.toggle(event)}
                    aria-controls="popup_menu_left"
                    aria-haspopup
                    className="ml-2 cursor-pointer"
                    image="/images/avatars/amyelsner.png"
                    shape="circle"/>
                <Menu model={items} popup ref={menu} popupAlignment="left"/>
            </>
        )
    }

    return (
        <>
            {isAuthenticated ? avatarWithMenu() : signInButton()}
            <UserSignInModal visible={signInModalVisible} setVisible={setSignInModalVisible}/>
        </>
    )
}