import styled from "styled-components"

const Badge = styled.span`
    display: inline-block;
    padding: 0.25em 0.4em;
    font-size: 95%;
    font-weight: 700;
    line-height: 1;
    text-align: center;
    white-space: nowrap;
    vertical-align: baseline;
    border-radius: 5px;
    color: #fff;
    background-color: ${(props) => props.color || "#CCCCCC"};
`

export default Badge
