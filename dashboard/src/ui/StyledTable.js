import styled from "styled-components";
import { Table, DataTable } from "grommet";

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

const StyledDataTable = styled(DataTable)`
  background: white;

  ${this} tbody tr th, td {
    background: none;
  }

  ${this} tbody tr:hover, focus {
    background: #fafafa;
  }

  ${this} thead th {
    color: #390099;
  }

  ${this} thead th div {
    border-bottom: none;
  }
`;
export { StyledDataTable, StyledTable };

export default StyledTable;
