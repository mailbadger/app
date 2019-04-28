import React, { Fragment } from "react";
import { FormField, Button, TextInput } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object, ref, addMethod, mixed } from "yup";
import axios from "axios";
import qs from "qs";

const changePassValidation = object().shape({
  access_key: string().required("Please enter your Amazon access key."),
  secret_key: string().required("Please enter your Amazon secret key."),
  region: string().required("Please enter the Amazon region")
});

const Form = ({ handleSubmit, handleChange, isSubmitting, errors }) => (
  <Fragment>
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

      <Button type="submit" primary disabled={isSubmitting} label="Submit" />
    </form>
  </Fragment>
);

const ChangePasswordForm = () => {
  const handleSubmit = (values, { setSubmitting, setErrors }) => {
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
    setTimeout(() => {
      callApi();
      setSubmitting(false);
    }, 3000);
    // callApi();

    // //done submitting, set submitting to false
    // setSubmitting(false);

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
