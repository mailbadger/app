import styled from "styled-components";
import PropTypes from "prop-types";
import { Button } from "grommet";

const StyledButton = styled(Button)`
  border-radius: 5px;
  border: 1px solid ${props => props.color || "#654FAA"};
  color: white;
  background: ${props => props.color || "#654FAA"};
  ${this}:hover, focus {
    box-shadow: 0 0 0 2px ${props => props.color || "#6FFFB0"};
  }
  ${this}:focus {
    box-shadow: 0 0 0 2px ${props => props.color || "#6FFFB0"};
  }
`;

StyledButton.propTypes = {
  color: PropTypes.string
};

export default StyledButton;
