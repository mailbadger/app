import React, { Fragment, useState } from "react";
import { Box, ResponsiveContext } from "grommet";

import ProtectedRoute from "./ProtectedRoute";
import Sidebar from "./Sidebar";
import Templates from "./Templates";
import Settings from "./Settings";

const Routes = React.memo(() => (
  <Box flex align="stretch" justify="start">
    <ProtectedRoute
      path="/dashboard/subscribers"
      component={() => <div>subs</div>}
    />
    <ProtectedRoute
      path="/dashboard/lists"
      component={() => <div>lists</div>}
    />
    <ProtectedRoute path="/dashboard/templates" component={Templates} />
    <ProtectedRoute
      path="/dashboard/campaigns"
      component={() => <div>campaigns</div>}
    />
    <ProtectedRoute path="/dashboard/settings" component={Settings} />
  </Box>
));

const Dashboard = () => {
  const [showSidebar, setSidebar] = useState(true);

  return (
    <ResponsiveContext.Consumer>
      {size => (
        <Fragment>
          {/*<AppBar>
            <Heading
              level="3"
              onClick={() => setSidebar(!showSidebar)}
              margin="none"
            >
              Mail Badger
            </Heading>
          </AppBar>*/}
          <Box
            direction="row"
            flex
            animation="fadeIn"
            overflow={{ horizontal: "hidden" }}
          >
            <Sidebar
              showSidebar={showSidebar}
              size={size}
              closeSidebar={() => setSidebar(false)}
            />
            <Routes />
          </Box>
        </Fragment>
      )}
    </ResponsiveContext.Consumer>
  );
};

export default Dashboard;
