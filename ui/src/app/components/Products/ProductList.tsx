import {IProduct} from "./IProduct.ts";
import {ReactElement} from "react";
import ProductCard from "./ProductCard.tsx";
import {Skeleton} from "primereact/skeleton";

interface IProductListProps {
    products: IProduct[],
}

export default function ProductList(props: IProductListProps): ReactElement {
    const productsList: ReactElement[] = props.products?.map((product: IProduct, index: number): ReactElement => (
        <ProductCard product={product} key={index}/>
    ));
    const skeletonList: ReactElement[] = new Array(10).fill(0).map((): ReactElement => (
        <Skeleton size="260px"></Skeleton>
    ));
    return (
        <>
            <section className="flex gap-4 flex-wrap">
                {props.products ? productsList : skeletonList}
            </section>
        </>
    );
}