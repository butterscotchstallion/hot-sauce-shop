import {ReactElement, Suspense, useEffect, useState} from "react";
import {Card} from "primereact/card";
import ProductFilterList from "../components/Products/ProductFilterList.tsx";
import {IProduct} from "../components/Products/IProduct.ts";
import {getProducts} from "../components/Products/ProductService.ts";
import {Subscription} from "rxjs";
import ProductList from "../components/Products/ProductList.tsx";
import Throbber from "../components/Shared/Throbber.tsx";
import {ITag} from "../components/Tag/ITag.ts";
import {getTags} from "../components/Tag/TagService.ts";
import {IProductsResults} from "../components/Products/IProductsResults.ts";

// Tags with the checked attribute layered on top for
// the purposes of this component
export interface IDisplayTag extends ITag {
    checked: boolean;
}

export default function ProductsPage(): ReactElement {
    const [products, setProducts] = useState<IProduct[]>([]);
    const [tags, setTags] = useState<IDisplayTag[]>([]);
    const [totalProducts, setTotalProducts] = useState(0);

    function toggleFilter(checked: boolean) {

    }

    useEffect(() => {
        const products$: Subscription = getProducts().subscribe({
            next: (results: IProductsResults) => {
                setProducts(results.inventory);
                setTotalProducts(results.total);
            },
            error: (err) => {
                console.error(err);
            }
        });
        const tags$: Subscription = getTags().subscribe({
            next: (results: ITag[]) => {
                const displayTags: IDisplayTag[] = [];
                results.map((tag: ITag) => {
                    displayTags.push({...tag, checked: false});
                });
                setTags(displayTags);
            },
            error: (err) => {
                console.error(err);
            }
        });
        return () => {
            products$.unsubscribe();
            tags$.unsubscribe();
        }
    }, []);

    return (
        <>
            <div className="flex gap-4">
                <menu className="w-[205px]">
                    <Card title="Filters">
                        <ProductFilterList tags={tags} toggleFilter={toggleFilter}/>
                    </Card>
                </menu>

                <section className="w-full">
                    <h2 className="font-bold text-lg mb-4">Products</h2>

                    <Suspense fallback={<Throbber/>}>
                        <ProductList products={products} total={totalProducts}/>
                    </Suspense>
                </section>
            </div>
        </>
    )
}