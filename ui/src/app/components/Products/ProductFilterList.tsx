import {ReactElement} from "react";
import {Checkbox, CheckboxChangeEvent} from "primereact/checkbox";
import {IDisplayTag} from "../../pages/ProductListPage.tsx";
import {ITag} from "../Tag/ITag.ts";

interface IProductFilterListProps {
    tags: IDisplayTag[],
    toggleFilter: (tag: ITag, checked: boolean) => void,
}

export default function ProductFilterList(props: IProductFilterListProps): ReactElement {
    return (
        <>
            <ul>
                {props.tags.map((tag: IDisplayTag, index: number) => (
                    <li key={index}>
                        <Checkbox
                            inputId={"filter-" + tag.slug}
                            onChange={(e: CheckboxChangeEvent) => props.toggleFilter(tag, !!e.checked)}
                            checked={tag.checked}
                        />
                        <label htmlFor={"filter-" + tag.slug}
                               className="pl-2 cursor-pointer"
                               title={tag.description}
                        >
                            {tag.name}
                        </label>
                    </li>
                ))}
            </ul>
        </>
    )
}