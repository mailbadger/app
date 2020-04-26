import React, { useState, useEffect, useContext } from "react";
import PropTypes from "prop-types";
import { Formik, ErrorMessage, FieldArray } from "formik";
import { string, object, array } from "yup";
import qs from "qs";
import { Add, Trash } from "grommet-icons";
import {
  Box,
  Button,
  Select,
  FormField,
  TextInput,
  Text,
  Heading,
} from "grommet";

import { mainInstance as axios } from "../axios";
import useApi from "../hooks/useApi";
import { NotificationsContext } from "../Notifications/context";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import StyledSpinner from "../ui/StyledSpinner";

const subscrValidation = object().shape({
  name: string().max(191, "The name must not exceed 191 characters."),
  metadata: array().of(
    object().shape({
      key: string().matches(
        /^[\w-]*$/,
        "The key must consist only of alphanumeric and hyphen characters."
      ),
      val: string()
        .max(191, "The value must not exceed 191 characters.")
        .required("Value is required."),
    })
  ),
});

const EditForm = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  hideModal,
  values,
  setFieldValue,
}) => {
  const [selected, setSelected] = useState(values.segments);
  const [options, setOptions] = useState({
    collection: values.segments,
    url: "/api/segments?per_page=40",
  });
  const callApi = async () => {
    const res = await axios(options.url);
    let url = "";
    if (res.data.links.next) {
      url = res.data.links.next;
    }

    // filter out the preselected segments to avoid duplicate items
    let col = [];
    for (let i = 0; i < res.data.collection.length; i++) {
      let found = false;
      for (let j = 0; j < values.segments.length; j++) {
        if (values.segments[j].id === res.data.collection[i].id) {
          found = true;
          break;
        }
      }

      if (!found) {
        col.push(res.data.collection[i]);
      }
    }

    setOptions({
      collection: [...options.collection, ...col],
      url: url,
    });
  };

  useEffect(() => {
    callApi();
  }, []);

  const onMore = () => {
    if (options.url) {
      callApi();
    }
  };

  const onChange = ({ value: nextSelected }) => {
    setSelected(nextSelected);
    setFieldValue("segments", nextSelected);
  };

  return (
    <Box
      direction="column"
      fill
      margin={{ left: "medium", right: "medium", bottom: "medium" }}
    >
      <form onSubmit={handleSubmit}>
        <Box>
          <FormField htmlFor="name" label="Subscriber Name (Optional)">
            <TextInput
              name="name"
              value={values.name}
              onChange={handleChange}
              placeholder="John Doe"
            />
            <ErrorMessage name="name" />
          </FormField>
          <FormField htmlFor="segments" label="Add to segments (Optional)">
            <Select
              multiple
              closeOnChange={false}
              placeholder="select an option..."
              value={selected}
              labelKey="name"
              valueKey="id"
              options={options.collection}
              dropHeight="medium"
              onMore={onMore}
              onChange={onChange}
            />
          </FormField>
          <FieldArray
            name="metadata"
            render={(arrayHelpers) => (
              <Box flex={true} overflow="auto" style={{ maxHeight: "200px" }}>
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

EditForm.propTypes = {
  hideModal: PropTypes.func,
  handleSubmit: PropTypes.func,
  handleChange: PropTypes.func,
  isSubmitting: PropTypes.bool,
  setFieldValue: PropTypes.func,
  values: PropTypes.shape({
    name: PropTypes.string,
    metadata: PropTypes.arrayOf(
      PropTypes.shape({
        key: PropTypes.string,
        val: PropTypes.string,
      })
    ),
    segments: PropTypes.arrayOf(
      PropTypes.shape({
        id: PropTypes.number,
      })
    ),
  }),
};

const EditSubscriber = ({ id, callApi, hideModal }) => {
  const { createNotification } = useContext(NotificationsContext);
  const [state] = useApi({
    url: `/api/subscribers/${id}`,
  });

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const postForm = async () => {
      try {
        let data = {
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

        await axios.put(
          `/api/subscribers/${id}`,
          qs.stringify(data, { arrayFormat: "brackets" })
        );
        createNotification("Subscriber has been edited successfully.");

        //done submitting, set submitting to false
        setSubmitting(false);
        await callApi({ url: "/api/subscribers" });

        hideModal();
      } catch (error) {
        if (error.response) {
          const { message, errors } = error.response.data;

          setErrors(errors);

          const msg = message
            ? message
            : "Unable to edit subscriber. Please try again.";

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
        <Heading level="3">Subscriber not found.</Heading>
      </Box>
    );
  }

  let m = [];
  if (!state.isLoading && state.data) {
    const { metadata } = state.data;

    for (var key in metadata) {
      if (Object.prototype.hasOwnProperty.call(metadata, key)) {
        m.push({ key: key, val: metadata[key] });
      }
    }
  }

  return (
    <Box direction="row">
      {!state.isLoading && state.data && (
        <Formik
          initialValues={{
            name: state.data.name,
            metadata: m,
            segments: state.data.segments,
          }}
          onSubmit={handleSubmit}
          validationSchema={subscrValidation}
        >
          {(props) => <EditForm {...props} hideModal={hideModal} />}
        </Formik>
      )}
    </Box>
  );
};

EditSubscriber.propTypes = {
  id: PropTypes.number,
  callApi: PropTypes.func,
  hideModal: PropTypes.func,
};

export default EditSubscriber;
