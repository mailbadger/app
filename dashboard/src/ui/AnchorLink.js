/* eslint-disable no-unused-vars */
import React from "react";
import PropTypes from "prop-types";
import { Anchor, ThemeContext } from "grommet";
import { Link } from "react-router-dom";

const AnchorLink = (props) => {
  const { hover } = props;
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

  return (
    <ThemeContext.Extend
      value={{
        anchor: {
          textDecoration: "none",
          fontWeight: fw,
          color: {
            dark: props.active ? "brand" : "dark-1",
            light: props.active ? "brand" : "dark-1",
          },
          extend: {
            boxShadow: "none",
          },
          hover: {
            textDecoration: "none",
            fontWeight: hfw,
            extend: {
              color: hcolor,
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
};

export default AnchorLink;
