import {PRODUCT_AUTOCOMPLETE_URL, PRODUCTS_URL} from "../Shared/Api.ts";
import {Subject} from "rxjs";
import {IProductsResults} from "./IProductsResults.ts";
import {IProduct} from "./IProduct.ts";
import {IAutocompleteSuggestion} from "./IAutocompleteSuggestion.ts";
import {IDisplayTag} from "../../pages/ProductListPage.tsx";

function getFilterTagsURLParameter(filters: IDisplayTag[]): string {
    if (filters.length > 0) {
        const tagIds: string = filters.map((filter: IDisplayTag) => filter.id).join(",");
        return "&tags=" + tagIds
    } else {
        return "";
    }
}

export function getProducts(offset: number = 0, perPage: number = 10, sort: string = "name", filters: IDisplayTag[]): Subject<IProductsResults> {
    const products$ = new Subject<IProductsResults>();
    let productsUrl: string = `${PRODUCTS_URL}?offset=${offset}&perPage=${perPage}&sort=${sort}`;
    productsUrl += getFilterTagsURLParameter(filters);
    fetch(productsUrl).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                products$.next(resp?.results || [])
            });
        } else {
            products$.error(res.json().then(resp => resp?.message || "Unknown error"));
        }
    }).catch((err: string) => {
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

export function getProductAutocompleteSuggestions(query: string): Subject<IAutocompleteSuggestion[]> {
    const autocomplete$ = new Subject<IAutocompleteSuggestion[]>();
    fetch(`${PRODUCT_AUTOCOMPLETE_URL}?q=${encodeURI(query)}`).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                autocomplete$.next(resp?.results?.suggestions || [])
            });
        } else {
            autocomplete$.error(res.statusText);
        }
    });
    return autocomplete$;
}

export function addOrUpdateItem(product: IProduct, isNewProduct: boolean): Subject<boolean> {
    const updateItem$ = new Subject<boolean>();
    fetch(`${PRODUCTS_URL}/${product.slug}`, {
        method: isNewProduct ? 'POST' : 'PUT',
        body: JSON.stringify(product),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    updateItem$.next(true);
                } else {
                    updateItem$.error(resp?.message || "Unknown error");
                }
            });
        } else {
            updateItem$.error(res.statusText);
        }
    }).catch((err) => {
        updateItem$.error(err);
    })
    return updateItem$;
}

export function addItem(product: IProduct): Subject<boolean> {
    const addItem$ = new Subject<boolean>();
    fetch(PRODUCTS_URL, {
        method: 'POST',
        body: JSON.stringify(product),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    addItem$.next(true);
                } else {
                    addItem$.error(resp?.message || "Unknown error");
                }
            });
        }
    });
    return addItem$;
}