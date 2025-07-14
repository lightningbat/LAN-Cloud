import './style.scss'
import { SunIcon, MoonIcon } from "../../../../icons"
import { useEffect, useState } from 'react';

export default function ThemeSwitch() {
    const [isToggled, setIsToggled] = useState(false);

    function toggleSwitch() {
        setIsToggled(!isToggled);
        document.documentElement.setAttribute('data-theme', isToggled ? 'light' : 'dark');
        localStorage.setItem('theme', isToggled ? 'light' : 'dark');
    }

    useEffect(() => {
        // get value from local storage
        const theme = localStorage.getItem('theme');
        if (theme) {
            document.documentElement.setAttribute('data-theme', theme);
            setIsToggled(theme === 'dark');
        }
    }, []);

    return (
        <button
            type="button"
            role="switch"
            aria-checked={isToggled}
            className="theme-switch"
            onClick={toggleSwitch}>

            <span className={`switch-thumb ${isToggled ? 'toggled' : ''}`}>
                { isToggled ? 
                <MoonIcon style={{ width: '70%', height: '70%' }} /> : 
                <SunIcon style={{ width: '70%', height: '70%' }} />}
            </span>

        </button>
    )
}