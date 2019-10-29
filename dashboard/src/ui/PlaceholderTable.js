import React from "react";
import PropTypes from "prop-types";
import { TableBody } from "grommet";
import StyledTable from "./StyledTable";
import PlaceholderRow from "./PlaceholderRow";

const PlaceholderTable = ({ header: Header, numRows, numCols, ...rest }) => {
  const rows = [];
  for (var i = 0; i < numRows; i++) {
    rows.push(<PlaceholderRow key={i} columns={numCols} />);
  }

  return (
    <StyledTable {...rest}>
      <Header />
      <TableBody>{rows}</TableBody>
    </StyledTable>
  );
};

PlaceholderTable.propTypes = {
  header: PropTypes.func,
  numCols: PropTypes.number,
  numRows: PropTypes.number
};

export default PlaceholderTable;
