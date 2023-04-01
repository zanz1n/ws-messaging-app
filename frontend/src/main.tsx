import { StrictMode } from "react";
import ReactDOM from "react-dom/client";
import { RouterProvider } from "react-router-dom";
import { router } from "./router";
import "./index.css";

ReactDOM.createRoot(document.body).render(
    <StrictMode>
        <RouterProvider router={router} />
    </StrictMode>,
);
