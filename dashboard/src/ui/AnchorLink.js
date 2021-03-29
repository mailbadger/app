/* eslint-disable no-unused-vars */
import React from "react";
import PropTypes from "prop-types";
import { Anchor, ThemeContext } from "grommet";
import { Link } from "react-router-dom";

const AnchorLink = (props) => {
  const { hover, color, active } = props;
  let hfw = "bold";
  let hcolor = "";
  if (hover) {
    hfw = hover.fontWeight;
    hcolor = hover.color;
  }

  const { fontWeight } = props;
  let fw = props.active ? "bold" : "500";
  if (fontWeight) {
    fw = fontWeight;
  }

  let dark = active ? "white" : "light-1";
  let light = active ? "white" : "dark-1";

  return (
    <ThemeContext.Extend
      value={{
        anchor: {
          textDecoration: "none",
          fontWeight: fw,
          color: {
            dark: dark,
            light: light,
          },
          extend: {
            boxShadow: "none",
          },
          hover: {
            textDecoration: "none",
            fontWeight: hfw,
            extend: {
              color: "accent-1",
            },
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
  fontWeight: PropTypes.string,
  hover: PropTypes.shape({
    fontWeight: PropTypes.string,
    color: PropTypes.string,
  }),
  color: PropTypes.shape({
    active: PropTypes.string,
    idle: PropTypes.string,
  }),
};

export default AnchorLink;
