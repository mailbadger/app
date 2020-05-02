import React, { Component } from "react";
import { Router, Route, Switch } from "react-router-dom";
import { Box, Grommet } from "grommet";
import { grommet } from "grommet/themes";
import { deepMerge } from "grommet/utils";

import { AuthProvider } from "./Auth/context";
import Landing from "./Landing";
import Dashboard from "./Dashboard";
import Logout from "./Auth/Logout";
import VerifyEmail from "./VerifyEmail";
import ProtectedRoute from "./ProtectedRoute";
import history from "./history";

const theme = deepMerge(grommet, {
  global: {
    font: {
      family:
        "-apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', 'Ubuntu', sans-serif;",
      size: "14px",
      height: "20px",
    },
    colors: {
      background: "light-3",
      brand: "#6650AA",
    },
  },
  tabs: {
    header: {
      background: "white",
    },
  },
  tab: {
    color: "#888888",
    active: {
      color: "brand",
    },
    border: false,
  },
  formField: {
    label: {
      color: "#ACACAC",
      size: "small",
    },
    border: false,
    borderColor: "#CACACA",
    margin: 0,
  },
  button: {
    border: {
      radius: "5px",
      color: "#6650AA",
    },
    primary: {
      color: "#6650AA",
    },
  },
  anchor: {
    primary: {
      color: "#999999",
    },
    color: "#6650AA",
  },
});

class App extends Component {
  render() {
    return (
      <Grommet theme={theme} full>
        <Router history={history}>
          <AuthProvider>
            <Box flex background="background">
              <Switch>
                <ProtectedRoute path="/dashboard" component={Dashboard} />
                <Route path="/logout" component={Logout} />
                <Route path="/verify-email/:token" component={VerifyEmail} />
                <Route path="/" component={Landing} />
              </Switch>
            </Box>
          </AuthProvider>
        </Router>
      </Grommet>
    );
  }
}

export default App;
