import {Subject} from "rxjs";
import {IProduct} from "../Products/IProduct.ts";
import {TAGS_URL} from "../Shared/Api.ts";

export function getTags() {
    const tags$ = new Subject<IProduct[]>();
    fetch(TAGS_URL).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => tags$.next(resp.results.tags));
        } else {
            tags$.error(res.statusText);
        }
    }).catch((err) => {
        tags$.error(err);
    });
    return tags$;
}