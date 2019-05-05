import React from "react";
import { Formik } from "formik";
import Login from "./Login";
import Register from "./Register";
import { string, object } from "yup";
import axios from "axios";
import qs from "qs";
import { Route, Switch } from "react-router-dom";

const loginValidation = object().shape({
  email: string().email("Please enter an email"),
  password: string().required("Please enter your password")
});

const registerValidation = object().shape({
  email: string().email("Please enter an email"),
  password: string().required("Please enter your password")
});

const Auth = props => {
  const handleSubmit = (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        const result = await axios.post(
          "/api/authenticate",
          qs.stringify({
            username: values.email,
            password: values.password
          })
        );

        result.data.token.expires_in =
          result.data.token.expires_in * 1000 + new Date().getTime();
        localStorage.setItem("token", JSON.stringify(result.data.token));
        localStorage.setItem("user", JSON.stringify(result.data.user));

        props.redirect();
      } catch (error) {
        setErrors(error.response.data);
      }
    };

    callApi();

    //done submitting, set submitting to false
    setSubmitting(false);

    return;
  };

  return (
    <Switch>
      <Route
        path="/login"
        component={() => (
          <Formik
            onSubmit={handleSubmit}
            validationSchema={loginValidation}
            render={props => <Login {...props} />}
          />
        )}
      />
      <Route
        path="/register"
        component={() => (
          <Formik
            onSubmit={handleSubmit}
            validationSchema={registerValidation}
            render={props => <Register {...props} />}
          />
        )}
      />
      <Route
        path="/"
        component={() => (
          <Formik
            onSubmit={handleSubmit}
            validationSchema={loginValidation}
            render={props => <Login {...props} />}
          />
        )}
      />
    </Switch>
  );
};

export default Auth;
