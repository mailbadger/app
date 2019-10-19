import React, { Fragment } from "react";
import { FormField, Button, TextInput, Box } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object, addMethod } from "yup";
import axios from "axios";
import qs from "qs";

import equalTo from "../utils/equalTo";
import { FormPropTypes } from "../PropTypes";

addMethod(string, "equalTo", equalTo);

const forgotPassValidation = object().shape({
  email: string()
    .email("The email must be a valid format")
    .required("Please enter your email")
});

const Form = ({ handleSubmit, handleChange, isSubmitting, errors }) => (
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
    >
      {errors && errors.message && <div>{errors.message}</div>}
      <form
        onSubmit={handleSubmit}
        style={{ color: "black", width: "90%", height: "100%" }}
      >
        <FormField label="Email" htmlFor="email">
          <TextInput
            placeholder="you@email.com"
            name="email"
            onChange={handleChange}
          />
          <ErrorMessage name="email" />
        </FormField>
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
          disabled={isSubmitting}
          type="submit"
          alignSelf="stretch"
          textAlign="center"
          primary
          label="Submit"
        />
      </form>
    </Box>
  </Fragment>
);

Form.propTypes = FormPropTypes;

const ForgotPasswordForm = () => {
  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.post(
          `/api/forgot-password`,
          qs.stringify({
            email: values.email
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
      onSubmit={handleSubmit}
      validationSchema={forgotPassValidation}
      render={props => <Form {...props} />}
    />
  );
};

export default ForgotPasswordForm;
