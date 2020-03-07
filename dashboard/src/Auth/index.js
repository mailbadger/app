import React from "react";
import Login from "./Login";
import Register from "./Register";
import NewPassword from "./NewPassword";
import { Route, Switch } from "react-router-dom";
import ForgotPassword from "./ForgotPassword";
import { AuthConsumer } from "./context";

export const socialAuthEnabled = () =>
  process.env.REACT_APP_ENABLE_SOCIAL_AUTH === "true";

const Auth = () => {
  return (
    <AuthConsumer>
      {({ fetchUser }) => (
        <Switch>
          <Route
            path="/login"
            component={() => <Login fetchUser={fetchUser} />}
          />
          <Route
            path="/signup"
            component={() => <Register fetchUser={fetchUser} />}
          />
          <Route path="/forgot-password/:token" component={NewPassword} />
          <Route path="/forgot-password" component={ForgotPassword} />
          <Route path="/" component={() => <Login fetchUser={fetchUser} />} />
        </Switch>
      )}
    </AuthConsumer>
  );
};

export default Auth;
