import { Route, createBrowserRouter, createRoutesFromElements } from "react-router-dom";
import App from "./App";
import SignInPage from "./pages/SignIn.page";
import SignUpPage from "./pages/SignUp.page";
import ChatPage from "./pages/Chat.page";
import { SocketProvider } from "./lib/SocketContext";

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
