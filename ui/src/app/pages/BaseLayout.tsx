import {PrimeReactProvider} from "primereact/api";
import * as React from "react";
import {useEffect} from "react";
import NavigationMenu from "../components/Navigation/NavigationMenu.tsx";
import {useDispatch} from "react-redux";
import {getCartItems, recalculateSubtotal} from "../components/Cart/CartService.ts";
import {ICart} from "../components/Cart/ICart.ts";
import {setCartItems, setCartSubtotal, setIdQuantityMap} from "../components/Cart/Cart.slice.ts";
import {Subscription} from "rxjs";
import Cookies from "js-cookie";
import {getUserDetailsBySessionId} from "../components/User/UserService.ts";
import {setSignedIn, setUser, setUserExperience, setUserLevel, setUserRoles} from "../components/User/User.slice.ts";
import {IUserDetails} from "../components/User/IUserDetails.ts";
import {ErrorBoundary} from "react-error-boundary";
import {Card} from "primereact/card";
import {ChatArea} from "../components/Chat/ChatArea.tsx";
import {ChatWindow} from "../components/Chat/ChatWindow.tsx";

type Props = {
    children: React.ReactNode
}

export default function BaseLayout({children}: Props) {
    const dispatch = useDispatch();
    useEffect(() => {
        let user$: Subscription | null = null;
        let cartItems$: Subscription | null = null;
        const sessionId: string | undefined = Cookies.get("sessionId");
        if (sessionId) {
            user$ = getUserDetailsBySessionId().subscribe({
                next: (userDetails: IUserDetails) => {
                    dispatch(setUser(userDetails.user));
                    dispatch(setUserRoles(userDetails.roles));
                    dispatch(setSignedIn(true));
                    dispatch(setUserLevel(userDetails.userLevelInfo.level));
                    dispatch(setUserExperience(userDetails.userLevelInfo.experience));
                },
                error: () => {
                    console.error('Error loading user');
                }
            });
            cartItems$ = getCartItems().subscribe({
                next: (cartItems: ICart[]) => {
                    dispatch(setCartItems(cartItems));
                    dispatch(setIdQuantityMap(cartItems));
                    dispatch(setCartSubtotal(recalculateSubtotal(cartItems)));
                },
                error: () => {
                    console.error('Error loading cart items');
                }
            });
        }
        return () => {
            cartItems$?.unsubscribe();
            user$?.unsubscribe();
        }
    }, [dispatch]);

    return (
        <PrimeReactProvider>
            <NavigationMenu/>
            <main className="container mx-auto max-w-7xl mb-10">
                <section className="mt-4">
                    <ErrorBoundary fallback={(
                        <>
                            <Card title="Something went wrong">
                                <section className="min-h-[344px]">
                                    <img src="/images/dumpster-fire.jpg" width="612" height="344"
                                         alt="Something went wrong"/>
                                </section>
                            </Card>
                        </>
                    )}>{children}</ErrorBoundary>
                </section>
            </main>
            <ChatArea>
                <ChatWindow/>
            </ChatArea>
        </PrimeReactProvider>
    )
}