import React, { Fragment } from "react";
import PropTypes from "prop-types";
import { Paragraph, FormField, Anchor, Box, Heading } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { Link } from "react-router-dom";
import ReCAPTCHA from "react-google-recaptcha";
import { Mail } from "grommet-icons";
import { string, object, ref, addMethod } from "yup";
import { mainInstance as axios } from "../axios";
import qs from "qs";

import equalTo from "../utils/equalTo";
import { socialAuthEnabled } from "../Auth";
import SocialButtons from "./SocialButtons";
import StyledTextInput from "../ui/StyledTextInput";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import { FormPropTypes } from "../PropTypes";

addMethod(string, "equalTo", equalTo);

const registerValidation = object().shape({
  email: string().email("Please enter your email"),
  password: string().required("Please enter your password").min(8),
  password_confirm: string()
    .equalTo(ref("password"), "Passwords don't match")
    .required("Confirm Password is required"),
});

const Form = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  setFieldValue,
  errors,
}) => (
  <Fragment>
    <Box
      direction="row"
      alignSelf="center"
      background="#ffffff"
      border={{ color: "#CFCFCF" }}
      animation="fadeIn"
      margin={{ top: "40px" }}
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
            marginBottom: "0px",
          }}
        >
          Create your Mailbadger account now
        </Heading>
        <FormField label="Email" htmlFor="email">
          <StyledTextInput
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
            onChange={(value) => setFieldValue("token_response", value, true)}
          />
        )}
        <Paragraph
          style={{ marginTop: "14px", paddingTop: "0px" }}
          size="small"
          textAlign="center"
          alignSelf="center"
          alignContent="center"
        >
          By clicking any of the Sign Up buttons, I agree to the
          <Anchor href=""> terms of service</Anchor>
        </Paragraph>
        <Box>
          <ButtonWithLoader
            margin={{ bottom: "small" }}
            icon={<Mail />}
            disabled={isSubmitting}
            type="submit"
            primary
            label="Sign Up"
          />
        </Box>
        {socialAuthEnabled() && (
          <Fragment>
            <Paragraph
              style={{
                borderTop: "1px solid #CACACA",
                marginTop: "14px",
                paddingTop: "0px",
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
      Already have an account? <Link to="/login">Sign in</Link>{" "}
    </Paragraph>
  </Fragment>
);

Form.propTypes = FormPropTypes;

const RegisterForm = (props) => {
  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        let params = {
          email: values.email,
          password: values.password,
        };

        if (process.env.REACT_APP_RECAPTCHA_SITE_KEY !== "") {
          if (!values.token_response) {
            setErrors({ message: "Invalid re-captcha response." });
            return;
          }

          params.token_response = values.token_response;
        }

        await axios.post("/api/signup", qs.stringify(params));
        props.fetchUser();
      } catch (error) {
        setErrors(error.response.data);
      }
    };

    await callApi();

    //done submitting, set submitting to false
    setSubmitting(false);
  };

  RegisterForm.propTypes = {
    fetchUser: PropTypes.func.isRequired,
  };

  return (
    <Formik
      initialValues={{ email: "", password: "", password_confirm: "" }}
      onSubmit={handleSubmit}
      validationSchema={registerValidation}
    >
      {(props) => <Form {...props} />}
    </Formik>
  );
};

export default RegisterForm;
