import {PRODUCTS_URL} from "../Shared/Api.ts";
import {Subject} from "rxjs";
import {IProductsResults} from "./IProductsResults.ts";

export function getProducts(offset: number = 0, perPage: number = 10): Subject<IProductsResults> {
    const products$ = new Subject<IProductsResults>();
    fetch(`${PRODUCTS_URL}?offset=${offset}&perPage=${perPage}`).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                products$.next(resp?.results || [])
            });
        } else {
            products$.error(res.statusText);
        }
    }).catch((err) => {
        products$.error(err);
    });
    return products$;
}