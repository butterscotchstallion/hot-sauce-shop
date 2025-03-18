import {Menu} from "primereact/menu";
import {Avatar} from "primereact/avatar";
import {MenuItem} from "primereact/menuitem";
import {RefObject, useRef} from "react";

export default function UserAvatarMenu() {
    const menu: RefObject<Menu | null> = useRef<Menu>(null);
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