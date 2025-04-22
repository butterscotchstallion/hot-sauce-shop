import {Subject} from "rxjs";
import {IShippingOption} from "./IShippingOption.ts";
import {ORDERS_SHIPPING_OPTIONS_URL} from "../Shared/Api.ts";

export function getShippingOptions() {
    const options$ = new Subject<IShippingOption[]>();
    fetch(ORDERS_SHIPPING_OPTIONS_URL).then((res: Response) => {
        if (res.ok) {
            res.json().then(resp => options$.next(resp?.results?.shippingOptions));
        } else {
            res.json().then(resp => options$.error(resp?.message || "Unknown error"));
        }
    }).catch((err) => {
        options$.error(err);
    });
    return options$;
}