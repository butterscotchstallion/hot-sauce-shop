import {PrimeReactProvider} from "primereact/api";
import * as React from "react";
import NavigationMenu from "../components/Navigation/NavigationMenu.tsx";

type Props = {
    children: React.ReactNode
}

export default function BaseLayout({children}: Props) {
    return (
        <PrimeReactProvider>
            <NavigationMenu/>
            <main className="container mx-auto max-w-7xl mb-10">
                <section className="mt-4">
                    {children}
                </section>
            </main>
        </PrimeReactProvider>
    )
}