import React from "react";
import PropTypes from "prop-types";
import { Route, Redirect } from "react-router-dom";
import { AuthConsumer } from "./Auth/context";

const componentOrRedirect = (isAuthenticated, Component) => {
  const WrappedComponent = props => {
    return isAuthenticated ? (
      <Component {...props} />
    ) : (
      <Redirect
        to={{
          pathname: "/login",
          state: { from: props.location }
        }}
      />
    );
  };

  WrappedComponent.propTypes = {
    location: PropTypes.any
  };

  return WrappedComponent;
};

const ProtectedRoute = ({ component: Component, ...rest }) => (
  <AuthConsumer>
    {({ isAuthenticated }) => (
      <Route
        render={componentOrRedirect(isAuthenticated, Component)}
        {...rest}
      />
    )}
  </AuthConsumer>
);

ProtectedRoute.propTypes = {
  component: PropTypes.elementType
};

export default ProtectedRoute;
