import './App.scss'
import BaseLayout from "./pages/BaseLayout.tsx";
import {Outlet} from "react-router";
import {ReactElement} from "react";
import {Provider} from "react-redux";
import {store} from "./store.ts";

function App(): ReactElement {
    return (
        <Provider store={store}>
            <BaseLayout>
                <Outlet/>
            </BaseLayout>
        </Provider>
    )
}

export default App
