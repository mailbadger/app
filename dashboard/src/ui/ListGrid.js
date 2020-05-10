import React, { useContext } from "react";
import PropTypes from "prop-types";
import { Grid, ResponsiveContext } from "grommet";

const ListGrid = ({ children }) => {
  const size = useContext(ResponsiveContext);

  let columns = ["medium", "medium", "small"];
  let areas = [
    ["nav", "nav", "nav"],
    ["main", "main", "main"],
  ];

  if (size === "large") {
    columns = ["small", "small", "large", "xsmall"];
    areas = [
      [".", "nav", "nav", "nav"],
      [".", "main", "main", "main"],
    ];
  }

  return (
    <Grid
      rows={["fill", "fill"]}
      columns={columns}
      gap="small"
      margin="medium"
      areas={areas}
    >
      {children}
    </Grid>
  );
};

ListGrid.propTypes = {
  children: PropTypes.element.isRequired,
};

export default ListGrid;
