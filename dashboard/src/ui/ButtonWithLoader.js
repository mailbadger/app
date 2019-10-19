import React from "react";
import PropTypes from "prop-types";
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

ButtonWithLoader.propTypes = {
  disabled: PropTypes.bool,
  icon: PropTypes.element,
  label: PropTypes.node
};

export default ButtonWithLoader;
