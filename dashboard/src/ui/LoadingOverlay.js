import React from "react"
import { Box } from "grommet"
import BarLoader from "./BarLoader"

const LoadingOverlay = () => (
    <Box margin="20%" alignSelf="center">
        <BarLoader />
    </Box>
)

export default LoadingOverlay
