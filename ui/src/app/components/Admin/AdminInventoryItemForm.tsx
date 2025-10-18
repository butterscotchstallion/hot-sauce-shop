import {InputText} from "primereact/inputtext";
import {IProduct} from "../Products/IProduct.ts";
import {FormEvent, ReactElement, RefObject, useEffect, useRef, useState} from "react";
import {InputTextarea} from "primereact/inputtextarea";
import SpiceRating from "../Products/SpiceRating.tsx";
import {Button} from "primereact/button";
import {Card} from "primereact/card";
import {NavigateFunction, useNavigate} from "react-router";
import {ProductSchema} from "../Products/ProductSchema.ts";
// ZodError is used in an exception not but detected for some reason
// eslint-disable-next-line @typescript-eslint/no-unused-vars
import {z, ZodError, ZodIssue} from "zod";
import {addOrUpdateItem} from "../Products/ProductService.ts";
import {Toast} from "primereact/toast";
import {MultiSelect, MultiSelectChangeEvent} from "primereact/multiselect";
import {ITag} from "../Tag/ITag.ts";
import {IFormErrata} from "../Shared/IFormErrata.ts";

interface IAdminInventoryItemFormProps {
    product: IProduct | undefined;
    productTags: ITag[];
    availableTags: ITag[];
    isNewProduct: boolean;
}

