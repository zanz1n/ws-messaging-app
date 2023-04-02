import { createContext, useContext, useMemo } from "react";
import { useNavigate } from "react-router-dom";
import jwtDecode from "jwt-decode";
import clientConfig from "../../env-settings.json";

export interface AuthedUserData {
    username: string;
    id: string;
}

export interface LoginProps {
    username: string;
    password: string;
}

export interface RegisterProps extends LoginProps {
    confirmPassword: string;
}

export interface AuthContext {
    get isAuthenticated(): boolean;
    get userData(): AuthedUserData | null;
    get token(): string | null;
    login(props: LoginProps): Promise<void>;
    register(props: RegisterProps): Promise<void>;
    logout: () => void;

}

const Context = createContext({} as AuthContext);

export function AuthProvider({ children }: { children: React.ReactElement | React.ReactElement[]}) {
    const navigate = useNavigate();

    function tokenK() {
        return localStorage.getItem("token");
    }

    function token() {
        const find = tokenK();
        return find ? localStorage.getItem(find) : null;
    }

    function isAuthenticated() {
        const find = token();
        if (!find) return false;
        try {
            const decoded = jwtDecode(find);
            if (decoded && typeof decoded == "object" && "username" in decoded && "id" in decoded) {
                return true;
            }
            return false;
        } catch (e) {
            return false;
        }
    }

    async function login(props: LoginProps) {
        const res = await fetch(`${clientConfig.ApiUri}/auth/signin`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(props)
        }).then(res => res.json()).catch(() => null) as unknown;

        if (res && typeof res == "object" && "token" in res && typeof res["token"] == "string") {
            const uid = crypto.randomUUID();
            localStorage.setItem("token", uid);
            localStorage.setItem(uid, res.token);
            return;
        }
        throw new Error("The username or password do not match.");
    }

    async function register({ username, password, confirmPassword }: RegisterProps) {
        if (password != confirmPassword) throw new Error("The passwords do not match.");

        const res = await fetch(`${clientConfig.ApiUri}/auth/signup`, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify({ username, password })
        }).then(res => res.json()).catch(() => null) as unknown;

        if (res && typeof res == "object" && "token" in res && typeof res["token"] == "string") {
            const uid = crypto.randomUUID();
            localStorage.setItem("token", uid);
            localStorage.setItem(uid, res.token);
            return;
        }

        throw new Error("The username is already taken.");
    }

    function getUserData() {
        const authenticated = isAuthenticated();
        if (!authenticated) return null;
        const find = token();
        if (!find) return null;
        try {
            const decoded = jwtDecode(find);
            if (decoded &&
                typeof decoded == "object" &&
                "username" in decoded &&
                "id" in decoded &&
                typeof decoded["username"] == "string" &&
                typeof decoded["id"] == "string"
            ) {
                return {
                    username: decoded["username"],
                    id: decoded["id"]
                };
            }
            return null;
        } catch (e) {
            return null;
        }
    }

    function logout() {
        const findK = tokenK();
        if (findK) {
            localStorage.removeItem(findK);
            localStorage.removeItem("token");
        }
        navigate("/auth/signin");
    }

    const ctx = useMemo(() => ({
        get isAuthenticated() { return isAuthenticated(); },
        get userData() { return getUserData(); },
        get token() { return token(); },
        login: (props: LoginProps) => login(props),
        logout: () => logout(),
        register: (props: RegisterProps) => register(props),
    } satisfies AuthContext), []);

    return (
        <Context.Provider value={ctx}>
            {children}
        </Context.Provider>
    );
}

export function useAuth() {
    return useContext(Context);
}
