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

export default function ActionBar() {
    const icon_style = { width: '1.2rem', height: '1.2rem' }
    return (
        <div className="action-bar">
            <button className="tooltip">
                <PasteIcon style={icon_style} />
                <span className="tooltip-text">Paste</span>
            </button>
            <button className="tooltip">
                <CheckBoxIcon style={icon_style} />
                <span className="tooltip-text">Select</span>
            </button>
            <button className="tooltip">
                <CopyIcon style={icon_style} />
                <span className="tooltip-text">Copy</span>
            </button>
            <button className="tooltip">
                <CutIcon style={icon_style} />
                <span className="tooltip-text">Cut</span>
            </button>
            <button className="tooltip">
                <DeleteIcon style={icon_style} />
                <span className="tooltip-text">Delete</span>
            </button>
            <button className="tooltip">
                <RenameIcon style={icon_style} />
                <span className="tooltip-text">Rename</span>
            </button>
            <button className="tooltip">
                <AddIcon style={icon_style} />
                <span className="tooltip-text">Add</span>
            </button>
            <button className="tooltip">
                <TagIcon style={icon_style} />
                <span className="tooltip-text">Tag</span>
            </button>
            <button className="tooltip">
                <LockIcon style={icon_style} />
                <span className="tooltip-text">Lock</span>
            </button>
        </div>
    )
}