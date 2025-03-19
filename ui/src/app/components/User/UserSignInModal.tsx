import {Dialog} from "primereact/dialog";
import {InputText} from "primereact/inputtext";
import {useState} from "react";
import {Button} from "primereact/button";

interface IUserSignInModalProps {
    visible: boolean;
    setVisible: (visible: boolean) => void;
}

export function UserSignInModal(props: IUserSignInModalProps) {
    const [username, setUsername] = useState<string>("");
    const [password, setPassword] = useState<string>("");

    return (
        <>
            <Dialog
                header="Sign In"
                visible={props.visible}
                style={{width: '20vw'}} onHide={() => {
                if (!props.visible) return;
                props.setVisible(false);
            }}>
                <div className="flex flex-col gap-10 justify-center m-0">
                    <div>
                        <label className="mb-2 block" htmlFor="username">Username</label>
                        <InputText size={15}
                                   className="w-full"
                                   onChange={(e) => {
                                       setUsername(e.target.value)
                                   }}/>
                    </div>
                    <div>
                        <label className="mb-2 block" htmlFor="password">Password</label>
                        <InputText size={15}
                                   className="w-full"
                                   onChange={(e) => {
                                       setPassword(e.target.value)
                                   }}/>
                    </div>

                    <div className="">
                        <Button
                            disabled={username.length === 0 || password.length === 0}
                            label="Sign In"
                            icon="pi pi-lock"
                            className="w-full"/>
                    </div>
                </div>
            </Dialog>
        </>
    )
}