import React, { Fragment } from "react";
import { Paragraph, FormField, Button, TextInput, Anchor } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { NavLink } from "react-router-dom";
import ReCAPTCHA from "react-google-recaptcha";
import { Mail } from "grommet-icons";
import { string, object, ref, addMethod } from "yup";
import axios from "axios";
import qs from "qs";
import equalTo from "../utils/equalTo";
import { socialAuthEnabled } from "../Auth";
import SocialButtons from "./SocialButtons";

addMethod(string, "equalTo", equalTo);

const registerValidation = object().shape({
  email: string().email("Please enter your email"),
  password: string()
    .required("Please enter your password")
    .min(8),
  password_confirm: string()
    .equalTo(ref("password"), "Passwords don't match")
    .required("Confirm Password is required")
});

const Form = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  setFieldValue,
  errors
}) => (
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
      {process.env.REACT_APP_RECAPTCHA_SITE_KEY && (
        <ReCAPTCHA
          sitekey={process.env.REACT_APP_RECAPTCHA_SITE_KEY}
          onChange={value => setFieldValue("token_response", value, true)}
        />
      )}
      Already have an account? <NavLink to="/login">Sign in</NavLink>
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
      {socialAuthEnabled() && (
        <Fragment>
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
          <SocialButtons />
        </Fragment>
      )}
    </form>
  </Fragment>
);

const RegisterForm = props => {
  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        let params = {
          email: values.email,
          password: values.password
        };

        if (process.env.REACT_APP_RECAPTCHA_SITE_KEY !== "") {
          if (!values.token_response) {
            setErrors({ message: "Invalid re-captcha response." });
            return;
          }

          params.token_response = values.token_response;
        }

        const result = await axios.post("/api/signup", qs.stringify(params));

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
