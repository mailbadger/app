import React, { Component } from "react";
import { Router, Route, Switch } from "react-router-dom";
import { Box, Grommet } from "grommet";
import { AuthProvider } from "./Auth/context";

import Landing from "./Landing";
import Dashboard from "./Dashboard";
import Logout from "./Auth/Logout";
import VerifyEmail from "./VerifyEmail";
import ProtectedRoute from "./ProtectedRoute";
import history from "./history";

const theme = {
  global: {
    font: {
      family:
        "-apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Helvetica Neue', 'Ubuntu', sans-serif;",
      size: "14px",
      height: "20px"
    },
    colors: {
      background: "#E3E8EE",
      brand: "#6650AA"
    }
  },
  tabs: {
    header: {
      background: "white"
    }
  },
  tab: {
    color: "#888888",
    active: {
      color: "brand"
    },
    border: false
  },
  select: {
    options: {
      text: {
        color: "#333",
        size: "small",
        margin: "0px 20px",
        border: "none"
      },
      container: {
        marginLeft: "-30px"
      }
    },
    connrol: {
      extend: {
        marginLeft: "-30px"
      }
    },
    container: {
      marginLeft: "-30px",
      extend: {
        boxShaddow:
          "0 7px 14px 0 rgba(60,66,87, 0.1), 0 3px 6px 0 rgba(0, 0, 0, .07)",
        border: "1px solid #e4e4e4"
      }
    }
  },

  formField: {
    label: {
      color: "#ACACAC",
      size: "small",
      margin: { vertical: "0", top: "small", horizontal: "0" },
      weight: 300
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
