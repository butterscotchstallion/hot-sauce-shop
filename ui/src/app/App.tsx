import './App.scss'
import BaseLayout from "./pages/BaseLayout.tsx";
import {Outlet} from "react-router";
import {ReactElement} from "react";
import {Provider} from "react-redux";
import {store} from "./store.ts";
import {AuthProvider} from "react-oidc-context";

function App(): ReactElement {
    const oidcConfig = {
        authority: "http://localhost:8080/realms/master",
        client_id: "netherrealm-realm",
        redirect_uri: "http://localhost:5173/products",
    };
    return (
        <AuthProvider {...oidcConfig}>
            <Provider store={store}>
                <BaseLayout>
                    <Outlet/>
                </BaseLayout>
            </Provider>
        </AuthProvider>
    )
}

export default App
