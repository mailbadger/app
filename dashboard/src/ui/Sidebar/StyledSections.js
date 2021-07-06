import styled, { css } from "styled-components";
import { Box } from "grommet";
import { AnchorLink } from "../index";
import { Link } from "react-router-dom";
import React from "react";

export const activeAnchor = css`
  background-color: #fadcff;
  color: #541388;
`;

export const StyledAnchorLink = styled((props) => <AnchorLink {...props} />)`
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  border-radius: 4px;
  padding: 13px 8.6px 13px 8px;
  ${(props) => (props.active ? `${activeAnchor}` : "color: #444444")};
  position: relative;
  justify-content: center;
`;

export const StyledLink = styled(Link)`
  ${(props) =>
    props.hovered
      ? `${activeAnchor}
border-top-right-radius: 0;
border-bottom-right-radius: 0
`
      : ""};

  /*
      Hack for Safari browser as position fixed is not behaving the same way as on other browsers 
      */
  @media not all and (min-resolution: 0.001dpcm) {
    @media {
      div:nth-of-type(2) {
        margin-left: 73px !important;
      }
    }
  }
`;

export const StyledMenuItemContainer = styled(Box)`
  width: 100%;
  display: flex;
  align-items: center;
  ${(props) =>
    props.active ? ` border-right: 3px solid #ffdaff;` : "border: none"};
  margin-bottom: 10px;
`;
