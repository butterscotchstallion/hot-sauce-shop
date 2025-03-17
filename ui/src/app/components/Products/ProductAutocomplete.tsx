import {useState} from "react";
import {
    AutoComplete,
    AutoCompleteChangeEvent,
    AutoCompleteCompleteEvent,
    AutoCompleteSelectEvent
} from "primereact/autocomplete";
import {getProductAutocompleteSuggestions} from "./ProductService.ts";
import {debounceTime} from "rxjs";
import {IAutocompleteSuggestion} from "./IAutocompleteSuggestion.ts";
import {NavigateFunction, useNavigate} from "react-router";

export default function ProductAutocomplete() {
    const [value, setValue] = useState<string>('');
    const [items, setItems] = useState<string[]>([]);
    const [nameSlugMap, setNameSlugMap] = useState<Map<string, string>>(new Map<string, string>());
    const navigate: NavigateFunction = useNavigate();

    const search = (event: AutoCompleteCompleteEvent) => {
        getProductAutocompleteSuggestions(event.query).pipe(
            debounceTime(250)
        ).subscribe((suggestions: IAutocompleteSuggestion[]) => {
            const slugMap = new Map<string, string>();
            setItems(suggestions.map((s: IAutocompleteSuggestion) => {
                slugMap.set(s.name, s.slug);
                return s.name;
            }));
            setNameSlugMap(slugMap);
        });
    }

    const onSelect = (productName: string) => {
        const slug: string | undefined = nameSlugMap.get(productName);
        if (slug) {
            setValue('');
            navigate(`/products/${slug}`);
        }
    }

    return (
        <AutoComplete value={value}
                      suggestions={items}
                      completeMethod={search}
                      placeholder="Search"
                      onSelect={(e: AutoCompleteSelectEvent) => onSelect(e.value)}
                      onChange={(e: AutoCompleteChangeEvent) => setValue(e.value)}
                      forceSelection/>
    )
}