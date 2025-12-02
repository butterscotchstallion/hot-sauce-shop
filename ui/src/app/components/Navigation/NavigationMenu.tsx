import {ReactElement} from 'react';
import {Menubar} from 'primereact/menubar';
import {MenuItem} from 'primereact/menuitem';
import {NavLink} from "react-router";
import CartSidebar from "../Cart/CartSidebar.tsx";
import ProductAutocomplete from "../Products/ProductAutocomplete.tsx";
import UserAvatarMenu from "../User/UserAvatarMenu.tsx";
import {UserLevel} from "../User/UserLevel.tsx";
import {WS} from "../Shared/WS.tsx";

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
            url: "/posts",
            template: itemRenderer,
        },
        {
            label: 'Products',
            icon: 'pi pi-gift',
            url: "/products",
            template: itemRenderer,
        },
        {
            label: 'Posts',
            icon: 'pi pi-envelope',
            template: itemRenderer,
            url: "/posts"
        }
    ];

    const start = <div
        className="text-2xl pl-2 font-bold w-[200px] all-small-caps transition-colors duration-300 hover:text-orange-500 hover:animate-pulse">
        <NavLink to="/">Caliente Corner</NavLink>
    </div>;
    const end = (
        <div className="flex align-items-center align-middle gap-4 pr-4">
            <ProductAutocomplete/>
            <CartSidebar/>
            <UserLevel/>
            <div style={{"lineHeight": 2.5}}><WS/></div>
            <UserAvatarMenu/>
        </div>
    );

    return (
        <div className="card">
            <Menubar model={items} start={start} end={end}/>
        </div>
    )
}
