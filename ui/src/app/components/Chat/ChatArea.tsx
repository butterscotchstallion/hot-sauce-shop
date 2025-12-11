export function ChatArea({children}) {
    return (
        <>
            <section
                className="fixed bottom-0 right-0 flex w-1/2 min-h-[350px]">
                {children}
            </section>
        </>
    )
}