import {BreadCrumb} from 'primereact/breadcrumb';
import {MenuItem, MenuItemOptions} from 'primereact/menuitem';
import {NavLink} from "react-router";

interface IAdminBreadcrumbsProps {
    icon: string;
    label: string;
}

export function AdminBreadcrumbs(props: IAdminBreadcrumbsProps) {
    const iconItemTemplate = (item: MenuItem, options: MenuItemOptions) => {
        return (
            <NavLink to={item.url || ''} className={options.className}>
                <span className={`${item.icon} pr-2`}></span> {item.label}
            </NavLink>
        );
    };
    const items: MenuItem[] = [
        {icon: props.icon, template: iconItemTemplate, label: props.label}
    ];
    const home: MenuItem = {icon: 'pi pi-home', url: '/admin', template: iconItemTemplate};

    return (
        <BreadCrumb model={items} home={home}/>
    )
}