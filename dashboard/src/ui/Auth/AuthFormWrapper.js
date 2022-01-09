import styled, { css } from "styled-components"
import { Box } from "grommet"
import PropTypes from "prop-types"

const absoluteCenterPosition = css`
    position: absolute;
    top: 50%;
    transform: translatey(-50%);
`

const AuthFormWrapper = styled(Box)`
    height: fit-content;
    ${({ isMobile }) =>
        isMobile
            ? `width: 335px`
            : `${absoluteCenterPosition}
			width: 503px`}

    @media (max-device-width: 1024px) and (orientation:landscape) {
        width: 503px;
    }

    @media only screen and (min-device-width: 768px) and (max-device-width: 1024px) and (orientation: portrait) {
        ${absoluteCenterPosition};
        width: 503px;
    }
`

AuthFormWrapper.propTypes = {
    isMobile: PropTypes.bool,
}

export default AuthFormWrapper
