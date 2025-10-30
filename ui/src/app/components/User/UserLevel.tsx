import {Tag} from "primereact/tag";
import {useSelector} from "react-redux";
import {RootState} from "../../store.ts";
import {IUser} from "./IUser.ts";
import {Badge} from "primereact/badge";

export function UserLevel() {
    const user: IUser | null = useSelector((state: RootState) => state.user.user);
    const level: number | null = useSelector((state: RootState) => state.user.level);
    const experience: number | null = useSelector((state: RootState) => state.user.experience);
    return (
        <>
            {user && (
                <Tag
                    severity="info"
                    icon="pi pi-user"
                    style={{width: "80px", height: "40px"}}
                    title={`Level ${level}`}>
                    <Badge value={level} severity="danger"></Badge>
                </Tag>
            )}
        </>
    )
}