import { useNavigate } from "react-router-dom";
import styles from "./Header.module.css";

export default function Header() {
    const isAuthenticated = true;
    const user = { username: "zanz1n", uuid: "22228f55-75e4-4466-ae5a-4adea203e2e3" };

    const navigate = useNavigate();

    return (
        <div className={styles.header}>
            <header className={styles.headerContainer}>
                <div className={styles.left}>
                    <h1 className={styles.title}>Ws App</h1>
                </div>
                <div className={styles.right}>
                    {
                        isAuthenticated ? (
                            <>
                                <button onClick={() => {
                                    localStorage.removeItem("token");
                                    navigate("/auth/signin");
                                }} className={styles.purple}>Logout</button>
                                <div>
                                    <p>Logged in as</p>
                                    <p>{user.username}</p>
                                </div>
                            </>
                        ) : (
                            <>
                                <button onClick={() => navigate("/auth/signin")} className={styles.green}>Sign In</button>
                                <button onClick={() => navigate("/auth/signup")} className={styles.pink}>Sign Up</button>
                            </>
                        )
                    }
                </div>
            </header>
        </div>
    );
}
