import styled from "styled-components";

const Badge = styled.span`
  display: inline-block;
  padding: 0.25em 0.4em;
  font-size: 95%;
  font-weight: 700;
  line-height: 1;
  text-align: center;
  white-space: nowrap;
  vertical-align: baseline;
  color: #541388;
  background-color: ${(props) => props.color || "rgba(84, 19, 136, 0.2)"};
  border: solid 1px #541388;
`;

export default Badge;
