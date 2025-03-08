import {ReactElement, Suspense, useEffect, useState} from "react";
import {Card} from "primereact/card";
import ProductFilterList from "../components/Products/ProductFilterList.tsx";
import {IProduct} from "../components/Products/IProduct.ts";
import {getProducts} from "../components/Products/ProductService.ts";
import {Subscription} from "rxjs";
import ProductList from "../components/Products/ProductList.tsx";
import Throbber from "../components/Shared/Throbber.tsx";

export interface IFilters {
    displayName: string;
    name: string;
    checked: boolean;
}

export default function ProductsPage(): ReactElement {
    const [products, setProducts] = useState<IProduct[]>([]);
    const [filters, setFilters] = useState<IFilters[]>([
        {displayName: "Habanero", name: "habanero", checked: false},
        {displayName: "Datil", name: "datil", checked: false},
        {displayName: "Jalapeño", name: "jalapeño", checked: false},
    ]);

    function toggleFilter(checked: boolean) {

    }

    useEffect(() => {
        const products$: Subscription = getProducts().subscribe({
            next: (products: IProduct[]) => {
                setProducts(products);
            },
            error: (err) => {
                console.error(err);
            }
        });
        return () => {
            products$.unsubscribe();
        }
    }, []);

    return (
        <>
            <div className="flex gap-4">
                <menu className="w-[200px]">
                    <Card title="Filters">
                        <ProductFilterList filters={filters} toggleFilter={toggleFilter}/>
                    </Card>
                </menu>

                <section className="w-full">
                    <h2 className="font-bold text-lg mb-4">Products</h2>

                    <Suspense fallback={<Throbber/>}>
                        <ProductList products={products}/>
                    </Suspense>
                </section>
            </div>
        </>
    )
}