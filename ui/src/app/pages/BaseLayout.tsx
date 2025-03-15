import {PrimeReactProvider} from "primereact/api";
import * as React from "react";
import {useEffect} from "react";
import NavigationMenu from "../components/Navigation/NavigationMenu.tsx";
import {useDispatch} from "react-redux";
import {getCartItems} from "../components/Cart/CartService.ts";
import {ICart} from "../components/Cart/ICart.ts";
import {setCartItems, setIdQuantityMap} from "../components/Cart/Cart.slice.ts";

type Props = {
    children: React.ReactNode
}

export default function BaseLayout({children}: Props) {
    const dispatch = useDispatch();
    useEffect(() => {
        const cartItems$ = getCartItems().subscribe({
            next: (cartItems: ICart[]) => {
                dispatch(setCartItems(cartItems));
                dispatch(setIdQuantityMap(cartItems));
            },
            error: () => {
                console.error('Error loading cart items');
            }
        });
        return () => {
            cartItems$.unsubscribe();
        }
    }, [dispatch]);

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