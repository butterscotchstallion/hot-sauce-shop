import {Tag} from "primereact/tag";
import {useDispatch, useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {IUser} from "./IUser.ts";
import {Badge} from "primereact/badge";
import {UserLevelModal} from "./UserLevelModal.tsx";
import {RefObject, useEffect, useRef, useState} from "react";
import useWebSocket from "react-use-websocket";
import {WS_URL} from "../Shared/WS.tsx";
import {parseNotification, WebsocketMessageType} from "../Shared/Notification.ts";
import {setUserExperience, setUserLevel} from "./User.slice.ts";
import {Toast} from "primereact/toast";

export function UserLevel() {
    const toast: RefObject<Toast | null> = useRef(null);
    const dispatch = useDispatch();
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const level: number | 1 = useSelector((state: RootState) => state.user.level);
    const experience: number | 0 = useSelector((state: RootState) => state.user.experience);
    const [userLevelModalVisible, setUserLevelModalVisible] = useState<boolean>(false);
    const {lastMessage} = useWebSocket(WS_URL, {
        shouldReconnect: () => true,
    });
    const [percentageOfLevelComplete, setPercentageOfLevelComplete] = useState<number>(0);

    useEffect(() => {
        if (lastMessage) {
            const data = parseNotification(lastMessage);
            if (data.messageType === WebsocketMessageType.USER_LEVEL_UPDATE) {
                dispatch(setUserLevel(data.updatedUserLevel));
                dispatch(setUserExperience(data.updatedUserExperience));
                setPercentageOfLevelComplete(data.percentageOfLevelComplete);
                if (toast.current) {
                    toast.current.show({
                        severity: 'success',
                        summary: 'Success',
                        detail: 'You gained experience!',
                        life: 3000,
                    })
                }
            }
        }
    }, [lastMessage]);

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
                percentageOfLevelComplete={percentageOfLevelComplete}
            />
            <Toast ref={toast}/>
        </>
    )
}