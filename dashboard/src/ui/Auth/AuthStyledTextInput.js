import styled from "styled-components";
import { TextInput } from "grommet";

const AuthStyledTextInput = styled(TextInput)`
  border: 1px solid #fadcff;
  box-shadow: none;
  color: #000;
  font-weight: 300;
  transition: 0.4s;
  height: 44px;
  border-radius: 20px;
  border: solid 1px #fadcff;
  background-color: #f5f5fa;
  ${this}:focus {
    outline: 0;
    border: 1px solid #8770cf;
    box-shadow: none;
    transition: 0.4s;
  }
`;

export default AuthStyledTextInput;
