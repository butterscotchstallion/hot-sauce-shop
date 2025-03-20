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

export default function UserAvatarMenu() {
    const dispatch = useDispatch();
    const isSignedIn = useSelector((state: RootState) => state.user.isSignedIn);
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const menu: RefObject<Menu | null> = useRef<Menu>(null);
    const [signInModalVisible, setSignInModalVisible] = useState<boolean>(false);
    const items: MenuItem[] = [
        {
            label: user ? user.username : "User menu",
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
                    image="/images/avatars/amyelsner.png"
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