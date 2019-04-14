import React, { Fragment, useState } from "react";
import { Box, Heading, ResponsiveContext } from "grommet";

import ProtectedRoute from "./ProtectedRoute";
import Sidebar from "./Sidebar";

const AppBar = props => (
  <Box
    tag="header"
    direction="row"
    align="center"
    justify="between"
    background="brand"
    pad={{ left: "medium", right: "small", vertical: "small" }}
    elevation="medium"
    style={{ zIndex: "1" }}
    {...props}
  />
);

const Dashboard = () => {
  const [showSidebar, setSidebar] = useState(true);

  return (
    <ResponsiveContext.Consumer>
      {size => (
        <Fragment>
          <AppBar>
            <Heading
              level="3"
              onClick={() => setSidebar(!showSidebar)}
              margin="none"
            >
              Mail Badger
            </Heading>
          </AppBar>
          <Box direction="row" flex overflow={{ horizontal: "hidden" }}>
            <Sidebar
              showSidebar={showSidebar}
              size={size}
              closeSidebar={() => setSidebar(false)}
            />
            <Box flex align="center" justify="center">
              <ProtectedRoute
                path="/dashboard/subscribers"
                component={() => <div>subs</div>}
              />
              <ProtectedRoute
                path="/dashboard/lists"
                component={() => <div>lists</div>}
              />
              <ProtectedRoute
                path="/dashboard/templates"
                component={() => <div>templates</div>}
              />
              <ProtectedRoute
                path="/dashboard/campaigns"
                component={() => <div>campaigns</div>}
              />
              <ProtectedRoute
                path="/dashboard/settings"
                component={() => <div>settings</div>}
              />
            </Box>
          </Box>
        </Fragment>
      )}
    </ResponsiveContext.Consumer>
  );
};

export default Dashboard;
