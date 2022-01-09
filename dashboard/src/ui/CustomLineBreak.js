import React from "react"
import PropTypes from "prop-types"
import { Text, Box } from "grommet"
import styled from "styled-components"

const StyledLineBreak = styled(Text)`
    border-top: ${(props) => (props.applyBorder ? "1px solid #000" : "none")};
    width: ${(props) => props.width};
    display: flex;
    flex-direction: column;
    align-self: center;
    text-align: center;
    font-size: 16px;
    line-height: 25px;
    color: #000;
`

const CustomLineBreak = ({ text }) => {
    return (
        <Box direction="row">
            <StyledLineBreak as="div" applyBorder width="208px" />
            <StyledLineBreak as="div" width="31px">
                {text}
            </StyledLineBreak>
            <StyledLineBreak as="div" applyBorder width="208px" />
        </Box>
    )
}

CustomLineBreak.propTypes = {
    text: PropTypes.string,
}

export default CustomLineBreak
