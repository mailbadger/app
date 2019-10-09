import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import CreateSegmentForm from "./Create";
import EditSegmentForm from "./Edit";

const Segments = () => (
  <Switch>
    <ProtectedRoute
      path="/dashboard/segments/new"
      component={CreateSegmentForm}
    />
    <ProtectedRoute
      path="/dashboard/segments/:id/edit"
      component={EditSegmentForm}
    />
    <ProtectedRoute exact path="/dashboard/segments" component={List} />
  </Switch>
);

export default Segments;
