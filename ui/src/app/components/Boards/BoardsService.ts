import {IBoard} from "./IBoard.ts";
import {Subject} from "rxjs";
import {BOARDS_URL} from "../Shared/Api.ts";

export function getBoards(): Subject<IBoard[]> {
    const boards$ = new Subject<IBoard[]>();
    fetch(`${BOARDS_URL}`, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                boards$.next(resp?.results?.boards || null);
            });
        } else {
            boards$.error(res.statusText);
        }
    }).catch((err) => {
        boards$.error(err);
    });
    return boards$;
}
