import React from "react"
import PropTypes from "prop-types"
import styled, { keyframes, css } from "styled-components"

const loading = keyframes`
  0% {
    transform: rotate(0deg);
  }
  100% {
    transform: rotate(360deg);
  }
`

const animation = (props) =>
    css`
        ${loading} ${props.duration}s infinite linear;
    `

function getColor(props) {
    const d = document.createElement("div")
    d.style.color = props.color
    document.body.appendChild(d)
    const rgbcolor = window.getComputedStyle(d).color
    const match =
        /rgba?\((\d+)\s*,\s*(\d+)\s*,\s*(\d+)\s*(,\s*\d+[.d+]*)*\)/g.exec(
            rgbcolor
        )
    const color = `${match[1]}, ${match[2]}, ${match[3]}`
    return color
}

const RotateSpin = styled.div`
    animation: ${animation};
    border: ${(props) => `1.1em solid rgba(${getColor(props)}, 0.2)`};
    border-left: ${(props) => `1.1em solid ${props.color}`};
    border-radius: 50%;
    font-size: ${(props) => `${props.size}px`};
    height: 6em;
    margin: 0px auto;
    position: relative;
    text-indent: -9999em;
    transform: translateZ(0);
    width: 6em;

    &:after {
        border-radius: 50%;
        height: 6em;
        width: 6em;
    }
`

const StyledSpinner = (props) => <RotateSpin {...props} />

StyledSpinner.propTypes = {
    color: PropTypes.string,
    duration: PropTypes.number,
    size: PropTypes.number,
}

StyledSpinner.defaultProps = {
    color: "#000",
    duration: 1.1,
    size: 3,
}

export default StyledSpinner
