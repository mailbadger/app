import React, { Fragment } from "react";
import { FormField, Button, TextInput } from "grommet";
import { ErrorMessage } from "formik";

const LoginForm = props => {
  const { handleSubmit, handleChange, errors } = props;

  return (
    <Fragment>
      {errors && errors.message && <div>{errors.message}</div>}
      <form onSubmit={handleSubmit}>
        <FormField label="Email" htmlFor="email">
          <TextInput name="email" onChange={handleChange} />
          <ErrorMessage name="email" />
        </FormField>
        <FormField label="Password" htmlFor="password">
          <TextInput name="password" type="password" onChange={handleChange} />
          <ErrorMessage name="password" />
        </FormField>

        <Button type="submit" primary label="Submit" />
      </form>
    </Fragment>
  );
};

export default LoginForm;
