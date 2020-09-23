import React, { useContext, useState } from "react";
import { Box, ResponsiveContext } from "grommet";

import Notification from "./Notifications";
import { NotificationsProvider } from "./Notifications/context";
import ProtectedRoute from "./ProtectedRoute";
import Sidebar from "./Sidebar";
import Subscribers from "./Subscribers";
import Templates from "./Templates";
import Segments from "./Segments";
import Campaigns from "./Campaigns";
import Settings from "./Settings";
import { SesKeysProvider } from "./Settings/SesKeysContext";

const Routes = React.memo(() => (
  <Box flex align="stretch" justify="start">
    <ProtectedRoute path="/dashboard/subscribers" component={Subscribers} />
    <ProtectedRoute path="/dashboard/segments" component={Segments} />
    <ProtectedRoute
      path="/dashboard/templates"
      component={() => (
        <SesKeysProvider>
          <Templates />
        </SesKeysProvider>
      )}
    />
    <ProtectedRoute
      path="/dashboard/campaigns"
      component={() => (
        <SesKeysProvider>
          <Campaigns />
        </SesKeysProvider>
      )}
    />
    <ProtectedRoute
      path="/dashboard/settings"
      component={() => (
        <SesKeysProvider>
          <Settings />
        </SesKeysProvider>
      )}
    />
  </Box>
));

Routes.displayName = "Routes";

const Dashboard = () => {
  const [showSidebar, setSidebar] = useState(true);
  const size = useContext(ResponsiveContext);

  return (
    <>
      <Box
        style={{ position: "fixed", top: 0, left: 0, height: "100%" }}
        background="brand"
      >
        <Sidebar
          showSidebar={showSidebar}
          size={size}
          closeSidebar={() => setSidebar(false)}
        />
      </Box>
      <Box
        direction="row"
        fill
        animation="fadeIn"
        overflow={{ horizontal: "hidden" }}
        margin={{ left: "18em" }}
      >
        <NotificationsProvider>
          <Routes />
          <Notification />
        </NotificationsProvider>
      </Box>
    </>
  );
};

export default Dashboard;
