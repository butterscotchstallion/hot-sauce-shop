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
import {OrderCheckoutPage} from "./pages/Checkout/OrderCheckoutPage.tsx";
import {AccountPage} from "./pages/Account/AccountPage.tsx";
import {AccountSignInPage} from "./pages/Account/AccountSignInPage.tsx";
import BoardsListPage from "./pages/Boards/BoardsListPage.tsx";
import PostsListPage from "./pages/Boards/PostsListPage.tsx";
import NewPostPage from "./pages/Boards/NewPostPage.tsx";
import UserProfilePage from "./pages/Users/UserProfilePage.tsx";
import {BoardSettingsPage} from "./pages/Boards/BoardSettingsPage.tsx";

const HomePage = React.lazy(() => import("./pages/HomePage.tsx"));
const ProductListPage = React.lazy(() => import("./pages/Products/ProductListPage.tsx"));
const ProductDetailPage = React.lazy(() => import("./pages/Products/ProductDetailPage.tsx"));

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
                    <Route path="boards/:boardSlug" element={<PostsListPage/>}/>
                    <Route path="posts" element={<PostsListPage/>}/>
                    <Route path="boards/:slug/posts/new" element={<NewPostPage/>}/>
                    <Route path="boards/:boardSlug/posts/:postSlug" element={<PostsListPage/>}/>
                    <Route path="users/:slug" element={<UserProfilePage/>}/>
                    <Route path="ws" element={<HomePage/>}/>
                    <Route path="boards/:boardSlug/settings" element={<BoardSettingsPage/>}/>
                </Route>
                <Route path="*" element={<Navigate to="/"/>}/>
            </Routes>
        </BrowserRouter>
    </Suspense>
)
