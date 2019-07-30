import styled from "styled-components";
import { TextInput } from "grommet";

const StyledTextInput = styled(TextInput)`
  background: #ffffff;
  border: 1px solid #bebebe;
  border-radius: 3px;
  box-shadow: none;
  color: #333333;
  font-weight: 300;
  transition: 0.4s;
  ${this}:focus {
    outline: 0;
    border: 1px solid #8770cf;
    box-shadow: none;
    transition: 0.4s;
  }
`;

export default StyledTextInput;
