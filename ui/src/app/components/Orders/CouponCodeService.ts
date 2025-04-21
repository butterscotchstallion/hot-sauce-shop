import {Subject} from "rxjs";
import {ICouponCode} from "./ICouponCode.ts";
import {ORDERS_COUPON_CODE_URL} from "../Shared/Api.ts";

export function getCouponCodeByCode(code: string): Subject<ICouponCode> {
    const code$: Subject<ICouponCode> = new Subject<ICouponCode>();
    fetch(`${ORDERS_COUPON_CODE_URL}/${code}`).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => code$.next(resp?.results?.couponCode));
        } else {
            res.json().then(resp => code$.error(resp?.message || "Unknown error"));
        }
    }).catch((err) => {
        code$.error(err);
    });
    return code$;
}