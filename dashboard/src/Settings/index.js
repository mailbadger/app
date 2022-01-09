import React from "react"
import { Box, Grid } from "grommet"
import AddSesKeys from "./AddSesKeys"
import ChangePassword from "./ChangePassword"
import { CustomTabs } from "../ui"

const tabs = [
    {
        title: "Email Transport",
        children: <AddSesKeys />,
    },
    {
        title: "Account",
        children: <ChangePassword />,
    },
]

const Settings = () => (
    <Grid
        rows={["fill"]}
        columns={["1fr", "1fr"]}
        gap="small"
        margin="medium"
        areas={[{ name: "tabs", start: [0, 0], end: [0, 0] }]}
    >
        <Box
            elevation="xsmall"
            round="small"
            margin={{ top: "medium" }}
            gridArea="tabs"
        >
            <CustomTabs tabs={tabs} />
        </Box>
    </Grid>
)

export default Settings
