import {Card} from "primereact/card";
import * as React from "react";
import {useState} from "react";
import {InputTextarea} from "primereact/inputtextarea";
import {Button} from "primereact/button";
import "./ChatWindow.scss";

interface IChatMessage {
    username: string;
    message: string;
}

export function ChatWindow() {
    const [closed, setClosed] = useState<boolean>(false);
    const [isMinimized, setIsMinimized] = useState<boolean>(false);
    const [minimizedStyles, setMinimizedStyles] = useState<string>('');
    const [recipient] = useState<string>("JalapeñoLover");
    const [outgoingMessage, setOutgoingMessage] = useState<string>("");
    const [messages] = React.useState<IChatMessage[]>([
        {
            username: "SauceBoss",
            message: "How are you enjoying the sauce?"
        },
        {
            username: recipient,
            message: "It's delicious! I love the balance between tangy and spicy!"
        },
        {
            username: "SauceBoss",
            message: "That's great to hear!"
        },
        {
            username: recipient,
            message: "Do you have any sauces with a bit more heat?"
        }
    ]);
    const header = (
        <div className="flex justify-between gap-x-2 cursor-pointer">
            <h2 className="p-4 m-0 text-lg font-bold">
                {recipient}
            </h2>
            <div className="mt-4 mr-4 max-w-[50px] w-full cursor-pointer flex justify-between">
                <i
                    onClick={() => minimize()}
                    className="pi pi-window-minimize hover:text-stone-600"
                    title="Minimize chat window"></i>
                <i
                    onClick={() => setClosed(true)}
                    className="pi pi-times hover:text-stone-600"
                    title="Close chat window"></i>
            </div>
        </div>
    );
    const minimize = () => {
        if (isMinimized) {
            setIsMinimized(false);
            setMinimizedStyles('');
        } else {
            setIsMinimized(true);
            setMinimizedStyles('h-[50px]');
        }
    }
    const footer = () => {
        return (
            <section className="w-full">
                <div className="w-full flex gap-x-2">
                    <InputTextarea
                        className="w-full"
                        value={outgoingMessage}
                        onChange={(e: React.ChangeEvent<HTMLTextAreaElement>) => setOutgoingMessage(e.target.value)}
                        rows={2}
                        cols={30}/>
                    <Button icon="pi pi-send">Send</Button>
                </div>
            </section>
        )
    }
    return (
        <>
            {!closed && (
                <div
                    className={`chat-window max-w-[300px] h-full min-w-[250px] mb-2 border-1 border-solid border-gray-600 ${minimizedStyles}`}>
                    <Card
                        header={header}
                        footer={footer}
                        className="h-full">
                        {!isMinimized && (
                            <>
                                <section
                                    className="w-full h-[225px] overflow-y-scroll bg-stone-800 p-2 mb-2 mb-0 pb-0">
                                    <ul className="list-style-none">
                                        {messages.map((item: IChatMessage, index: number) => (
                                            <li key={`chat-message-${index}`} className="mb-4 text-sm">
                                                <div className="flex gap-x-2">
                                                    <aside className="w-[50px] h-[50px]">
                                                        <img src="/images/hot-pepper-thumbnail.png"
                                                             alt="hot pepper"
                                                             width="50"
                                                             height="50"
                                                        />
                                                    </aside>
                                                    <div>
                                                        <strong className="block mb-1">{item.username}</strong>
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