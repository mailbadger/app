import styled from "styled-components";
import { Table } from "grommet";

const StyledTable = styled(Table)`
  background: white;
  border-radius: 4px;

  ${this} tbody tr:hover, focus {
    background: #fafafa;
  }

  ${this} thead {
    color: #390099;
  }

  ${this} thead th div {
    border-bottom: none;
  }
`;

export default StyledTable;
