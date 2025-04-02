import {Menu} from "primereact/menu";
import {Avatar} from "primereact/avatar";
import {MenuItem} from "primereact/menuitem";
import {RefObject, useRef, useState} from "react";
import {Button} from "primereact/button";
import {UserSignInModal} from "./UserSignInModal.tsx";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {IUser} from "./IUser.ts";
import {confirmDialog} from "primereact/confirmdialog";
import {removeSessionCookie} from "./UserService.ts";
import {setSignedOut} from "./User.slice.ts";
import {Badge} from "primereact/badge";
import {NavLink} from "react-router";

interface IUserAvatarMenuItem extends MenuItem {
    url: string;
    badge: string;
    shortcut: string;
    label: string;
}

export default function UserAvatarMenu() {
    const dispatch = useDispatch();
    const isSignedIn = useSelector((state: RootState) => state.user.isSignedIn);
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const menu: RefObject<Menu | null> = useRef<Menu>(null);
    const [signInModalVisible, setSignInModalVisible] = useState<boolean>(false);
    const itemRenderer = (item: IUserAvatarMenuItem) => (
        <div className='p-menuitem-content'>
            <NavLink
                to={item.url}
                className="flex align-items-center p-menuitem-link"
                onClick={(event) => menu?.current?.toggle(event)}
            >
                <span className={item.icon}/>
                <span className="mx-2">{item.label}</span>
                {item.badge && <Badge className="ml-auto" value={item.badge}/>}
                {item.shortcut && <span
                    className="ml-auto border-1 surface-border border-round surface-100 text-xs p-1">{item.shortcut}</span>}
            </NavLink>
        </div>
    );
    const items: MenuItem[] = [
        {
            label: user ? user.username : "User menu",
            items: [
                {
                    label: 'Admin',
                    icon: 'pi pi-lock',
                    url: '/admin',
                    template: itemRenderer
                },
                {
                    label: 'Account Settings',
                    icon: 'pi pi-user',
                    url: '/account',
                    template: itemRenderer
                },
                {
                    label: 'Sign Out',
                    icon: 'pi pi-sign-out',
                    command: () => {
                        confirmDialog({
                            header: "Sign Out",
                            message: 'Are you sure you want to sign out?',
                            icon: 'pi pi-exclamation-triangle',
                            defaultFocus: 'accept',
                            accept: () => {
                                removeSessionCookie();
                                dispatch(setSignedOut(null));
                            },
                            reject: () => {
                            }
                        });
                    }
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
                    image={`/images/avatars/${user?.avatarFilename}`}
                    shape="circle"/>
                <Menu model={items} popup ref={menu} popupAlignment="left"/>
            </>
        )
    }

    return (
        <>
            {isSignedIn ? avatarWithMenu() : signInButton()}
            <UserSignInModal visible={signInModalVisible} setVisible={setSignInModalVisible}/>
        </>
    )
}