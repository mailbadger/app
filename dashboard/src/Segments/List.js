import React, { useState, useContext, Fragment } from "react";
import PropTypes from "prop-types";
import { parseISO, formatRelative } from "date-fns";
import { More, Add } from "grommet-icons";
import { mainInstance as axios } from "../axios";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import qs from "qs";

import { useApi } from "../hooks";
import {
  TableBody,
  TableRow,
  TableCell,
  Box,
  Button,
  Heading,
  Select,
  FormField,
  TextInput,
  ResponsiveContext,
  } from "grommet";
import history from "../history";
import {
  StyledTable,
  ButtonWithLoader,
  Modal,
  AnchorLink,
} from "../ui";
import { NotificationsContext } from "../Notifications/context";
import DeleteSegment from "./Delete";
import { DashboardDataTable } from "../ui/DashboardDataTable";
import DashboardPlaceholderTable from "../ui/DashboardPlaceholderTable";
import {
  StyledHeaderWrapper,
  StyledHeaderButtons,
  StyledHeaderTitle,
  StyledHeaderButton,
  StyledActions,
} from "../Subscribers/StyledSections";
import { StyledTableHeader } from "../ui/DashboardStyledTable";

const Row = ({ segment, setShowDelete }) => {
  const ca = parseISO(segment.created_at);
  const ua = parseISO(segment.updated_at);
  return (
    <TableRow>
      <TableCell scope="row" size="medium">
        <AnchorLink
          size="medium"
          fontWeight="bold"
          to={`/dashboard/segments/${segment.id}`}
          label={segment.name}
        />
      </TableCell>
      <TableCell scope="row" size="medium">
        <strong>{segment.subscribers_in_segment}</strong>
      </TableCell>
      <TableCell scope="row" size="medium">
        {formatRelative(ca, new Date())}
      </TableCell>
      <TableCell scope="row" size="medium">
        {formatRelative(ua, new Date())}
      </TableCell>
      <TableCell scope="row" size="xsmall" align="end">
        <StyledActions>
        <Select
          alignSelf="center"
          plain
          icon={<More />}
          options={["View", "Delete"]}
          onChange={({ option }) => {
            (function () {
              switch (option) {
                case "View":
                  history.push(`/dashboard/segments/${segment.id}`);
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
        </StyledActions>
      </TableCell>
    </TableRow>
  );
};

Row.propTypes = {
  segment: PropTypes.shape({
    name: PropTypes.string,
    id: PropTypes.number,
    subscribers_in_segment: PropTypes.number,
    created_at: PropTypes.string,
    updated_at: PropTypes.string,
  }),
  setShowDelete: PropTypes.func,
};

const Header = () => (
  <StyledTableHeader>
    <TableRow>
      <TableCell scope="col" border="bottom" size="small">
        <strong>Name</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="small">
        <strong>Total Subscribers</strong>
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
  </StyledTableHeader>
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
  name: string()
    .required("Please enter a segment name.")
    .max(191, "The name must not exceed 191 characters."),
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

const getData = (segmentsData, setShowDelete) => {
  const data = [];

  for (let i = 0; i < segmentsData.length; i += 1) {
    const { name, subscribers_in_segment, created_at, updated_at, id } = segmentsData[i];

    const dateCreatedAt = new Date(created_at);
    const dateUpdatedAt = parseISO(updated_at);

    data.push({
      name,
      subscribers_in_segment,
      created: dateCreatedAt.toLocaleDateString("en-US"),
      updated: formatRelative(dateUpdatedAt, new Date()),
      actions: (
        <StyledActions>
          <Select
            alignSelf="center"
            plain
            defaultValue="View"
            icon={<More />}
            options={["View","Delete"]}
            onChange={({ option }) => {
              (() => {
                switch (option) {
                  case "View":
                    history.push(`/dashboard/segments/${id}`);
                    break;
                  case "Delete":
                    setShowDelete({
                      show: true,
                      name,
                      id,
                    });
                    break;
                  default:
                    return "null";
                }
              })();
            }}
          />
        </StyledActions>
      ),
    });
  }
  return data;
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
  const contextSize = useContext(ResponsiveContext);

  const data = getData(
    state.data.collection,
    setShowDelete
  );
   const columns = [
    { property: "name", header: "Name", size: "small" },
    { property: "subscribers_in_segment", header: "Subscribers in Group", size: "small" },
    { property: "created", header: "Created At", size: "small" },
    { property: "updated", header: "Updated At", size: "small" },
    { property: "actions", header: "Actions", size: "small", align: "center" },
  ];

  let table = null;
  if (state.isLoading) {
    table = (
      <DashboardPlaceholderTable
      columns={columns}
      numCols={columns.length}
      numRows={10}
    />
    
    );
  } else if (state.data && state.data.collection.length > 0) {
    table = (
      <DashboardDataTable
        columns={columns}
        data={data}
        isLoading={state.isLoading}
        setShowDelete={setShowDelete}     
        prevLinks={state.data.links.previous}
        nextLinks={state.data.links.next}
      />
    );
  }

  return (
    <>
      {showDelete.show && (
        <Modal
          title={`Delete segment ${showDelete.name} ?`}
          hideModal={hideModal}
          form={
            <DeleteSegment
              id={showDelete.id}
              onSuccess={async () => {
                await callApi({ url: "/api/segments" });
                hideModal();
              }}
              onCancel={hideModal}
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

    <StyledHeaderWrapper
        size={contextSize}
        gridArea="nav"
        margin={{ left: "40px", right: "100px", bottom: "22px", top: "40px" }}
      >
        <StyledHeaderTitle size={contextSize}>Segments</StyledHeaderTitle>
        <StyledHeaderButtons size={contextSize} margin={{ left: "auto" }}>
          <Fragment>
           
            <StyledHeaderButton
              width="154"
              margin={{ right: "small" }}
              label="Create New"
              color="status-ok"
              icon={<Add />}
              onClick={() => openCreateModal(true)}
            />
          </Fragment>
        </StyledHeaderButtons>
      </StyledHeaderWrapper>
      <Box gridArea="main">
        <Box animation="fadeIn">
          {table}
          {!state.isLoading && !state.isError && state.data.collection.length === 0 ? (
            <Box align="center" margin={{ top: "large" }}>
              <Heading level="2">Create your first segment.</Heading>
            </Box>
          ) : null}
        </Box>
      </Box>


      {/* <Box gridArea="nav" direction="row" border={{ side: 'bottom', color: 'light-4' }}>
        <Box margin={{ right: "small" }} alignSelf="center">
          <Heading level="2">Segments</Heading>
        </Box>
        <Box alignSelf="center">
          <Button
            primary
            color="status-ok"
            label="Create new"
            icon={<Add />}
            reverse
            onClick={() => openCreateModal(true)}
          />
        </Box>
      </Box> */}
      {/* <Box gridArea="main">
        <Box animation="fadeIn">
          {table}

          {!state.isLoading && !state.isError && state.data.collection.length === 0 ? (
            <Box align="center" margin={{ top: "small" }}>
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
      </Box> */}
    </>
  );
};

export default List;
