import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List, { Row, Header, SubscriberTable } from "./List";
import BulkDelete from "./BulkDelete";

const Subscribers = () => {
  return (
    <Switch>
      <ProtectedRoute
        exact
        path="/dashboard/subscribers"
        component={() => <List />}
      />
      <ProtectedRoute
        exact
        path="/dashboard/subscribers/bulk-delete"
        component={BulkDelete}
      />
    </Switch>
  );
};

export { Row, Header, SubscriberTable as Table };

export default Subscribers;
