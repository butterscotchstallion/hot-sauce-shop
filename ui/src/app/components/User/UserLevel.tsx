import {Tag} from "primereact/tag";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {IUser} from "./IUser.ts";
import {Badge} from "primereact/badge";
import {UserLevelModal} from "./UserLevelModal.tsx";
import {useState} from "react";

export function UserLevel() {
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const level: number | null = useSelector((state: RootState) => state.user.level);
    const experience: number | null = useSelector((state: RootState) => state.user.experience);
    const [userLevelModalVisible, setUserLevelModalVisible] = useState<boolean>(false);

    return (
        <>
            {user && (
                <Tag
                    onClick={() => setUserLevelModalVisible(!userLevelModalVisible)}
                    icon="pi pi-user"
                    style={{width: "60px", height: "40px", "cursor": "pointer"}}
                    title={`Level ${level} - ${experience} experience gained`}>
                    <Badge value={level} severity="contrast"></Badge>
                </Tag>
            )}

            <UserLevelModal
                visible={userLevelModalVisible}
                setVisible={setUserLevelModalVisible}
                level={level}
                experience={experience}
            />
        </>
    )
}