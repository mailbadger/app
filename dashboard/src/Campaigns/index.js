import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import Send from "./Send";

const Campaigns = () => (
  <Switch>
    <ProtectedRoute exact path="/dashboard/campaigns" component={List} />
    <ProtectedRoute
      exact
      path="/dashboard/campaigns/send/:id"
      component={Send}
    />
  </Switch>
);

export default Campaigns;
