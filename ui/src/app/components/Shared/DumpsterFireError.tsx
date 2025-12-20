import {Card} from "primereact/card";

export function DumpsterFireError() {
    return (
        <>
            <Card title="Something went wrong">
                <section className="min-h-[344px]">
                    <img src="/images/dumpster-fire.jpg" width="612" height="344"
                         alt="Something went wrong"/>
                </section>
            </Card>
        </>
    )
}