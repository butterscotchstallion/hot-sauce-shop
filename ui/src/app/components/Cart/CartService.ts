import {Subject} from "rxjs";
import {CART_URL} from "../Shared/Api.ts";
import {ICart} from "./ICart.ts";
import {IAddCartItemRequest} from "./IAddCartItemRequest.ts";

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

export function addCartItem(addCartItemRequest: IAddCartItemRequest): Subject<boolean> {
    const addCartItem$ = new Subject<boolean>();
    fetch(CART_URL, {
        method: 'POST',
        body: JSON.stringify(addCartItemRequest),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    addCartItem$.next(resp.status === "OK");
                } else {
                    addCartItem$.error(resp?.message || "Unknown error");
                }
            });
        } else {
            addCartItem$.error(res.statusText);
        }
    }).catch((err) => {
        addCartItem$.error(err);
    })
    return addCartItem$;
}