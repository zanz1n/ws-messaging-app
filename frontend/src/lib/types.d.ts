export type WithChildren<T = unknown> = {
    children?: React.ReactElement | React.ReactElement[] | string
} & T

export interface BaseMessage {
    id: string;
    createdAt: number;
    updatedAt: number;
    user: {
        id: string;
        username: string;
    }
    image: string | null;
    content: string | null;
}
