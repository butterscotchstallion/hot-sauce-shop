import {Subject} from "rxjs";
import {TAGS_URL} from "../Shared/Api.ts";
import {ITag} from "./ITag.ts";

export function getTags(): Subject<ITag[]> {
    const tags$ = new Subject<ITag[]>();
    fetch(TAGS_URL).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => tags$.next(resp?.results?.tags || []));
        } else {
            tags$.error(res.json().then(resp => resp?.message || "Unknown error"));
        }
    }).catch((err) => {
        tags$.error(err);
    });
    return tags$;
}