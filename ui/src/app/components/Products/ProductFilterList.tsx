import * as React from "react";
import {ReactElement, useEffect, useState} from "react";
import {Checkbox, CheckboxChangeEvent} from "primereact/checkbox";
import {IDisplayTag} from "../../pages/Products/ProductListPage.tsx";
import {ITag} from "../Tag/ITag.ts";
import {Subscription} from "rxjs";
import {getTags} from "../Tag/TagService.ts";
import {Toast} from "primereact/toast";

interface IProductFilterListProps {
    toast: React.RefObject<Toast | null>;
    onFiltersChanged: (filters: IDisplayTag[]) => void;
}

export default function ProductFilterList(props: IProductFilterListProps): ReactElement {
    const [tags, setTags] = useState<IDisplayTag[]>([]);

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
                if (props.toast.current) {
                    props.toast.current.show({
                        severity: 'error',
                        summary: 'Error',
                        detail: 'Error loading filters: ' + err,
                        life: 3000,
                    })
                }
            }
        });
        return () => {
            tags$.unsubscribe();
        }
    }, [props.toast]);

    function toggleFilter(checked: boolean, index: number) {
        tags[index].checked = checked;
        setTags([...tags]);
        props.onFiltersChanged(tags.filter((tag: IDisplayTag): boolean => tag.checked));
    }

    return (
        <>
            <ul>
                {tags.map((tag: IDisplayTag, index: number) => (
                    <li key={index}>
                        <Checkbox
                            inputId={"filter-" + tag.slug}
                            onChange={(e: CheckboxChangeEvent) => toggleFilter(!!e.checked, index)}
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