import {PrimeReactProvider} from "primereact/api";
import * as React from "react";
import {useEffect} from "react";
import NavigationMenu from "../components/Navigation/NavigationMenu.tsx";
import {useDispatch} from "react-redux";
import {getCartItems} from "../components/Cart/CartService.ts";
import {ICart} from "../components/Cart/ICart.ts";
import {setCartItems, setIdQuantityMap} from "../components/Cart/Cart.slice.ts";
import {Subscription} from "rxjs";
import Cookies from "js-cookie";
import {getUserBySessionId} from "../components/User/UserService.ts";
import {IUser} from "../components/User/IUser.ts";
import {setSignedIn, setUser} from "../components/User/User.slice.ts";

type Props = {
    children: React.ReactNode
}

export default function BaseLayout({children}: Props) {
    const dispatch = useDispatch();
    useEffect(() => {
        let user$: Subscription | null = null;
        const sessionId: string | undefined = Cookies.get("sessionId");
        if (sessionId) {
            user$ = getUserBySessionId().subscribe({
                next: (user: IUser) => {
                    dispatch(setUser(user));
                    dispatch(setSignedIn(true));
                },
                error: () => {
                    console.error('Error loading user');
                }
            });
        }

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
            user$?.unsubscribe();
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