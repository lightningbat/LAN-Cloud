import './style.scss'
// import PropTypes from 'prop-types'

// Spinner.propTypes = {
//     scale: PropTypes.number,
//     color: PropTypes.string,
//     thickness: PropTypes.number
// }
/**
 * 
 * @param {number} scale size/scale of the spinner
 * @param {string} color color of the spinner
 * @param {number} thickness thickness of the spinner (in pixels)
 * @returns {JSX.Element} spinner
 */
export default function Spinner({scale = 1, thickness}) {
    const defaultScale = 50 * scale
    thickness = thickness ? `${thickness}px` : '0.03em'
    return (
        <div className="loading-animation spinner" style={{
            fontSize: `${defaultScale}px`,
            borderWidth: thickness
        }}/>
    )
}