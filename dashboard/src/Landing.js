import React, { Component, Fragment } from "react";
import { Box } from "grommet";
import { Redirect } from "react-router-dom";

import { AuthContext } from "./Auth/AuthContext";
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

    if (auth.isAuthenticated() || this.state.redirectToReferrer) {
      return <Redirect to={from} />;
    }

    return (
      <Fragment>
        <Box
          direction="row"
          flex="grow"
          alignSelf="center"
          background="#ffffff"
          border={{ color: "#CFCFCF" }}
          animation="fadeIn"
          margin={{ top: "40px", bottom: "40px" }}
          elevation="medium"
          width="medium"
          gap="small"
          pad="medium"
          align="center"
          justify="center"
        >
          <Auth
            setSession={auth.setSession}
            redirect={() => this.setState({ redirectToReferrer: true })}
          />
        </Box>
      </Fragment>
    );
  }
}

export default Landing;
