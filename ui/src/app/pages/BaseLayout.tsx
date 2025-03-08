import {PrimeReactProvider} from "primereact/api";
import * as React from "react";
import NavigationMenu from "../components/Navigation/NavigationMenu.tsx";

type Props = {
    children: React.ReactNode
}

export default function BaseLayout({children}: Props) {
    return (
        <PrimeReactProvider>
            <section className="bg-black p-4">
                {/*
                <div className="flex justify-end pr-4">
                    <CartMenu/>
                </div>
                */}
                <NavigationMenu/>
            </section>
            <main className="container mx-auto max-w-7xl">
                {/*
                <header className="flex justify-between mt-4">
                    <nav>
                        <ul className="w-full flex justify-end gap-10 text-xl mt-10">
                            <li>
                                <NavLink to="/"
                                         className="transition-colors duration-300 hover:text-orange-500 hover:animate-pulse">
                                    Home
                                </NavLink>
                            </li>
                            <li>
                                <NavLink to="/products"
                                         className="transition-colors duration-300 hover:text-orange-500 hover:animate-pulse">
                                    Products
                                </NavLink>
                            </li>
                            <li>
                                <NavLink to="/merchandise"
                                         className="transition-colors duration-300 hover:text-orange-500 hover:animate-pulse">
                                    Merchandise
                                </NavLink>
                            </li>
                            <li>
                                <NavLink to="/contact"
                                         className="transition-colors duration-300 hover:text-orange-500 hover:animate-pulse">
                                    Contact
                                </NavLink>
                            </li>
                        </ul>
                    </nav>
                </header>
                */}
                <section className="mt-4">
                    {children}
                </section>
            </main>
        </PrimeReactProvider>
    )
}