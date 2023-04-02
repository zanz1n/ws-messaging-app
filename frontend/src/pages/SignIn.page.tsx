import styles from "../components/auth/Form.module.css";
import { useState } from "react";
import Form from "../components/auth/Form";
import InputLabel from "../components/auth/InputLabel";
import SubmitButton from "../components/auth/SubmitButton";
import SwitchPages from "../components/auth/SwitchPage";
import Header from "../components/Header";

function validate(target: unknown) {
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
        typeof target["password"]["value"] == "string") {
        return true;
    }
    return false;
}

function handleValueUpdate(e: React.ChangeEvent<HTMLInputElement>, setSendable: React.Dispatch<React.SetStateAction<boolean>>) {
    const username = document.getElementById("username") as HTMLInputElement;
    const password = document.getElementById("password") as HTMLInputElement;
    if (!username.value || username.value == "" || !password.value || password.value == "") {
        setSendable(false);
        return;
    }
    else setSendable(true);
}

export default function SignInPage() {
    const [error, setError] = useState<string | null>(null);

    const [sendable, setSendable] = useState<boolean>(false);

    return (
        <>
            <Header/>
            <main className={styles.main}>
                <div className={styles.formContainer}>
                    <Form error={error} onSubmit={(e) => {
                        e.preventDefault();
                        (async(target: unknown) => {
                            if (validate(target)) {
                                setError("The username or password is invalid.");
                            } else {
                                setError("The username or password is invalid.");
                            }
                        })(e.target);
                    }}>
                        <InputLabel onChange={(e) => handleValueUpdate(e, setSendable)} required identifier="username" type="text">
                        Username
                        </InputLabel>

                        <InputLabel onChange={(e) => handleValueUpdate(e, setSendable)} required identifier="password" type="password">
                        Password
                        </InputLabel>

                        <SubmitButton enabled={sendable} >Log In</SubmitButton>
                        <SwitchPages plain="New here?" to="/auth/signup">Create an account</SwitchPages>
                    </Form>
                </div>
            </main>
        </>
    );
}
