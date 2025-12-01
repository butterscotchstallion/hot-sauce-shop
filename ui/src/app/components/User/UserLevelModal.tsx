import {Dialog} from "primereact/dialog";
import {ProgressBar} from "primereact/progressbar";

interface IUserLevelModalProps {
    visible: boolean;
    setVisible: (visible: boolean) => void;
    level: number;
    experience: number;
}

export function UserLevelModal({visible, setVisible, user}: IUserLevelModalProps) {
    return (
        <>
            <Dialog
                header="User Level"
                visible={visible}
                style={{width: '20vw'}} onHide={() => {
                if (!visible) return;
                setVisible(false);
            }}>
                <div className="flex flex-col gap-5 justify-center m-0">
                    <div className="card">
                        <ProgressBar value={50}></ProgressBar>
                    </div>
                </div>
            </Dialog>
        </>
    )
}