import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import Details from "./Details";

const Segments = () => (
  <Switch>
    <ProtectedRoute exact path="/dashboard/segments" component={List} />
    <ProtectedRoute exact path="/dashboard/segments/:id" component={Details} />
  </Switch>
);

export default Segments;
