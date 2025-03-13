import './App.scss'
import BaseLayout from "./pages/BaseLayout.tsx";
import {Outlet} from "react-router";
import {ReactElement} from "react";
import {ToastContextProvider} from "./components/Shared/ToastContext.tsx";

function App(): ReactElement {
    return (
        <ToastContextProvider>
            <BaseLayout>
                <Outlet/>
            </BaseLayout>
        </ToastContextProvider>
    )
}

export default App
