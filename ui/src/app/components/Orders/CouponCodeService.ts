import {Subject} from "rxjs";
import {ICouponCode} from "./ICouponCode.ts";

export function getCouponCodeByCode(code: string): Subject<ICouponCode> {
    const code$: Subject<ICouponCode> = new Subject<ICouponCode>();

    // TODO: implement coupon code API
    code$.next({
        code: "SPICY20",
        description: "20% off your order",
        expirationDate: new Date(2030, 1, 1),
    });

    return code$;
}