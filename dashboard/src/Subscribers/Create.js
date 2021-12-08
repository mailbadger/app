import React, { useState, useEffect, useContext, useReducer } from "react";
import PropTypes from "prop-types";
import { Formik, ErrorMessage, FieldArray } from "formik";
import { string, object, array } from "yup";
import qs from "qs";
import { Add, Trash } from "grommet-icons";
import { Box, Button, Select, FormField, TextInput, Text } from "grommet";

import { mainInstance as axios } from "../axios";
import { NotificationsContext } from "../Notifications/context";
import { ButtonWithLoader } from "../ui";
import { useApi } from "../hooks";
import { endpoints } from "../network/endpoints";

const subscrValidation = object().shape({
  email: string().email().required("Please enter a subscriber email."),
  name: string().max(191, "The name must not exceed 191 characters."),
  metadata: array().of(
    object().shape({
      key: string()
        .matches(
          /^[\w-]*$/,
          "The key must consist only of alphanumeric and hyphen characters."
        )
        .required("Key is required."),
      val: string()
        .max(191, "The value must not exceed 191 characters.")
        .required("Value is required."),
    })
  ),
});

const reducer = (state, action) => {
  switch (action.type) {
    case "append":
      return [...state, ...action.payload];
    default:
      throw new Error("invalid action type.");
  }
};

const CreateForm = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  hideModal,
  values,
  setFieldValue,
}) => {
  const [selected, setSelected] = useState("");

  const [segments, callApi] = useApi(
    {
      url: `${endpoints.getGroups}?per_page=40`,
    },
    {
      collection: [],
      links: {
        next: null,
      },
    }
  );

  const [options, dispatch] = useReducer(reducer, []);

  useEffect(() => {
    if (segments.isError || segments.isLoading) {
      return;
    }

    if (segments && segments.data) {
      dispatch({ type: "append", payload: segments.data.collection });
    }
  }, [segments]);

  const onMore = () => {
    if (segments.isError || segments.isLoading) {
      return;
    }

    let url = "";
    if (segments && segments.data && segments.data.links) {
      url = segments.data.links.next;
    }

    if (!url) {
      return;
    }

    callApi({
      url: url,
    });
  };

  const onChange = ({ value: nextSelected }) => {
    setFieldValue("segments", nextSelected);
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
          <FormField htmlFor="email" label="Subscriber Email">
            <TextInput
              name="email"
              onChange={handleChange}
              placeholder="john.doe@example.com"
            />
            <ErrorMessage name="email" />
          </FormField>
          <FormField htmlFor="name" label="Subscriber Name (Optional)">
            <TextInput
              name="name"
              onChange={handleChange}
              placeholder="John Doe"
            />
            <ErrorMessage name="name" />
          </FormField>
          <FormField htmlFor="segments" label="Add to groups (Optional)">
            <Select
              multiple
              closeOnChange={false}
              placeholder="select an option..."
              value={selected}
              labelKey="name"
              valueKey="id"
              options={options}
              dropHeight="medium"
              onMore={onMore}
              onChange={onChange}
            />
          </FormField>
          <FieldArray
            name="metadata"
            render={(arrayHelpers) => (
              <Box>
                <Button
                  margin={{ top: "small", bottom: "small" }}
                  alignSelf="start"
                  hoverIndicator="light-1"
                  onClick={() => arrayHelpers.push({ key: "", val: "" })}
                >
                  <Box pad="small" direction="row" align="center" gap="small">
                    <Text>Add field</Text>
                    <Add />
                  </Box>
                </Button>
                <Box flex={true} overflow="auto" style={{ maxHeight: "200px" }}>
                  {values.metadata && values.metadata.length > 0
                    ? values.metadata.map((m, i) => (
                        <Box key={i} direction="row" style={{ flexShrink: 0 }}>
                          <FormField htmlFor={`metadata[${i}].key`} label="Key">
                            <TextInput
                              name={`metadata[${i}].key`}
                              onChange={handleChange}
                              value={m.key}
                            />
                            <ErrorMessage name={`metadata[${i}].key`} />
                          </FormField>
                          <FormField
                            margin={{ left: "small" }}
                            htmlFor={`metadata[${i}].val`}
                            label="Value"
                          >
                            <TextInput
                              name={`metadata[${i}].val`}
                              onChange={handleChange}
                              value={m.val}
                            />
                            <ErrorMessage name={`metadata[${i}].val`} />
                          </FormField>
                          <Button
                            margin={{ left: "small" }}
                            alignSelf="end"
                            hoverIndicator="light-1"
                            onClick={() => arrayHelpers.remove(i)}
                          >
                            <Box pad="small" direction="row" align="center">
                              <Trash />
                            </Box>
                          </Button>
                        </Box>
                      ))
                    : null}
                </Box>
              </Box>
            )}
          />

          <Box direction="row" alignSelf="end" margin={{ top: "large" }}>
            <Box margin={{ right: "small" }}>
              <Button label="Cancel" onClick={() => hideModal()} />
            </Box>
            <Box>
              <ButtonWithLoader
                type="submit"
                primary
                disabled={isSubmitting}
                label="Save Subscriber"
              />
            </Box>
          </Box>
        </Box>
      </form>
    </Box>
  );
};

CreateForm.propTypes = {
  hideModal: PropTypes.func,
  handleSubmit: PropTypes.func,
  handleChange: PropTypes.func,
  isSubmitting: PropTypes.bool,
  setFieldValue: PropTypes.func,
  values: PropTypes.shape({
    metadata: PropTypes.arrayOf(
      PropTypes.shape({
        key: PropTypes.string,
        val: PropTypes.string,
      })
    ),
  }),
};

const CreateSubscriber = ({ callApi, hideModal }) => {
  const { createNotification } = useContext(NotificationsContext);

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const postForm = async () => {
      try {
        let data = {
          email: values.email,
          segments: values.segments,
        };
        if (values.name !== "") {
          data.name = values.name;
        }
        if (values.metadata.length > 0) {
          data.metadata = values.metadata.reduce((map, meta) => {
            map[meta.key] = meta.val;
            return map;
          }, {});
        }

        if (values.segments.length > 0) {
          data.segments = values.segments.map((s) => s.id);
        }

        await axios.post(
          endpoints.postSubscribers,
          qs.stringify(data, { arrayFormat: "brackets" })
        );
        createNotification("Subscriber has been created successfully.");

        //done submitting, set submitting to false
        setSubmitting(false);
        await callApi({ url: endpoints.getSubscribers });

        hideModal();
      } catch (error) {
        if (error.response) {
          const { message, errors } = error.response.data;

          setErrors(errors);

          const msg = message
            ? message
            : "Unable to create subscriber. Please try again.";

          createNotification(msg, "status-error");

          //done submitting, set submitting to false
          setSubmitting(false);
        }
      }
    };

    await postForm();

    return;
  };

  return (
    <Box direction="row">
      <Formik
        initialValues={{ email: "", name: "", metadata: [], segments: [] }}
        onSubmit={handleSubmit}
        validationSchema={subscrValidation}
      >
        {(props) => <CreateForm {...props} hideModal={hideModal} />}
      </Formik>
    </Box>
  );
};

CreateSubscriber.propTypes = {
  callApi: PropTypes.func,
  hideModal: PropTypes.func,
};

export default CreateSubscriber;
