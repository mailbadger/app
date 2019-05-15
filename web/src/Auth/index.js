import React from "react";
import Login from "./Login";
import Register from "./Register";
import NewPassword from "./NewPassword";
import { Route, Switch } from "react-router-dom";
import ForgotPassword from "./ForgotPassword";

const Auth = props => {
  return (
    <Switch>
      <Route
        path="/login"
        component={() => <Login redirect={props.redirect} />}
      />
      <Route
        path="/signup"
        component={() => <Register redirect={props.redirect} />}
      />
      <Route path="/forgot-password/:token" component={NewPassword} />
      <Route path="/forgot-password" component={ForgotPassword} />
      <Route path="/" component={() => <Login redirect={props.redirect} />} />
    </Switch>
  );
};

export default Auth;
