import React, { useContext } from "react";
import PropTypes from "prop-types";
import {
  TableBody,
  Box,
  ResponsiveContext,
  TableRow,
  TableCell,
} from "grommet";
import StyledTable, { StyledTableHeader } from "./DashboardStyledTable";
import PlaceholderRow from "./PlaceholderRow";
import { DashboardSearchPlaceholder } from "./DashboardDataTable";

const DashboardPlaceholderTable = ({ columns, numRows, numCols, ...rest }) => {
  const rows = [];
  for (var i = 0; i < numRows; i++) {
    rows.push(<PlaceholderRow key={i} columns={numCols} />);
  }
  const size = useContext(ResponsiveContext);

  return (
    <Box
      fill="horizontal"
      style={{
        padding: size === "large" ? "0 100px 15px" : "20px",
        display: size === "large" ? "flex" : "table",
        overflow: "auto",
      }}
    >
      <DashboardSearchPlaceholder pad={{ bottom: "10px" }} searchInput="" />
      <StyledTable {...rest}>
        <PlaceholderHeader columns={columns} />
        <TableBody>{rows}</TableBody>
      </StyledTable>
    </Box>
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
