import {InputText} from "primereact/inputtext";
import {IProduct} from "../Products/IProduct.ts";
import {ReactElement, useEffect, useState} from "react";
import {InputTextarea} from "primereact/inputtextarea";

interface IAdminInventoryItemFormProps {
    product: IProduct | undefined;
}

export default function AdminInventoryItemForm(props: IAdminInventoryItemFormProps): ReactElement {
    const [productName, setProductName] = useState<string>("");
    const [productPrice, setProductPrice] = useState<number>(0);
    const [productShortDescription, setProductShortDescription] = useState<string>("");
    const [productDescription, setProductDescription] = useState<string>("");

    useEffect(() => {
        if (props.product) {
            setProductName(props.product.name);
            setProductPrice(props.product.price);
            setProductShortDescription(props.product.shortDescription);
            setProductDescription(props.product.description);
        }
    }, [props.product]);

    return (
        <>
            <section className="flex flex-col gap-4 w-full">
                <section className="flex w-full">
                    <div className="flex gap-10">
                        <div>
                            <label className="mb-2 block" htmlFor="name">Name</label>
                            <InputText value={productName} onChange={(e) => {
                                setProductName(e.target.value)
                            }}/>
                        </div>

                        <div>
                            <label className="mb-2 block" htmlFor="shortDescription">Short Description</label>
                            <InputTextarea rows={5} cols={30} value={productShortDescription} onChange={(e) => {
                                setProductShortDescription(e.target.value)
                            }}/>
                        </div>
                    </div>
                </section>

                <section className="flex w-full">
                    <div className="flex gap-10">


                        <div>
                            <label className="mb-2 block" htmlFor="price">Price</label>
                            <InputText value={productPrice.toString()} onChange={(e) => {
                                setProductPrice(Number(e.target.value))
                            }}/>
                        </div>

                        <div>
                            <label className="mb-2 block" htmlFor="description">Description</label>
                            <InputTextarea rows={5} cols={30} value={productDescription} onChange={(e) => {
                                setProductDescription(e.target.value)
                            }}/>
                        </div>
                    </div>
                </section>
            </section>
        </>
    )
}