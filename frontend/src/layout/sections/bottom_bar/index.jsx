import "./style.scss"
import { MenuIcon, ActionbarIcon, ToolbarIcon, ClipboardIcon } from "../../../icons"
import Toolbar from "../../small_ui/toolbar"
import ActionBar from "../../small_ui/action_bar"
import { useState } from "react";
import { useUiContext } from "../../../context/ui_context";

export default function BottomBar() {
    const [showOptionsBar, setShowOptionsBar] = useState(false);
    const [selectedBar, setSelectedBar] = useState(null);

    const { setIsFolderTreeOpen } = useUiContext();

    function handleClick(target) {
        if (target == selectedBar) setShowOptionsBar(!showOptionsBar);
        else setShowOptionsBar(true), setSelectedBar(target);
    }
    return (
        <div className="bottom-bar">
            <div className={`options-bar ${showOptionsBar ? "active" : ""}`}>
                {selectedBar === "toolbar" && <Toolbar />}
                {selectedBar === "action-bar" && <ActionBar />}
            </div>
            <div className="bottom-bar-btns">
                <div className="menu-open-btn" onClick={() => setIsFolderTreeOpen(true)}><MenuIcon style={{ width: '1.5rem', height: '1.5rem' }} /></div>
                <div className={`toolbar-open-btn ${showOptionsBar && selectedBar === "toolbar" ? "active" : ""}`} onClick={() => handleClick("toolbar")}>
                    <ToolbarIcon style={{ width: '1.3rem', height: '1.3rem'  }} />
                    <p className="label">Tools</p>
                </div>
                <div className={`action-bar-open-btn ${showOptionsBar && selectedBar === "action-bar" ? "active" : ""}`} onClick={() => handleClick("action-bar")}>
                    <ActionbarIcon style={{ width: '1.3rem', height: '1.3rem' }} />
                    <p className="label">Actions</p>
                </div>
                <div className="clipboard-open-btn"><ClipboardIcon style={{ width: '1.3rem', height: '1.3rem' }} /></div>
            </div>
        </div>
    )
}