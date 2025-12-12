import {ChatBuddyList} from "./ChatBuddyList.tsx";

export function ChatArea({children}) {
    return (
        <>
            <section
                className="fixed w-full bottom-0 right-0 min-h-[350px] flex flex-wrap align-bottom gap-2 m-2 justify-end">
                <div className="w-3/4 flex gap-2 pl-4 pr-4 justify-end">
                    {children}
                </div>
                <ChatBuddyList/>
            </section>
        </>
    )
}