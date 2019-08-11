import React from "react";
import styled from "styled-components";
import { RotateSpinLoader } from "react-css-loaders";
import StyledButton from "./StyledButton";

const StyledSpinner = styled(RotateSpinLoader)`
  margin: 0 auto !important;
  font-size: 0.2em !important;
`;

const ButtonWithLoader = ({ disabled, icon, label, ...rest }) => {
  return (
    <StyledButton
      icon={!disabled ? icon : null}
      label={!disabled ? label : <StyledSpinner size={3} color="#fff" />}
      disabled={disabled}
      {...rest}
    />
  );
};

export default ButtonWithLoader;
