import {Subject} from "rxjs";
import {VOTE_ADD_UPDATE_URL} from "../Shared/Api.ts";
import {IVoteRequest} from "./IVoteRequest.ts";

export function addUpdateVote(postId: number, voteValue: number): Subject<number> {
    const addVote$ = new Subject<number>();
    const voteRequest: IVoteRequest = {
        postId,
        voteValue,
    }
    fetch(`${VOTE_ADD_UPDATE_URL}/${postId}`, {
        credentials: 'include',
        method: 'POST',
        body: JSON.stringify(voteRequest),
        headers: {
            'Content-Type': 'application/json'
        }
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                addVote$.next(resp.results.voteId);
            });
        } else {
            res.json().then(resp => {
                addVote$.error(resp?.message || "Unknown error");
            });
        }
    }).catch((err) => {
        addVote$.error(err);
    });
    return addVote$;
}