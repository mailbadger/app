import React, { Component } from "react";
import { Router, Route, Switch } from "react-router-dom";
import { Box, Grommet } from "grommet";
import Landing from "./Landing";
import { AuthProvider } from "./Auth/AuthContext";
import Dashboard from "./Dashboard";
import Logout from "./Auth/Logout";
import ProtectedRoute from "./ProtectedRoute";
import history from "./history";

const theme = {
  global: {
    font: {
      family: "Roboto",
      size: "14px",
      height: "20px"
    }
  }
};

class App extends Component {
  render() {
    return (
      <Grommet theme={theme} full>
        <Router history={history}>
          <AuthProvider>
            <Box fill>
              <Switch>
                <ProtectedRoute path="/dashboard" component={Dashboard} />
                <Route path="/logout" component={Logout} />
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
