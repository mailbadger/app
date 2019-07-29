import React from "react";
import { Box } from "grommet";

const StyledFormBox = () => (
  <Box
    direction="row"
    flex="grow"
    alignSelf="center"
    background="#ffffff"
    border={{ color: "#CFCFCF" }}
    animation="fadeIn"
    margin={{ top: "40px", bottom: "10px" }}
    elevation="medium"
    width="medium"
    gap="small"
    pad="medium"
    align="center"
    justify="center"
    style={{ borderRadius: "5px" }}
  />
);

export default StyledFormBox;
