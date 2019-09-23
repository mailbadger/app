import React from "react";
import { Box } from "grommet";
import { BarLoader } from "react-css-loaders";

const LoadingOverlay = () => (
  <Box margin="20%" alignSelf="center" animation="fadeOut">
    <BarLoader />
  </Box>
);

export default LoadingOverlay;
