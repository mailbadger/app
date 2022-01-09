import styled, { css } from "styled-components"
import { Box } from "grommet"
import { AnchorLink } from "../index"
import { Link } from "react-router-dom"
import React from "react"

export const activeAnchor = css`
    background-color: #fadcff;
    color: #541388;
`

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
`

export const StyledLink = styled(Link)`
    ${(props) =>
        props.hovered
            ? `${activeAnchor}
borderTopRightRadius: '0',
borderBottomRightRadius: '0'
`
            : ""};
`

export const StyledMenuItemContainer = styled(Box)`
    width: 100%;
    display: flex;
    align-items: center;
    ${(props) =>
        props.active ? ` border-right: 3px solid #ffdaff;` : "border: none"};
    margin-bottom: 10px;
`
