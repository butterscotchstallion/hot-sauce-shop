import {CartItemsDataTable} from "../components/Cart/CartItemsDataTable.tsx";
import {DataTable} from "primereact/datatable";
import {IDeliveryOption} from "../components/Orders/IDeliveryOption.ts";
import {Column} from "primereact/column";
import {useEffect, useState} from "react";
import {Button} from "primereact/button";
import {useSelector} from "react-redux";
import {RootState} from "../store.ts";
import dayjs, {Dayjs} from 'dayjs';
import {Tooltip} from 'primereact/tooltip';
import {Card} from "primereact/card";
import {Link, NavigateFunction, useNavigate} from "react-router";
import {InputText} from "primereact/inputtext";

interface IOrderTotalItems {
    name: string;
    price: number;
}

export function OrderCheckoutPage() {
    const navigate: NavigateFunction = useNavigate();
    const [convenienceFee] = useState(Math.random() * 20);
    const [couponCodes, setCouponCodes] = useState<string[]>([]);
    const [couponCode, setCouponCode] = useState<string>("");
    const cartSubtotal: number = useSelector((state: RootState) => state.cart.cartSubtotal);
    const today: Dayjs = dayjs();
    const twoDay: Dayjs = today.add(1, "days");
    const threeDay: Dayjs = today.add(2, "days");
    const whenever: Dayjs = today.add(7, "days");
    const instantTransmission: Dayjs = today.add(1, "hours");
    const deliveryDateFormat: string = "ddd, MMM D";
    const deliveryOptions: IDeliveryOption[] = [
        {
            name: "Instant Transmission",
            price: 99.99,
            deliveryDate: instantTransmission.format(deliveryDateFormat),
            description: "Teleported via black hole after packaging"
        },
        {
            name: "Two Day",
            price: 9.99,
            deliveryDate: twoDay.format(deliveryDateFormat)
        },
        {
            name: "Three Day",
            price: 4.99,
            deliveryDate: threeDay.format(deliveryDateFormat)
        },
        {
            name: "Whenever",
            price: 0.00,
            deliveryDate: whenever.format(deliveryDateFormat),
            description: "Usually about a week"
        },
    ];
    const [selectedDeliveryOption, setSelectedDeliveryOption] = useState<IDeliveryOption>(deliveryOptions[0]);
    const [orderTotal, setOrderTotal] = useState<string>(cartSubtotal.toFixed(2));
    const getTaxAmount = (): number => {
        return (parseFloat(orderTotal) * 0.06)
    }
    const orderTotalItems: IOrderTotalItems[] = [
        {name: "Subtotal", price: cartSubtotal},
        {name: "Shipping & Handling", price: selectedDeliveryOption.price},
        {name: "Estimated Taxes", price: getTaxAmount()},
        {name: "Convenience Fee", price: convenienceFee}
    ];
    const priceFormatted = (row: IDeliveryOption) => {
        return row.price > 0 ? `$${row.price.toFixed(2)}` : <strong className="text-yellow-200">FREE</strong>;
    }
    const deliveryOptionName = (row: IDeliveryOption) => {
        return <>
            <p>
                {row.name} {row?.description &&
                <i className="pl-2 cursor-pointer pi pi-question-circle custom-target-icon text-yellow-200"
                   data-pr-tooltip={row.description}
                   data-pr-position="right"
                   data-pr-at="right+5 top"
                   data-pr-my="left center-2"></i>}
            </p>
        </>
    }
    useEffect(() => {
        const newOrderTotal: string = (parseFloat(String(cartSubtotal)) + selectedDeliveryOption.price).toFixed(2);
        setOrderTotal(newOrderTotal);
    }, [cartSubtotal, selectedDeliveryOption]);
    return (
        <>
            <h1 className="text-2xl font-bold mb-4">Checkout</h1>

            <section className="flex justify-between gap-4">
                <CartItemsDataTable hideSubtotal={true}/>

                <aside className="w-2/3">
                    <section className="flex flex-col gap-4">
                        <section>
                            <h4 className="text-xl font-bold mb-2">Delivery Options</h4>
                            <DataTable
                                value={deliveryOptions}
                                selection={selectedDeliveryOption}
                                onSelectionChange={(e) => setSelectedDeliveryOption(e.value || deliveryOptions[0])}>
                                <Column selectionMode="single" headerStyle={{width: '3rem'}}></Column>
                                <Column field="name" header="Name" body={deliveryOptionName}/>
                                <Column field="price" header="Price" body={priceFormatted}/>
                                <Column field="deliveryDate" header="Delivery Date"/>
                            </DataTable>
                        </section>

                        <Card>
                            <section className="flex justify-between gap-4 items-center">
                                <div>
                                    <p>Delivering to John HotSauceLover</p>
                                    <p>123 MAPLE ST, Boulder CO, 12345, United States</p>
                                </div>
                                <Button
                                    onClick={() => navigate("/account/addresses")}
                                    link
                                    label="Change"
                                    icon="pi pi-address-book"/>
                            </section>
                        </Card>

                        <section>
                            <Card>
                                <section className="flex justify-between items-center gap-4">
                                    <div className="flex flex-col gap-6">
                                        <section>
                                            Paying with
                                            <Link to={`/account/payment-methods/1`}>
                                                <i className="pi pi-credit-card mr-1 ml-2"></i> Visa 1234
                                            </Link>
                                        </section>

                                        <section>
                                            <div className="flex gap-4">
                                                <InputText
                                                    type="text"
                                                    className="p-inputtext-sm"
                                                    placeholder="Enter coupon code"
                                                    maxLength={20}
                                                    value={couponCode}
                                                    onChange={(e) => setCouponCode(e.target.value)}
                                                />
                                                <Button
                                                    label="Apply"
                                                    icon="pi pi-plus"
                                                    size="small"
                                                    disabled={couponCode.length === 0}
                                                />
                                            </div>
                                        </section>
                                    </div>
                                    <Button link label="Change" icon="pi pi-wallet"
                                            onClick={() => navigate("/account/payment-methods")}/>
                                </section>
                            </Card>
                        </section>

                        <section>
                            <DataTable value={orderTotalItems}>
                                <Column field="name" header="Item"/>
                                <Column field="price" header="Cost" body={priceFormatted}/>
                            </DataTable>
                        </section>

                        <section>
                            <section className="flex justify-between gap-4">
                                <h4 className="text-xl font-bold mb-2">
                                    Order total: ${orderTotal}
                                </h4>
                                <Button label="Place Order" icon="pi pi-cart-arrow-down"/>
                            </section>
                        </section>
                    </section>
                </aside>
            </section>

            <Tooltip target=".custom-target-icon"/>
        </>
    )
}