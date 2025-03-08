import {IProduct} from "./IProduct.ts";
import {ReactElement} from "react";
import ProductCard from "./ProductCard.tsx";

interface IProductListProps {
    products: IProduct[]
}

export default function ProductList(props: IProductListProps): ReactElement {
    const productsList: ReactElement[] = props.products?.map((product: IProduct, index: number): ReactElement => (
        <ProductCard product={product} key={index}/>
    ))
    return (
        <>
            <section className="flex gap-4 flex-wrap">
                {productsList}
            </section>
        </>
    );
}