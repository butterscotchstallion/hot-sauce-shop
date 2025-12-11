export function ChatArea({children}) {
    return (
        <>
            <section
                className="fixed bottom-0 right-0 w-1/2 min-h-[350px] grid grid-flow-col justify-items-end-safe mr-4">
                {children}
            </section>
        </>
    )
}