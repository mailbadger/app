import React, { Component, Fragment } from "react";
import { Box } from "grommet";
import { Redirect } from "react-router-dom";

import { AuthContext } from "./Auth/context";
import Auth from "./Auth";

class Landing extends Component {
  state = {
    redirectToReferrer: false
  };

  static contextType = AuthContext;

  render() {
    const { from } = this.props.location.state || {
      from: { pathname: "/dashboard" }
    };

    let auth = this.context;
    if (auth.isLoading) {
      return <div>Loading...</div>;
    }

    if (auth.isAuthenticated || this.state.redirectToReferrer) {
      return <Redirect to={from} />;
    }

    return (
      <Fragment>
        <Box>
          <Auth />
        </Box>
      </Fragment>
    );
  }
}

export default Landing;
