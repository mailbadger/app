import React, { useState, useContext } from "react";
import PropTypes from "prop-types";
import { parseISO } from "date-fns";
import { More, Add } from "grommet-icons";
import axios from "axios";
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
  Layer,
  Heading,
  Select,
  FormField,
  TextInput
} from "grommet";
import history from "../history";
import StyledTable from "../ui/StyledTable";
import StyledButton from "../ui/StyledButton";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import PlaceholderRow from "../ui/PlaceholderRow";
import { NotificationsContext } from "../Notifications/context";

const Row = ({ segment, setShowDelete }) => {
  const res = parseISO(segment.created_at);
  return (
    <TableRow>
      <TableCell scope="row" size="xlarge">
        {segment.name}
      </TableCell>
      <TableCell scope="row" size="medium">
        {res.toUTCString()}
      </TableCell>
      <TableCell scope="row" size="xsmall">
        <Select
          alignSelf="center"
          plain
          icon={<More />}
          options={["Edit", "Delete"]}
          onChange={({ option }) => {
            (function() {
              switch (option) {
                case "Edit":
                  history.push(`/dashboard/segments/${segment.id}/edit`);
                  break;
                case "Delete":
                  setShowDelete({
                    show: true,
                    name: segment.name,
                    id: segment.id
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
    created_at: PropTypes.string
  }),
  setShowDelete: PropTypes.func
};

const Header = () => (
  <TableHeader>
    <TableRow>
      <TableCell scope="col" border="bottom" size="medium">
        <strong>Name</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="medium">
        <strong>Date</strong>
      </TableCell>
      <TableCell
        style={{ textAlign: "right" }}
        align="end"
        scope="col"
        border="bottom"
        size="small"
      />
    </TableRow>
  </TableHeader>
);
const PlaceholderTable = () => (
  <StyledTable caption="Segments">
    <Header />
    <TableBody>
      <PlaceholderRow columns={3} />
      <PlaceholderRow columns={3} />
      <PlaceholderRow columns={3} />
      <PlaceholderRow columns={3} />
      <PlaceholderRow columns={3} />
    </TableBody>
  </StyledTable>
);

const SegmentTable = React.memo(({ list, setShowDelete }) => (
  <StyledTable caption="Segments">
    <Header />
    <TableBody>
      {list.map(s => (
        <Row segment={s} key={s.id} setShowDelete={setShowDelete} />
      ))}
    </TableBody>
  </StyledTable>
));

SegmentTable.displayName = "SegmentTable";
SegmentTable.propTypes = {
  list: PropTypes.array,
  setShowDelete: PropTypes.func
};

const segmentValidation = object().shape({
  name: string().required("Please enter a segment name.")
});

const CreateForm = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  hideModal
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
  isSubmitting: PropTypes.bool
};

const CreateSegment = ({ callApi, hideModal }) => {
  const { createNotification } = useContext(NotificationsContext);

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const postForm = async () => {
      try {
        await axios.post(
          "/api/segments",
          qs.stringify({
            name: values.name
          })
        );
        createNotification("Segment has been created successfully.");

        await callApi({ url: "/api/segments" });

        hideModal();
      } catch (error) {
        if (error.response) {
          const { message, errors } = error.response.data;

          setErrors(errors);

          const msg = message
            ? message
            : "Unable to create segment. Please try again.";

          createNotification(msg, "status-error");
        }
      }
    };

    await postForm();

    //done submitting, set submitting to false
    setSubmitting(false);

    return;
  };

  return (
    <Box direction="row">
      <Formik
        onSubmit={handleSubmit}
        validationSchema={segmentValidation}
        render={props => <CreateForm {...props} hideModal={hideModal} />}
      />
    </Box>
  );
};

CreateSegment.propTypes = {
  callApi: PropTypes.func,
  hideModal: PropTypes.func
};

const DeleteForm = ({ id, callApi, hideModal }) => {
  const deleteSegment = async id => {
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
  id: PropTypes.string,
  callApi: PropTypes.func,
  hideModal: PropTypes.func
};

const Modal = ({ hideModal, title, form }) => {
  return (
    <Layer onEsc={() => hideModal()} onClickOutside={() => hideModal()}>
      <Box width="30em">
        <Heading margin="small" level="3">
          {title}
        </Heading>
        {form}
      </Box>
    </Layer>
  );
};

Modal.propTypes = {
  hideModal: PropTypes.func,
  title: PropTypes.string,
  form: PropTypes.element.isRequired
};

const List = () => {
  const [showDelete, setShowDelete] = useState({ show: false, name: "" });
  const [showCreate, openCreateModal] = useState(false);
  const hideModal = () => setShowDelete({ show: false, name: "", id: "" });

  const [state, callApi] = useApi(
    {
      url: "/api/segments"
    },
    {
      collection: [],
      init: true
    }
  );

  let table = null;
  if (state.isLoading) {
    table = <PlaceholderTable />;
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
        { name: "main", start: [0, 1], end: [1, 1] }
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
          <ButtonWithLoader
            label="Create new"
            icon={<Add color="#ffffff" />}
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
              <StyledButton
                label="Previous"
                onClick={() => console.log("Previous")}
              />
            </Box>
            <Box>
              <StyledButton label="Next" onClick={() => console.log("Next")} />
            </Box>
          </Box>
        ) : null}
      </Box>
    </Grid>
  );
};

export default List;
