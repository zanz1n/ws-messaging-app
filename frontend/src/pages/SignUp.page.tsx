import { useState } from "react";
import Header from "../components/Header";
import Form from "../components/auth/Form";
import styles from "../components/auth/Form.module.css";
import InputLabel from "../components/auth/InputLabel";
import SubmitButton from "../components/auth/SubmitButton";
import SwitchPages from "../components/auth/SwitchPage";
import { useAuth } from "../lib/AuthContext";
import { useNavigate } from "react-router-dom";

interface RegisterDomData {
    username: {
        value: string;
    };
    password: {
        value: string;
    };
    confirmPassword: {
        value: string;
    };
}

function validate(target: unknown): target is RegisterDomData {
    if (target &&
        typeof target == "object" &&
        "username" in target &&
        target["username"] &&
        typeof target["username"] == "object" &&
        "value" in target["username"] &&
        target["username"]["value"] &&
        typeof target["username"]["value"] == "string" &&
        "password" in target &&
        target["password"] &&
        typeof target["password"] == "object" &&
        "value" in target["password"] &&
        target["password"]["value"] &&
        typeof target["password"]["value"] == "string" &&
        "confirmPassword" in target &&
        target["confirmPassword"] &&
        typeof target["confirmPassword"] == "object" &&
        "value" in target["confirmPassword"] &&
        target["confirmPassword"]["value"] &&
        typeof target["confirmPassword"]["value"] == "string"
    ) {
        return true;
    } else {
        return false;
    }
}

export default function SignUpPage() {
    const [password , setPassword] = useState<string>("");

    const [confirmPassword , setConfirmPassword] = useState<string>("");

    const [error, setErrorRaw] = useState<string | null>(null);
    
    const [sendable, setSendable] = useState<boolean>(false);

    const { isAuthenticated } = useAuth();

    const navigate = useNavigate();

    if (isAuthenticated) {
        navigate("/");
    }
    
    function setError(e: string | null) {
        setErrorRaw(e);
        if (e == null) setSendable(true);
        else setSendable(false);
    }

    // const navigate = useNavigate();

    function handlePasswordUpdate(t: "password" | "confirmPassword") {
        return function(e: { target: { value: string; }; }) {
            const value = e.target.value;
            if (t == "password") {
                if (confirmPassword == "" || confirmPassword != value) {
                    setError("The passwords do not match.");
                } else {
                    setError(null);
                }
                setPassword(value);
            } else if (t == "confirmPassword") {
                if (password == "" || password != value) {
                    setError("The passwords do not match.");
                } else {
                    setError(null);
                }
                setConfirmPassword(value);
            }
        };
    }
    
    return <>
        <Header/>
        <main className={styles.main}>
            <div className={styles.formContainer}>
                <div className={styles.formTitle}>
                    <h1>Sign Up</h1>
                </div>
                <Form error={error}
                    onSubmit={(e) => {
                        e.preventDefault();
                        (async(target: unknown) => {
                            if (validate(target)) {
                                if (target["password"]["value"] != target["confirmPassword"]["value"]) {
                                    setError("The passwords do not match.");
                                    return;
                                }
                                if (!target["username"]["value"] || target["username"]["value"] == "") {
                                    setError("Please enter a username.");
                                    return;
                                }
                                // const result = await register({
                                //     username: target["username"]["value"],
                                //     password: target["password"]["value"],
                                //     confirmPassword: target["confirmPassword"]["value"]
                                // });

                                // if (result) {
                                //     setError(null);
                                //     navigate("/");
                                //     return;
                                // } else {
                                setError("An error occurred while creating your account.");
                            // }
                            }
                        })(e.target);
                    }}>

                    <InputLabel required identifier="username" type="text"
                        onChange={(e) => {
                            if (!e.target.value || e.target.value == "") {
                                setError("Please enter a username.");
                                return;
                            } else {
                                setError(null);
                            }
                        }}>
                    Username
                    </InputLabel>

                    <InputLabel required identifier="password" type="password"
                        onChange={handlePasswordUpdate("password")}>
                        Password
                    </InputLabel>

                    <InputLabel required identifier="confirmPassword" type="password"
                        onChange={handlePasswordUpdate("confirmPassword")}>
                        Confirm Password
                    </InputLabel>

                    <SubmitButton enabled={sendable} >Create Account</SubmitButton>

                    <SwitchPages plain="Already have an account?" to="/auth/signin">Login</SwitchPages>

                </Form>
            </div>
        </main>
    </>;
}
