import {Subject} from "rxjs";
import {CART_URL} from "../Shared/Api.ts";
import {ICart} from "./ICart.ts";
import {IAddCartItemRequest} from "./IAddCartItemRequest.ts";
import {IDeleteCartItemRequest} from "./IDeleteCartItemRequest.ts";

export function calculateCartItemsTotal(cartItems: ICart[]): number {
    return cartItems.reduce(
        (sum: number, item: ICart) => sum + item.quantity,
        0,
    );
}

export function recalculateSubtotal(cartItems: ICart[]): number {
    return cartItems.reduce(
        (acc: number, item: ICart) =>
            acc + item.price * item.quantity,
        0,
    );
}

export function getCartItems(): Subject<ICart[]> {
    const cartItems$ = new Subject<ICart[]>();
    fetch(CART_URL, {
        credentials: 'include',
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                cartItems$.next(resp?.results.cartItems || null)
            }).catch((err) => {
                cartItems$.error(err);
            })
        } else {
            cartItems$.error(res.statusText);
        }
    }).catch((err) => {
        cartItems$.error(err);
    })
    return cartItems$;
}

// Cannot serialize Map type :(
export function getIdQuantityMap(cartItems: ICart[]) {
    const idQuantityMap: any = {};
    cartItems.forEach((item: ICart) => {
        idQuantityMap[item.inventoryItemId] = item.quantity;
    })
    return idQuantityMap;
}

export function addCartItem(addCartItemRequest: IAddCartItemRequest): Subject<boolean> {
    const addCartItem$ = new Subject<boolean>();
    fetch(CART_URL, {
        method: 'POST',
        body: JSON.stringify(addCartItemRequest),
        credentials: 'include',
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

export function deleteCartItem(deleteCartItemRequest: IDeleteCartItemRequest): Subject<boolean> {
    const deleteCartItem$ = new Subject<boolean>();
    fetch(CART_URL, {
        method: 'DELETE',
        body: JSON.stringify(deleteCartItemRequest),
        credentials: 'include',
        headers: {
            'Content-Type': 'application/json'
        }
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    deleteCartItem$.next(true);
                } else {
                    deleteCartItem$.error(resp?.message || "Unknown error");
                }
            });
        }
    });
    return deleteCartItem$;
}