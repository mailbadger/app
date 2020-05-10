import React, { useContext } from "react";
import { Switch } from "react-router-dom";
import { Box, Heading, Button } from "grommet";
import { Add } from "grommet-icons";
import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import CreateTemplateForm from "./Create";
import EditTemplateForm from "./Edit";
import { SesKeysContext } from "../Settings/SesKeysContext";
import { BarLoader, ListGrid } from "../ui";
import history from "../history";

const Templates = () => {
  const { keys, isLoading } = useContext(SesKeysContext);

  if (isLoading) {
    return (
      <Box gridArea="nav" alignSelf="center" margin="20%">
        <BarLoader size={15} />
      </Box>
    );
  }

  if (keys) {
    return (
      <Switch>
        <ProtectedRoute
          path="/dashboard/templates/new"
          component={CreateTemplateForm}
        />
        <ProtectedRoute
          path="/dashboard/templates/:id/edit"
          component={EditTemplateForm}
        />
        <ProtectedRoute
          exact
          path="/dashboard/templates"
          component={() => (
            <ListGrid>
              <List />
            </ListGrid>
          )}
        />
      </Switch>
    );
  }

  return (
    <Box>
      <Box align="center" margin={{ top: "large" }}>
        <Heading level="2">Please provide your AWS SES keys first.</Heading>
      </Box>
      <Box align="center" margin={{ top: "medium" }}>
        <Button
          primary
          color="status-ok"
          label="Add SES Keys"
          icon={<Add />}
          reverse
          onClick={() => history.push("/dashboard/settings")}
        />
      </Box>
    </Box>
  );
};

export default Templates;
