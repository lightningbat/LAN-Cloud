import "./style.scss"
import {
    PasteIcon,
    CheckBoxIcon,
    CopyIcon,
    CutIcon,
    DeleteIcon,
    RenameIcon,
    AddIcon, 
    TagIcon, 
    LockIcon
} from "../../../icons"

import { useExplorerContext } from "../../../context/explorer_context"
import { useImmer } from "use-immer"
import { useEffect } from "react"

export default function ActionBar() {
    const icon_style = { width: '1.2rem', height: '1.2rem' }
    const { selectionMode, setSelectionMode, selectedItems, renameState, setRenameState } = useExplorerContext()
    const [activeState, setActiveState] = useImmer({
        paste: false,
        select: false,
        copy: false,
        cut: false,
        delete: false,
        rename: false,
        add: false,
        tag: false,
        lock: false
    })

    useEffect(() => {
        const item_count = Object.keys(selectedItems).length;
        if (item_count === 1) {
            setActiveState(draft => { draft.rename = true })
        } else {
            setActiveState(draft => { draft.rename = false })
            setRenameState(false)
        }
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [selectedItems])
    
    return (
        <div className="action-bar">
            <div className="tooltip">
                <button className="tooltip-btn"><PasteIcon style={icon_style} /></button>
                <span className="tooltip-text">Paste</span>
            </div>
            <div className="tooltip">
                <button className="tooltip-btn" onClick={()=> setSelectionMode(!selectionMode)}><CheckBoxIcon style={icon_style} /></button>
                <span className="tooltip-text">Select</span>
            </div>
            <div className="tooltip">
                <button className="tooltip-btn"><CopyIcon style={icon_style} /></button>
                <span className="tooltip-text">Copy</span>
            </div>
            <div className="tooltip">
                <button className="tooltip-btn"><CutIcon style={icon_style} /></button>
                <span className="tooltip-text">Cut</span>
            </div>
            <div className="tooltip">
                <button className="tooltip-btn"><DeleteIcon style={icon_style} /></button>
                <span className="tooltip-text">Delete</span>
            </div>
            <div className="tooltip">
                <button className="tooltip-btn" disabled={!activeState.rename} onClick={() => setRenameState(!renameState)}>
                    <RenameIcon style={icon_style} />
                </button>
                <span className="tooltip-text">Rename</span>
            </div>
            <div className="tooltip">
                <button className="tooltip-btn"><AddIcon style={icon_style} /></button>
                <span className="tooltip-text">Add</span>
            </div>
            <div className="tooltip">
                <button className="tooltip-btn"><TagIcon style={icon_style} /></button>
                <span className="tooltip-text">Tag</span>
            </div>
            <div className="tooltip">
                <button className="tooltip-btn"><LockIcon style={icon_style} /></button>
                <span className="tooltip-text">Lock</span>
            </div>
        </div>
    )
}