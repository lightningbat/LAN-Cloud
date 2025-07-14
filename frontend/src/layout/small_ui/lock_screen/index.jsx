import "./style.scss"
import { EyeOpenIcon, EyeCloseIcon } from "../../../icons"
import { useState } from "react";
import { scrypt } from "scrypt-js";
import { decryptAESGCM, setKey } from "../../../utils/crypto"
import { Spinner } from "../../../loading_animations";
import { HMAC } from "@stablelib/hmac";
import { SHA256 } from "@stablelib/sha256";


export default function LockScreen({ setAuthenticated }) {
    const [inputType, setInputType] = useState("password");
    const [inputValue, setInputValue] = useState("");
    const [processing, setProcessing] = useState(false)
    const [errMsg, setErrMsg] = useState("")

    async function handleSubmit() {
        event.preventDefault();
        setProcessing(true)
        setErrMsg("")
        await connect();
    }

    async function connect() {
        // random string to verify authenticity of server responsex
        const randomString = Math.random().toString(36).substring(2, 15);

        const res = await fetch(`${import.meta.env.VITE_SERVER_URL}/handshake/connect`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ client_data: randomString })
        });

        if (res.status !== 200) {
            setProcessing(false)
            setErrMsg(res.statusText);
            return;
        }

        const data = await res.json();
        const { session_id, salt, scrypt_params, iv_base64, ciphertext_base64 } = data;

        const inputBuffer = new TextEncoder().encode(inputValue);
        const saltBuffer = new TextEncoder().encode(salt);

        const key = await scrypt(
            inputBuffer,
            saltBuffer,
            scrypt_params.iterations,
            scrypt_params.block_size,
            scrypt_params.parallelism,
            scrypt_params.hash_len
        );

        let decrypted;
        try {
            decrypted = await decryptAESGCM(key, iv_base64, ciphertext_base64);
        } catch {
            setProcessing(false)
            setErrMsg("Invalid password or server error (decryption)");
            return;
        }


        if (!decrypted.includes(randomString)) {
            setProcessing(false)
            setErrMsg("Invalid password or server error (missing random string)");
            return;
        }

        const nonceStr = decrypted.replace(randomString, "");

        let responseHash;
        if (window.crypto?.subtle?.importKey) {
            const hmacKey = await window.crypto.subtle.importKey(
                "raw",
                key,
                { name: "HMAC", hash: "SHA-256" },
                false,
                ["sign", "verify"]
            );

            responseHash = await window.crypto.subtle.sign(
                "HMAC",
                hmacKey,
                new TextEncoder().encode(nonceStr)
            );

        } else {
            const hmac = new HMAC(SHA256, key);
            hmac.update(new TextEncoder().encode(nonceStr));
            responseHash = hmac.digest();
        }

        const responseHashBase64 = btoa(String.fromCharCode(...new Uint8Array(responseHash)));
        authenticate(session_id, responseHashBase64, key);
    }

    async function authenticate(session_id, responseHashBase64, key) {
        const res1 = await fetch(`${import.meta.env.VITE_SERVER_URL}/handshake/authenticate`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ session_id, auth_key: responseHashBase64 })
        });

        if (res1.status != 200) {
            setProcessing(false)
            setErrMsg("Invalid password or server error (authentication)");
        }

        const data1 = await res1.json();

        const { session_id: new_session_id, iv_base64: new_iv_base64, ciphertext_base64: new_ciphertext_base64 } = data1;

        // decrypt the data
        let decrypted;
        try {
            decrypted = await decryptAESGCM(key, new_iv_base64, new_ciphertext_base64);
        } catch {
            setProcessing(false)
            setErrMsg("Invalid password or server error (decryption)");
            return;
        }

        // convert decrypted from base64 to buffer
        const newSessionKeyBuffer = Uint8Array.from(atob(decrypted), c => c.charCodeAt(0));

        setKey(new_session_id, newSessionKeyBuffer);

        setProcessing(false)
        setAuthenticated(true);
    }

    return (
        <div className="lock-screen-cont">
            <div className="lock-screen-window">
                <h2 className="title">Lock Screen</h2>
                <form onSubmit={handleSubmit}>
                    <div className="input">
                        <input type={inputType} value={inputValue} onChange={(e) => setInputValue(e.target.value)} placeholder="Enter server password" required />
                        <div className="input-icon" onClick={() => setInputType(inputType === "password" ? "text" : "password")}>
                            {inputType === "password" ?
                                <EyeOpenIcon style={{ width: "100%", height: "100%" }} /> :
                                <EyeCloseIcon style={{ width: "100%", height: "100%" }} />}
                        </div>
                    </div>
                    <p className="error-msg">{errMsg}</p>
                    <div className="cont">
                        {!processing && <button className="unlock-btn" type="submit">Unlock</button>}
                        {processing && <Spinner scale={0.5} />}
                    </div>
                </form>
            </div>
        </div>
    )
}