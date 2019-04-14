import React, { Component } from "react";
import Auth from "./Auth";
import { Box } from "grommet";
import { Redirect } from "react-router-dom";
import { AuthContext } from "./Auth/AuthContext";

class Landing extends Component {
  state = {
    redirectToReferrer: false
  };

  render() {
    const { from } = this.props.location.state || {
      from: { pathname: "/dashboard" }
    };
    let auth = this.context;

    if (auth.isAuthenticated || this.state.redirectToReferrer) {
      return <Redirect to={from} />;
    }

    return (
      <Box flex align="center" justify="center">
        <Auth
          setSession={auth.setSession}
          redirect={() => this.setState({ redirectToReferrer: true })}
        />
      </Box>
    );
  }
}

Landing.contextType = AuthContext;

export default Landing;
