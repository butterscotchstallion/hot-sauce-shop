import {ChatBuddyList} from "./ChatBuddyList.tsx";
import {ChatWindow} from "./ChatWindow.tsx";
import {IChatMessage} from "./IChatMessage.ts";
import {useState} from "react";

export function ChatArea() {
    const [chatAreaStyles, setChatAreaStyles] = useState<string>('w-full')
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
        console.log(`Conversation with ${recipient} closed.`);
        const newConversations = conversations.filter((conversation: IChatMessage[]) => {
            return conversation[0].recipient !== recipient
        });
        setConversations(newConversations);
        console.log(conversations);

        if (newConversations.length === 0) {
            setChatAreaStyles('');
        }
    }
    // useEffect(() => {
    //     const conversations: IChatMessage[] = [];
    //     for (let j = 0; j < 3; j++) {
    //         conversations.push(makeConversation())
    //     }
    //     setConversations(conversations);
    // }, []);
    return (
        <>
            <section
                className={`fixed bottom-0 right-0 min-h-[350px] flex flex-wrap align-bottom m-2 flex gap-2 pl-4 pr-4 justify-end ${chatAreaStyles}`}>
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
            </section>
        </>
    )
}