import React, { Fragment } from "react";
import PropTypes from "prop-types";
import { FormField, Button, TextInput } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object, ref, addMethod } from "yup";
import { mainInstance as axios } from "../axios";
import qs from "qs";

import equalTo from "../utils/equalTo";
import history from "../history";
import { FormPropTypes } from "../PropTypes";

addMethod(string, "equalTo", equalTo);

const passwordValidation = object().shape({
  password: string().required("Please enter a password").min(8),
  password_confirm: string()
    .equalTo(ref("password"), "Passwords don't match")
    .required("Confirm Password is required"),
});

const Form = ({ handleSubmit, handleChange, isSubmitting, errors }) => (
  <Fragment>
    {errors && errors.message && <div>{errors.message}</div>}
    <form
      onSubmit={handleSubmit}
      style={{ color: "black", width: "90%", height: "100%" }}
    >
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

      <Button
        plain
        style={{
          marginTop: "10px",
          marginBottom: "10px",
          borderRadius: "5px",
          padding: "8px",
          background: "#654FAA",
          width: "100%",
          textAlign: "center",
        }}
        disabled={isSubmitting}
        type="submit"
        alignSelf="stretch"
        textAlign="center"
        primary
        label="Change Password"
      />
    </form>
  </Fragment>
);

Form.propTypes = FormPropTypes;

const NewPasswordForm = (props) => {
  const {
    match: { params },
  } = props;

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.put(
          `/api/forgot-password/${params.token}`,
          qs.stringify({
            password: values.password,
          })
        );

        history.replace("/login");
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
      initialValues={{ password: "", password_confirm: "" }}
      onSubmit={handleSubmit}
      validationSchema={passwordValidation}
    >
      {(props) => <Form {...props} />}
    </Formik>
  );
};

NewPasswordForm.propTypes = {
  match: PropTypes.shape({
    params: PropTypes.shape({
      token: PropTypes.string,
    }),
  }),
};

export default NewPasswordForm;
