import {ReactElement, RefObject, SetStateAction, Suspense, useEffect, useRef, useState} from "react";
import {Card} from "primereact/card";
import {IProduct} from "../components/Products/IProduct.ts";
import {Subscription} from "rxjs";
import ProductList from "../components/Products/ProductList.tsx";
import Throbber from "../components/Shared/Throbber.tsx";
import {ITag} from "../components/Tag/ITag.ts";
import {IProductsResults} from "../components/Products/IProductsResults.ts";
import {Paginator} from "primereact/paginator";
import {getProducts} from "../components/Products/ProductService.ts";
import {Toast} from "primereact/toast";
import {Dropdown} from "primereact/dropdown";
import ProductFilterList from "../components/Products/ProductFilterList.tsx";

// Tags with the checked attribute layered on top for
// the purposes of this component
export interface IDisplayTag extends ITag {
    checked: boolean;
}

interface IProductListSortOptions {
    label: string,
    value: string,
}

export default function ProductListPage(): ReactElement {
    const toast: RefObject<Toast | null> = useRef(null);
    const [offset, setOffset] = useState<number>(0);
    const [perPage, setPerPage] = useState<number>(10);
    const [products, setProducts] = useState<IProduct[]>([]);
    const [totalProducts, setTotalProducts] = useState<number>(0);
    const [productListSortKey, setProductListSortKey] = useState<string>("name");
    const productListSortOptions: IProductListSortOptions[] = [
        {label: "Sort by Name", value: "name"},
        {label: "Sort by Price", value: "price"},
        {label: "Sort by Spice Rating", value: "spice_rating"},
        {label: "Sort by Date Added", value: "created_at"},
        {label: "Sort by Review Count", value: "review_count"},
    ];
    const [filters, setFilters] = useState<IDisplayTag[]>([]);

    const onPageChange = (event: { first: SetStateAction<number>; rows: SetStateAction<number>; }) => {
        setOffset(event.first);
        setPerPage(event.rows);
    };

    /**
     * Only enabled filters
     * @param filters
     */
    const onFiltersChanged = (filters: IDisplayTag[]) => {
        setFilters(filters);
    };

    useEffect(() => {
        const products$: Subscription = getProducts(offset, perPage, productListSortKey, filters).subscribe({
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
            }
        });
        return () => {
            products$.unsubscribe();
        }
    }, [filters, offset, perPage, productListSortKey]);

    return (
        <>
            <div className="flex gap-4">
                <menu className="w-[205px]">
                    <Card title="Filters">
                        {<ProductFilterList onFiltersChanged={onFiltersChanged} toast={toast}/>}
                    </Card>
                </menu>

                <section className="w-full">
                    <section className="flex justify-between mb-4 pr-1">
                        <h2 className="font-bold text-lg mb-4" data-testid="page-header">Products</h2>
                        <Dropdown value={productListSortKey}
                                  onChange={(e) => setProductListSortKey(e.value)}
                                  options={productListSortOptions}
                                  optionLabel="label"
                                  className="w-[10rem]"/>
                    </section>

                    {products.length > 0 ? (
                        <Suspense fallback={<Throbber/>}>
                            <ProductList products={products} toast={toast}/>

                            {products.length >= perPage ? (
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
                        <>No products found.</>
                    )}
                </section>
            </div>

            <Toast ref={toast}/>
        </>
    )
}