import React, { Fragment, useState } from "react";
import {
  Layer,
  Box,
  FormField,
  Button,
  TextInput,
  Select,
  Heading
} from "grommet";
import { Trash } from "grommet-icons";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import axios from "axios";
import qs from "qs";

import regions from "../regions/regions.json";
import useApi from "../hooks/useApi";

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

const SesKey = ({ sesKey, setShowDelete }) => (
  <Box direction="column" margin={{ top: "small" }}>
    <Box direction="row">
      <Box margin={{ right: "small" }}>
        <strong>Region:</strong>
      </Box>
      <Box>{sesKey.region}</Box>
    </Box>
    <Box direction="row">
      <Box margin={{ right: "small" }}>
        <strong>Access key:</strong>
      </Box>
      <Box margin={{ right: "small" }}>{sesKey.access_key}</Box>
      <Box>
        <Button plain icon={<Trash />} onClick={() => setShowDelete(true)} />
      </Box>
    </Box>
  </Box>
);

const deleteKeys = async () => {
  await axios.delete(`/api/ses-keys`);
};

const DeleteLayer = ({ setShowDelete, callApi }) => {
  const hideModal = () => setShowDelete(false);
  return (
    <Layer onEsc={() => hideModal()} onClickOutside={() => hideModal()}>
      <Heading margin="small" level="4">
        Delete key ?
      </Heading>
      <Box direction="row" alignSelf="end" pad="small">
        <Box margin={{ right: "small" }}>
          <Button label="Cancel" onClick={() => hideModal()} />
        </Box>
        <Box>
          <Button
            label="Delete"
            onClick={() => {
              deleteKeys();
              callApi({ url: "/api/ses-keys" });
              hideModal();
            }}
          />
        </Box>
      </Box>
    </Layer>
  );
};

const AddSesKeysForm = () => {
  const [showDelete, setShowDelete] = useState(false);
  const [state, callApi] = useApi({
    url: `/api/ses-keys`
  });

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const addKeys = async () => {
      try {
        await axios.post(
          "/api/ses-keys",
          qs.stringify({
            access_key: values.access_key,
            secret_key: values.secret_key,
            region: values.region.code
          })
        );

        await callApi({ url: `/api/ses-keys` });
      } catch (error) {
        setErrors(error.response.data);
      }
    };

    addKeys();

    //done submitting, set submitting to false
    setSubmitting(false);

    return;
  };

  let body = (
    <Formik
      onSubmit={handleSubmit}
      initialValues={{
        region: { code: "", name: "" }
      }}
      validationSchema={addSesKeysValidation}
      render={Form}
    />
  );

  if (state.isLoading) {
    body = <div>Loading...</div>;
  }

  if (!state.isError && state.data) {
    body = (
      <SesKey
        callApi={callApi}
        setShowDelete={setShowDelete}
        sesKey={state.data}
      />
    );
  }

  return (
    <Fragment>
      {showDelete && (
        <DeleteLayer setShowDelete={setShowDelete} callApi={callApi} />
      )}
      <Heading level="3">Amazon SES Keys</Heading>
      {body}
    </Fragment>
  );
};

export default AddSesKeysForm;
