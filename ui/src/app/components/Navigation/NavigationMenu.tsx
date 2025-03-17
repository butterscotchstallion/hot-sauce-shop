import {ReactElement} from 'react';
import {Menubar} from 'primereact/menubar';
import {MenuItem} from 'primereact/menuitem';
import {Avatar} from 'primereact/avatar';
import {NavLink} from "react-router";
import CartSidebar from "../Cart/CartSidebar.tsx";
import ProductAutocomplete from "../Products/ProductAutocomplete.tsx";

export default function NavigationMenu(): ReactElement {
    const itemRenderer: (item: MenuItem) => ReactElement = (item: MenuItem) => (
        <NavLink className="flex align-items-center p-menuitem-link" to={item.url || '/'}>
            <span className={item.icon}/>
            <span className="mx-2">{item.label}</span>
        </NavLink>
    );
    const items: MenuItem[] = [
        {
            label: 'Home',
            icon: 'pi pi-home',
            url: "/",
            template: itemRenderer,
        },
        {
            label: 'Products',
            icon: 'pi pi-gift',
            url: "/products",
            template: itemRenderer,
        },
        {
            label: 'Contact',
            icon: 'pi pi-envelope',
            template: itemRenderer,
            url: "/contact"
        }
    ];

    const start = <div
        className="text-2xl pl-2 font-bold w-[200px] all-small-caps transition-colors duration-300 hover:text-orange-500 hover:animate-pulse">
        <NavLink to="/">Caliente Corner</NavLink>
    </div>;
    const end = (
        <div className="flex align-items-center gap-4 pr-4">
            <ProductAutocomplete/>
            <CartSidebar/>
            <Avatar className="ml-2 cursor-pointer" image="/images/avatars/amyelsner.png" shape="circle"/>
        </div>
    );

    return (
        <div className="card">
            <Menubar model={items} start={start} end={end}/>
        </div>
    )
}
