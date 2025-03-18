import {PrimeReactProvider} from "primereact/api";
import * as React from "react";
import {useEffect} from "react";
import NavigationMenu from "../components/Navigation/NavigationMenu.tsx";
import {useDispatch} from "react-redux";
import {getCartItems} from "../components/Cart/CartService.ts";
import {ICart} from "../components/Cart/ICart.ts";
import {setCartItems, setIdQuantityMap} from "../components/Cart/Cart.slice.ts";
import {Subscription} from "rxjs";
import {AuthContextProps, useAuth} from "react-oidc-context";

type Props = {
    children: React.ReactNode
}

export default function BaseLayout({children}: Props) {
    const dispatch = useDispatch();
    const auth: AuthContextProps = useAuth();

    useEffect(() => {
        const cartItems$: Subscription = getCartItems().subscribe({
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

    switch (auth.activeNavigator) {
        case "signinSilent":
            return <div>Signing you in...</div>;
        case "signoutRedirect":
            return <div>Signing you out...</div>;
    }

    if (auth.isLoading) {
        return <div>Loading...</div>;
    }

    if (auth.error) {
        return <div>Oops... {auth.error.message}</div>;
    }
    
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