import {Card} from "primereact/card";
import * as React from "react";
import {InputTextarea} from "primereact/inputtextarea";
import {Button} from "primereact/button";

interface IChatMessage {
    username: string;
    message: string;
}

export function ChatWindow() {
    const [outgoingMessage, setOutgoingMessage] = React.useState<string>("");
    const [messages, setMessages] = React.useState<IChatMessage[]>([
        {
            username: "SauceBoss",
            message: "How are you enjoying the sauce?"
        },
        {
            username: "JalapeñoLover",
            message: "It's delicious! I love the balance between tangy and spicy!"
        },
        {
            username: "SauceBoss",
            message: "That's great to hear!"
        }
    ]);

    return (
        <>
            <Card title="JalapeñoLover" className="max-w-2/3 min-w-[300px] mb-2 border-1 border-solid border-gray-500">
                <section className="w-full h-[200px] overflow-y-scroll bg-stone-800 p-2 mb-2">
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
                                        {item.message}
                                    </div>
                                </div>
                            </li>
                        ))}
                    </ul>
                </section>
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
            </Card>
        </>
    )
}