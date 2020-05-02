import React, { useState } from "react";
import PropTypes from "prop-types";
import { parseISO, formatRelative } from "date-fns";
import {
  More,
  Add,
  UserAdd,
  Upload,
  SubtractCircle,
  FormPreviousLink,
  FormNextLink,
} from "grommet-icons";
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
} from "grommet";

import { useApi } from "../hooks";
import { StyledTable, PlaceholderTable, Modal, SecondaryButton } from "../ui";
import CreateSubscriber from "./Create";
import DeleteSubscriber from "./Delete";
import EditSubscriber from "./Edit";

const Row = ({ subscriber, setShowDelete, setShowEdit }) => {
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
                  setShowEdit({
                    show: true,
                    id: subscriber.id,
                  });
                  break;
                case "Delete":
                  setShowDelete({
                    show: true,
                    email: subscriber.email,
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
  setShowEdit: PropTypes.func,
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

const SubscriberTable = React.memo(({ list, setShowDelete, setShowEdit }) => (
  <StyledTable>
    <Header />
    <TableBody>
      {list.map((s) => (
        <Row
          subscriber={s}
          key={s.id}
          setShowDelete={setShowDelete}
          setShowEdit={setShowEdit}
        />
      ))}
    </TableBody>
  </StyledTable>
));

SubscriberTable.displayName = "SubscriberTable";
SubscriberTable.propTypes = {
  list: PropTypes.array,
  setShowDelete: PropTypes.func,
  setShowEdit: PropTypes.func,
};

const ActionButtons = () => (
  <>
    <SecondaryButton
      margin={{ right: "small" }}
      icon={<UserAdd size="20px" />}
      label="Import from file"
    />
    <SecondaryButton
      margin={{ right: "small" }}
      icon={<SubtractCircle size="20px" />}
      label="Delete from file"
    />
    <SecondaryButton icon={<Upload size="20px" />} label="Export" />
  </>
);

const List = () => {
  const [showDelete, setShowDelete] = useState({
    show: false,
    email: "",
    id: "",
  });
  const [showEdit, setShowEdit] = useState({ show: false, id: "" });
  const [showCreate, openCreateModal] = useState(false);
  const hideDeleteModal = () =>
    setShowDelete({ show: false, email: "", id: "" });
  const hideEditModal = () => setShowEdit({ show: false, id: "" });

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
        setShowEdit={setShowEdit}
      />
    );
  }

  return (
    <Grid
      rows={["fill", "fill"]}
      columns={["1fr", "1fr"]}
      gap="small"
      margin="medium"
      areas={[
        { name: "nav", start: [0, 0], end: [1, 1] },
        { name: "main", start: [0, 1], end: [1, 1] },
      ]}
    >
      {showDelete.show && (
        <Modal
          title={`Delete subscriber ${showDelete.email} ?`}
          hideModal={hideDeleteModal}
          form={
            <DeleteSubscriber
              id={showDelete.id}
              callApi={callApi}
              hideModal={hideDeleteModal}
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
      {showEdit.show && (
        <Modal
          title={`Edit subscriber`}
          hideModal={hideEditModal}
          form={
            <EditSubscriber
              id={showEdit.id}
              callApi={callApi}
              hideModal={hideEditModal}
            />
          }
        />
      )}
      <Box gridArea="nav" direction="row">
        <Box alignSelf="center" margin={{ right: "small" }}>
          <Heading level="2">Subscribers</Heading>
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
        <Box margin={{ left: "auto" }} alignSelf="center" direction="row">
          <ActionButtons />
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
