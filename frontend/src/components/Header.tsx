import { Link, useNavigate } from "react-router-dom";
import styles from "./Header.module.css";
import { useAuth } from "../lib/AuthContext";

export default function Header() {
    const { isAuthenticated, userData: user, logout } = useAuth();

    const navigate = useNavigate();

    return (
        <div className={styles.header}>
            <header className={styles.headerContainer}>
                <div className={styles.left}>
                    <Link to="/">
                        <h1 className={styles.title}>Ws App</h1>
                    </Link>
                </div>
                <div className={styles.right}>
                    {
                        isAuthenticated ? (
                            <>
                                <button onClick={() => {
                                    logout();
                                }} className={styles.purple}>Logout</button>
                                <div>
                                    <p>Logged in as</p>
                                    <p>{user?.username}</p>
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
