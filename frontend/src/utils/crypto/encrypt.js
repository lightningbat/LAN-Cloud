import { AES } from "@stablelib/aes";
import { GCM } from "@stablelib/gcm";
import { randomBytes } from "@stablelib/random";

function encodeUint8ArrayToBase64(uint8Array) {
    let binary = '';
    const len = uint8Array.byteLength;
    const chunkSize = 16384;

    for (let i = 0; i < len; i += chunkSize) {
        const chunk = uint8Array.subarray(i, Math.min(i + chunkSize, len));
        binary += String.fromCharCode(...chunk);
    }
    return btoa(binary);
}

async function encryptAESGCM(keyByte, plaintext) {
    
    if (window.crypto?.subtle?.importKey) {
        const iv = new Uint8Array(16); // 128-bit IV for GCM mode
        const cryptoKey = await crypto.subtle.importKey(
            "raw",
            keyByte,
            { name: "AES-GCM" },
            false,
            ["encrypt"]
        );
        const encrypted = await crypto.subtle.encrypt(
            {
                name: "AES-GCM",
                iv: iv
            },
            cryptoKey,
            new TextEncoder().encode(plaintext)
        );
        return {
            iv: btoa(String.fromCharCode(...iv)),
            ciphertext: encodeUint8ArrayToBase64(new Uint8Array(encrypted))
        };
    } else {
        const iv = randomBytes(12);
        const aes = new AES(keyByte);
        const gcm = new GCM(aes);
        const encrypted = gcm.seal(iv, new TextEncoder().encode(plaintext));
        return {
            iv: btoa(String.fromCharCode(...iv)),
            ciphertext: encodeUint8ArrayToBase64(encrypted)
        };
    }
}

async function encryptJSON(keyByte, obj) {
    const { iv, ciphertext } = await encryptAESGCM(keyByte, JSON.stringify(obj));
    return { iv, ciphertext }
}

export { encryptAESGCM, encryptJSON };