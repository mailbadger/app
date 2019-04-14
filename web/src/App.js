import React, { Component } from "react";
import { BrowserRouter as Router, Route, Switch } from "react-router-dom";
import { Box, Heading, Grommet } from "grommet";
import Landing from "./Landing";
import { AuthProvider } from "./Auth/AuthContext";
import Dashboard from "./Dashboard";
import ProtectedRoute from "./ProtectedRoute";

const theme = {
  global: {
    font: {
      family: "Roboto",
      size: "14px",
      height: "20px"
    }
  }
};

const AppBar = props => (
  <Box
    tag="header"
    direction="row"
    align="center"
    justify="between"
    background="brand"
    pad={{ left: "medium", right: "small", vertical: "small" }}
    elevation="medium"
    style={{ zIndex: "1" }}
    {...props}
  />
);

class App extends Component {
  state = {
    showSidebar: true
  };

  render() {
    const { showSidebar } = this.state;
    return (
      <Grommet theme={theme} full>
        <Router>
          <AuthProvider>
            <Box fill>
              <AppBar>
                <Heading
                  level="3"
                  onClick={() =>
                    this.setState({ showSidebar: !this.state.showSidebar })
                  }
                  margin="none"
                >
                  Mail Badger
                </Heading>
              </AppBar>
              <Box direction="row" flex overflow={{ horizontal: "hidden" }}>
                <Switch>
                  <ProtectedRoute
                    exact
                    path="/dashboard"
                    component={() => (
                      <Dashboard
                        showSidebar={showSidebar}
                        closeSidebar={() =>
                          this.setState({ showSidebar: false })
                        }
                      />
                    )}
                  />
                  <Route path="/" component={Landing} />
                </Switch>
              </Box>
            </Box>
          </AuthProvider>
        </Router>
      </Grommet>
    );
  }
}

export default App;
