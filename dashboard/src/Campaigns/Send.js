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
import { Redirect } from "react-router-dom";
import DOMPurify from "dompurify";
import qs from "qs";

import { useApi } from "../hooks";
import history from "../history";
import {
  LoadingOverlay,
  SecondaryButton,
  Notice,
  ButtonWithLoader,
  Modal,
  Badge,
  AnchorLink,
} from "../ui";
import EditCampaign from "./Edit";
import DeleteCampaign from "./Delete";
import { mainInstance as axios } from "../axios";
import { NotificationsContext } from "../Notifications/context";
import { endpoints } from "../network/endpoints";

const sendValidation = object().shape({
  source: string()
    .email("Email is not valid.")
    .required("Please enter a valid email address."),
  from_name: string()
    .required("Please enter the sender's name.")
    .max(191, "The name must not exceed 191 characters."),
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
  segments: array().min(1, "Please choose atleast one group."),
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
  quota,
}) => {
  const [totalSelectedSubs, setTotalSubs] = useState(0);
  const [selected, setSelected] = useState(values.segments);
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
    let totalSubs = 0;
    for (let i = 0; i < nextSelected.length; i++) {
      totalSubs += nextSelected[i].subscribers_in_segment;
    }

    setTotalSubs(totalSubs);
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
              placeholder="Select a group..."
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
          <Box margin={{ top: "large", bottom: "small" }}>
            <Text alignSelf="end">Total recipients: {totalSelectedSubs}</Text>
            <Text alignSelf="end">
              Send quota:{" "}
              {quota &&
                quota.data &&
                quota.data.max_24_hour_send - quota.data.sent_last_24_hours}
            </Text>
          </Box>
          <Box direction="row" alignSelf="end">
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
  quota: PropTypes.shape({
    isLoading: PropTypes.bool,
    isError: PropTypes.bool,
    data: PropTypes.shape({
      max_24_hour_send: PropTypes.number,
      sent_last_24_hours: PropTypes.number,
    }),
  }),
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

const handleSubmit = (id, setSuccess, createNotification) => async (
  values,
  { setSubmitting, setErrors }
) => {
  const postForm = async () => {
    try {
      let data = {
        from_name: values.from_name,
        source: values.source,
      };

      if (values.metadata.length > 0) {
        data.default_template_data = values.metadata.reduce((map, meta) => {
          map[meta.key] = meta.val;
          return map;
        }, {});
      }

      if (values.segments.length > 0) {
        data.segment_id = values.segments.map((s) => s.id);
      }

      await axios.post(
        endpoints.postCampaignsStart,
        qs.stringify(data, { arrayFormat: "brackets" })
      );
      createNotification(
        "The campaign has started. We will begin queueing e-mails shortly."
      );

      setSuccess(true);
    } catch (error) {
      if (error.response) {
        const { message, errors } = error.response.data;

        setErrors(errors);

        const msg = message
          ? message
          : "Unable to send campaign. Please try again.";

        createNotification(msg, "status-error");
      }
    }
  };

  await postForm();

  setSubmitting(false);
};

const PreviewTemplate = React.memo(({ name }) => {
  const [template] = useApi({
    url: `${endpoints.getTemplates}/${name}`,
  });

  if (template.isLoading) {
    return <LoadingOverlay />;
  }

  return (
    <Box direction="column">
      <Box direction="row" margin={{ bottom: "small" }}>
        <Box direction="column" align="start">
          <Box>
            <Text>
              Name{" "}
              <Badge color="#00b4d8">
                {template && template.data && template.data.name}
              </Badge>
            </Text>
          </Box>
          <Box margin={{ top: "xsmall" }}>
            <Text>
              Subject{" "}
              {template && template.data && (
                <Badge color="#00b4d8">{template.data.subject_part}</Badge>
              )}
            </Text>
          </Box>
        </Box>
        <Box margin={{ left: "auto", top: "auto" }}>
          {template && template.data && (
            <AnchorLink
              size="medium"
              to={`/dashboard/templates/${template.data.name}/edit`}
            >
              Edit template <Edit fontWeight="bold" size="18px" />
            </AnchorLink>
          )}
        </Box>
      </Box>
      {template && template.data && (
        <Box elevation="small">
          <iframe
            frameBorder="0"
            height="550px"
            title="preview-template"
            srcDoc={DOMPurify.sanitize(template.data.html_part, {
              USE_PROFILES: { html: true },
            })}
          />
        </Box>
      )}
    </Box>
  );
});

PreviewTemplate.propTypes = {
  name: PropTypes.string.isRequired,
};

const SendCampaign = ({ match }) => {
  const { createNotification } = useContext(NotificationsContext);
  const [success, setSuccess] = useState(false);
  const [showEdit, setShowEdit] = useState(false);
  const [showDelete, setShowDelete] = useState(false);

  const [quota] = useApi({
    url: "/api/ses/quota",
  });

  const [campaign, callApi] = useApi({
    url: endpoints.getCampaigns(match.params.id),
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

  if (
    success ||
    (campaign && campaign.data && campaign.data.status !== "draft")
  ) {
    return <Redirect to={`/dashboard/campaigns/${match.params.id}/report`} />;
  }

  return (
    <DetailsGrid>
      {campaign && campaign.data && (
        <>
          {showEdit && (
            <Modal
              title={`Edit group`}
              hideModal={() => setShowEdit(false)}
              form={
                <EditCampaign
                  id={campaign.data.id}
                  hideModal={() => setShowEdit(false)}
                  onSuccess={() =>
                    callApi({ url: endpoints.getCampaign(campaign.data.id) })
                  }
                />
              }
            />
          )}
          {showDelete && (
            <Modal
              title={`Delete campaign ${campaign.data.name} ?`}
              hideModal={() => setShowDelete(false)}
              form={
                <DeleteCampaign
                  id={campaign.data.id}
                  onSuccess={() => history.replace("/dashboard/campaigns")}
                  hideModal={() => setShowDelete(false)}
                />
              }
            />
          )}
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
              onSubmit={handleSubmit(
                campaign.data.id,
                setSuccess,
                createNotification
              )}
              initialValues={{
                from_name: "",
                source: "",
                segments: [],
                metadata: [],
              }}
            >
              {(props) => <Form {...props} quota={quota} />}
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
          <Box gridArea="main" margin={{ left: "medium" }}>
            <PreviewTemplate name={campaign.data.template_name} />
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
