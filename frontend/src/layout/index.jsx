import './style.scss'
import { QuickAccess, Explorer, Tags, FolderTree, BottomBar } from "./sections";
import {Toolbar, LockScreen} from "./small_ui";
import { useUiContext } from '../context/ui_context';

function Layout() {
    const { isFolderTreeOpen, authenticated, setAuthenticated } = useUiContext();
    
    return (
        <div className="layout">
            {!authenticated && <LockScreen setAuthenticated={setAuthenticated} />}
            <div className={`side-menu ${isFolderTreeOpen ? "active" : "closed"}`}>
                <FolderTree />
            </div>
            <div className="content">
                <div className="pc-toolbar"><Toolbar /></div>
                <QuickAccess />
                <Explorer />
                <Tags />
                <BottomBar />
            </div>
        </div>
    )
}

export default Layout