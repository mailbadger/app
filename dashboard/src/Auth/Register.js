import React, { Fragment } from "react";
import { Paragraph, FormField, Anchor, Box, Heading } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { NavLink } from "react-router-dom";
import ReCAPTCHA from "react-google-recaptcha";
import { Mail, StatusCriticalSmall } from "grommet-icons";
import { string, object, ref, addMethod } from "yup";
import axios from "axios";
import qs from "qs";
import equalTo from "../utils/equalTo";
import { socialAuthEnabled } from "../Auth";
import SocialButtons from "./SocialButtons";
import StyledTextInput from "../ui/StyledTextInput";
import StyledButton from "../ui/StyledButton";

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
    <Box
      direction="row"
      flex="grow"
      alignSelf="center"
      background="#ffffff"
      border={{ color: "#CFCFCF" }}
      animation="fadeIn"
      margin={{ top: "40px", bottom: "10px" }}
      elevation="medium"
      width="medium"
      gap="small"
      pad="medium"
      align="center"
      justify="center"
      style={{ borderRadius: "5px" }}
    >
      {errors && errors.message && <div>{errors.message}</div>}
      <form onSubmit={handleSubmit} style={{}}>
        <Heading
          textAlign="center"
          level="3"
          color="#564392"
          style={{
            fontWeight: "400",
            marginTop: "0px",
            paddingBottom: "0px",
            marginBottom: "0px"
          }}
        >
          Create your MailBadger account now
        </Heading>
        <FormField label="Email" htmlFor="email">
          <StyledTextInput
            style={{ border: "1px solid #CACACA" }}
            placeholder="you@email.com"
            name="email"
            onChange={handleChange}
          />
          <Paragraph
            style={{ margin: "0", padding: "0" }}
            size="small"
            color="#D85555"
          >
            <ErrorMessage name="email" />
          </Paragraph>
        </FormField>
        <FormField label="Password" htmlFor="password">
          <StyledTextInput
            placeholder="****"
            name="password"
            type="password"
            onChange={handleChange}
          />
          <Paragraph
            style={{ margin: "0", padding: "0" }}
            size="small"
            color="#D85555"
          >
            <ErrorMessage name="password" />
          </Paragraph>
        </FormField>
        <FormField label="Confirm Password" htmlFor="password_confirm">
          <StyledTextInput
            placeholder="****"
            name="password_confirm"
            type="password"
            onChange={handleChange}
          />
          <Paragraph
            style={{ margin: "0", padding: "0" }}
            size="small"
            color="#D85555"
          >
            <ErrorMessage name="password_confirm" />
          </Paragraph>
        </FormField>
        {process.env.REACT_APP_RECAPTCHA_SITE_KEY && (
          <ReCAPTCHA
            sitekey={process.env.REACT_APP_RECAPTCHA_SITE_KEY}
            onChange={value => setFieldValue("token_response", value, true)}
          />
        )}
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
        <StyledButton
          style={{
            width: "100%",
            marginBottom: "4px"
          }}
          icon={<Mail />}
          disabled={isSubmitting}
          type="submit"
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
    </Box>
    <Paragraph alignSelf="center" size="small" textAlign="center">
      {" "}
      Already have an account? <NavLink to="/login">Sign in</NavLink>{" "}
    </Paragraph>
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
