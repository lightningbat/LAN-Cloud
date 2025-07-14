import { useEffect, useState } from "react";
import ExplorerContext from "./explorerContext";
import { useUiContext } from "../ui_context";
import { getKey, decryptAESGCM } from "../../utils/crypto";
import { useImmer } from "use-immer";

const ExplorerProvider = ({ children }) => {
    const [loading, setLoading] = useState(null);
    const [foldersData, updateFoldersData] = useImmer({});
    const [filesData, updateFilesData] = useImmer({});
    const [selectedFolderId, setSelectedFolderId] = useState(null);
    const [tagsInfo, updateTagsInfo] = useImmer({
        SystemTags: {
            images: 0,
            videos: 0,
            audios: 0,
            documents: 0
        },
        UserTags: {
            /*
            "id": {
                name: "Tag 1",
                color: "#ff234f",
                count: 0
            } */
        }
    })

    // useEffect(() => {
    //     console.log("selectedFolderId", selectedFolderId);
    // }, [selectedFolderId])

    
    // useEffect(() => {
    //     console.log("foldersData", foldersData);
    //     console.log("filesData", filesData);
    // }, [foldersData, filesData])
    
    // useEffect(() => {
    //     console.table({ loading, selectedFolderId });
    // }, [loading]);
    
    const { authenticated } = useUiContext();

    function _fetch(route, body) {
        return fetch(`${import.meta.env.VITE_SERVER_URL}/${route}`, {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(body)
        });
    }

    useEffect(() => {
        (async () => {
            if (authenticated) {

                await loadFolder("root", true);
                setLoading(null);

                const { Key, SessionId } = await getKey();

                const tags_resp = await _fetch("getTags", { session_id: SessionId });

                if (tags_resp.status !== 200) {
                    alert("Failed to get tags");
                    return;
                }

                const { iv_base64: tags_iv, ciphertext_base64: tags_ciphertext } = await tags_resp.json();

                try {
                    const decryptedData = await decryptAESGCM(Key, tags_iv, tags_ciphertext);
                    const data = JSON.parse(decryptedData);
                    updateTagsInfo(draft => {
                        draft.SystemTags.images = data.image_count;
                        draft.SystemTags.videos = data.video_count;
                        draft.SystemTags.audios = data.audio_count;
                        draft.SystemTags.documents = data.document_count;
                        draft.UserTags = data.custom_tags || {};
                    })

                } catch (err) {
                    alert("Failed to decrypt tags");
                    console.error(err);
                }
            }
        })();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [authenticated]);

    useEffect(() => {
        if (selectedFolderId) {
            const subfolder_ids = Object.keys(foldersData[selectedFolderId]?.sub_folders || {});
            const file_ids = Object.keys(foldersData[selectedFolderId]?.files || {});

            const missing_folders_count = subfolder_ids.filter(folder_id => !foldersData[folder_id]).length;
            const missing_files_count = file_ids.filter(file_id => !filesData[file_id]).length;

            if (missing_folders_count > 0 || missing_files_count > 0) {
                (async () => {
                    // load current folder
                    await loadFolder(selectedFolderId);
                    setLoading(null);

                    // climb up till root folder and load all missing parent folders
                    let parent_folder_id = foldersData[selectedFolderId].parent_id;
                    async function recurse(parent_folder_id) {
                        // return if parent_folder_id is root
                        if (parent_folder_id === "root") return;
                        // recurse if parent folder is already loaded and also it's not root
                        // use case scenario: root: loaded -> folder1: not loaded -> folder2(current folder in check): loaded
                        if (foldersData[parent_folder_id] && parent_folder_id !== "root") {
                            return await recurse(foldersData[parent_folder_id].parent_id);
                        }
                        await loadFolder(parent_folder_id); // load parent folder
                        return await recurse(foldersData[parent_folder_id].parent_id);
                    }
                    await recurse(parent_folder_id);
                    setLoading(null);
                })();
            }
        }
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedFolderId])

    async function loadFolder(folder_id, setActiveFolder = false) {
        setLoading(folder_id);

        const { Key, SessionId } = getKey();

        let response;
        try {
            response = await _fetch("getFolder", { session_id: SessionId, folder_id });
        } catch {
            alert("Failed to get folder data");
            return false;
        }

        if (response.status !== 200) {
            alert("Failed to get folder data");
            return false;
        }

        const { iv_base64: folder_iv, ciphertext_base64: folder_ciphertext } = await response.json();

        try {
            const decryptedData = await decryptAESGCM(Key, folder_iv, folder_ciphertext);
            const { folder, subfolders, files } = JSON.parse(decryptedData);
            updateFoldersData(draft => {
                draft[folder_id] = folder;
                for (const [folder_id, folder] of Object.entries(subfolders)) {
                    draft[folder_id] = folder;
                }
            });
            updateFilesData(draft => {
                for (const [file_id, file] of Object.entries(files)) {
                    draft[file_id] = file;
                }
            });
            if (setActiveFolder) setSelectedFolderId(folder_id);
            return true;
        } catch (err) {
            alert("Failed to decrypt folder data");
            console.error(err);
            return false;
        }
    }

    return (
        <ExplorerContext.Provider value={{
            loading,
            foldersData,
            updateFoldersData,
            filesData,
            updateFilesData,
            selectedFolderId,
            setSelectedFolderId,
            tagsInfo,
            loadFolder
        }}>
            {children}
        </ExplorerContext.Provider>
    );
};

export default ExplorerProvider;