import { Route, createBrowserRouter, createRoutesFromElements, useNavigate } from "react-router-dom";
import App from "./App";
import SignInPage from "./pages/SignIn.page";
import SignUpPage from "./pages/SignUp.page";
import ChatPage from "./pages/Chat.page";
import { useEffect } from "react";
import { ErrorBoundary } from "react-error-boundary";
import { SocketProvider } from "./lib/SocketContext";

function RedirectFallback() {
    const navigate = useNavigate();
    useEffect(() => {
        navigate("/auth/signin");
    }, []);
    return <></>;
}

function ProtectedRoute({ children }: { children: JSX.Element }) {
    return <>
        <ErrorBoundary FallbackComponent={RedirectFallback}>
            {children}
        </ErrorBoundary>
    </>;
}

export const router = createBrowserRouter(createRoutesFromElements(
    <Route path="/" element={<App/>}>
        <Route path="/" element={
            <SocketProvider>
                <ChatPage/>
            </SocketProvider>
        }/>
        <Route path="auth/signin" element={<SignInPage/>}/>
        <Route path="auth/signup" element={<SignUpPage/>}/>
    </Route>
));
