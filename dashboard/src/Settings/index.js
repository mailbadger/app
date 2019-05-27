import React from "react";
import { Grid, Box } from "grommet";

import ChangePassword from "./ChangePassword";
import AddSesKeys from "./AddSesKeys";

const Settings = () => (
  <Grid
    rows={["small", "medium"]}
    columns={["1/3", "30%"]}
    gap="medium"
    areas={[{ name: "main", start: [1, 0], end: [1, 1] }]}
  >
    <Box gridArea="main">
      <AddSesKeys />
      <ChangePassword />
    </Box>
  </Grid>
);

export default Settings;
