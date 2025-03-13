import {ReactElement, RefObject, SetStateAction, Suspense, useEffect, useRef, useState} from "react";
import {Card} from "primereact/card";
import ProductFilterList from "../components/Products/ProductFilterList.tsx";
import {IProduct} from "../components/Products/IProduct.ts";
import {Subscription} from "rxjs";
import ProductList from "../components/Products/ProductList.tsx";
import Throbber from "../components/Shared/Throbber.tsx";
import {ITag} from "../components/Tag/ITag.ts";
import {getTags} from "../components/Tag/TagService.ts";
import {IProductsResults} from "../components/Products/IProductsResults.ts";
import {Paginator} from "primereact/paginator";
import {getProducts} from "../components/Products/ProductService.ts";
import {Toast} from "primereact/toast";

// Tags with the checked attribute layered on top for
// the purposes of this component
export interface IDisplayTag extends ITag {
    checked: boolean;
}

export default function ProductListPage(): ReactElement {
    const toast: RefObject<Toast | null> = useRef(null);
    const [offset, setOffset] = useState<number>(0);
    const [perPage, setPerPage] = useState<number>(10);
    const [products, setProducts] = useState<IProduct[]>([]);
    const [tags, setTags] = useState<IDisplayTag[]>([]);
    const [totalProducts, setTotalProducts] = useState<number>(0);
    const [isErrorLoadingTags, setIsErrorLoadingTags] = useState<boolean>(false);
    const [isErrorLoadingProducts, setIsErrorLoadingProducts] = useState<boolean>(false);

    function toggleFilter(checked: boolean) {
        console.log(checked);
    }

    const onPageChange = (event: { first: SetStateAction<number>; rows: SetStateAction<number>; }) => {
        setOffset(event.first);
        setPerPage(event.rows);
    };

    useEffect(() => {
        const products$: Subscription = getProducts(offset, perPage).subscribe({
            next: (results: IProductsResults) => {
                setProducts(results.inventory);
                setTotalProducts(results.total);
            },
            error: (err: string) => {
                if (toast.current) {
                    toast.current.show({
                        severity: 'error',
                        summary: 'Error',
                        detail: 'Error loading products: ' + err,
                        life: 3000,
                    })
                }
                setIsErrorLoadingProducts(true);
            }
        });
        return () => {
            products$.unsubscribe();
        }
    }, [offset, perPage]);

    useEffect(() => {
        const tags$: Subscription = getTags().subscribe({
            next: (results: ITag[]) => {
                const displayTags: IDisplayTag[] = [];
                results.map((tag: ITag) => {
                    displayTags.push({...tag, checked: false});
                });
                setTags(displayTags);
            },
            error: (err) => {
                if (toast.current) {
                    toast.current.show({
                        severity: 'error',
                        summary: 'Error',
                        detail: 'Error loading filters: ' + err,
                        life: 3000,
                    })
                }
                setIsErrorLoadingTags(true);
            }
        });
        return () => {
            tags$.unsubscribe();
        }
    }, []);

    return (
        <>
            <div className="flex gap-4">
                <menu className="w-[205px]">
                    {tags.length > 0 ? (
                        <Card title="Filters">
                            <ProductFilterList tags={tags} toggleFilter={toggleFilter}/>
                        </Card>
                    ) : ""}
                </menu>

                <section className="w-full">
                    <h2 className="font-bold text-lg mb-4">Products</h2>

                    {products.length > 0 ? (
                        <Suspense fallback={<Throbber/>}>
                            <ProductList products={products} toast={toast}/>

                            {products.length > 0 ? (
                                <div className="card mt-4 mb-4">
                                    <Paginator first={offset}
                                               rows={perPage}
                                               totalRecords={totalProducts}
                                               rowsPerPageOptions={[10, 20, 30]}
                                               onPageChange={onPageChange}/>
                                </div>
                            ) : ""}
                        </Suspense>
                    ) : (
                        <Throbber/>
                    )}
                </section>
            </div>

            <Toast ref={toast}/>
        </>
    )
}