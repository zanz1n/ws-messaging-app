import { createContext, useContext, useMemo } from "react";
import clientConfig from "../../env-settings.json";
import { useAuth } from "./AuthContext";

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

export type ListennerCallback = (this: WebSocket, e: MessageEvent<IncomingChatMessage>) => unknown | Promise<unknown>

export interface SocketContext {
    socket: WebSocket;
    onMessage(callback: ListennerCallback): void;
    sendMessage(message: ChatMessagePayload): void;
}

const Context = createContext({} as SocketContext);

export function SocketProvider({ children }: { children: React.ReactElement | React.ReactElement[] }) {
    const { isAuthenticated, token } = useAuth();
    if (!isAuthenticated || !token) throw new Error("Chat gatweway connection requires authentication");

    const socket = new WebSocket(clientConfig.WsEndpoint, token);

    function sendMessage(message: ChatMessagePayload) {
        socket.send(JSON.stringify(message));
    }

    function onMessage(callback: ListennerCallback) {
        socket.onmessage = callback;
    }

    const value = useMemo(() => ({
        onMessage,
        sendMessage,
        socket
    } satisfies SocketContext), []);

    return <Context.Provider value={value}>
        {children}
    </Context.Provider>;
}

export function useSocket() {
    return useContext(Context);
}
