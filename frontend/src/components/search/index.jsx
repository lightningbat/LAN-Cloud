import './style.scss'
import { SearchIcon } from '../../icons'
import { useRef } from 'react'

export default function Search() {
    const input = useRef();

    return (
        <div className="search-box">
            <div className='search-icon-cont' onClick={() => input.current.focus()}><SearchIcon style={{ width: '1.5rem', height: '1.5rem' }} /></div>
            <input type="text" placeholder="Search..." ref={input} />
            <button><SearchIcon style={{ width: '1.5rem', height: '1.5rem' }} /></button>
        </div>
    )
}