import React, { Fragment } from "react";
import { FormField, Button, TextInput, Select, Heading } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import axios from "axios";
import qs from "qs";

import regions from "../regions/regions.json";

const addSesKeysValidation = object().shape({
  access_key: string().required("Please enter your Amazon access key."),
  secret_key: string().required("Please enter your Amazon secret key."),
  region: string().required("Please enter the Amazon region")
});

const opts = regions.filter(r => r.public);

const Form = ({
  handleSubmit,
  values,
  handleChange,
  setFieldValue,
  isSubmitting,
  errors
}) => (
  <Fragment>
    <Heading level="3">Add Amazon SES Keys</Heading>
    {errors && errors.message && <div>{errors.message}</div>}
    <form onSubmit={handleSubmit}>
      <FormField label="Access key" htmlFor="access_key">
        <TextInput name="access_key" onChange={handleChange} />
        <ErrorMessage name="access_key" />
      </FormField>
      <FormField label="Secret key" htmlFor="secret_key">
        <TextInput name="secret_key" onChange={handleChange} />
        <ErrorMessage name="secret_key" />
      </FormField>
      <FormField label="Region" htmlFor="region">
        <Select
          options={opts}
          value={values.region}
          name="region"
          onChange={({ option }) => setFieldValue("region", option, true)}
          valueKey="code"
          labelKey="name"
        />
        <ErrorMessage name="region" />
      </FormField>

      <Button type="submit" disabled={isSubmitting} label="Add keys" />
    </form>
  </Fragment>
);

const AddSesKeysForm = () => {
  const handleSubmit = (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.post(
          "/api/ses-keys",
          qs.stringify({
            access_key: values.access_key,
            secret_key: values.secret_key,
            region: values.region.code
          })
        );
      } catch (error) {
        setErrors(error.response.data);
      }
    };

    callApi();

    //done submitting, set submitting to false
    setSubmitting(false);

    return;
  };

  return (
    <Formik
      onSubmit={handleSubmit}
      initialValues={{
        region: { code: "", name: "" }
      }}
      validationSchema={addSesKeysValidation}
      render={Form}
    />
  );
};

export default AddSesKeysForm;
