import { createContext, useContext, useMemo } from "react";
import { WithChildren } from "./types";
import { useAuth } from "./AuthContext";
import clientConfig from "../../env-settings.json";

export interface IncomingChatMessage {
    type: "messageCreated";
    id: string
    createdAt: number
    updatedAt: number
    user: {
        id: string;
        username: string;
    }
    image: string | null
    content: string | null
}

export interface IncomingError {
    error: string
}

export interface SocketContext {
    ws: WebSocket;
    onMessage: (callback: (message: IncomingChatMessage) => void | Promise<void>) => void;
    onError: (callback: (err: IncomingError) => void | Promise<void>) => void;
    onClose: (callback: () => void | Promise<void>) => void;
    close: () => void;
}

const Context = createContext({} as SocketContext);

export function SocketProvider({ children }: WithChildren) {
    const { token } = useAuth();

    let errrorCallback = (_: IncomingError): void | Promise<void> => {/**/};
    let messageCallback = (_: IncomingChatMessage): void | Promise<void> => {/**/};
    let closeCallback = (): void | Promise<void> => {/**/};

    const ws = new WebSocket(`${clientConfig.WsEndpoint}?auth_token=${token}`);

    ws.onmessage = (e) => {
        const data = JSON.parse(e.data);

        if ("error" in data) {
            errrorCallback(data);
        } else if (data["type"] == "messageCreated") {
            messageCallback(data);
        }
    };

    ws.onclose = () => {
        closeCallback();
    };

    setInterval(() => {
        console.log("Ping sent to ws");
        ws.send("{\"type\":\"ping\"}");
    }, 30 * 1000);

    const ctx = useMemo(() => ({
        ws,
        close() {
            ws.close();
        },
        onClose(callback) { closeCallback = callback; },
        onError(callback) { errrorCallback = callback; },
        onMessage(callback) { messageCallback = callback; },
    } satisfies SocketContext), []);

    return <Context.Provider value={ctx}>
        {children}
    </Context.Provider>;
}

export function useSocket() {
    return useContext(Context);
}
