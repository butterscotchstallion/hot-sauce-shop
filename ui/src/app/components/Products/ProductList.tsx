import {IProduct} from "./IProduct.ts";
import {ReactElement, useState} from "react";
import ProductCard from "./ProductCard.tsx";
import {Paginator} from "primereact/paginator";

interface IProductListProps {
    products: IProduct[],
    total: number,
}

export default function ProductList(props: IProductListProps): ReactElement {
    const [first, setFirst] = useState(0);
    const [rows, setRows] = useState(10);
    const productsList: ReactElement[] = props.products?.map((product: IProduct, index: number): ReactElement => (
        <ProductCard product={product} key={index}/>
    ));
    const onPageChange = (event) => {
        setFirst(event.first);
        setRows(event.rows);
    };
    return (
        <>
            <section className="flex gap-4 flex-wrap">
                {productsList}
            </section>
            <div className="card mt-4 mb-4">
                <Paginator first={first}
                           rows={rows}
                           totalRecords={props.total}
                           rowsPerPageOptions={[10, 20, 30]}
                           onPageChange={onPageChange}/>
            </div>
        </>
    );
}