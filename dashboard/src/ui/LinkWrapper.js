import { Link } from "react-router-dom"
import styled from "styled-components"

export const LinkWrapper = styled(Link)`
    text-decoration: none;
    color: black;
    &:hover {
        color: #444444;
    }
    &:focus,
    &:visited,
    &:link,
    &:active {
        text-decoration: none;
    }
`
