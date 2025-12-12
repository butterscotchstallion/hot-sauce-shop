import {Card} from "primereact/card";
import {Badge} from "primereact/badge";
import * as React from "react";

export function ChatBuddyList() {
    const [buddies, setBuddies] = React.useState<string[]>([
        "JalapeñoLover",
        "SauceBoss",
        "BaconJamEnjoyer",
        "SweetAndSmokey",
        "TangyHot",
    ])
    const header = () => {
        return (
            <div className="flex justify-between gap-x-2 bg-stone-800 cursor-pointer"
                 onClick={() => minimize()}>
                <h2 className="p-2 m-0 text-lg font-bold  line-height-2">Buddy List</h2>
                <i className="pi pi-window-minimize hover:text-stone-600 mr-4 mt-4"
                   title="Minimize chat window"></i>
            </div>
        )
    }
    return (
        <>
            <section
                className="chat-buddy-list w-[250px] bg-black-200 border-1 border-solid border-gray-600">
                <Card header={header} className="h-full">
                    <ul className="list-style-none">
                        {buddies.map((buddy: string, _: number) => (
                            <li className="block p-4 mb-2 cursor-pointer hover:bg-stone-700">
                                <Badge value="" severity="success" className="mr-2"></Badge> {buddy}
                            </li>
                        ))}
                    </ul>
                </Card>
            </section>
        </>
    )
}