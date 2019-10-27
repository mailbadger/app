import React from "react";
import PropTypes from "prop-types";
import styled, { keyframes, css } from "styled-components";

const loading = keyframes`
0%,
80%,
100% {
  box-shadow: 0 0;
  height: 4em;
}
40% {
  box-shadow: 0 -2em;
  height: 5em;
}
`;

const animation = props =>
  css`
    ${loading} ${props.duration}s infinite ease-in-out
  `;

const Bar = styled.div`
  animation: ${animation};
  animation-delay: ${props => `${props.duration * -0.16}s`};
  background: ${props => props.color};
  color: ${props => props.color};
  font-size: ${props => `${props.size}px`};
  height: 4em;
  margin: 88px auto;
  position: relative;
  text-indent: -9999em;
  transform: translateZ(0);
  width: 1em;

  &:before {
    animation: ${animation};
    animation-delay: ${props => `${props.duration * -0.32}s`};
    background: ${props => props.color};
    content: "";
    height: 4em;
    left: -1.5em;
    position: absolute;
    top: 0;
    width: 1em;
  }

  &:after {
    animation: ${animation};
    animation-delay: ${props => `${props.duration * 0.08}s`};
    background: ${props => props.color};
    content: "";
    height: 4em;
    left: 1.5em;
    position: absolute;
    top: 0;
    width: 1em;
  }
`;

const BarLoader = props => <Bar {...props} />;

BarLoader.propTypes = {
  color: PropTypes.string,
  duration: PropTypes.number,
  size: PropTypes.number
};

BarLoader.defaultProps = {
  color: "#3e3e3e",
  duration: 1,
  size: 11
};

export default BarLoader;
