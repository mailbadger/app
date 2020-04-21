import React from "react";
import { Switch } from "react-router-dom";
import { Box, Heading, Button } from "grommet";
import { Add } from "grommet-icons";
import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import CreateTemplateForm from "./Create";
import EditTemplateForm from "./Edit";
import { SesKeysProvider, SesKeysConsumer } from "../Settings/SesKeysContext";
import BarLoader from "../ui/BarLoader";
import history from "../history";

const Templates = () => (
  <SesKeysProvider>
    <SesKeysConsumer>
      {({ keys, isLoading }) => {
        if (isLoading) {
          return (
            <Box alignSelf="center" margin="20%">
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
                component={List}
              />
            </Switch>
          );
        }

        return (
          <Box>
            <Box align="center" margin={{ top: "large" }}>
              <Heading level="2">
                Please provide your AWS SES keys first.
              </Heading>
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
      }}
    </SesKeysConsumer>
  </SesKeysProvider>
);

export default Templates;
