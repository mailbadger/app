import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List from "./List";

const Segments = () => (
  <Switch>
    <ProtectedRoute
      path="/dashboard/segments/new"
      component={CreateTemplateForm}
    />
    <ProtectedRoute
      path="/dashboard/segments/:id/edit"
      component={EditTemplateForm}
    />
    <ProtectedRoute exact path="/dashboard/segments" component={List} />
  </Switch>
);

export default Segments;
