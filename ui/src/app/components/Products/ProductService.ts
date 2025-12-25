import {PRODUCT_AUTOCOMPLETE_URL, PRODUCTS_URL} from "../Shared/Api.ts";
import {Subject} from "rxjs";
import {IProductsResults} from "./types/IProductsResults.ts";
import {IProduct} from "./types/IProduct.ts";
import {IAutocompleteSuggestion} from "./types/IAutocompleteSuggestion.ts";
import {IDisplayTag} from "../../pages/Products/ProductListPage.tsx";
import {IProductDetail} from "./types/IProductDetail.ts";
import {IAddProductReviewRequest} from "./types/IAddProductReviewRequest.ts";
import {IProductReviewResponse} from "./types/IProductReviewResponse.ts";

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
            res.json().then(resp => {
                products$.error(resp?.message || "Unknown error")
            });
        }
    }).catch((err: string) => {
        products$.error(err);
    });
    return products$;
}

export function getProductDetail(slug: string): Subject<IProductDetail> {
    const product$ = new Subject<IProductDetail>();
    fetch(`${PRODUCTS_URL}/${slug}`).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                product$.next({
                    product: resp.results.product,
                    tags: resp.results.tags
                });
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

export function addReview(review: IAddProductReviewRequest, productSlug: string): Subject<boolean> {
    const review$ = new Subject<boolean>();
    fetch(`${PRODUCTS_URL}/${productSlug}/reviews`, {
        method: 'POST',
        body: JSON.stringify(review),
        headers: {
            'Content-Type': 'application/json'
        },
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    review$.next(true);
                } else {
                    review$.error(resp?.message || "Unknown error");
                }
            });
        } else {
            review$.error(res.statusText);
        }
    }).catch((err) => {
        review$.error(err);
    });
    return review$;
}

export function getProductReviews(productSlug: string): Subject<IProductReviewResponse> {
    const reviews$ = new Subject<IProductReviewResponse>();
    fetch(`${PRODUCTS_URL}/${productSlug}/reviews`).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    reviews$.next({
                        reviews: resp.results.reviews,
                        reviewRatingDistributions: resp.results.ratingDistribution,
                    });
                } else {
                    reviews$.error(resp?.message || "Unknown error");
                }
            });
        } else {
            reviews$.error(res.statusText);
        }
    }).catch((err) => {
        reviews$.error(err);
    });
    return reviews$;
}

export function deleteProduct(slug: string): Subject<boolean> {
    const deleteItem$ = new Subject<boolean>();
    fetch(`${PRODUCTS_URL}/${slug}`, {
        method: 'DELETE',
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                if (resp?.status === "OK") {
                    deleteItem$.next(true);
                } else {
                    deleteItem$.error(resp?.message || "Unknown error");
                }
            });
        } else {
            deleteItem$.error(res.statusText);
        }
    }).catch((err) => {
        deleteItem$.error(err);
    });
    return deleteItem$;
}