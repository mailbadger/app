import React, { Fragment, useState } from "react";
import { Box, ResponsiveContext } from "grommet";

import Notification from "./Notifications";
import { NotificationsProvider } from "./Notifications/context";
import ProtectedRoute from "./ProtectedRoute";
import Sidebar from "./Sidebar";
import Templates from "./Templates";
import Segments from "./Segments";
import Campaigns from "./Campaigns";
import Settings from "./Settings";

const Routes = React.memo(() => (
  <Box flex align="stretch" justify="start">
    <ProtectedRoute
      path="/dashboard/subscribers"
      component={() => <div>subs</div>}
    />
    <ProtectedRoute path="/dashboard/segments" component={Segments} />
    <ProtectedRoute path="/dashboard/templates" component={Templates} />
    <ProtectedRoute path="/dashboard/campaigns" component={Campaigns} />
    <ProtectedRoute path="/dashboard/settings" component={Settings} />
  </Box>
));

Routes.displayName = "Routes";

const Dashboard = () => {
  const [showSidebar, setSidebar] = useState(true);

  return (
    <ResponsiveContext.Consumer>
      {size => (
        <Fragment>
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
            <NotificationsProvider>
              <Routes />
              <Notification />
            </NotificationsProvider>
          </Box>
        </Fragment>
      )}
    </ResponsiveContext.Consumer>
  );
};

export default Dashboard;
