import { useEffect, useState } from "react";
import Header from "../components/Header";
import { useAuth } from "../lib/AuthContext";
import { useNavigate } from "react-router-dom";
import styles from "./Chat.module.css";
import { useSocket } from "../lib/SocketContext";
import clientConfig from "../../env-settings.json";
import { BaseMessage } from "../lib/types";
import Message from "../components/Message";

export interface ChatMessagePayload {
    content: string | null;
    image: string | null;
}

export default function ChatPage() {
    const { isAuthenticated, token, userData, logout } = useAuth();
    const { onMessage, close } = useSocket();

    const navigate = useNavigate();

    const [messages, setMessages] = useState<BaseMessage[]>([]);

    useEffect(() => {
        if (!isAuthenticated || !token) {
            close();
            navigate("/auth/signin");
            return;
        }

        fetch(`${clientConfig.ApiUri}/messages?t=16816144697510&l=100`, {
            headers: {
                authorization: token
            }
        }).then(res => {
            if (!res.ok) throw new Error("failed to fetch messages");

            return res.json();
        }).then((data: unknown) => {
            if (data && typeof data == "object" && "data" in data && typeof data["data"] == "object") {
                const msgs = data["data"] as BaseMessage[];

                setMessages(msgs);
                return;
            }
            throw new Error("received incomplete or corrupted data from the server");
        }).catch((e => console.error(e)));
    }, []);

    onMessage((message: BaseMessage) => {
        if ("type" in message) {
            delete message["type"];
        }
        setMessages([ ...messages, message]);
    });

    const [_, setErr] = useState(null as string | null);

    function handleMessageSubmit(e: React.FormEvent<HTMLFormElement>) {
        e.preventDefault();

        if (!token) {
            return logout();
        }

        const text = document.getElementById("content") as HTMLInputElement;
    
        if (!text || !text.value || text.value == "") {
            setErr("the message content must not be empty!");
            return;
        }

        fetch(clientConfig.ApiUri + "/messages", {
            method: "POST",
            headers: {
                "Authorization": token,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                content: text.value
            })
        }).then(res => {
            if (!res.ok) throw new Error();
            return res.json();
        }).then(() => {
            text.value = "";
        }).catch(() => {
            setErr("something went wrong sending your message");
        });
    }
    
    return <>
        <Header/>
        <main className={styles.main}>
            <div className={styles.chatContainer}>
                <div className={styles.messagesContainer}>
                    <div className={styles.messages}>
                        {messages.map((m) => {
                            const date = new Date(m.createdAt);
                            const timeString = `${date.getHours()}:${date.getMinutes()}:` + 
                                (date.getSeconds() > 10 ? date.getSeconds() : `0${date.getSeconds()}`);
                            return <Message
                                content={m.content} 
                                image={m.image}
                                self={userData?.id == m.user.id}
                                timeFmt={timeString}
                                userId={m.user.id}
                                username={m.user.username}
                                key={m.id}
                            />;
                        })}
                    </div>
                </div>
                <form className={styles.form} onSubmit={handleMessageSubmit}>
                    <div className={styles.formInput}>
                        <input type="text" id="content" name="content" />
                    </div>
                    <button type="submit">Enviar</button>
                </form>
            </div>
        </main>
    </>;
}
