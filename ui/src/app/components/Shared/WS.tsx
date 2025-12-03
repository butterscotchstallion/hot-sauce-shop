import {RefObject, useEffect, useRef, useState} from "react";
import useWebSocket, {ReadyState} from "react-use-websocket";
import {Tooltip} from "primereact/tooltip";

export const WS_URL = "ws://localhost:8081/ws"

export function WS() {
    const {readyState, lastMessage} = useWebSocket(WS_URL, {
        shouldReconnect: () => true,
    });
    const connectionStatus: string = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Connected',
        [ReadyState.CLOSING]: 'Disconnecting',
        [ReadyState.CLOSED]: 'Disconnected',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];
    const populateStatusColorMap = () => {
        const statusColorMap = new Map<string, string>()
        statusColorMap.set('Connecting', "gray");
        statusColorMap.set('Connected', "yellow");
        statusColorMap.set('Disconnecting', "red");
        statusColorMap.set('Disconnected', "red");
        statusColorMap.set('Unintantiated', "gray");
        return statusColorMap;
    }
    const statusColorMap: RefObject<Map<string, string>> = useRef<Map<string, string>>(
        populateStatusColorMap()
    );
    const [statusColor, setStatusColor] = useState<string>('gray');

    useEffect(() => {
        const color = statusColorMap.current.get(connectionStatus) || "gray";
        setStatusColor(color);
    }, [connectionStatus]);

    useEffect(() => {
        if (lastMessage !== null) {
            console.log("WS message recv: ", lastMessage);
        }
    }, [lastMessage]);

    return (
        <>
            <i
                className="pi pi-bolt ws-icon"
                style={{"color": statusColor}}
                data-pr-tooltip={connectionStatus}
                data-pr-position="bottom"
            />
            <Tooltip target=".ws-icon"/>
        </>
    )
}