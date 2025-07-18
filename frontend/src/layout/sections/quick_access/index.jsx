import "./style.scss"
import { ImageIcon, VideoIcon, AudioIcon, DocumentIcon, ArrowHeadUpIcon } from "../../../icons"
import { useState, useEffect } from "react";
import { useExplorerContext } from "../../../context/explorer_context"

export default function QuickAccess() {
    const btns = [
        [ <ImageIcon style={{ width: '1.5rem', height: '1.5rem' }} />, "Images" ],
        [ <VideoIcon style={{ width: '1.5rem', height: '1.5rem' }} />, "Videos" ],
        [ <AudioIcon style={{ width: '1.5rem', height: '1.5rem' }} />, "Audios" ],
        [ <DocumentIcon style={{ width: '1.5rem', height: '1.5rem' }} />, "Documents" ],
    ]

    const [itemCounts, setItemCounts] = useState([0, 0, 0, 0]);
    const [collapsed, setCollapsed] = useState(false);
    const { tagsInfo, selectedTagState, setSelectedTagState } = useExplorerContext();

    useEffect(() => {
        setItemCounts([
            tagsInfo.SystemTags.images, 
            tagsInfo.SystemTags.videos, 
            tagsInfo.SystemTags.audios, 
            tagsInfo.SystemTags.documents
        ]);
    }, [tagsInfo]);

    return (
        <div className={`quick-access ${collapsed ? "collapsed" : ""}`}>
            <div className="quick-access-content">
                <h2 className="quick-access-title">Quick Access</h2>
                <div className="quick-access-btns">
                    {btns.map((btn, index) => <Button key={index} icon={btn[0]} name={btn[1]} itemCounts={itemCounts[index]} selectedTagState={selectedTagState} setSelectedTagState={setSelectedTagState} />)}
                </div>
            </div>
            <button className="quick-access-expand-btn"
                onClick={() => setCollapsed(!collapsed)}>
                <div className="icon"><ArrowHeadUpIcon style={{ width: '100%', height: '100%' }} /></div>
            </button>
        </div>
    )
}

function Button({ icon, name, itemCounts, selectedTagState, setSelectedTagState }) {
    return (
        <button className={`quick-access-btn ${name} ${selectedTagState?.type === "SystemTags" && selectedTagState?.id === name.toLowerCase() ? "selected" : ""}`} onClick={()=> setSelectedTagState({type: "SystemTags", id: name.toLowerCase()})}>
            {icon}
            <p>{name}</p>
            <span>{itemCounts}</span>
        </button>
    )
}