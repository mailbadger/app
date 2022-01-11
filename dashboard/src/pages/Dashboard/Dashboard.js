import React, { useContext, useState } from "react"
import { Box, ResponsiveContext, Grid } from "grommet"

import Notification from "../../Notifications"
import { NotificationsProvider } from "../../Notifications/context"
import ProtectedRoute from "../../ProtectedRoute"
import Sidebar from "../../ui/Sidebar/Sidebar"
import Subscribers from "../Subscribers"
import Templates from "../Templates"
import Groups from "../Groups"
import Campaigns from "../Campaigns"
import Settings from "../Settings"
import { SesKeysProvider } from "../Settings/SesKeysContext"

const Routes = React.memo(() => (
    <Box flex align="stretch" justify="start">
        <ProtectedRoute path="/dashboard/subscribers" component={Subscribers} />
        <ProtectedRoute path="/dashboard/groups" component={Groups} />
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
))

Routes.displayName = "Routes"

const Dashboard = () => {
    const [showSidebar, setSidebar] = useState(true)
    const size = useContext(ResponsiveContext)

    return (
        <Grid
            fill
            rows={["auto"]}
            columns={["84px", "flex"]}
            areas={[
                { name: "sidebar", start: [0, 0], end: [0, 0] },
                { name: "main", start: [1, 0], end: [1, 0] },
            ]}
        >
            <Sidebar
                gridArea="sidebar"
                showSidebar={showSidebar}
                size={size}
                closeSidebar={() => setSidebar(false)}
            />
            <Box
                gridArea="main"
                overflow="auto"
                animation="fadeIn"
                background="#f5f5fa"
            >
                <NotificationsProvider>
                    <Routes />
                    <Notification />
                </NotificationsProvider>
            </Box>
        </Grid>
    )
}

export default Dashboard
