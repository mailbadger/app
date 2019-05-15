import React, { Fragment } from "react";
import { FormField, Button, TextInput, Paragraph } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { Facebook, Google, Github, Mail } from "grommet-icons";
import { string, object } from "yup";
import { NavLink } from "react-router-dom";
import axios from "axios";
import qs from "qs";

const Form = ({ handleSubmit, handleChange, isSubmitting, errors }) => (
  <Fragment>
    {errors && errors.message && <div>{errors.message}</div>}
    <form onSubmit={handleSubmit} style={{ width: "90%" }}>
      <FormField label="Email" htmlFor="email">
        <TextInput
          placeholder="you@email.com"
          name="email"
          onChange={handleChange}
        />
        <ErrorMessage name="email" />
      </FormField>
      <FormField label="Password" htmlFor="password">
        <TextInput
          placeholder="****"
          name="password"
          type="password"
          onChange={handleChange}
        />
        <ErrorMessage name="password" />
      </FormField>

      <NavLink to="/forgot-password">Forgot my password?</NavLink>

      <Button
        icon={<Mail />}
        style={{
          marginTop: "18px",
          marginBottom: "4px",
          background: "#654FAA"
        }}
        disabled={isSubmitting}
        type="submit"
        primary
        label="Login with email"
      />

      <Paragraph
        style={{
          borderTop: "1px solid #CACACA",
          marginTop: "14px",
          paddingTop: "0px"
        }}
        size="small"
        textAlign="center"
        alignSelf="center"
        alignContent="center"
      >
        or
      </Paragraph>

      <Button
        style={{
          marginTop: "0px",
          marginBottom: "4px",
          background: "#4267B2"
        }}
        type="button"
        icon={<Facebook />}
        primary
        label=" Login with facebook"
      />
      <Button
        style={{
          marginTop: "4px",
          marginBottom: "4px",
          background: "#4285F4"
        }}
        type="button"
        icon={<Google />}
        primary
        label=" Login with google"
      />
      <Button
        style={{
          marginTop: "4px",
          marginBottom: "4px",
          background: "#333333"
        }}
        type="button"
        icon={<Github />}
        primary
        label=" Login with github"
      />
    </form>
  </Fragment>
);

const loginValidation = object().shape({
  email: string().required("Please enter your email"),
  password: string().required("Please enter your password")
});

const LoginForm = props => {
  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
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

    await callApi();

    //done submitting, set submitting to false
    setSubmitting(false);
  };

  return (
    <Formik
      onSubmit={handleSubmit}
      validationSchema={loginValidation}
      render={props => <Form {...props} />}
    />
  );
};

export default LoginForm;
