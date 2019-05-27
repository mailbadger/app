import React from "react";
import { Route, Redirect } from "react-router-dom";
import { AuthConsumer } from "./Auth/AuthContext";

const ProtectedRoute = ({ component: Component, ...rest }) => (
  <AuthConsumer>
    {({ isAuthenticated }) => (
      <Route
        render={props =>
          isAuthenticated() ? (
            <Component {...props} />
          ) : (
            <Redirect
              to={{
                pathname: "/login",
                state: { from: props.location }
              }}
            />
          )
        }
        {...rest}
      />
    )}
  </AuthConsumer>
);

export default ProtectedRoute;
