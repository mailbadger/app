import React, { Fragment } from "react";

import Sidebar from "./Sidebar";
import { Box, ResponsiveContext } from "grommet";

const Dashboard = props => (
  <ResponsiveContext.Consumer>
    {size => (
      <Fragment>
        <Sidebar {...props} size={size} />
        <Box flex align="center" justify="center">
          center
        </Box>
      </Fragment>
    )}
  </ResponsiveContext.Consumer>
);

export default Dashboard;
