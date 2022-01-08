import React from "react"
import PropTypes from "prop-types"
import { CheckBox as GrommetCheckBox, Anchor, Box } from "grommet"
import styled from "styled-components"

const StyledCheckbox = styled(Box)`
    color: #000;
    display: block;
    font-size: 16px;
    line-height: 25px;
`

const Checkbox = ({ checked, name, label, handleChange, optionalText }) => (
    <GrommetCheckBox
        onChange={handleChange}
        name={name}
        checked={checked}
        label={
            <StyledCheckbox>
                {label}˝
                {optionalText && (
                    <Anchor as="span" href="" color="#000">
                        {" "}
                        {optionalText}
                    </Anchor>
                )}
            </StyledCheckbox>
        }
    />
)

Checkbox.propTypes = {
    checked: PropTypes.bool,
    name: PropTypes.string,
    label: PropTypes.string,
    handleChange: PropTypes.func,
    optionalText: PropTypes.string,
}

export default Checkbox
