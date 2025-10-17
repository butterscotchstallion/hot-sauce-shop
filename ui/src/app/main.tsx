import * as React from 'react'
import {Suspense} from 'react'
import {createRoot} from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import {BrowserRouter, Navigate, Route, Routes} from "react-router";
import Throbber from "./components/Shared/Throbber.tsx";
import AdminPage from "./pages/Admin/AdminPage.tsx";
import AdminInventoryPage from "./pages/Admin/AdminInventoryPage.tsx";
import {AdminUserDetailPage} from "./pages/Admin/AdminUserDetailPage.tsx";
import {OrderCheckoutPage} from "./pages/OrderCheckoutPage.tsx";
import {AccountPage} from "./pages/Account/AccountPage.tsx";
import {AccountSignInPage} from "./pages/Account/AccountSignInPage.tsx";
import BoardsListPage from "./pages/Boards/BoardsListPage.tsx";
import BoardPostListPage from "./pages/Boards/BoardPostListPage.tsx";
import PostsListPage from "./pages/Boards/PostsListPage.tsx";

const HomePage = React.lazy(() => import("./pages/HomePage.tsx"));
const ProductListPage = React.lazy(() => import("./pages/ProductListPage.tsx"));
const ProductDetailPage = React.lazy(() => import("./pages/ProductDetailPage.tsx"));

createRoot(document.getElementById('root')!).render(
    <Suspense fallback={<Throbber/>}>
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<App/>}>
                    <Route path="" element={<HomePage/>}/>
                    <Route path="products" element={<ProductListPage/>}/>
                    <Route path="products/:slug" element={<ProductDetailPage/>}/>
                    <Route path="admin" element={<AdminPage/>}/>
                    <Route path="admin/products/edit/:slug" element={<AdminInventoryPage isNewProduct={false}/>}/>
                    <Route path="admin/products/add" element={<AdminInventoryPage isNewProduct={true}/>}/>
                    <Route path="admin/users/edit/:slug" element={<AdminUserDetailPage isNewUser={false}/>}/>
                    <Route path="admin/users/add" element={<AdminUserDetailPage isNewUser={true}/>}/>
                    <Route path="orders/checkout" element={<OrderCheckoutPage/>}/>
                    <Route path="account" element={<AccountPage/>}/>
                    <Route path="account/sign-in" element={<AccountSignInPage/>}/>
                    <Route path="boards" element={<BoardsListPage/>}/>
                    <Route path="boards/:slug" element={<BoardPostListPage/>}/>
                    <Route path="posts" element={<PostsListPage/>}/>
                </Route>
                <Route path="*" element={<Navigate to="/"/>}/>
            </Routes>
        </BrowserRouter>
    </Suspense>
)
