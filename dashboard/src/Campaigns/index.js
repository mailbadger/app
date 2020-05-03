import React from "react";
import { Switch } from "react-router-dom";

import ProtectedRoute from "../ProtectedRoute";
import List from "./List";
import { SesKeysProvider } from "../Settings/SesKeysContext";

const Segments = () => (
  <Switch>
    <SesKeysProvider>
      <ProtectedRoute exact path="/dashboard/campaigns" component={List} />
    </SesKeysProvider>
  </Switch>
);

export default Segments;
