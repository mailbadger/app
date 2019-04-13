import React, { Component } from "react";
import { Box, Heading, Grommet, ResponsiveContext } from "grommet";
import Auth from "./Auth";
import { AuthProvider, AuthConsumer } from "./Auth/AuthContext";
import Sidebar from "./Sidebar";

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
        <AuthProvider>
          <ResponsiveContext.Consumer>
            {size => (
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
                <AuthConsumer>
                  {({ isAuth }) => (
                    <Box
                      direction="row"
                      flex
                      overflow={{ horizontal: "hidden" }}
                    >
                      {isAuth && (
                        <Sidebar
                          showSidebar={showSidebar}
                          closeSidebar={() =>
                            this.setState({ showSidebar: false })
                          }
                          size={size}
                        />
                      )}
                      <Box flex align="center" justify="center">
                        <Auth />
                      </Box>
                    </Box>
                  )}
                </AuthConsumer>
              </Box>
            )}
          </ResponsiveContext.Consumer>
        </AuthProvider>
      </Grommet>
    );
  }
}

export default App;
