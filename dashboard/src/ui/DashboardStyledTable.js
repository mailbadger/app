import styled from "styled-components";
import { Table, DataTable, TableHeader } from "grommet";
import { tableHeading } from "./DashboardDataTable";

const StyledTable = styled(Table)`
  background: white;

  ${this} tbody {
    tr {
      height: 52px;
    }

    tr:hover,
    focus {
      background: #fafafa;
    }
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
    color: #541388;
  }

  ${this} thead th div {
    border-bottom: none;
  }
`;

const StyledTableHeader = styled(TableHeader)`
  ${this} tr {
    width: 100%;
    font-size: 13px;
    ${tableHeading};
  }

  th {
    border: none;
    font-size: 18px;
    line-height: 24px;

    &:first-of-type {
      padding-left: 30px;
    }
  }
`;
export { StyledDataTable, StyledTable, StyledTableHeader };

export default StyledTable;
