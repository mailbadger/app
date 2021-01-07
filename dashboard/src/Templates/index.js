import React from "react";
import { Switch } from "react-router-dom";
import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import CreateTemplateForm from "./Create";
import EditTemplateForm from "./Edit";
import { ListGrid } from "../ui";

const Templates = () => {
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
};

export default Templates;
