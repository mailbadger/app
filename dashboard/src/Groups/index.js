import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import Details from "./Details";
import { ListGrid } from "../ui";

const Groups = () => (
  <Switch>
    <ProtectedRoute
      exact
      path="/dashboard/groups"
      component={() => (
        <ListGrid>
          <List />
        </ListGrid>
      )}
    />
    <ProtectedRoute exact path="/dashboard/groups/:id" component={Details} />
  </Switch>
);

export default Groups;
