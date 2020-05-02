/* eslint-disable no-unused-vars */
import React from "react";
import PropTypes from "prop-types";
import { Anchor, ThemeContext } from "grommet";
import { Link } from "react-router-dom";

const AnchorLink = (props) => {
  return (
    <ThemeContext.Extend
      value={{
        anchor: {
          textDecoration: "none",
          fontWeight: props.active ? 800 : 500,
          color: {
            dark: props.active ? "brand" : "dark-1",
            light: props.active ? "brand" : "dark-1",
          },
          extend: {
            boxShadow: "none",
          },
          hover: {
            fontWeight: 800,
          },
        },
      }}
    >
      <Anchor
        as={({ active, colorProp, hasIcon, hasLabel, focus, ...rest }) => (
          <Link {...rest} />
        )}
        {...props}
      />
    </ThemeContext.Extend>
  );
};

AnchorLink.propTypes = {
  active: PropTypes.bool,
};

export default AnchorLink;
