import * as React from "react";
import {Button} from "primereact/button";
import {Menu} from "primereact/menu";

export default function CartMenu() {
    const cartMenu = React.useRef<Menu>(null);
    const items = [
        {
            label: 'Cart',
            items: [
                {
                    label: 'Refresh',
                    icon: 'pi pi-refresh'
                },
                {
                    label: 'Export',
                    icon: 'pi pi-upload'
                }
            ]
        }
    ];
    return (
        <>
            <Menu model={items} ref={cartMenu} popup id="popup_menu_right" popupAlignment="right"/>
            <Button
                label="Cart"
                icon="pi pi-shopping-cart"
                className="mr-2"
                onClick={(event) => cartMenu?.current.toggle(event)}
                aria-controls="popup_menu_right"
                aria-haspopup/>
        </>
    )
}