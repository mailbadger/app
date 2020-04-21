import React, { Component, Fragment } from "react";
import PropTypes from "prop-types";
import { Box } from "grommet";
import { Redirect } from "react-router-dom";

import { AuthContext } from "./Auth/context";
import Auth from "./Auth";
import LoadingOverlay from "./ui/LoadingOverlay";

class Landing extends Component {
  constructor(props) {
    super(props);
    this.state = {
      redirectToReferrer: false,
    };
  }

  componentDidMount() {
    if (localStorage.getItem("force_login")) {
      let auth = this.context;
      auth.clear();
    }
  }

  render() {
    const { from } = this.props.location.state || {
      from: { pathname: "/dashboard" },
    };

    let auth = this.context;

    if (auth.isLoading) {
      return <LoadingOverlay />;
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

Landing.propTypes = {
  location: PropTypes.shape({
    state: PropTypes.shape({
      from: PropTypes.shape({
        pathname: PropTypes.string,
      }),
    }),
  }),
};

Landing.contextType = AuthContext;

export default Landing;
