import {useCallback, useEffect, useState} from "react";
import useWebSocket, {ReadyState} from "react-use-websocket";
import {Button} from "primereact/button";

export function WS() {
    const [socketUrl] = useState('ws://localhost:8081/ws');
    const [messageHistory, setMessageHistory] = useState<MessageEvent<never>[]>([]);
    const {sendMessage, lastMessage, readyState} = useWebSocket(socketUrl, {
        shouldReconnect: () => true,
    });
    const handleClickSendMessage = useCallback(() => sendMessage('Hello'), [sendMessage]);
    const connectionStatus: string = {
        [ReadyState.CONNECTING]: 'Connecting',
        [ReadyState.OPEN]: 'Open',
        [ReadyState.CLOSING]: 'Closing',
        [ReadyState.CLOSED]: 'Closed',
        [ReadyState.UNINSTANTIATED]: 'Uninstantiated',
    }[readyState];

    useEffect(() => {
        if (lastMessage !== null) {
            setMessageHistory((prev: MessageEvent<never>[]) => prev.concat(lastMessage));
        }
    }, [lastMessage, setMessageHistory]);

    return (
        <>
            <div>
                <Button
                    onClick={handleClickSendMessage}
                    disabled={readyState !== ReadyState.OPEN}
                >
                    Click Me to send 'Hello'
                </Button>
                <h1 className="text-2xl font-bold ">The WebSocket is currently {connectionStatus}</h1>
                {lastMessage ? <span>Last message: {lastMessage.data}</span> : null}
                <ul>
                    {messageHistory.map((message, idx) => (
                        <li key={idx}>{message ? message.data : null}</li>
                    ))}
                </ul>
            </div>
        </>
    )
}