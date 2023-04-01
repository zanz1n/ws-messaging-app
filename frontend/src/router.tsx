import { Route, createBrowserRouter, createRoutesFromElements } from "react-router-dom";
import App from "./App";
import SignInPage from "./pages/SignIn.page";
import SignUpPage from "./pages/SignUp.page";
import ChatPage from "./pages/Chat.page";

function ProtectedRoute({ children }: { children: JSX.Element }) {
    return children;
}

function ErrorBoundary() {
    return <div>{"Erro"}</div>;
}

export const router = createBrowserRouter(createRoutesFromElements(
    <Route ErrorBoundary={ErrorBoundary} path="/" element={<App/>}>
        <Route path="/" element={<ProtectedRoute><ChatPage/></ProtectedRoute>}/>
        <Route path="auth/signin" element={<SignInPage/>}/>
        <Route path="auth/signup" element={<SignUpPage/>}/>
    </Route>
));
