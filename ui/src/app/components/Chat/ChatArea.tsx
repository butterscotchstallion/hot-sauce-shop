import {ChatBuddyList} from "./ChatBuddyList.tsx";

export function ChatArea({children}) {
    return (
        <>
            <section
                className="fixed w-full bottom-0 right-0 min-h-[350px] flex flex-wrap justify-between gap-2 m-2">
                <div className="w-3/4 flex gap-2 pl-4 pr-4">
                    {children}
                </div>
                <ChatBuddyList/>
            </section>
        </>
    )
}