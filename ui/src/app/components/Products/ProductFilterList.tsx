import {ReactElement} from "react";
import {Checkbox} from "primereact/checkbox";
import {IFilters} from "../../pages/ProductsPage.tsx";

interface IProductFilterListProps {
    filters: IFilters[],
    toggleFilter: (checked: boolean) => void,
}

export default function ProductFilterList(props: IProductFilterListProps): ReactElement {
    return (
        <>
            <ul>
                {props.filters.map((filter: IFilters, index: number) => (
                    <li key={index}>
                        <Checkbox
                            inputId={"filter-" + filter.name}
                            onChange={e => props.toggleFilter(!!e.checked)}
                            checked={filter.checked}
                        />
                        <label htmlFor={"filter-" + filter.name} className="pl-2 cursor-pointer">
                            {filter.displayName}
                        </label>
                    </li>
                ))}
            </ul>
        </>
    )

}