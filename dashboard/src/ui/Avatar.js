import React from "react"
import PropTypes from "prop-types"
import { Box } from "grommet"

const Avatar = ({ name, ...rest }) => (
    <Box
        alignContent="end"
        a11yTitle={`${name} avatar`}
        height="30px"
        width="30px"
        round="full"
        margin={{ top: "0", right: "3px", bottom: "0" }}
        style={{
            borderRadius: "5px",
        }}
        background="url(https://www.gravatar.com/avatar/94d093eda664addd6e450d7e9881bcad?s=80&d=identicon)"
        {...rest}
    />
)

Avatar.propTypes = {
    name: PropTypes.string,
}

export default Avatar
