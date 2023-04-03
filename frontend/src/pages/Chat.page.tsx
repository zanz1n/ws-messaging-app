import { useEffect, useState } from "react";
import Header from "../components/Header";
import { useAuth } from "../lib/AuthContext";
import { useNavigate } from "react-router-dom";
import clientConfig from "../../env-settings.json";

export interface IncomingChatMessage {
    id: string;
    content: string | null;
    image: string | null;
    author: {
        id: string;
        username: string;
    };
}

export interface ChatMessagePayload {
    content: string | null;
    image: string | null;
}

export default function ChatPage() {
    const { isAuthenticated, token } = useAuth();

    const navigate = useNavigate();

    const [messages, setMessages] = useState([] as IncomingChatMessage[]);

    useEffect(() => {
        if (!isAuthenticated || !token) {
            navigate("/auth/signin");
            return;
        }
    }, []);
    
    return <>
        <script>
        window.alert("A")
        </script>
        <Header/>
    </>;
}
