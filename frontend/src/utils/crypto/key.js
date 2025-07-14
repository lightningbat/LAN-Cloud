import { decryptAESGCM } from "./decrypt";

let Key = null;
let SessionId = null;

export function setKey(sessionId, key) {
    Key = key;
    SessionId = sessionId;

    // auto refresh session key
    setTimeout(refreshKey, 15 * 60 * 1000);
}

export function getKey() {
    return { Key, SessionId };
}

async function refreshKey() {
    try {
        const res = await fetch(`${import.meta.env.VITE_SERVER_URL}/refreshSessionKey`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ session_id: SessionId })
        });

        if (res.status !== 200) {
            alert("Failed to refresh session key");
            return;
        }

        const data = await res.json();

        // decrypts data
        let newSessionKey;
        try {
            newSessionKey = await decryptAESGCM(Key, data.iv_base64, data.ciphertext_base64);
        } catch (err) {
            alert("Failed to refresh session key (decrypt)");
            console.error(err);
            return;
        }

        // convert from base64 to buffer
        newSessionKey = Uint8Array.from(atob(newSessionKey), c => c.charCodeAt(0));

        setKey(data.session_id, newSessionKey);
    }
    catch (err) {
        alert("Failed to refresh session key (catch)");
        console.error(err);
    }
}