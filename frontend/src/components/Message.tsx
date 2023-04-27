import styles from "./Message.module.css";

export interface MessageProps {
    self: boolean;
    image: string | null;
    content: string | null;
    userId: string;
    username: string;
    timeFmt: string;
}

export default function Message({ content, image, userId, username, timeFmt, self }: MessageProps) {
    return <>
        <div className={`${styles.message} ${self ? styles.self : styles.other}`}>
            <div className={styles.body}>
                { content ? <p>{content}</p> : undefined }
                { image ? <img src={image} alt={image} /> : undefined }
            </div>
            <div className={styles.metadata}>
                <span id={userId}>{username} • {timeFmt}</span>
            </div>
        </div>
    </>;
}
