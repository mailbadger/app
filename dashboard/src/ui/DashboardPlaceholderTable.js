import React, { useContext } from "react";
import PropTypes from "prop-types";
import { TableBody, ResponsiveContext, TableRow, TableCell } from "grommet";
import StyledTable, { StyledTableHeader } from "./DashboardStyledTable";
import PlaceholderRow from "./PlaceholderRow";
import {
  DashboardSearchPlaceholder,
  DashboardWrapper,
} from "./DashboardDataTable";

const DashboardPlaceholderTable = ({ columns, numRows, numCols, ...rest }) => {
  const rows = [];
  for (var i = 0; i < numRows; i++) {
    rows.push(<PlaceholderRow key={i} columns={numCols} />);
  }
  const size = useContext(ResponsiveContext);

  return (
    <DashboardWrapper fill="horizontal" overflow="auto" contextSize={size}>
      <DashboardSearchPlaceholder pad={{ bottom: "10px" }} searchInput="" />
      <StyledTable {...rest}>
        <PlaceholderHeader columns={columns} />
        <TableBody>{rows}</TableBody>
      </StyledTable>
    </DashboardWrapper>
  );
};

export const PlaceholderHeader = ({ columns }) => (
  <StyledTableHeader>
    <TableRow>
      {columns.map((column) => {
        const { header, size, align } = column;
        return (
          <TableCell
            key={Math.random()}
            scope="col"
            border="bottom"
            align={align ? align : "start"}
            size={size}
          >
            <strong>{header}</strong>
          </TableCell>
        );
      })}
    </TableRow>
  </StyledTableHeader>
);

PlaceholderHeader.propTypes = {
  columns: PropTypes.array,
};

DashboardPlaceholderTable.propTypes = {
  numCols: PropTypes.number,
  numRows: PropTypes.number,
  columns: PropTypes.array,
};

export default DashboardPlaceholderTable;
