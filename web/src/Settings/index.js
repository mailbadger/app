import React from "react";
import { Grid, Box } from "grommet";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import ChangePassword from "./ChangePassword";

const Settings = () => (
  <Grid
    rows={["small", "medium"]}
    columns={["1/3", "30%"]}
    gap="medium"
    areas={[{ name: "main", start: [1, 1], end: [1, 1] }]}
  >
    <Box gridArea="main">
      <Switch>
        <ProtectedRoute path="/dashboard/settings" component={ChangePassword} />
      </Switch>
    </Box>
  </Grid>
);

export default Settings;
