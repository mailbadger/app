import styled from "styled-components";
import { Button } from "grommet";

const StyledButton = styled(Button)`
  border-radius: 5px;
  border: 1px soid ${props => props.color || "#6FFFB0"};
  color: white;
  background: ${props => props.color || "#654FAA"};
  ${this}:hover, focus {
    box-shadow: 0 0 0 2px ${props => props.color || "#6FFFB0"};
  }
  ${this}:focus {
    box-shadow: 0 0 0 2px ${props => props.color || "#6FFFB0"};
  }
`;

export default StyledButton;
