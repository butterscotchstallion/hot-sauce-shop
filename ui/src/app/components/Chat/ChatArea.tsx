import {ChatBuddyList} from "./ChatBuddyList.tsx";
import {ChatWindow} from "./ChatWindow.tsx";
import {IChatMessage} from "./IChatMessage.ts";
import {useEffect, useState} from "react";

export function ChatArea() {
    const [conversations, setConversations] = useState<IChatMessage[]>([])
    const getMessage = () => {
        const messages: string[] = [
            "How are you enjoying the sauce?",
            "It's delicious! I love the balance between tangy and spicy!",
            "That's great to hear!",
            "Do you have any sauces with a bit more heat?"
        ];
        return messages[Math.floor(Math.random() * messages.length)];
    }
    const makeConversation = (): IChatMessage[] => {
        const recipients: string[] = ["JalapeñoLover", "SeranoGal", "SauceBoss", "TangyAndSweet", "SpiceDemon"];
        const messages: IChatMessage[] = [];
        for (let j = 0; j < Math.floor(Math.random() * 3) + 1; j++) {
            messages.push({
                recipient: recipients[Math.floor(Math.random() * recipients.length)],
                message: getMessage()
            });
        }
        return messages;
    }
    const onNewConversation = (recipient: string) => {
        const newConversations: IChatMessage[] = conversations;
        newConversations.push({
            recipient: recipient,
            message: ''
        });
        setConversations(newConversations);
        console.log(conversations);
    }
    const onConversationClosed = (recipient: string) => {
        const newConversations: IChatMessage[] = conversations.filter((conversation: IChatMessage) => conversation.recipient !== recipient);
        setConversations(newConversations);
    }
    useEffect(() => {
        const conversations: IChatMessage[] = [];
        for (let j = 0; j < 3; j++) {
            conversations.push(makeConversation())
        }
        setConversations(conversations);
    }, []);
    return (
        <>
            {conversations.length > 0 && (
                <section
                    className="fixed w-full bottom-0 right-0 min-h-[350px] flex flex-wrap align-bottom gap-2 m-2 justify-end">
                    <div className="w-3/4 flex gap-2 pl-4 pr-4 justify-end">
                        {conversations.map((conversation: IChatMessage, index: number) => (
                            <ChatWindow
                                key={index}
                                conversation={conversation}
                                onConversationClosed={onConversationClosed}
                            />
                        ))}
                        <ChatBuddyList
                            key="buddy-list"
                            onNewConversation={onNewConversation}/>
                    </div>
                </section>
            )}
        </>
    )
}