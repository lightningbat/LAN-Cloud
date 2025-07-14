async function encryptAESGCM(keyByte, plaintext) {
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
        ciphertext: btoa(String.fromCharCode(...new Uint8Array(encrypted)))
    };
}

async function encryptJSON(obj, keyByte) {
    const { iv, ciphertext } = await encryptAESGCM(keyByte, JSON.stringify(obj));
    return { iv, ciphertext }
}

export { encryptAESGCM, encryptJSON };