import {PRODUCTS_URL} from "../Shared/Api.ts";
import {Subject} from "rxjs";
import {IProductsResults} from "./IProductsResults.ts";
import {IProduct} from "./IProduct.ts";

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

export function getProductDetail(slug: string): Subject<IProduct> {
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