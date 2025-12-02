import {RefObject, useEffect, useRef, useState} from "react";
import useWebSocket, {ReadyState} from "react-use-websocket";

export function WS() {
    const [socketUrl] = useState('ws://localhost:8081/ws');
    const {readyState} = useWebSocket(socketUrl, {
        shouldReconnect: () => true,
    });
    const connectionStatus: string = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];
    const populateStatusColorMap = () => {
        const statusColorMap = new Map<string, string>()
        statusColorMap.set('Connecting', "gray");
        statusColorMap.set('Open', "yellow");
        statusColorMap.set('Closing', "red");
        statusColorMap.set('Closed', "red");
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

    return (
        <>
            <i className="pi pi-bolt" style={{"color": statusColor}} title={connectionStatus}/>
        </>
    )
}