export default function AdminInventoryItemForm(props: IAdminInventoryItemFormProps): ReactElement {
    const toast: RefObject<Toast | null> = useRef<Toast>(null);
    const [productName, setProductName] = useState<string>("");
    const [productPrice, setProductPrice] = useState<number>(0);
    const [productTags, setProductTags] = useState<ITag[]>(props.productTags);
    const [productShortDescription, setProductShortDescription] = useState<string>("");
    const [productDescription, setProductDescription] = useState<string>("");
    const [spiceRating, setSpiceRating] = useState<number>(3);
    const [productSlug, setProductSlug] = useState<string>("");
    const navigate: NavigateFunction = useNavigate();
    const defaultErrata: IFormErrata = {
        name: '',
        description: '',
        shortDescription: '',
        price: '',
        spiceRating: ''
    };
    const [formErrata, setFormErrata] = useState<IFormErrata>(defaultErrata);
    const resetErrata = () => {
        setFormErrata(defaultErrata);
    }
    const validate = (): boolean => {
        try {
            ProductSchema.parse({
                name: productName,
                price: productPrice,
                shortDescription: productShortDescription,
                description: productDescription,
                spiceRating: spiceRating,
                slug: productSlug,
            });
            resetErrata();

            return true;
        } catch (err: ZodError | unknown) {
            if (err instanceof z.ZodError) {
                const newErrata: IFormErrata = {...formErrata};
                err.issues.forEach((issue: ZodIssue) => {
                    newErrata[issue.path[0]] = issue.message;
                });
                setFormErrata(newErrata);
            }
            return false;
        }
    }

    const onSubmit = (event: FormEvent<HTMLElement>) => {
        event.preventDefault();
        const valid: boolean = validate();

        if (valid) {
            const product: IProduct = {
                slug: productSlug,
                name: productName,
                description: productDescription,
                shortDescription: productShortDescription,
                price: productPrice,
                spiceRating: spiceRating,
                tagIds: productTags.map((tag: ITag) => tag.id),
            }
            addOrUpdateItem(product, props.isNewProduct).subscribe({
                next: () => {
                    if (toast.current) {
                        toast?.current.show({
                            severity: 'success',
                            summary: 'Success',
                            detail: 'Product saved successfully',
                            life: 3000,
                        })
                    }
                },
                error: (err) => {
                    console.log(err);
                    if (toast.current) {
                        toast?.current.show({
                            severity: 'error',
                            summary: 'Error',
                            detail: 'Error saving product: ' + err + '.',
                            life: 3000,
                        })
                    }
                }
            });
        }
    }

    useEffect(() => {
        if (!props.isNewProduct && props.product) {
            setProductName(props.product.name);
            setProductPrice(props.product.price);
            setProductShortDescription(props.product.shortDescription);
            setProductDescription(props.product.description);
            setSpiceRating(props.product.spiceRating);
            setProductSlug(props.product.slug);
            setProductTags(props.productTags);
        } else {
            setProductName('');
            setProductPrice(9.99);
            setProductShortDescription('');
            setProductDescription('');
            setSpiceRating(3);
            setProductSlug('');
            resetErrata();
        }
    }, [props.isNewProduct, props.product]);

    const goToProductPage = () => {
        if (productSlug) {
            navigate(`/products/${productSlug}`);
        } else {
            console.error("Product slug is not set");
        }
    }

    const goToNewProductPage = () => {
        navigate('/admin/products/add');
    }

    return (
        <>
            <form onSubmit={onSubmit} className="w-full m-0 p-0">
                <section className="flex flex-col gap-4 w-full">
                    {/* header */}
                    <section className="flex w-full justify-between">
                        <h1 className="font-bold text-2xl mb-4">
                            Editing {props.isNewProduct ? "New Product" : productName}
                        </h1>
                        <div className="flex justify-end w-[500px] gap-4">
                            <Button type="button"
                                    label="Add Product"
                                    icon="pi pi-plus"
                                    severity="info"
                                    onClick={() => goToNewProductPage()}/>
                            <Button type="button"
                                    label="View Product"
                                    icon="pi pi-eye"
                                    severity="info"
                                    onClick={() => goToProductPage()}/>
                            <Button type="submit" label="Save" icon="pi pi-save"/>
                        </div>
                    </section>

                    {/* content */}
                    <section className="flex flex-col gap-4 w-full">
                        <Card>
                            {/* outer content container */}
                            <section className="flex w-full gap-4 justify-between">
                                <section className="left-col">
                                    <section>
                                        <div className="flex w-full justify-between gap-10">
                                            <div className="w-full mb-4">
                                                <label className="mb-2 block" htmlFor="name">Name</label>
                                                <InputText
                                                    className="w-full"
                                                    value={productName}
                                                    invalid={!!formErrata.name}
                                                    onChange={(e) => {
                                                        setProductName(e.target.value);
                                                        validate();
                                                    }}
                                                />
                                                <p className="text-red-500 pt-2">{formErrata.name}</p>
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
                                                <InputText
                                                    type="number"
                                                    invalid={!!formErrata.price}
                                                    value={productPrice.toString()}
                                                    step={.01}
                                                    onChange={(e) => {
                                                        setProductPrice(Number(e.target.value));
                                                        validate();
                                                    }}/>
                                                <p className="text-red-500 pt-2">{formErrata.price}</p>
                                            </div>
                                        </div>
                                    </section>

                                    <section className="w-1/2 flex gap-4 justify-between">
                                        <div>
                                            <label className="mb-2 block" htmlFor="shortDescription">Short
                                                Description</label>
                                            <InputTextarea
                                                invalid={!!formErrata.shortDescription}
                                                rows={5}
                                                cols={30}
                                                value={productShortDescription}
                                                onChange={(e) => {
                                                    setProductShortDescription(e.target.value);
                                                    validate();
                                                }}/>
                                            <p className="text-red-500 pt-2">{formErrata.shortDescription}</p>
                                        </div>
                                        <div>
                                            <label className="mb-2 block" htmlFor="description">Description</label>
                                            <InputTextarea
                                                invalid={!!formErrata.description}
                                                rows={5}
                                                cols={30}
                                                value={productDescription}
                                                onChange={(e) => {
                                                    setProductDescription(e.target.value);
                                                    validate();
                                                }}/>
                                            <p className="text-red-500 pt-2">{formErrata.description}</p>
                                        </div>
                                    </section>
                                </section>
                                <section className="right-col w-1/2">
                                    <h3 className="text-1xl font-bold mb-4">Tags</h3>
                                    <MultiSelect value={productTags}
                                                 onChange={(e: MultiSelectChangeEvent) => setProductTags(e.value)}
                                                 options={props.availableTags}
                                                 optionLabel="name"
                                                 display="chip"
                                                 placeholder="Select Tags"
                                                 className="w-full md:w-20rem"/>
                                </section>
                            </section>
                        </Card>
                    </section>
                </section>
            </form>
            <Toast ref={toast}/>
        </>
    )
}