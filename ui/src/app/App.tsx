import './App.scss'
import BaseLayout from "./pages/BaseLayout.tsx";
import {Outlet} from "react-router";
import {ReactElement} from "react";

function App(): ReactElement {
    return (
        <BaseLayout>
            <Outlet/>
        </BaseLayout>
    )
}

export default App
