import React, { Fragment } from "react";
import PropTypes from "prop-types";
import { FormField, Box } from "grommet";
import { Formik } from "formik";
import { string, object, addMethod } from "yup";
import { mainInstance as axios } from "../axios";
import qs from "qs";

import equalTo from "../utils/equalTo";
import { FormPropTypes } from "../PropTypes";
import {
  AuthFormWrapper,
  AuthStyledTextLabel,
  AuthStyledTextInput,
  AuthFormFieldError,
  AuthStyledButton,
  AuthFormSubmittedError,
} from "../ui";
import { endpoints } from "../network/endpoints";

addMethod(string, "equalTo", equalTo);

const forgotPassValidation = object().shape({
  email: string()
    .email("The email must be a valid format")
    .required("Please enter your email"),
});

const Form = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  errors,
  isMobile,
}) => (
  <Fragment>
    <Box
      flex={true}
      direction="row"
      alignSelf="center"
      justify="center"
      align="center"
    >
      <AuthFormWrapper isMobile={isMobile}>
        <form onSubmit={handleSubmit}>
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
          <Box>
            <AuthStyledButton
              margin={{ top: "medium", bottom: "medium" }}
              disabled={isSubmitting}
              type="submit"
              primary
              label="Submit"
              alignSelf="start"
            />
          </Box>
        </form>
      </AuthFormWrapper>
    </Box>
  </Fragment>
);

Form.propTypes = FormPropTypes;

const ForgotPasswordForm = ({ isMobile }) => {
  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.post(
          endpoints.forgotPassword,
          qs.stringify({
            email: values.email,
          })
        );
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
      initialValues={{ email: "" }}
      onSubmit={handleSubmit}
      validationSchema={forgotPassValidation}
    >
      {(props) => <Form isMobile={isMobile} {...props} />}
    </Formik>
  );
};

ForgotPasswordForm.propTypes = {
  isMobile: PropTypes.bool,
};

export default ForgotPasswordForm;
