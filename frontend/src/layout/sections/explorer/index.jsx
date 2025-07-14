import './style.scss'
import Path from "./components/path"
import Content from './components/content'
import { GridIcon } from '../../../icons'
import ActionBar from '../../small_ui/action_bar'

export default function Explorer() {
    return (
        <div className="explorer">
            <div className='explorer-top-bar'>
                <Path />
                <GridIcon style={{ width: '1.5rem', height: '1.5rem' }} />
            </div>
            <div className="left-right">
                <Content />
                <div className='pc-action-bar'><ActionBar /></div>
            </div>
        </div>
    )
}