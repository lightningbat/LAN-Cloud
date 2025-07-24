import { useEffect, useState } from "react";
import ExplorerContext from "./explorerContext";
import { useUiContext } from "../ui_context";
import { getKey, decryptAESGCM, decryptJSON } from "../../utils/crypto";
import { useImmer } from "use-immer";

const ExplorerProvider = ({ children }) => {
    const { authenticated } = useUiContext();

    const [rootFolderId, setRootFolderId] = useState(null);
    const [loading, setLoading] = useState(null);
    const [foldersData, updateFoldersData] = useImmer({});
    const [filesData, updateFilesData] = useImmer({});
    const [selectedFolderId, _setSelectedFolderId] = useState(null);
    const [selectedTagState, _setSelectedTagState] = useState(null); // null | { type: "SystemTags" | "User", id: tag_id | "image" | "video" | "audio" | "document" }
    const [tagsInfo, updateTagsInfo] = useImmer({
        SystemTags: {
            images: 0,
            videos: 0,
            audios: 0,
            documents: 0
        },
        UserTags: {} // tag id => { name, color, item_count }
    })
    const [tagsItems, updateTagsItems] = useImmer({
        SystemTags: {
            images: {}, // file id => {}
            videos: {},
            audios: {},
            documents: {}
        },
        UserTags: {} // tag id => { files: [], folders: [] }
    });
    const [selectionMode, setSelectionMode] = useState(false)
    const [selectedItems, updateSelectedItems] = useImmer({}); // id => "file" | "folder"

    // clear selected items if selection mode is disabled
    useEffect(() => {
        if (!selectionMode) {
            updateSelectedItems((draft) => {
                for (const id in draft) {
                    delete draft[id];
                }
            });
        }
    })

    function setSelectedFolderId(id) {
        if (id === null) {
            _setSelectedFolderId(null);
        } else {
            _setSelectedFolderId(id);
            _setSelectedTagState(null);
        }
    }

    function setSelectedTagState({ type, id }) {
        if (type === null) {
            _setSelectedTagState(null);
        } else {
            _setSelectedTagState({ type, id });
            _setSelectedFolderId(null);
        }
    };
    

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

                // fetch root folder id
                const resp = await fetch(`${import.meta.env.VITE_SERVER_URL}/getRootDirId`);
                if (resp.status !== 200) {
                    alert("Failed to get root folder id");
                    return;
                }
                const { root_dir_id } = await resp.json();
                setRootFolderId(root_dir_id);

                // fetch tags info
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
                    // create empty objects for all user tag items
                    if (data.custom_tags) {
                        updateTagsItems(draft => {
                            for (const [tag_id] of Object.entries(data.custom_tags)) {
                                draft.UserTags[tag_id] = {files: [], folders: []};
                            }
                        });
                    }
                } catch (err) {
                    alert("Failed to decrypt tags");
                    console.error(err);
                }
            }
        })();
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [authenticated]);

    // load root folder
    useEffect(() => {
        if (!rootFolderId) return;
        (async () => {
            await loadFolder(rootFolderId, true);
        })();
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [rootFolderId])

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

                    // climb up till root folder and load all missing parent folders
                    let parent_folder_id = foldersData[selectedFolderId].parent_id;
                    async function recurse(parent_folder_id) {
                        // return if parent_folder_id is root
                        if (parent_folder_id === rootFolderId) return;
                        // skip loading if parent_folder_id is already loaded
                        if (foldersData[parent_folder_id]) {
                            return await recurse(foldersData[parent_folder_id].parent_id);
                        }
                        const res = await loadFolder(parent_folder_id); // load parent folder
                        if (!res) return;
                        return await recurse(res.parent_id);
                    }
                    await recurse(parent_folder_id);
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
            setLoading(null);
            return false;
        }

        if (response.status !== 200) {
            alert("Failed to get folder data");
            setLoading(null);
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
            setLoading(null);
            return folder;
        } catch (err) {
            alert("Failed to decrypt folder data");
            setLoading(null);
            console.error(err);
            return false;
        }
    }

    useEffect(() => {
        if (selectedTagState === null) return;
        if (selectedTagState.type === "SystemTags") {
            // check if system tag is already loaded
            const current_items_count = Object.keys(tagsItems.SystemTags[selectedTagState.id]).length;
            if (current_items_count !== tagsInfo.SystemTags[selectedTagState.id]) {
                (async () => {
                    const { Key, SessionId } = getKey();
                    const resp = await _fetch("getTagItems", { session_id: SessionId, tag: { type: "System", id: selectedTagState.id } });
                    if (resp.status !== 200) {
                        alert("Failed to get tag items");
                        return;
                    }
                    const { iv_base64: iv, ciphertext_base64: ciphertext } = await resp.json();
                    try {
                        const decryptedData = await decryptJSON(Key, iv, ciphertext);
                        updateTagsItems(draft => {
                            draft.SystemTags[selectedTagState.id] = decryptedData;
                        });
                        loadFilesMetaData(Object.keys(decryptedData) || []);
                    } catch (err) {
                        alert("Failed to decrypt tag items");
                        console.error(err);
                    }
                })()
            }
        } else {
            // check if user tag is already loaded
            const current_items_count = tagsItems.UserTags[selectedTagState.id].files.length + tagsItems.UserTags[selectedTagState.id].folders.length;
            if (current_items_count !== tagsInfo.UserTags[selectedTagState.id].item_count) {
                (async () => {
                    const { Key, SessionId } = getKey();
                    const resp = await _fetch("getTagItems", { session_id: SessionId, tag: { type: "User", id: selectedTagState.id } });
                    if (resp.status !== 200) {
                        alert("Failed to get tag items");
                        return;
                    }
                    const { iv_base64: iv, ciphertext_base64: ciphertext } = await resp.json();
                    try {
                        const decryptedData = await decryptJSON(Key, iv, ciphertext);
                        // separating files and folders ids
                        const file_ids = Object.keys(decryptedData).filter(file_id => decryptedData[file_id] === "file");
                        const folder_ids = Object.keys(decryptedData).filter(file_id => decryptedData[file_id] === "folder");
                        updateTagsItems(draft => {
                            draft.UserTags[selectedTagState.id].files = file_ids;
                            draft.UserTags[selectedTagState.id].folders = folder_ids;
                        })
                        await loadFilesMetaData(file_ids);
                        await loadFoldersMetaData(folder_ids);
                    } catch (err) {
                        alert("Failed to decrypt tag items");
                        console.error(err);
                    }
                })()
            }
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedTagState])

    async function loadFilesMetaData(fileIds) {
        if (!fileIds) return;
        // checking if files are already loaded
        const missing_files = fileIds.filter(file_id => !filesData[file_id]);
        if (missing_files.length > 0) {
            const { Key, SessionId } = getKey();
            const resp = await _fetch("getFilesMetaData", { session_id: SessionId, file_ids: missing_files });
            if (resp.status !== 200) {
                alert("Failed to get files meta data");
                return;
            }
            const { iv_base64: iv, ciphertext_base64: ciphertext } = await resp.json();
            try {
                const decryptedData = await decryptJSON(Key, iv, ciphertext);
                updateFilesData(draft => {
                    for (const [file_id, file] of Object.entries(decryptedData)) {
                        draft[file_id] = file;
                    }
                });
            } catch (err) {
                alert("Failed to decrypt files meta data");
                console.error(err);
            }
        }
    }
    async function loadFoldersMetaData(folderIds) {
        if (!folderIds) return;
        const missing_folders = folderIds.filter(folder_id => !foldersData[folder_id]);
        if (missing_folders.length > 0) {
            const { Key, SessionId } = getKey();
            const resp = await _fetch("getFoldersMetaData", { session_id: SessionId, folder_ids: missing_folders });
            if (resp.status !== 200) {
                alert("Failed to get folders meta data");
                return;
            }
            const { iv_base64: iv, ciphertext_base64: ciphertext } = await resp.json();
            try {
                const decryptedData = await decryptJSON(Key, iv, ciphertext);
                updateFoldersData(draft => {
                    for (const [folder_id, folder] of Object.entries(decryptedData)) {
                        draft[folder_id] = folder;
                    }
                });
            } catch (err) {
                alert("Failed to decrypt folders meta data");
                console.error(err);
            }
        }
    }

    return (
        <ExplorerContext.Provider value={{
            rootFolderId,
            loading,
            foldersData,
            updateFoldersData,
            filesData,
            updateFilesData,
            selectedFolderId,
            setSelectedFolderId,
            tagsInfo,
            loadFolder,
            selectedTagState,
            setSelectedTagState,
            tagsItems,
            selectionMode,
            setSelectionMode,
            selectedItems,
            updateSelectedItems,
        }}>
            {children}
        </ExplorerContext.Provider>
    );
};

export default ExplorerProvider;