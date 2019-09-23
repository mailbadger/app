import React, { Fragment } from "react";
import { FormField, Paragraph, Heading, Box } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { Mail } from "grommet-icons";
import { string, object } from "yup";
import { NavLink } from "react-router-dom";
import axios from "axios";
import qs from "qs";

import SocialButtons from "./SocialButtons";
import { socialAuthEnabled } from "../Auth";
import StyledTextInput from "../ui/StyledTextInput";
import ButtonWithLoader from "../ui/ButtonWithLoader";

const Form = ({ handleSubmit, handleChange, isSubmitting, errors }) => (
  <Fragment>
    <Box
      direction="row"
      alignSelf="center"
      background="#ffffff"
      border={{ color: "#CFCFCF" }}
      animation="fadeIn"
      margin={{ top: "40px", bottom: "40px" }}
      elevation="medium"
      width="medium"
      gap="small"
      pad="medium"
      align="center"
      justify="center"
      style={{ borderRadius: "5px" }}
    >
      <form onSubmit={handleSubmit}>
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
          Welcome back!
        </Heading>
        <Paragraph
          style={{ paddingTop: "0px", marginTop: "5px", fontSize: "17px" }}
          textAlign="center"
          color="#ACACAC"
        >
          We are so excited to see you again!
        </Paragraph>
        <Paragraph textAlign="center" size="small" color="#D85555">
          {errors && errors.message}
        </Paragraph>
        <FormField label="Email" htmlFor="email">
          <StyledTextInput
            placeholder="you@email.com"
            name="email"
            onChange={handleChange}
          />
          <ErrorMessage name="email" />
        </FormField>
        <FormField label="Password" htmlFor="password">
          <StyledTextInput
            placeholder="****"
            name="password"
            type="password"
            onChange={handleChange}
          />
          <ErrorMessage name="password" />
        </FormField>

        <NavLink to="/forgot-password">Forgot your password?</NavLink>

        <ButtonWithLoader
          icon={<Mail />}
          style={{
            width: "100%",
            marginTop: "18px",
            marginBottom: "4px"
          }}
          disabled={isSubmitting}
          type="submit"
          primary
          label={"Login with email"}
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
        <Paragraph
          style={{
            borderTop: "1px solid #CACACA",
            marginTop: "14px",
            paddingTop: "10px"
          }}
          size="small"
          textAlign="center"
          alignSelf="center"
          alignContent="center"
        >
          Don't have an account? <NavLink to="/signup">Sign up</NavLink>
        </Paragraph>
      </form>
    </Box>
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
        props.setUser(result.data.user);
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
