import "./style.scss"
import { CloseIcon, ArrowHeadRightIcon, ArrowHeadDownIcon } from "../../../icons"
import { useEffect } from "react"
import { useUiContext } from "../../../context/ui_context"
import { useExplorerContext } from "../../../context/explorer_context"
import { useImmer } from "use-immer"

export default function FolderTree() {
    // use ui context
    const { isFolderTreeOpen, setIsFolderTreeOpen } = useUiContext();
    const { foldersData, selectedFolderId, setSelectedFolderId, loadFolder } = useExplorerContext();
    const [nodeState, setNodeState] = useImmer({});

    // console.log(nodeState["root"]);

    useEffect(() => {
        if (foldersData["root"]) {
            const _nodeState = {};
            if (nodeState["root"] === undefined) _nodeState["root"] = true;
            else _nodeState["root"] = nodeState["root"];
            const recurse = (parent_id) => {
                if (!foldersData[parent_id]) return;
                const child_ids = Object.keys(foldersData[parent_id].sub_folders);
                for (const child_id of child_ids) {
                    if (nodeState[child_id] === undefined) _nodeState[child_id] = false;
                    else _nodeState[child_id] = nodeState[child_id];
                    recurse(child_id);
                }
            }
            recurse("root");
            setNodeState(_nodeState);
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [foldersData])

    // set node state to open for selected folder
    useEffect(() => {
        if (selectedFolderId) {
            setNodeState(draft => {
                draft[selectedFolderId] = true;
                // if (selectedFolderId !== "root"){
                //     // recurse back to root/missing folder to open all parent folders
                // }
            })
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedFolderId])

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
        <div className={`folder-tree ${isFolderTreeOpen ? "active" : "closed"}`}>
            <div className="top-bar">
                <h2 className="title">Folder Tree</h2>
                <button className="close-btn" onClick={() => setIsFolderTreeOpen(false)}>
                    <CloseIcon style={{ width: '1.2rem', height: '1.2rem' }} />
                </button>
            </div>
            <FolderNode folder_id="root" opened={nodeState} setOpened={setNodeState} />
        </div>
    )
}

function FolderNode({ folder_id, opened, setOpened, }) {

    const { foldersData, selectedFolderId, setSelectedFolderId, loadFolder } = useExplorerContext();

    let name = "", childNodes = [];
    if (foldersData[folder_id]) {
        name = foldersData[folder_id].name;
        childNodes = Object.keys(foldersData[folder_id].sub_folders);
    } else {
        return null
    }

    const toggleOpen = () => {
        setOpened(draft => { draft[folder_id] = !draft[folder_id] })
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
            <p onClick={() => setSelectedFolderId(folder_id)} className="name">{name}</p>
        </div>
        {opened[folder_id] && <div className="left-right-child">
            <div className="white-space"></div>
            <div className="child-nodes">
                {childNodes.map((child, index) => <FolderNode key={index} folder_id={child} opened={opened} setOpened={setOpened} />)}
            </div>
        </div>}
    </div>)
}