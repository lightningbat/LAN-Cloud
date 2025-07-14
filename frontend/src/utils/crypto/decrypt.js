import { AES } from "@stablelib/aes";
import { GCM } from "@stablelib/gcm";

/**
 * 
 * @param {Uint8Array} keyByte 
 * @param {String} base64Iv 
 * @param {String} base64Ciphertext 
 * @returns {Promise<String>}
 */
async function decryptAESGCM(keyByte, base64Iv, base64Ciphertext) {

    const iv = Uint8Array.from(atob(base64Iv), c => c.charCodeAt(0));
    const ciphertext = Uint8Array.from(atob(base64Ciphertext), c => c.charCodeAt(0));

    if (window.crypto?.subtle?.importKey) {

        const iv = Uint8Array.from(atob(base64Iv), c => c.charCodeAt(0));
        const ciphertext = Uint8Array.from(atob(base64Ciphertext), c => c.charCodeAt(0));

        const cryptoKey = await crypto.subtle.importKey(
            "raw",
            keyByte,
            { name: "AES-GCM" },
            false,
            ["decrypt"]
        );


        const decrypted = await crypto.subtle.decrypt(
            {
                name: "AES-GCM",
                iv: iv
            },
            cryptoKey,
            ciphertext
        );

        return new TextDecoder().decode(decrypted); // plaintext string

    } else {
        const aes = new AES(keyByte);
        const gcm = new GCM(aes);
        const decrypted = gcm.open(iv, ciphertext);

        return new TextDecoder().decode(decrypted);
    }
}

async function decryptJSON(keyByte, base64Iv, base64Ciphertext) {
    const plaintext = await decryptAESGCM(keyByte, base64Iv, base64Ciphertext);
    return JSON.parse(plaintext);
}

export { decryptAESGCM, decryptJSON };