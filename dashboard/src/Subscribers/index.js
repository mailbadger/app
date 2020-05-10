import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List, { Row, Header, SubscriberTable } from "./List";
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
    </Switch>
  );
};

export { Row, Header, SubscriberTable as Table };

export default Subscribers;
