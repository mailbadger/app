import styled from "styled-components";
import { Table } from "grommet";

const StyledTable = styled(Table)`
  background: white;
  color: #4f566b;
  border-radius: 4px;

  ${this} tbody tr:hover, focus {
    background: #fafafa;
  }

  ${this} thead {
    color: #6650aa;
  }

  ${this} thead th div {
    border-bottom: none;
  }

  ${this} tbody {
    color: #4f566b;
  }
`;

export default StyledTable;
