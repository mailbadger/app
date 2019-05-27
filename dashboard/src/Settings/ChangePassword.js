import React, { Fragment } from "react";
import { FormField, Button, TextInput, Heading } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object, ref, addMethod } from "yup";
import axios from "axios";
import qs from "qs";
import equalTo from "../utils/equalTo";

addMethod(string, "equalTo", equalTo);

const changePassValidation = object().shape({
  password: string().required("Please enter your password."),
  new_password: string()
    .min(8, "Your password must be atleast 8 characters.")
    .required("Password must not be empty."),
  new_password_confirm: string()
    .equalTo(ref("new_password"), "Passwords don't match")
    .required("Confirm Password is required")
});

const Form = ({ handleSubmit, handleChange, isSubmitting, errors }) => (
  <Fragment>
    <Heading level="3">Change password</Heading>
    {errors && errors.message && <div>{errors.message}</div>}
    <form onSubmit={handleSubmit}>
      <FormField label="Old password" htmlFor="password">
        <TextInput name="password" type="password" onChange={handleChange} />
        <ErrorMessage name="password" />
      </FormField>
      <FormField label="New password" htmlFor="new_password">
        <TextInput
          name="new_password"
          type="password"
          onChange={handleChange}
        />
        <ErrorMessage name="new_password" />
      </FormField>
      <FormField label="Confirm new password" htmlFor="new_password_confirm">
        <TextInput
          name="new_password_confirm"
          type="password"
          onChange={handleChange}
        />
        <ErrorMessage name="new_password_confirm" />
      </FormField>

      <Button type="submit" disabled={isSubmitting} label="Update password" />
    </form>
  </Fragment>
);

const ChangePasswordForm = () => {
  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.post(
          "/api/users/password",
          qs.stringify({
            password: values.password,
            new_password: values.new_password
          })
        );
      } catch (error) {
        setErrors(error.response.data);
      }
    };

    await callApi();

    //done submitting, set submitting to false
    setSubmitting(false);

    return;
  };

  return (
    <Formik
      onSubmit={handleSubmit}
      validationSchema={changePassValidation}
      render={Form}
    />
  );
};

export default ChangePasswordForm;
