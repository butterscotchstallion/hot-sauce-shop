import {PRODUCTS_URL} from "../Shared/Api.ts";
import {Subject} from "rxjs";
import {IProduct} from "./IProduct.ts";

export function getProducts() {
    const products$ = new Subject<IProduct[]>();
    fetch(PRODUCTS_URL).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => products$.next(resp.results.inventory));
        } else {
            products$.error(res.statusText);
        }
    }).catch((err) => {
        products$.error(err);
    });
    return products$;
}