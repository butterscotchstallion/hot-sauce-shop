import {ReactElement} from "react";
import {Checkbox} from "primereact/checkbox";
import {IDisplayTag} from "../../pages/ProductsPage.tsx";

interface IProductFilterListProps {
    tags: IDisplayTag[],
    toggleFilter: (checked: boolean) => void,
}

export default function ProductFilterList(props: IProductFilterListProps): ReactElement {
    return (
        <>
            <ul>
                {props.tags.map((filter: IDisplayTag, index: number) => (
                    <li key={index}>
                        <Checkbox
                            inputId={"filter-" + filter.slug}
                            onChange={e => props.toggleFilter(!!e.checked)}
                            checked={filter.checked}
                        />
                        <label htmlFor={"filter-" + filter.slug} className="pl-2 cursor-pointer">
                            {filter.name}
                        </label>
                    </li>
                ))}
            </ul>
        </>
    )

}