import React from "react";
import styled from "styled-components";
import { Box } from "grommet";
import { BarLoader } from "react-css-loaders";

const LoadingContainer = styled.div`
  width: 100%;
  height: 100%;
  display: flex;
  margin: 20% 0;
  flex-direction: column;
  align-items: center;
  justify-content: center;
`;

const LoadingOverlay = () => (
  <Box margin="20%" alignSelf="center" animation="fadeOut">
    <BarLoader />
  </Box>
);

export default LoadingOverlay;
