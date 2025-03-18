import {InputText} from "primereact/inputtext";

export default function AdminInventoryItemForm() {

    return (
        <>
            <div className="flex flex-col gap-4">
                <div>
                    <label className="mb-4 block" htmlFor="name">Name</label>
                    <InputText value={product?.name} onChange={(e) => setValue("name", e.target.value)}/>
                </div>
            </div>
        </>
    )
}