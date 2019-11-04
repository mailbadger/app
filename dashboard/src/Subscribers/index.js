import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List from "./List";

const Subscribers = () => (
  <Switch>
    <ProtectedRoute exact path="/dashboard/subscribers" component={List} />
  </Switch>
);

export default Subscribers;
