import {Subject} from "rxjs";
import {CART_URL} from "../Shared/Api.ts";
import {ICart} from "./ICart.ts";

export function getCartItems(): Subject<ICart[]> {
    const cartItems$ = new Subject<ICart[]>();
    fetch(CART_URL).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                cartItems$.next(resp?.results.cartItems || null)
            });
        } else {
            cartItems$.error(res.statusText);
        }
    })
    return cartItems$;
}