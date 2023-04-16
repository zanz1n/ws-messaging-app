import { useEffect, useState } from "react";
import Header from "../components/Header";
import { useAuth } from "../lib/AuthContext";
import { useNavigate } from "react-router-dom";
import clientConfig from "../../env-settings.json";
import styles from "./Chat.module.css";
import { IncomingChatMessage, useSocket } from "../lib/SocketContext";

export interface ChatMessagePayload {
    content: string | null;
    image: string | null;
}

export default function ChatPage() {
    const { isAuthenticated, token } = useAuth();
    const { onMessage, close } = useSocket();

    const navigate = useNavigate();

    const [messages, setMessages] = useState<IncomingChatMessage[]>([]);

    useEffect(() => {
        if (!isAuthenticated || !token) {
            close();
            navigate("/auth/signin");
            return;
        }
    }, []);

    onMessage((message) => {
        console.log(message);
        setMessages([ ...messages, message]);
    });
    
    return <>
        {console.log(messages)}
        <Header/>
        <main className={styles.main}>
            {messages.map((m) => {
                return JSON.stringify(m);
            })}
        </main>
    </>;
}
