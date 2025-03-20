import {Dialog} from "primereact/dialog";
import {InputText} from "primereact/inputtext";
import * as React from "react";
import {ChangeEvent, RefObject, useRef, useState} from "react";
import {Button} from "primereact/button";
import {ISignInResponse, ValidateUsernameAndPassword} from "./UserService.ts";
import {useDispatch} from "react-redux";
import {setSignedIn, setUser} from "./User.slice.ts";
import {Messages} from "primereact/messages";

interface IUserSignInModalProps {
    visible: boolean;
    setVisible: (visible: boolean) => void;
}

export function UserSignInModal(props: IUserSignInModalProps) {
    const dispatch = useDispatch();
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");
    const [isSigningIn, setIsSigningIn] = useState<boolean>(false);
    const msgs: RefObject<Messages | null> = useRef<Messages>(null);
    const [signInButtonLabel, setSignInButtonLabel] = useState<string>("Sign In");
    const signIn = () => {
        setIsSigningIn(true);
        setSignInButtonLabel("Signing In...");
        ValidateUsernameAndPassword(username, password).subscribe({
            next: (results: ISignInResponse) => {
                setSignInButtonLabel("Sign In");
                props.setVisible(false);
                dispatch(setUser(results.user));
                dispatch(setSignedIn(true));
                console.log("Signed in as " + results.user.username);
            },
            error: (err) => {
                showSignInMessage(err);
                console.error(err);
                setIsSigningIn(false);
                setSignInButtonLabel("Sign In");
            }
        });
    }
    const isFormValid = () => {
        return username.length > 0 && password.length > 0;
    }
    const handleKeyPress = (e: React.KeyboardEvent<HTMLInputElement>) => {
        if (e.key === 'Enter' && isFormValid()) {
            signIn();
        }
    }
    const showSignInMessage = (message: string) => {
        msgs.current?.clear();
        msgs.current?.show({
            id: '1',
            sticky: true,
            severity: 'error',
            summary: 'Error',
            detail: message,
            closable: false,
            content: (
                <React.Fragment>
                    <div className="ml-2"><i className="pi pi-info-circle"></i> {message}</div>
                </React.Fragment>
            )
        });
    }

    return (
        <>
            <Dialog
                header="Sign In"
                visible={props.visible}
                style={{width: '20vw'}} onHide={() => {
                if (!props.visible) return;
                props.setVisible(false);
            }}>
                <div className="flex flex-col gap-5 justify-center m-0">
                    <section>
                        <Messages ref={msgs}/>
                    </section>
                    <div>
                        <label className="mb-2 block" htmlFor="username">Username</label>
                        <InputText
                            onKeyDown={handleKeyPress}
                            size={15}
                            value={username}
                            id="username"
                            className="w-full"
                            onChange={(e: ChangeEvent<HTMLInputElement>) => {
                                setUsername(e.target.value)
                            }}/>
                    </div>
                    <div>
                        <label className="mb-2 block" htmlFor="password">Password</label>
                        <InputText
                            type="password"
                            onKeyDown={handleKeyPress}
                            size={15}
                            value={password}
                            id="password"
                            className="w-full"
                            onChange={(e: ChangeEvent<HTMLInputElement>) => {
                                setPassword(e.target.value)
                            }}/>
                    </div>

                    <div className="">
                        <Button
                            onClick={() => {
                                signIn();
                            }}
                            disabled={!isFormValid() || isSigningIn}
                            label={signInButtonLabel}
                            icon="pi pi-lock"
                            className="w-full"/>
                    </div>
                </div>
            </Dialog>
        </>
    )
}