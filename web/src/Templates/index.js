import React from "react";
import { Grid, Box, Button } from "grommet";
import { Switch } from "react-router-dom";
import { Add } from "grommet-icons";

import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import CreateTemplateForm from "./Create";
import history from "../history";
import EditTemplateForm from "./Edit";

const Templates = () => {
  return (
    <Grid
      rows={["small", "fill"]}
      columns={["15%", "4fr", "1fr"]}
      gap="medium"
      areas={[
        { name: "nav", start: [0, 0], end: [0, 0] },
        { name: "main", start: [1, 1], end: [1, 1] }
      ]}
    >
      <Box
        gridArea="nav"
        direction="row"
        margin={{ top: "medium", left: "medium" }}
      >
        <Box>
          <Button
            label="Create new"
            icon={<Add />}
            reverse
            onClick={() => history.push("/dashboard/templates/new")}
          />
        </Box>
      </Box>
      <Box gridArea="main">
        <Switch>
          <ProtectedRoute
            path="/dashboard/templates/new"
            component={CreateTemplateForm}
          />
          <ProtectedRoute
            path="/dashboard/templates/:id/edit"
            component={EditTemplateForm}
          />
          <ProtectedRoute exact path="/dashboard/templates" component={List} />
        </Switch>
      </Box>
    </Grid>
  );
};

export default Templates;
