import React from "react";
import PropTypes from "prop-types";
import ContentLoader from "react-content-loader";
import { TableCell, TableRow } from "grommet";

const PlaceholderRow = ({ columns }) => {
  const random = Math.random() * (1 - 0.7) + 0.7;
  let cols = [];

  for (let i = 0; i < columns; i++) {
    cols.push(
      <TableCell key={i}>
        <ContentLoader
          height={45}
          width={700}
          speed={2}
          primaryColor="#d9d9d9"
          secondaryColor="#ecebeb"
        >
          <rect x="0" y="13" rx="4" ry="4" width={36 * random} height="12" />
          <rect x="34" y="13" rx="6" ry="6" width={240 * random} height="12" />
        </ContentLoader>
      </TableCell>
    );
  }

  return <TableRow>{cols}</TableRow>;
};

PlaceholderRow.propTypes = {
  columns: PropTypes.number
};

export default PlaceholderRow;
