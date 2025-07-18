import "./style.scss"
import { ArrowHeadDownIcon } from "../../../icons"
import { useState } from "react";
import { useExplorerContext } from "../../../context/explorer_context";

export default function Tags() {
    const { tagsInfo, selectedTagState, setSelectedTagState } = useExplorerContext();
    const [collapsed, setCollapsed] = useState(false);

    return (
        <div className={`tags ${collapsed ? "collapsed" : ""}`}>
            <div className="tag-expand-btn" onClick={()=> setCollapsed(!collapsed)}>
                <div className="icon"><ArrowHeadDownIcon style={{width: "100%", height: "100%"}} /></div>
            </div>
            <div className="collapsable">
                <h2 className="label">Tags</h2>
                <div className="tag-list">
                    {Object.keys(tagsInfo.UserTags).map((tag_id) => 
                        <Tag 
                            key={tagsInfo.UserTags[tag_id]}
                            id={tag_id}
                            name={tagsInfo.UserTags[tag_id].name} 
                            bgcolor={tagsInfo.UserTags[tag_id].color} 
                            selectedTagState={selectedTagState}
                            setSelectedTagState={setSelectedTagState} />
                        )}
                </div>
            </div>
        </div>
    )
}

function Tag({ id, bgcolor, name, selectedTagState, setSelectedTagState }) {
    return (
        <div className={`tag ${selectedTagState?.id === id ? "selected" : ""}`} onClick={() => setSelectedTagState({ type: "UserTags", id })}>
            <div className="circle" style={{ backgroundColor: bgcolor }}></div>
            <p className="name">{name}</p>
        </div>
    )
}

