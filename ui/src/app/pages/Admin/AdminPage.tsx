import {TabMenu} from "primereact/tabmenu";
import {ReactElement, useState} from "react";
import {MenuItem} from "primereact/menuitem";

export default function AdminPage() {
    const [activeIndex, setActiveIndex] = useState<number>(0);
    const items: MenuItem[] = [
        {label: 'Users', icon: 'pi pi-chart-line'},
        {label: 'Orders', icon: 'pi pi-gift'},
    ];
    const adminPages: ReactElement[] = [];

    return (
        <>
            <h1 className="text-3xl font-bold mb-4">Admin</h1>

            <TabMenu
                model={items}
                activeIndex={activeIndex}
                onTabChange={(e) => setActiveIndex(e.index)}
            />
            <div className="flex pt-4">
                {adminPages[activeIndex]}
            </div>
        </>
    );
}