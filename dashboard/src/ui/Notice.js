import React, { useState } from "react"
import PropTypes from "prop-types"
import { Box, Button, Text } from "grommet"
import { FormClose } from "grommet-icons"

import StatusIcons from "./StatusIcons"

const Notice = ({ message, status, ...rest }) => {
    const [closed, setClosed] = useState(false)
    if (closed) {
        return null
    }

    return (
        <Box
            border={{
                side: "all",
                size: "small",
                color: rest.borderColor || "dark-2",
            }}
            align="center"
            direction="row"
            gap="small"
            justify="between"
            elevation="medium"
            pad={{ vertical: "xsmall", horizontal: "small" }}
            background={rest.color}
        >
            <Box alignSelf="start">{StatusIcons[status]}</Box>
            <Box alignSelf="center" margin={{ left: "small" }}>
                <Text>{message}</Text>
            </Box>
            <Button
                alignSelf="start"
                icon={<FormClose />}
                onClick={() => setClosed(true)}
                plain
            />
        </Box>
    )
}

Notice.propTypes = {
    message: PropTypes.string,
    status: PropTypes.string,
    borderColor: PropTypes.string,
    color: PropTypes.string,
}

export default Notice
