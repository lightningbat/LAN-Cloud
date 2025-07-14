import "./style.scss"
import { ArrowHeadDownIcon } from "../../../icons"
import { useState } from "react";
import { useExplorerContext } from "../../../context/explorer_context";

export default function Tags() {
    const { tagsInfo } = useExplorerContext();
    const [collapsed, setCollapsed] = useState(false);

    return (
        <div className={`tags ${collapsed ? "collapsed" : ""}`}>
            <div className="tag-expand-btn" onClick={()=> setCollapsed(!collapsed)}>
                <div className="icon"><ArrowHeadDownIcon style={{width: "100%", height: "100%"}} /></div>
            </div>
            <div className="collapsable">
                <h2 className="label">Tags</h2>
                <div className="tag-list">
                    {/* {tags.map((tag, index) => <Tag key={index} name={tag[1]} bgcolor={tag[0]} />)} */}
                    {Object.keys(tagsInfo.UserTags).map((tag_id) => 
                        <Tag 
                            key={tagsInfo.UserTags[tag_id]} 
                            name={tagsInfo.UserTags[tag_id].name} 
                            bgcolor={tagsInfo.UserTags[tag_id].color} />)}
                </div>
            </div>
        </div>
    )
}

function Tag({ bgcolor, name }) {
    return (
        <div className="tag">
            <div className="circle" style={{ backgroundColor: bgcolor }}></div>
            <p className="name">{name}</p>
        </div>
    )
}

