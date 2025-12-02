import {Dialog} from "primereact/dialog";
import {ProgressBar} from "primereact/progressbar";
import {Column} from "primereact/column";
import {DataTable} from "primereact/datatable";

interface IUserLevelModalProps {
    visible: boolean;
    setVisible: (visible: boolean) => void;
    level: number;
    experience: number;
    percentageOfLevelComplete: number;
}

export function UserLevelModal({
                                   visible,
                                   setVisible,
                                   level,
                                   experience,
                                   percentageOfLevelComplete
                               }: IUserLevelModalProps) {
    const userLevelDetails = [{level, experience: experience.toLocaleString()}];
    return (
        <>
            <Dialog
                header="User Level Details"
                visible={visible}
                style={{width: '20vw'}} onHide={() => {
                if (!visible) return;
                setVisible(false);
            }}>
                <div className="flex flex-col gap-5 justify-center m-0">
                    <div className="card">
                        <section className="flex flex-col gap-y-1">
                            <DataTable value={userLevelDetails}>
                                <Column field="level" header="Level"></Column>
                                <Column field="experience" header="Experience"></Column>
                            </DataTable>
                            {experience > 0 && <ProgressBar value={percentageOfLevelComplete || 0}/>}
                        </section>
                    </div>
                </div>
            </Dialog>
        </>
    )
}