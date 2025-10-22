import {Subject} from "rxjs";
import {USER_VOTE_MAP_URL, VOTE_ADD_UPDATE_URL} from "../Shared/Api.ts";
import {IVoteRequest} from "./IVoteRequest.ts";
import {IVoteMap} from "./IVoteMap.ts";

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

export function getUserVoteMap(): Subject<Map<number, number>> {
    const voteMap$ = new Subject<Map<number, number>>
    fetch(USER_VOTE_MAP_URL, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                const resultMap: IVoteMap = resp.results.voteMap;
                const voteMap: Map<number, number> = new Map<number, number>();
                for (const postId in resultMap) {
                    voteMap.set(parseInt(postId), resultMap[postId]);
                }
                voteMap$.next(voteMap);
            });
        }
    }).catch((err) => {
        voteMap$.error(err);
    });
    return voteMap$;
}

export function getPostIdVoteValueMap(): Subject<Map<number, number>> {
    const postIdVoteValueMap$ = new Subject<Map<number, number>>
    fetch(USER_VOTE_MAP_URL, {
        credentials: 'include'
    }).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => {
                const resultMap: IVoteMap = resp.results.voteMap;
                const voteMap: Map<number, number> = new Map<number, number>();
                for (const postId in resultMap) {
                    voteMap.set(parseInt(postId), resultMap[postId]);
                }
                postIdVoteValueMap$.next(voteMap);
            });
        }
    }).catch((err) => {
        postIdVoteValueMap$.error(err);
    });
    return postIdVoteValueMap$;
}