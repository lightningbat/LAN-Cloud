import "./style.scss"
import { CloseIcon, ArrowHeadRightIcon, ArrowHeadDownIcon } from "../../../icons"
import { useEffect, useRef } from "react"
import { useUiContext } from "../../../context/ui_context"
import { useExplorerContext } from "../../../context/explorer_context"
import { useImmer } from "use-immer"

export default function FolderTree() {
    // use ui context
    const { isFolderTreeOpen, setIsFolderTreeOpen } = useUiContext();
    const { foldersData, selectedFolderId, rootFolderId } = useExplorerContext();
    const [nodeState, setNodeState] = useImmer({});
    const folderTreeRef = useRef(null);

    // set node state
    useEffect(() => {
        if (foldersData[rootFolderId]) {
            const _nodeState = {};
            if (nodeState[rootFolderId] === undefined) _nodeState[rootFolderId] = true;
            else _nodeState[rootFolderId] = nodeState[rootFolderId];
            const recurse = (parent_id) => {
                if (!foldersData[parent_id]) return;
                const child_ids = Object.keys(foldersData[parent_id].sub_folders);
                for (const child_id of child_ids) {
                    if (nodeState[child_id] === undefined) _nodeState[child_id] = false;
                    else _nodeState[child_id] = nodeState[child_id];
                    recurse(child_id);
                }
            }
            recurse(rootFolderId);
            setNodeState(_nodeState);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [foldersData])

    // recurse from selected folder to root to set node state to open
    useEffect(() => {
        if (selectedFolderId) {
            setNodeState(draft => {
                draft[selectedFolderId] = true;
                function recurse(parent_id) {
                    if (parent_id === rootFolderId) {
                        draft[rootFolderId] = true;
                        return;
                    }
                    if (!foldersData[parent_id]) return;
                    draft[parent_id] = true;
                    recurse(foldersData[parent_id].parent_id);
                }
                if (selectedFolderId !== rootFolderId) recurse(foldersData[selectedFolderId].parent_id);
            })
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedFolderId, foldersData])

    // scrolls to the selected node
    useEffect(() => {
        if (selectedFolderId) {
            const selectedFolder = folderTreeRef.current.querySelector(`.node-${selectedFolderId}`);
            if (!selectedFolder) return;
            const rect = selectedFolder.getBoundingClientRect();
            if (rect.top >= 0 && rect.bottom <= window.innerHeight) return;
            selectedFolder.scrollIntoView({ behavior: "smooth", block: "center" });
        }
    }, [selectedFolderId, nodeState])

    useEffect(() => {
        const mediaQuery = window.matchMedia("(max-width: 700px)");

        const handleResize = () => {
            if (mediaQuery.matches) {
                setIsFolderTreeOpen(false);
            } else {
                setIsFolderTreeOpen(true);
            }
        };

        handleResize();

        mediaQuery.addEventListener("change", handleResize);

        return () => {
            mediaQuery.removeEventListener("change", handleResize);
        };
        // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    return (
        <div className={`folder-tree ${isFolderTreeOpen ? "active" : "closed"}`} ref={folderTreeRef}>
            <div className="top-bar">
                <h2 className="title">Folder Tree</h2>
                <button className="close-btn" onClick={() => setIsFolderTreeOpen(false)}>
                    <CloseIcon style={{ width: '1.2rem', height: '1.2rem' }} />
                </button>
            </div>
            <FolderNode folder_id={rootFolderId} opened={nodeState} setOpened={setNodeState} />
        </div>
    )
}

function FolderNode({ folder_id, opened, setOpened, }) {

    const { foldersData, filesData, selectedFolderId, setSelectedFolderId, loadFolder } = useExplorerContext();

    let name = "", childNodes = [];
    if (foldersData[folder_id]) {
        name = foldersData[folder_id].name;
        childNodes = Object.keys(foldersData[folder_id].sub_folders);
    } else {
        return null
    }

    const toggleOpen = () => {
        setOpened(draft => { draft[folder_id] = !draft[folder_id] })
        // check if folder is already loaded
        if (foldersData[folder_id]) {
            // check if folder items are already loaded
            const subfolders = Object.keys(foldersData[folder_id].sub_folders);
            const files = Object.keys(foldersData[folder_id].files);
            
            const missingSubfolders = subfolders.filter(subfolder => !foldersData[subfolder]);
            const missingFiles = files.filter(file => !filesData[file]);
            if (missingSubfolders.length === 0 && missingFiles.length === 0) return;
        }
        loadFolder(folder_id);
    }

    return (<div className="folder-node">
        <div className={`left-right-parent ${selectedFolderId === folder_id ? "active" : ""}`}>
            <button className="drop-btn" onClick={() => toggleOpen()}>
                {opened[folder_id] ?
                    <ArrowHeadDownIcon style={{ width: '0.9rem', height: '0.9rem' }} /> :
                    <ArrowHeadRightIcon style={{ width: '0.9rem', height: '0.9rem' }} />
                }
            </button>
            <p onClick={() => setSelectedFolderId(folder_id)} className={`name node-${folder_id}`}>{name}</p>
        </div>
        {opened[folder_id] && <div className="left-right-child">
            <div className="white-space"></div>
            <div className="child-nodes">
                {childNodes.map((child, index) => <FolderNode key={index} folder_id={child} opened={opened} setOpened={setOpened} />)}
            </div>
        </div>}
    </div>)
}