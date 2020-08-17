import React, { useState, useReducer, useEffect, useContext } from "react";
import PropTypes from "prop-types";
import {
  Box,
  Heading,
  Text,
  Button,
  FormField,
  TextInput,
  Select,
  ResponsiveContext,
  Grid,
} from "grommet";
import { Edit, Trash, Add, Send, LinkPrevious } from "grommet-icons";
import { Formik, ErrorMessage, FieldArray } from "formik";
import { string, object, array } from "yup";

import { useApi } from "../hooks";
import history from "../history";
import {
  LoadingOverlay,
  SecondaryButton,
  Notice,
  ButtonWithLoader,
} from "../ui";

const sendValidation = object().shape({
  source: string()
    .email("Email is not valid.")
    .required("Please enter a valid email address."),
  from_name: string()
    .required("Please enter the sender's name.")
    .max(191, "The name must not exceed 191 characters."),
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
  segments: array().min(1, "Please choose atleast one segment."),
});

const DetailsGrid = ({ children }) => {
  const size = useContext(ResponsiveContext);

  let cols = ["small", "small", "large", "xsmall"];
  let areas = [
    ["title", "title", "title", "title"],
    ["info", "info", "main", "main"],
    ["info", "info", "main", "main"],
  ];

  if (size === "medium") {
    cols = ["120px", "240px", "600px", "xsmall"];
    areas = [
      ["title", "title", "title", "."],
      ["info", "info", "main", "main"],
      ["info", "info", "main", "main"],
    ];
  }

  return (
    <Grid
      rows={["xsmall", "1fr", "1fr"]}
      columns={cols}
      margin="medium"
      gap="small"
      areas={areas}
    >
      {children}
    </Grid>
  );
};

DetailsGrid.displayName = "DetailsGrid";
DetailsGrid.propTypes = {
  children: PropTypes.element,
};

const Form = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  values,
  setFieldValue,
}) => {
  const [selected, setSelected] = useState(values.segments);
  const [segments, callApi] = useApi(
    {
      url: `/api/segments?per_page=40`,
    },
    {
      collection: [],
      links: {
        next: null,
      },
    }
  );

  const reducer = (segments) => (state, action) => {
    let col = [];

    switch (action.type) {
      case "append":
        // filter out the preselected segments to avoid duplicate items
        for (let i = 0; i < action.payload.length; i++) {
          let found = false;
          for (let j = 0; j < segments.length; j++) {
            if (segments[j].id === action.payload[i].id) {
              found = true;
              break;
            }
          }

          if (!found) {
            col.push(action.payload[i]);
          }
        }

        return [...state, ...col];
      default:
        throw new Error("invalid action type.");
    }
  };

  const [options, dispatch] = useReducer(
    reducer(values.segments),
    values.segments
  );

  useEffect(() => {
    if (segments.isError || segments.isLoading) {
      return;
    }

    dispatch({ type: "append", payload: segments.data.collection });
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
    setSelected(nextSelected);
    setFieldValue("segments", nextSelected);
  };

  return (
    <Box
      round={{ corner: "top", size: "small" }}
      background="white"
      pad={{ horizontal: "medium", top: "small", bottom: "medium" }}
      margin={{ bottom: "small" }}
    >
      <form onSubmit={handleSubmit}>
        <Box>
          <FormField htmlFor="from_name" label="From Name">
            <TextInput
              name="from_name"
              placeholder="Acme Inc."
              value={values.from_name}
              onChange={handleChange}
            />
            <ErrorMessage name="from_name" />
          </FormField>
          <FormField htmlFor="source" label="From Email">
            <TextInput
              name="source"
              placeholder="john@example.com"
              value={values.source}
              onChange={handleChange}
            />
            <ErrorMessage name="source" />
          </FormField>
          <FormField htmlFor="segments" label="Choose subscribers">
            <Select
              multiple
              name="segments"
              closeOnChange={false}
              placeholder="select a segment..."
              value={selected}
              labelKey="name"
              valueKey="id"
              options={options}
              dropHeight="medium"
              onMore={onMore}
              onChange={onChange}
            />
            <ErrorMessage name="segments" />
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
                    <Text>Add default field</Text>
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
            <Box>
              <Button
                alignSelf="start"
                onClick={() => history.push("/dashboard/campaigns")}
              >
                <Box
                  direction="row"
                  pad={{ vertical: "xsmall" }}
                  align="center"
                >
                  <LinkPrevious size="20px" />
                  <Text margin={{ left: "xsmall" }}>Back</Text>
                </Box>
              </Button>
            </Box>
            <Box margin={{ left: "small" }}>
              <ButtonWithLoader
                type="submit"
                primary
                disabled={isSubmitting}
                label="Send"
                icon={<Send />}
                color="status-ok"
              />
            </Box>
          </Box>
        </Box>
      </form>
    </Box>
  );
};

Form.propTypes = {
  handleSubmit: PropTypes.func,
  handleChange: PropTypes.func,
  isSubmitting: PropTypes.bool,
  setFieldValue: PropTypes.func,
  values: PropTypes.shape({
    from_name: PropTypes.string,
    source: PropTypes.string,
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

const SendCampaign = ({ match }) => {
  const [showEdit, setShowEdit] = useState(false);
  const [showDelete, setShowDelete] = useState(false);

  const [campaign] = useApi({
    url: `/api/campaigns/${match.params.id}`,
  });

  if (campaign.isLoading) {
    return <LoadingOverlay />;
  }

  if (campaign.isError) {
    return (
      <Box margin="15%" alignSelf="center">
        <Heading>Campaign not found.</Heading>
      </Box>
    );
  }

  return (
    <DetailsGrid>
      {campaign && campaign.data && (
        <>
          <Box gridArea="title" direction="row">
            <Heading level="2" alignSelf="center">
              {campaign.data.name}
            </Heading>
            <Box direction="row" margin={{ left: "auto" }}>
              <SecondaryButton
                margin={{ right: "small" }}
                a11yTitle="edit campaign"
                alignSelf="center"
                icon={<Edit a11yTitle="edit campaign" color="dark-1" />}
                label="Edit"
                onClick={() => setShowEdit(true)}
              />
              <SecondaryButton
                a11yTitle="delete campaign"
                alignSelf="center"
                icon={<Trash a11yTitle="delete campaign" color="dark-1" />}
                label="Delete"
                onClick={() => setShowDelete(true)}
              />
            </Box>
          </Box>
          <Box gridArea="info" direction="column" alignSelf="start">
            <Formik
              validationSchema={sendValidation}
              initialValues={{
                from_name: "",
                source: "",
                segments: [],
              }}
            >
              {Form}
            </Formik>
            <Notice
              message={`
                        Default fields are used as a fallback when replacing the template's
                        variables before sending the newsletter.
                        `}
              status="status-info"
              color="white"
              borderColor="neutral-2"
            />
          </Box>
          <Box gridArea="main" margin={{ left: "small" }}>
            Sup
          </Box>
        </>
      )}
    </DetailsGrid>
  );
};

SendCampaign.propTypes = {
  match: PropTypes.shape({
    params: PropTypes.shape({
      id: PropTypes.string,
    }),
  }),
};

export default SendCampaign;
