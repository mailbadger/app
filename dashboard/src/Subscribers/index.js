import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List, { Row, Header, SubscriberTable } from "./List";
import Import from "./Import";
import BulkDelete from "./BulkDelete";
import { ListGrid } from "../ui";

const Subscribers = () => {
  return (
    <Switch>
      <ProtectedRoute
        exact
        path="/dashboard/subscribers"
        component={() => (
          <ListGrid>
            <List />
          </ListGrid>
        )}
      />
      <ProtectedRoute
        exact
        path="/dashboard/subscribers/import"
        component={Import}
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
