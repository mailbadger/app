import React, { useState, useContext } from "react";
import PropTypes from "prop-types";
import { parseISO, formatRelative } from "date-fns";
import { More, Add, FormPreviousLink, FormNextLink } from "grommet-icons";
import { mainInstance as axios } from "../axios";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
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
} from "grommet";
import history from "../history";
import StyledTable from "../ui/StyledTable";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import PlaceholderTable from "../ui/PlaceholderTable";
import Modal from "../ui/Modal";
import { NotificationsContext } from "../Notifications/context";

const Row = ({ segment, setShowDelete }) => {
  const ca = parseISO(segment.created_at);
  const ua = parseISO(segment.updated_at);
  return (
    <TableRow>
      <TableCell scope="row" size="xlarge">
        <strong>{segment.name}</strong>
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
                  history.push(`/dashboard/segments/${segment.id}/edit`);
                  break;
                case "Delete":
                  setShowDelete({
                    show: true,
                    name: segment.name,
                    id: segment.id,
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
  segment: PropTypes.shape({
    name: PropTypes.string,
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
        <strong>Name</strong>
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

const SegmentTable = React.memo(({ list, setShowDelete }) => (
  <StyledTable>
    <Header />
    <TableBody>
      {list.map((s) => (
        <Row segment={s} key={s.id} setShowDelete={setShowDelete} />
      ))}
    </TableBody>
  </StyledTable>
));

SegmentTable.displayName = "SegmentTable";
SegmentTable.propTypes = {
  list: PropTypes.array,
  setShowDelete: PropTypes.func,
};

const segmentValidation = object().shape({
  name: string().required("Please enter a segment name."),
});

const CreateForm = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  hideModal,
}) => (
  <Box
    direction="column"
    fill
    margin={{ left: "medium", right: "medium", bottom: "medium" }}
  >
    <form onSubmit={handleSubmit}>
      <Box>
        <FormField htmlFor="name" label="Segment Name">
          <TextInput
            name="name"
            onChange={handleChange}
            placeholder="My segment"
          />
          <ErrorMessage name="name" />
        </FormField>
        <Box direction="row" alignSelf="end" margin={{ top: "medium" }}>
          <Box margin={{ right: "small" }}>
            <Button label="Cancel" onClick={() => hideModal()} />
          </Box>
          <Box>
            <ButtonWithLoader
              type="submit"
              primary
              disabled={isSubmitting}
              label="Save Segment"
            />
          </Box>
        </Box>
      </Box>
    </form>
  </Box>
);

CreateForm.propTypes = {
  hideModal: PropTypes.func,
  handleSubmit: PropTypes.func,
  handleChange: PropTypes.func,
  isSubmitting: PropTypes.bool,
};

const CreateSegment = ({ callApi, hideModal }) => {
  const { createNotification } = useContext(NotificationsContext);

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const postForm = async () => {
      try {
        await axios.post(
          "/api/segments",
          qs.stringify({
            name: values.name,
          })
        );
        createNotification("Segment has been created successfully.");

        await callApi({ url: "/api/segments" });

        //done submitting, set submitting to false
        setSubmitting(false);

        hideModal();
      } catch (error) {
        if (error.response) {
          const { message, errors } = error.response.data;

          setErrors(errors);

          const msg = message
            ? message
            : "Unable to create segment. Please try again.";

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
        initialValues={{ name: "" }}
        onSubmit={handleSubmit}
        validationSchema={segmentValidation}
      >
        {(props) => <CreateForm {...props} hideModal={hideModal} />}
      </Formik>
    </Box>
  );
};

CreateSegment.propTypes = {
  callApi: PropTypes.func,
  hideModal: PropTypes.func,
};

const DeleteForm = ({ id, callApi, hideModal }) => {
  const deleteSegment = async (id) => {
    await axios.delete(`/api/segments/${id}`);
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
            await deleteSegment(id);
            await callApi({ url: "/api/segments" });
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
      url: "/api/segments",
    },
    {
      collection: [],
      init: true,
    }
  );

  let table = null;
  if (state.isLoading) {
    table = <PlaceholderTable header={Header} numCols={3} numRows={8} />;
  } else if (state.data.collection.length > 0) {
    table = (
      <SegmentTable
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
          title={`Delete segment ${showDelete.name} ?`}
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
          title={`Create segment`}
          hideModal={() => openCreateModal(false)}
          form={
            <CreateSegment
              callApi={callApi}
              hideModal={() => openCreateModal(false)}
            />
          }
        />
      )}
      <Box gridArea="nav" direction="row">
        <Box>
          <Heading level="2" margin={{ bottom: "xsmall" }}>
            Segments
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
              <Heading level="2">Create your first segment.</Heading>
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
