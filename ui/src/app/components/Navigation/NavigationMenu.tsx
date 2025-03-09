import {ReactElement} from 'react';
import {Menubar} from 'primereact/menubar';
import {InputText} from 'primereact/inputtext';
import {MenuItem} from 'primereact/menuitem';
import {Avatar} from 'primereact/avatar';
import {NavLink} from "react-router";

export default function NavigationMenu(): ReactElement {
    const itemRenderer = (item) => (
        <NavLink className="flex align-items-center p-menuitem-link" to={item.url}>
            <span className={item.icon}/>
            <span className="mx-2">{item.label}</span>
        </NavLink>
    );
    const items: MenuItem[] = [
        {
            label: 'Home',
            icon: 'pi pi-home',
            url: "/"
        },
        {
            label: 'Products',
            icon: 'pi pi-gift',
            url: "/products"
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
        <div className="flex align-items-center gap-2 pr-4">
            <InputText placeholder="Search" type="text" className="w-8rem sm:w-auto"/>
            <Avatar className="ml-2" image="/images/avatars/amyelsner.png" shape="circle"/>
        </div>
    );

    return (
        <div className="card">
            <Menubar model={items} start={start} end={end}/>
        </div>
    )
}
