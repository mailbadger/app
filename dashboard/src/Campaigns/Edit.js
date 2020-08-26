import React, { useState, useEffect, useContext, useReducer } from "react";
import PropTypes from "prop-types";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import qs from "qs";
import { Box, Button, Select, FormField, TextInput, Heading } from "grommet";

import { mainInstance as axios } from "../axios";
import { NotificationsContext } from "../Notifications/context";
import { ButtonWithLoader, StyledSpinner } from "../ui";
import { useApi } from "../hooks";

const campaignValidation = object().shape({
  name: string()
    .required()
    .max(191, "The name must not exceed 191 characters."),
  template_name: string()
    .required()
    .max(191, "The template name must not exceed 191 characters."),
});

const reducer = (templateName) => (state, action) => {
  const filtered = action.payload.filter(
    (templates) => templates.name !== templateName
  );

  switch (action.type) {
    case "append":
      return [...state, ...filtered];
    default:
      throw new Error("invalid action type.");
  }
};

const EditForm = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  hideModal,
  values,
  setFieldValue,
}) => {
  const [selected, setSelected] = useState({ name: values.template_name });
  const [templates, callApi] = useApi(
    {
      url: `/api/templates`,
    },
    {
      collection: [],
      next_token: "",
    }
  );

  const [options, dispatch] = useReducer(reducer(values.template_name), [
    { name: values.template_name },
  ]);

  useEffect(() => {
    if (templates.isError || templates.isLoading) {
      return;
    }

    if (templates && templates.data) {
      dispatch({ type: "append", payload: templates.data.collection });
    }
  }, [templates]);

  const onMore = () => {
    if (templates.isError || templates.isLoading) {
      return;
    }

    let next_token = "";
    if (templates && templates.data) {
      next_token = templates.data.next_token;
    }

    if (!next_token) {
      return;
    }

    callApi({
      url: `/api/templates?next_token=${encodeURIComponent(next_token)}`,
    });
  };

  const onChange = ({ value: nextSelected }) => {
    setFieldValue("template_name", nextSelected.name);
    setSelected(nextSelected);
  };

  return (
    <Box
      direction="column"
      fill
      margin={{ left: "medium", right: "medium", bottom: "medium" }}
    >
      <form onSubmit={handleSubmit}>
        <Box>
          <FormField htmlFor="name" label="Name">
            <TextInput
              name="name"
              value={values.name}
              onChange={handleChange}
              placeholder="Campaign name"
            />
            <ErrorMessage name="name" />
          </FormField>
          <FormField htmlFor="template_name" label="Choose template">
            <Select
              placeholder="select a template..."
              value={selected}
              labelKey="name"
              valueKey="name"
              options={options}
              dropHeight="medium"
              onMore={onMore}
              onChange={onChange}
            />

            <ErrorMessage name="template_name" />
          </FormField>

          <Box direction="row" alignSelf="end" margin={{ top: "large" }}>
            <Box margin={{ right: "small" }}>
              <Button label="Cancel" onClick={() => hideModal()} />
            </Box>
            <Box>
              <ButtonWithLoader
                type="submit"
                primary
                disabled={isSubmitting}
                label="Save Campaign"
              />
            </Box>
          </Box>
        </Box>
      </form>
    </Box>
  );
};

EditForm.propTypes = {
  hideModal: PropTypes.func,
  handleSubmit: PropTypes.func,
  handleChange: PropTypes.func,
  isSubmitting: PropTypes.bool,
  setFieldValue: PropTypes.func,
  values: PropTypes.shape({
    template_name: PropTypes.string,
    name: PropTypes.string,
  }),
};

const EditCampaign = ({ id, onSuccess, hideModal }) => {
  const { createNotification } = useContext(NotificationsContext);
  const [state] = useApi({
    url: `/api/campaigns/${id}`,
  });

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const postForm = async () => {
      try {
        let data = {
          name: values.name,
          template_name: values.template_name,
        };

        await axios.put(`/api/campaigns/${id}`, qs.stringify(data));

        createNotification("Campaign has been updated successfully.");

        //done submitting, set submitting to false
        setSubmitting(false);
        onSuccess();

        hideModal();
      } catch (error) {
        if (error.response) {
          const { message, errors } = error.response.data;

          setErrors(errors);

          const msg = message
            ? message
            : "Unable to edit campaign. Please try again.";

          createNotification(msg, "status-error");

          //done submitting, set submitting to false
          setSubmitting(false);
        }
      }
    };

    await postForm();

    return;
  };

  if (state.isLoading) {
    return (
      <Box margin="15%" alignSelf="center">
        <StyledSpinner color="#2e2e2e" size={8} />
      </Box>
    );
  }

  if (state.isError) {
    return (
      <Box margin="15%" alignSelf="center">
        <Heading level="3">Campaign not found.</Heading>
      </Box>
    );
  }

  return (
    <Box direction="row">
      {!state.isLoading && state.data && (
        <Formik
          initialValues={{
            name: state.data.name,
            template_name: state.data.template_name,
          }}
          onSubmit={handleSubmit}
          validationSchema={campaignValidation}
        >
          {(props) => <EditForm {...props} hideModal={hideModal} />}
        </Formik>
      )}
    </Box>
  );
};

EditCampaign.propTypes = {
  onSuccess: PropTypes.func,
  hideModal: PropTypes.func,
  id: PropTypes.number,
};

export default EditCampaign;
