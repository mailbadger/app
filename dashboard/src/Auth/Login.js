import React, { Fragment } from "react";
import { FormField, Paragraph, Box } from "grommet";
import { Formik } from "formik";
import { string, object } from "yup";
import { NavLink } from "react-router-dom";
import { mainInstance as axios } from "../axios";
import qs from "qs";
import SocialButtons from "./SocialButtons";
import { socialAuthEnabled } from "../Auth";
import {
  AuthStyledTextInput,
  AuthStyledTextLabel,
  AuthStyledButton,
  AuthStyledHeader,
  AuthFormFieldError,
  AuthStyledRedirectLink,
  CustomLineBreak,
  AuthFormWrapper,
  AuthFormSubmittedError,
} from "../ui";
import { FormPropTypes, AuthFormPropTypes } from "../PropTypes";

const Form = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  errors,
  isMobile,
}) => {
  const style = isMobile ? { height: "fit-content", padding: "20px 10px" } : {};
  return (
    <Box flex={true} direction="column" style={style}>
      <AuthStyledRedirectLink
        text="Don't have an account?"
        redirectLink="/signup"
        redirectLabel="Sign up"
      />
      <Box
        flex={true}
        direction="row"
        alignSelf="center"
        justify="center"
        align="center"
      >
        <AuthFormWrapper isMobile={isMobile}>
          <form onSubmit={handleSubmit}>
            <AuthStyledHeader isMobile={isMobile}>
              Welcome back !
            </AuthStyledHeader>
            <Paragraph
              margin={{ top: "10px", right: "0" }}
              color="#000"
              style={{
                fontSize: "23px",
                lineHeight: "38px",
                fontFamily: "Poppins Medium",
              }}
            >
              We are so excited to see you again !
            </Paragraph>
            <AuthFormSubmittedError>
              {errors && errors.message}
            </AuthFormSubmittedError>
            <FormField
              htmlFor="email"
              label={<AuthStyledTextLabel>Email</AuthStyledTextLabel>}
            >
              <AuthStyledTextInput name="email" onChange={handleChange} />
              <AuthFormFieldError name="email" />
            </FormField>
            <FormField
              htmlFor="password"
              label={<AuthStyledTextLabel>Password</AuthStyledTextLabel>}
            >
              <AuthStyledTextInput
                name="password"
                type="password"
                onChange={handleChange}
              />
              <AuthFormFieldError name="password" />
            </FormField>
            <NavLink
              style={{
                color: "#000",
                fontSize: "13px",
                fontFamily: "Poppins Medium",
              }}
              to="/forgot-password"
            >
              Forgot your password?
            </NavLink>
            <Box>
              <AuthStyledButton
                margin={{ top: "medium", bottom: "medium" }}
                disabled={isSubmitting}
                type="submit"
                primary
                label="Login with email"
              />
            </Box>
            {!socialAuthEnabled() && (
              <Fragment>
                <CustomLineBreak text="or" />
                <SocialButtons />
              </Fragment>
            )}
          </form>
        </AuthFormWrapper>
      </Box>
    </Box>
  );
};

Form.propTypes = FormPropTypes;

const loginValidation = object().shape({
  email: string().required("Please enter your email"),
  password: string().required("Please enter your password"),
});

const LoginForm = ({ isMobile, fetchUser }) => {
  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.post(
          "/api/authenticate",
          qs.stringify({
            username: values.email,
            password: values.password,
          })
        );

        setSubmitting(false);

        fetchUser();
      } catch (error) {
        setSubmitting(false);
        setErrors(error.response.data);
      }
    };

    await callApi();
  };

  return (
    <Formik
      initialValues={{ email: "", password: "" }}
      onSubmit={handleSubmit}
      validationSchema={loginValidation}
    >
      {(props) => <Form isMobile={isMobile} {...props} />}
    </Formik>
  );
};

LoginForm.propTypes = AuthFormPropTypes;

export default LoginForm;
