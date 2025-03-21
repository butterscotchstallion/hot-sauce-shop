import {InputText} from "primereact/inputtext";
import {IProduct} from "../Products/IProduct.ts";
import {ReactElement, useEffect, useState} from "react";
import {InputTextarea} from "primereact/inputtextarea";
import SpiceRating from "../Products/SpiceRating.tsx";
import {Button} from "primereact/button";
import {Card} from "primereact/card";
import {NavigateFunction, useNavigate} from "react-router";
import {Controller, SubmitHandler, useForm} from "react-hook-form";
import {ProductSchema} from "../Products/ProductSchema.ts";
import {zodResolver} from "@hookform/resolvers/zod";

interface IAdminInventoryItemFormProps {
    product: IProduct | undefined;
}

type IAdminInventoryItemFormInput = {
    name: string
    price: number
    shortDescription: string
    description: string
    spiceRating: number
}

export default function AdminInventoryItemForm(props: IAdminInventoryItemFormProps): ReactElement {
    const [productName, setProductName] = useState<string>("");
    const [productPrice, setProductPrice] = useState<number>(0);
    const [productShortDescription, setProductShortDescription] = useState<string>("");
    const [productDescription, setProductDescription] = useState<string>("");
    const [spiceRating, setSpiceRating] = useState<number>(3);
    const [productSlug, setProductSlug] = useState<string>("");
    const navigate: NavigateFunction = useNavigate();
    const {
        handleSubmit,
        control,
        formState: {errors},
    } = useForm<IAdminInventoryItemFormInput>({
        resolver: zodResolver(ProductSchema),
    });
    const onSubmit: SubmitHandler<IAdminInventoryItemFormInput> = (data: IAdminInventoryItemFormInput) => {
        console.log(data)
    }

    useEffect(() => {
        if (props.product) {
            setProductName(props.product.name);
            setProductPrice(props.product.price);
            setProductShortDescription(props.product.shortDescription);
            setProductDescription(props.product.description);
            setSpiceRating(props.product.spiceRating);
            setProductSlug(props.product.slug);
        }
    }, [props.product]);

    const goToProductPage = () => {
        if (productSlug) {
            navigate(`/products/${productSlug}`);
        } else {
            console.error("Product slug is not set");
        }
    }

    return (
        <>
            <form onSubmit={handleSubmit(onSubmit)} className="w-full m-0 p-0">
                <section className="flex flex-col gap-4 w-full">
                    <section className="flex w-full justify-between">
                        <h1 className="font-bold text-2xl mb-4">Editing {productName}</h1>
                        <div className="flex justify-end w-[300px] gap-4">
                            <Button label="View Product"
                                    icon="pi pi-eye"
                                    severity="info"
                                    onClick={() => goToProductPage()}/>
                            <Button type="submit" label="Save" icon="pi pi-save"/>
                        </div>
                    </section>

                    <section className="flex flex-col gap-4 w-full">
                        <Card>
                            <section className="flex w-full">
                                <div className="flex w-1/2 justify-between gap-10">
                                    <div className="w-full mb-4">
                                        <label className="mb-2 block" htmlFor="name">Name</label>
                                        <Controller
                                            name="name"
                                            control={control}
                                            render={({field}) => {
                                                return <>
                                                    <InputText
                                                        className="w-full"
                                                        {...field}
                                                        value={productName}
                                                        invalid={!!errors.name?.message}
                                                        onChange={(e) => {
                                                            setProductName(e.target.value)
                                                        }}/>
                                                </>
                                            }}
                                        />
                                        <p>{errors.name && errors.name.message}</p>
                                    </div>

                                    <div className="w-full">
                                        <label className="mb-2 block">Spice Rating</label>
                                        <SpiceRating
                                            rating={spiceRating}
                                            readOnly={false}
                                            onChange={(rating: number) => setSpiceRating(rating)}
                                        />
                                    </div>
                                </div>
                            </section>

                            <section className="flex w-full">
                                <div className="flex gap-10">
                                    <div className="w-full mb-4">
                                        <label className="mb-2 block" htmlFor="price">Price</label>
                                        <InputText type="number" value={productPrice.toString()} onChange={(e) => {
                                            setProductPrice(Number(e.target.value))
                                        }}/>
                                    </div>
                                </div>
                            </section>

                            <section className="w-1/2 flex justify-between">
                                <div>
                                    <label className="mb-2 block" htmlFor="shortDescription">Short Description</label>
                                    <InputTextarea rows={5} cols={30} value={productShortDescription} onChange={(e) => {
                                        setProductShortDescription(e.target.value)
                                    }}/>
                                </div>
                                <div>
                                    <label className="mb-2 block" htmlFor="description">Description</label>
                                    <InputTextarea rows={5} cols={30} value={productDescription} onChange={(e) => {
                                        setProductDescription(e.target.value)
                                    }}/>
                                </div>
                            </section>
                        </Card>
                    </section>
                </section>
            </form>
        </>
    )
}