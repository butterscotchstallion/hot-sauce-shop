import {Subject} from "rxjs";
import {ICouponCode} from "./ICouponCode.ts";

export function getCouponCodeByCode(code: string): Subject<ICouponCode> {
    const code$: Subject<ICouponCode> = new Subject<ICouponCode>();

    // TODO: implement coupon code API
    code$.next();

    return code$;
}