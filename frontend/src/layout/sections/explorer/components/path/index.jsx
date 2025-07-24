import "./style.scss"
import { HomeIcon, ArrowHeadRightIcon, ImageIcon, VideoIcon, AudioIcon, DocumentIcon } from "../../../../../icons";
import { useExplorerContext } from "../../../../../context/explorer_context";
import { useEffect, useRef, useState } from "react";

export default function Path(props) {
    const { loading, foldersData, selectedFolderId, setSelectedFolderId, selectedTagState, tagsInfo, rootFolderId } = useExplorerContext();
    const [path, setPath] = useState([]);
    const scrollablePathPartRef = useRef();

    useEffect(() => {
        if (selectedFolderId === null || selectedFolderId === rootFolderId) {
            setPath([]);
            return;
        }
        // load current path
        const _path = [{ id: selectedFolderId, name: foldersData[selectedFolderId]?.name, loading: loading === selectedFolderId  }];
        
        if (!foldersData[selectedFolderId]) return;
        // climb up till root is reached or parent folder is not loaded 
        let currentFolderId = foldersData[selectedFolderId].parent_id;
        while (foldersData[currentFolderId] && currentFolderId !== rootFolderId) {
            _path.unshift({ id: currentFolderId, name: foldersData[currentFolderId]?.name, loading: loading === currentFolderId });
            currentFolderId = foldersData[currentFolderId]?.parent_id;
        }
        setPath(_path);
    }, [selectedFolderId, foldersData, loading, rootFolderId])

    useEffect(() => {
        if (scrollablePathPartRef.current) scrollablePathPartRef.current.scrollLeft = scrollablePathPartRef.current.scrollWidth
    }, [path])
    
    return (
        <div className="path" {...props} >
            {foldersData[rootFolderId] && selectedFolderId && <div className="node root" onClick={()=> setSelectedFolderId(rootFolderId)}><HomeIcon style={{ width: '1.2rem', height: '1.2rem' }} /></div>}
            {selectedTagState && selectedTagState.type === "SystemTags" && (
                <div>
                    {selectedTagState.id === "images" && <div className="system-tag"><ImageIcon style={{ width: '1.2rem', height: '1.2rem' }} /> Images</div>}
                    {selectedTagState.id === "videos" && <div className="system-tag"><VideoIcon style={{ width: '1.2rem', height: '1.2rem' }} /> Videos</div>}
                    {selectedTagState.id === "audios" && <div className="system-tag"><AudioIcon style={{ width: '1.2rem', height: '1.2rem' }} /> Audios</div>}
                    {selectedTagState.id === "documents" && <div className="system-tag"><DocumentIcon style={{ width: '1.2rem', height: '1.2rem' }} /> Documents</div>}
                </div>
            )}
            {selectedTagState && selectedTagState.type === "UserTags" && (
                <div className="user-tag">
                    <div className="circle" style={{ backgroundColor: tagsInfo.UserTags[selectedTagState.id].color }}></div>
                    <p>{tagsInfo.UserTags[selectedTagState.id].name}</p>
                </div>
            )}
            {path != "" &&
                <div className="sub-nodes" ref={scrollablePathPartRef}>
                    {path.map((node) => <>
                        <ArrowHeadRightIcon style={{ width: '1rem', height: '1rem' }} />
                        <div className="node" onClick={()=> setSelectedFolderId(node.id)} key={node.id}>{node.name}</div>
                    </>)}
                </div>
            }
        </div>
    )
}