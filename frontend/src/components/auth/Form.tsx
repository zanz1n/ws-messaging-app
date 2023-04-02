import styles from "./Form.module.css";

export interface FormProps {
    type: "login" | "register";
    error?: string | null;
    onSubmit?: (e: React.FormEvent<HTMLFormElement>) => void | Promise<void>;
    children: React.ReactElement | React.ReactElement[];
}

export default function Form({ type, error, onSubmit, children }: FormProps) {
    return(
        <form className={type} onSubmit={(e) => {
            e.preventDefault();
            onSubmit?.(e);
        }}>
            <div className={`${styles.topError} ${error ? "" : styles.invisible}`}>
                <p>{error ?? "-"}</p>
            </div>
            {children}
            <div className={`${styles.topError} ${styles.invisible}`}>
                <p>-</p>
            </div>
        </form>
    );
}
