import './style.scss'

import ThemeSwitch from "./theme_switch"
import { Search } from "../../../components"
import { MenuIcon, UploadIcon, DownloadIcon, ClipboardIcon } from "../../../icons"
import { useUiContext } from '../../../context/ui_context'

export default function Toolbar() {
    const { isFolderTreeOpen, setIsFolderTreeOpen } = useUiContext();
    return (
        <div className="toolbar">
            {!isFolderTreeOpen ? <div className="menu-open-btn" onClick={() => setIsFolderTreeOpen(true)}>
                <MenuIcon style={{ width: '1.5rem', height: '1.5rem' }} />
            </div> : <div />}
            <div className='flex'>
                <ThemeSwitch />
                <div className="search-box-cont"><Search /></div>
                <div className="upload-btn"><UploadIcon style={{ width: '1.3rem', height: '1.3rem' }} /></div>
                <div className="download-btn"><DownloadIcon style={{ width: '1.3rem', height: '1.3rem' }} /></div>
                <div className="clipboard-btn"><ClipboardIcon style={{ width: '1.3rem', height: '1.3rem' }} /></div>
            </div>
        </div>
    )
}