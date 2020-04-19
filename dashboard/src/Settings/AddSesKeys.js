import React, { useState, useContext } from "react";
import PropTypes from "prop-types";
import {
  Layer,
  Box,
  FormField,
  Button,
  TextInput,
  Select,
  Heading,
} from "grommet";
import { Trash } from "grommet-icons";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import { mainInstance as axios } from "../axios";
import qs from "qs";

import { NotificationsContext } from "../Notifications/context";
import regions from "../regions/regions.json";
import useApi from "../hooks/useApi";
import useInterval from "../hooks/useInterval";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import { FormPropTypes } from "../PropTypes";
import StyledSpinner from "../ui/StyledSpinner";

const addSesKeysValidation = object().shape({
  access_key: string().required("Please enter your Amazon access key."),
  secret_key: string().required("Please enter your Amazon secret key."),
  region: string().required("Please enter the Amazon region"),
});

const opts = regions.filter((r) => r.public);

const Form = ({
  handleSubmit,
  values,
  handleChange,
  setFieldValue,
  isSubmitting,
}) => (
  <Box width="medium">
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
          placeholder="Select region"
        />
        <ErrorMessage name="region" />
      </FormField>

      <Box margin={{ top: "medium" }}>
        <ButtonWithLoader
          type="submit"
          primary
          disabled={isSubmitting}
          label="Add keys"
        />
      </Box>
    </form>
  </Box>
);

Form.propTypes = FormPropTypes;

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

SesKey.propTypes = {
  setShowDelete: PropTypes.func,
  sesKey: PropTypes.shape({
    region: PropTypes.string,
    access_key: PropTypes.string,
  }),
};

const deleteKeys = async () => {
  await axios.delete(`/api/ses-keys`);
};

const DeleteLayer = ({ setShowDelete, callApi }) => {
  const hideModal = () => setShowDelete(false);
  const [isSubmitting, setSubmitting] = useState(false);

  return (
    <Layer onEsc={() => hideModal()} onClickOutside={() => hideModal()}>
      <Box width="30em">
        <Heading margin="small" level="4">
          Delete Amazon SES key?
        </Heading>
        <Box direction="row" alignSelf="end" pad="small">
          <Box margin={{ right: "small" }}>
            <Button label="Cancel" onClick={() => hideModal()} />
          </Box>
          <Box>
            <ButtonWithLoader
              primary
              label="Delete"
              color="#FF4040"
              disabled={isSubmitting}
              onClick={async () => {
                setSubmitting(true);
                await deleteKeys();
                await callApi({ url: "/api/ses-keys" });
                setSubmitting(false);
                hideModal();
              }}
            />
          </Box>
        </Box>
      </Box>
    </Layer>
  );
};

DeleteLayer.propTypes = {
  setShowDelete: PropTypes.func,
  callApi: PropTypes.func,
};

const AddSesKeysForm = () => {
  const [showDelete, setShowDelete] = useState(false);
  const [state, callApi] = useApi({
    url: `/api/ses-keys`,
  });
  const { createNotification } = useContext(NotificationsContext);
  const [retries, setRetries] = useState(-1);

  useInterval(
    async () => {
      await callApi({ url: `/api/ses-keys` });
      setRetries(retries - 1);
    },
    retries > 0 ? 1000 : null
  );

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const addKeys = async () => {
      try {
        await axios.post(
          "/api/ses-keys",
          qs.stringify({
            access_key: values.access_key,
            secret_key: values.secret_key,
            region: values.region.code,
          })
        );

        setRetries(5); //reset retries
      } catch (error) {
        if (error.response) {
          setErrors(error.response.data);
          const { message } = error.response.data;
          const msg = message ? message : "Unable to add SES keys";

          createNotification(msg, "status-error");
        }
      }
    };

    await addKeys();

    //done submitting, set submitting to false
    setSubmitting(false);

    return;
  };

  let body = (
    <Formik
      onSubmit={handleSubmit}
      initialValues={{
        region: { code: "", name: "" },
      }}
      validationSchema={addSesKeysValidation}
    >
      {Form}
    </Formik>
  );

  if (retries === 0) {
    setRetries(-1);
    createNotification(
      "Unable to add SES keys, check the IAM permissions and try again.",
      "status-error"
    );
  }
  if (state.isLoading || retries > 0) {
    body = <StyledSpinner size={4} />;
  }

  if (!state.isError && state.data) {
    if (retries > 0) {
      setRetries(-1);
    }

    body = <SesKey setShowDelete={setShowDelete} sesKey={state.data} />;
  }

  return (
    <Box
      pad="large"
      alignSelf="start"
      background="#ffffff"
      elevation="medium"
      animation="fadeIn"
    >
      {showDelete && (
        <DeleteLayer setShowDelete={setShowDelete} callApi={callApi} />
      )}
      <Heading level="4" color="#564392" style={{ marginTop: "0px" }}>
        Amazon SES Keys
      </Heading>
      {body}
    </Box>
  );
};

export default AddSesKeysForm;
