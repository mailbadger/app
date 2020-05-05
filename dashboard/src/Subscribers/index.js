import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List, { Row, Header, SubscriberTable } from "./List";

const Subscribers = () => (
  <Switch>
    <ProtectedRoute exact path="/dashboard/subscribers" component={List} />
  </Switch>
);

export { Row, Header, SubscriberTable as Table };

export default Subscribers;
