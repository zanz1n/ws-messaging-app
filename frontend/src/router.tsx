import { createBrowserRouter } from "react-router-dom";
import App from "./App";
import SignInPage from "./pages/SignIn.page";
import SignUpPage from "./pages/SignUp.page";

function ProtectedRoute({ children }: { children: any }) {
    return children;
}

export const router = createBrowserRouter([
    {
        path: "/",
        element: <App/>,
        children: [
            {
                path: "/",
                element: <ProtectedRoute>{""}</ProtectedRoute>
            }
        ]
    },
    {
        path: "/auth/signin",
        element: <App/>,
        children: [
            {
                path: "/auth/signin",
                element: <SignInPage/>
            }
        ]
    },
    {
        path: "/auth/signup",
        element: <App/>,
        children: [
            {
                path: "/auth/signup",
                element: <SignUpPage/>
            }
        ]
    }
]);
