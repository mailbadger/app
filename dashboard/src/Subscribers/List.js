import React, { useState, useContext, useEffect } from "react";
import PropTypes from "prop-types";
import { parseISO, formatRelative } from "date-fns";
import {
  More,
  Add,
  FormPreviousLink,
  FormNextLink,
  Trash,
} from "grommet-icons";
import { mainInstance as axios } from "../axios";
import { Formik, ErrorMessage, FieldArray } from "formik";
import { string, object, array } from "yup";
import qs from "qs";

import useApi from "../hooks/useApi";
import {
  Grid,
  TableHeader,
  TableBody,
  TableRow,
  TableCell,
  Box,
  Button,
  Heading,
  Select,
  FormField,
  TextInput,
  Text,
} from "grommet";
import history from "../history";
import StyledTable from "../ui/StyledTable";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import PlaceholderTable from "../ui/PlaceholderTable";
import Modal from "../ui/Modal";
import { NotificationsContext } from "../Notifications/context";

const Row = ({ subscriber, setShowDelete }) => {
  const ca = parseISO(subscriber.created_at);
  const ua = parseISO(subscriber.updated_at);
  return (
    <TableRow>
      <TableCell scope="row" size="xlarge">
        <strong>{subscriber.email}</strong>
      </TableCell>
      <TableCell scope="row" size="medium">
        {formatRelative(ca, new Date())}
      </TableCell>
      <TableCell scope="row" size="medium">
        {formatRelative(ua, new Date())}
      </TableCell>
      <TableCell scope="row" size="xxsmall" align="end">
        <Select
          alignSelf="center"
          plain
          icon={<More />}
          options={["Edit", "Delete"]}
          onChange={({ option }) => {
            (function () {
              switch (option) {
                case "Edit":
                  history.push(`/dashboard/subscribers/${subscriber.id}/edit`);
                  break;
                case "Delete":
                  setShowDelete({
                    show: true,
                    name: subscriber.email,
                    id: subscriber.id,
                  });
                  break;
                default:
                  return null;
              }
            })();
          }}
        />
      </TableCell>
    </TableRow>
  );
};

Row.propTypes = {
  subscriber: PropTypes.shape({
    email: PropTypes.string,
    id: PropTypes.number,
    created_at: PropTypes.string,
    updated_at: PropTypes.string,
  }),
  setShowDelete: PropTypes.func,
};

const Header = () => (
  <TableHeader>
    <TableRow>
      <TableCell scope="col" border="bottom" size="small">
        <strong>Email</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="small">
        <strong>Created At</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="small">
        <strong>Updated At</strong>
      </TableCell>
      <TableCell align="end" scope="col" border="bottom" size="small">
        <strong>Action</strong>
      </TableCell>
    </TableRow>
  </TableHeader>
);

const SubscriberTable = React.memo(({ list, setShowDelete }) => (
  <StyledTable>
    <Header />
    <TableBody>
      {list.map((s) => (
        <Row subscriber={s} key={s.id} setShowDelete={setShowDelete} />
      ))}
    </TableBody>
  </StyledTable>
));

SubscriberTable.displayName = "SubscriberTable";
SubscriberTable.propTypes = {
  list: PropTypes.array,
  setShowDelete: PropTypes.func,
};

const subscrValidation = object().shape({
  email: string().email().required("Please enter a subscriber email."),
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

const CreateForm = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  hideModal,
  values,
  setFieldValue,
}) => {
  const [selected, setSelected] = useState("");
  const [options, setOptions] = useState({
    collection: [],
    url: "/api/segments?per_page=40",
  });
  const callApi = async () => {
    const res = await axios(options.url);
    setOptions({
      collection: [...options.collection, ...res.data.collection],
      url: res.data.links.next,
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
          "/api/subscribers",
          qs.stringify(data, { arrayFormat: "brackets" })
        );
        createNotification("Subscriber has been created successfully.");

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

const DeleteForm = ({ id, callApi, hideModal }) => {
  const deleteSubscriber = async (id) => {
    await axios.delete(`/api/subscribers/${id}`);
  };

  const [isSubmitting, setSubmitting] = useState(false);
  return (
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
            await deleteSubscriber(id);
            await callApi({ url: "/api/subscribers" });
            setSubmitting(false);
            hideModal();
          }}
        />
      </Box>
    </Box>
  );
};

DeleteForm.propTypes = {
  id: PropTypes.number,
  callApi: PropTypes.func,
  hideModal: PropTypes.func,
};

const List = () => {
  const [showDelete, setShowDelete] = useState({ show: false, name: "" });
  const [showCreate, openCreateModal] = useState(false);
  const hideModal = () => setShowDelete({ show: false, name: "", id: "" });

  const [state, callApi] = useApi(
    {
      url: "/api/subscribers",
    },
    {
      collection: [],
      init: true,
    }
  );

  let table = null;
  if (state.isLoading) {
    table = <PlaceholderTable header={Header} numCols={3} numRows={3} />;
  } else if (state.data.collection.length > 0) {
    table = (
      <SubscriberTable
        isLoading={state.isLoading}
        list={state.data.collection}
        setShowDelete={setShowDelete}
      />
    );
  }

  return (
    <Grid
      rows={["fill", "fill"]}
      columns={["1fr", "1fr"]}
      gap="medium"
      margin="medium"
      areas={[
        { name: "nav", start: [0, 0], end: [0, 1] },
        { name: "main", start: [0, 1], end: [1, 1] },
      ]}
    >
      {showDelete.show && (
        <Modal
          title={`Delete subscriber ${showDelete.email} ?`}
          hideModal={hideModal}
          form={
            <DeleteForm
              id={showDelete.id}
              callApi={callApi}
              hideModal={hideModal}
            />
          }
        />
      )}
      {showCreate && (
        <Modal
          title={`Create subscriber`}
          hideModal={() => openCreateModal(false)}
          form={
            <CreateSubscriber
              callApi={callApi}
              hideModal={() => openCreateModal(false)}
            />
          }
        />
      )}
      <Box gridArea="nav" direction="row">
        <Box>
          <Heading level="2" margin={{ bottom: "xsmall" }}>
            Subscribers
          </Heading>
        </Box>
        <Box margin={{ left: "medium", top: "medium" }}>
          <Button
            primary
            color="status-ok"
            label="Create new"
            icon={<Add />}
            reverse
            onClick={() => openCreateModal(true)}
          />
        </Box>
      </Box>
      <Box gridArea="main">
        <Box animation="fadeIn">
          {table}

          {!state.isLoading && state.data.collection.length === 0 ? (
            <Box align="center" margin={{ top: "large" }}>
              <Heading level="2">Create your first subscriber.</Heading>
            </Box>
          ) : null}
        </Box>
        {!state.isLoading && state.data.collection.length > 0 ? (
          <Box direction="row" alignSelf="end" margin={{ top: "medium" }}>
            <Box margin={{ right: "small" }}>
              <Button
                icon={<FormPreviousLink />}
                label="Previous"
                disabled={state.data.links.previous === null}
                onClick={() => {
                  callApi({
                    url: state.data.links.previous,
                  });
                }}
              />
            </Box>
            <Box>
              <Button
                icon={<FormNextLink />}
                reverse
                label="Next"
                disabled={state.data.links.next === null}
                onClick={() => {
                  callApi({
                    url: state.data.links.next,
                  });
                }}
              />
            </Box>
          </Box>
        ) : null}
      </Box>
    </Grid>
  );
};

export default List;
