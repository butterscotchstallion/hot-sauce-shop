import {Card} from "primereact/card";
import * as React from "react";
import {useEffect, useState} from "react";
import {InputTextarea} from "primereact/inputtextarea";
import {Button} from "primereact/button";
import "./ChatWindow.scss";
import {IChatMessage} from "./IChatMessage.ts";

interface IChatWindowProps {
    conversation: IChatMessage;
}

export function ChatWindow({conversation}: IChatWindowProps) {
    const [messagesArea, setMessagesArea] = useState<HTMLElement>();
    const [closed, setClosed] = useState<boolean>(false);
    const [isMinimized, setIsMinimized] = useState<boolean>(false);
    const [minimizedStyles, setMinimizedStyles] = useState<string>('');
    const [recipient, setRecipient] = useState<string>("JalapeñoLover");
    const [outgoingMessage, setOutgoingMessage] = useState<string>("");
    const scrollToBottom = () => {
        if (messagesArea) {
            messagesArea.scrollIntoView({behavior: "smooth"});
        }
    }
    const [messages] = useState<IChatMessage[]>(conversation);
    const header = (
        <div className="flex justify-between gap-x-2 cursor-pointer">
            <h2 className="p-4 pt-2 m-0 text-lg font-bold">
                {recipient}
            </h2>
            <div className="mt-4 mr-4 max-w-[50px] w-full cursor-pointer flex justify-between">
                <i
                    onClick={() => minimize()}
                    className="pi pi-window-minimize hover:text-yellow-200"
                    title="Minimize chat window"></i>
                <i
                    onClick={() => setClosed(true)}
                    className="pi pi-times hover:text-yellow-200"
                    title="Close chat window"></i>
            </div>
        </div>
    );
    const footer = () => {
        return !isMinimized && (
            <section className="w-full">
                <div className="w-full flex gap-x-2">
                    <InputTextarea
                        className="w-full"
                        value={outgoingMessage}
                        onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setOutgoingMessage(e.target.value)}
                        rows={2}
                        cols={30}/>
                    <Button icon="pi pi-send" className="min-w-[50px]"/>
                </div>
            </section>
        )
    }
    const minimize = () => {
        if (isMinimized) {
            setIsMinimized(false);
            setMinimizedStyles('h-full');
        } else {
            setIsMinimized(true);
            setMinimizedStyles('h-[50px] mt-auto');
        }
    }
    useEffect(() => {
        scrollToBottom();
        setRecipient(messages[0].recipient);
    }, [messages]);
    return (
        <>
            {!closed && (
                <div
                    className={`chat-window max-w-[350px] min-w-[250px] border-1 border-solid border-gray-600 ${minimizedStyles}`}>
                    <Card
                        header={header}
                        footer={footer}
                        className="h-full">
                        {!isMinimized && (
                            <>
                                {/* Chat messages */}
                                <section
                                    className="w-full h-[225px] overflow-y-scroll bg-stone-800 p-2 mb-2 mb-0 pb-0"
                                    ref={(el) => {
                                        if (el) {
                                            setMessagesArea(el)
                                        }
                                    }}>
                                    <ul className="list-style-none">
                                        {messages.map((item: IChatMessage, index: number) => (
                                            <li key={`chat-message-${index}`} className="mb-4 text-sm">
                                                <div className="flex gap-x-2">
                                                    <aside
                                                        className="w-[25px] h-[25px]">
                                                        <img src="/images/hot-pepper-buddy-icon.png"
                                                             alt="hot pepper"
                                                             width="25"
                                                             height="25"
                                                             className="rounded-2xl"
                                                        />
                                                    </aside>
                                                    <div className="w-3/4">
                                                        <strong className="block mb-1">{item.recipient}</strong>
                                                        <div className="pr-2">{item.message}</div>
                                                    </div>
                                                </div>
                                            </li>
                                        ))}
                                    </ul>
                                </section>
                            </>
                        )}
                    </Card>
                </div>
            )}
        </>
    )
}