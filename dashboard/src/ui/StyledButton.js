import styled from "styled-components";
import { Button } from "grommet";

const StyledButton = styled(Button)`
  text-transform: uppercase;
  border-radius: 5px;
  border: 1px soid ${props => props.inputColor || "#654FAA"};
  color: white;
  background: ${props => props.inputColor || "#654FAA"};
  ${this}:hover, focus {
    box-shadow: 0 0 0 2px ${props => props.inputColor || "#654FAA"};
  }
  ${this}:focus {
    box-shadow: 0 0 0 2px ${props => props.inputColor || "#654FAA"};
  }
`;

export default StyledButton;
