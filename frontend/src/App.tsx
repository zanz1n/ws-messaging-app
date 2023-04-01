import { Outlet } from "react-router-dom";
import { AuthProvider } from "./lib/AuthContext";

export default function App() {
    return (
        <Outlet/>
    );
}
