import {Subject} from "rxjs";
import {IProduct} from "../Products/IProduct.ts";
import {PRODUCTS_URL} from "../Shared/Api.ts";

export function getCartItems(slug: string): Subject<IProduct> {
    const product$ = new Subject<IProduct>();
    fetch(`${PRODUCTS_URL}/${slug}`).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                product$.next(resp?.results.product || null)
            });
        } else {
            product$.error(res.statusText);
        }
    })
    return product$;
}