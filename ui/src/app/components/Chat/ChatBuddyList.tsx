import {Card} from "primereact/card";
import {Badge} from "primereact/badge";
import * as React from "react";
import {useState} from "react";

interface IChatBuddyListProps {
    onNewConversation: (recipient: string) => void;
}

export function ChatBuddyList({onNewConversation}: IChatBuddyListProps) {
    const [isMinimized, setIsMinimized] = useState<boolean>(false);
    const [minimizedStyles, setMinimizedStyles] = useState<string>('');
    const [buddies] = React.useState<string[]>([
        "JalapeñoLover",
        "SauceBoss",
        "BaconJamEnjoyer",
        "SweetAndSmokey",
        "TangyHot",
    ])
    const minimize = () => {
        if (isMinimized) {
            setIsMinimized(false);
            setMinimizedStyles('h-full');
        } else {
            setIsMinimized(true);
            setMinimizedStyles('h-[50px] mt-auto');
        }
    }
    const header = () => {
        return (
            <div className="flex justify-between gap-x-2 bg-stone-800 cursor-pointer"
                 onClick={() => minimize()}>
                <h2 className="p-2 m-0 text-lg font-bold  line-height-2">Buddy List</h2>
                <i className="pi pi-window-minimize hover:text-yellow-200 mr-4 mt-4"
                   title="Minimize chat window"></i>
            </div>
        )
    }
    return (
        <>
            <section
                className={`chat-buddy-list min-w-[250px] bg-black-200 border-1 border-solid border-gray-600 ml-2 ${minimizedStyles}`}>
                <Card header={header} className="h-full">
                    {!isMinimized && (
                        <ul className="list-style-none block max-h-[100vw] overflow-y-scroll">
                            {buddies.map((buddy: string, index: number) => (
                                <li onClick={() => onNewConversation(buddy)}
                                    className="block p-4 mb-2 cursor-pointer hover:bg-stone-700"
                                    key={index}>
                                    <Badge value="" severity="success" className="mr-2"></Badge> {buddy}
                                </li>
                            ))}
                        </ul>
                    )}
                </Card>
            </section>
        </>
    )
}