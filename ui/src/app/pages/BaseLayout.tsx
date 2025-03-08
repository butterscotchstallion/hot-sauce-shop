import {NavLink} from "react-router";
import {PrimeReactProvider} from "primereact/api";
import * as React from "react";

type Props = {
    children: React.ReactNode
}

export default function BaseLayout({children}: Props) {
    return (
        <PrimeReactProvider>
            <main className="container mx-auto max-w-7xl">
                <header className="flex justify-between mt-4">
                    <div
                        className="text-4xl font-bold w-[200px] all-small-caps transition-colors duration-300 hover:text-orange-500 hover:animate-pulse">
                        <NavLink to="/">Caliente Corner</NavLink>
                    </div>

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

                <section className="mt-4">
                    {children}
                </section>
            </main>
        </PrimeReactProvider>
    )
}