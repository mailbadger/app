import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import Details from "./Details";
import { ListGrid } from "../ui";

const Segments = () => (
  <Switch>
    <ProtectedRoute
      exact
      path="/dashboard/segments"
      component={() => (
        <ListGrid>
          <List />
        </ListGrid>
      )}
    />
    <ProtectedRoute exact path="/dashboard/segments/:id" component={Details} />
  </Switch>
);

export default Segments;
