import {useState} from "react";
import {AutoComplete, AutoCompleteChangeEvent, AutoCompleteCompleteEvent} from "primereact/autocomplete";
import {getProductAutocompleteSuggestions} from "./ProductService.ts";
import {debounceTime} from "rxjs";
import {IAutocompleteSuggestion} from "./IAutocompleteSuggestion.ts";

export default function ProductAutocomplete() {
    const [value, setValue] = useState<string>('');
    const [items, setItems] = useState<string[]>([]);

    const search = (event: AutoCompleteCompleteEvent) => {
        getProductAutocompleteSuggestions(event.query).pipe(
            debounceTime(250)
        ).subscribe((suggestions: IAutocompleteSuggestion[]) => {
            setItems(suggestions.map((s: IAutocompleteSuggestion) => s.name));
        });
    }

    return (
        <AutoComplete value={value}
                      suggestions={items}
                      completeMethod={search}
                      placeholder="Search"
                      onChange={(e: AutoCompleteChangeEvent) => setValue(e.value)} forceSelection/>
    )
}