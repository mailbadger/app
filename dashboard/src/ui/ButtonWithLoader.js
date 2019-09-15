import React from "react";
import StyledButton from "./StyledButton";
import StyledSpinner from "./StyledSpinner";

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
