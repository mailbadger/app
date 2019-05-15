import React, { Fragment } from "react";
import { Paragraph, FormField, Button, TextInput, Anchor } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { Facebook, Google, Github, Mail } from "grommet-icons";
import { string, object, ref, addMethod } from "yup";
import axios from "axios";
import qs from "qs";
import equalTo from "../utils/equalTo";

addMethod(string, "equalTo", equalTo);

const registerValidation = object().shape({
  email: string().email("Please enter your email"),
  password: string().required("Please enter your password"),
  password_confirm: string()
    .equalTo(ref("password"), "Passwords don't match")
    .required("Confirm Password is required")
});

const Form = ({ handleSubmit, handleChange, isSubmitting, errors }) => (
  <Fragment>
    {errors && errors.message && <div>{errors.message}</div>}
    <form
      onSubmit={handleSubmit}
      style={{ color: "black", width: "90%", height: "100%" }}
    >
      <FormField label="Email" htmlFor="email">
        <TextInput
          style={{ border: "1px solid #CACACA" }}
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
      <FormField label="Confirm Password" htmlFor="password_confirm">
        <TextInput
          placeholder="****"
          name="password_confirm"
          type="password"
          onChange={handleChange}
        />
        <ErrorMessage name="password_confirm" />
      </FormField>

      <Paragraph
        style={{ marginTop: "14px", paddingTop: "0px" }}
        size="small"
        textAlign="center"
        alignSelf="center"
        alignContent="center"
      >
        By clicking any of the Sign Up buttons, I agree to the{" "}
        <Anchor href="">terms of service</Anchor>
      </Paragraph>

      <Button
        plain
        style={{
          marginTop: "10px",
          marginBottom: "10px",
          borderRadius: "5px",
          padding: "8px",
          background: "#654FAA",
          width: "100%",
          textAlign: "center"
        }}
        icon={<Mail />}
        disabled={isSubmitting}
        type="submit"
        alignSelf="stretch"
        textAlign="center"
        primary
        label="Sign UP with email"
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
        plain
        style={{
          marginTop: "7px",
          marginBottom: "7px",
          borderRadius: "5px",
          padding: "8px",
          background: "#4267B2",
          width: "100%",
          textAlign: "center"
        }}
        type="button"
        icon={<Facebook />}
        alignSelf="center"
        textAlign="center"
        primary
        label=" Sign Up with facebook"
      />
      <br />
      <Button
        plain
        style={{
          marginTop: "7px",
          marginBottom: "7px",
          borderRadius: "5px",
          padding: "8px",
          background: "#4285F4",
          width: "100%",
          textAlign: "center"
        }}
        type="button"
        icon={<Google />}
        alignSelf="stretch"
        textAlign="center"
        primary
        label=" Sign Up with google"
      />
      <br />
      <Button
        plain
        style={{
          marginTop: "7px",
          marginBottom: "7px",
          borderRadius: "5px",
          padding: "8px",
          background: "#333333",
          width: "100%",
          textAlign: "center"
        }}
        type="button"
        icon={<Github />}
        alignSelf="stretch"
        textAlign="center"
        primary
        label=" Sign Up with github"
      />
    </form>
  </Fragment>
);

const RegisterForm = props => {
  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        const result = await axios.post(
          "/api/signup",
          qs.stringify({
            email: values.email,
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
      validationSchema={registerValidation}
      render={props => <Form {...props} />}
    />
  );
};

export default RegisterForm;
