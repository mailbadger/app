import React, { Component } from "react";
import { Router, Route, Switch } from "react-router-dom";
import { Box, Grommet } from "grommet";
import Landing from "./Landing";
import { AuthProvider } from "./Auth/AuthContext";
import Dashboard from "./Dashboard";
import Logout from "./Auth/Logout";
import Register from "./Auth/Register";
import ProtectedRoute from "./ProtectedRoute";
import history from "./history";

const theme = {
  global: {
    font: {
      family: "Segoe UI",
      size: "14px",
      height: "20px"
    },
    colors: {
      background: "#F5F7F9",
      brand: "#6650AA"
    }
  },
  formField: {
    label: {
      color: "dark-3",
      size: "small",
      margin: { vertical: "0", top: "small", horizontal: "0" },
      weight: 600
    },
    border: false,
    borderColor: "#CACACA",
    margin: 0
  },

  button: {
    border: {
      radius: "5px",
      color: "#6650AA"
    },
    padding: {
      vertical: "7px",
      horizontal: "24px"
    },
    primary: {
      color: "#6650AA"
    },
    extend: props => {
      let extraStyles = "";
      if (props.primary) {
        extraStyles = `
            text-transform: uppercase;
          `;
      }
      return `
          color: white;
          font-size: 12px;
          font-weight: bold;
          width: 100%;
          border: 0px;
          border-radius:5px;
          ${extraStyles}
        `;
    }
  },
  anchor: {
    primary: {
      color: "#999999"
    },
    color: "#6650AA"
  }
};

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
                <Route path="/" component={Landing} />
                <Route path="/register" component={Register} />
              </Switch>
            </Box>
          </AuthProvider>
        </Router>
      </Grommet>
    );
  }
}

export default App;
