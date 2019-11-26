import React, { useState } from "react";
import PropTypes from "prop-types";
import { parseISO, formatRelative } from "date-fns";
import { More, Add, FormPreviousLink, FormNextLink } from "grommet-icons";
import axios from "axios";
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
  Select
} from "grommet";
import history from "../history";
import StyledTable from "../ui/StyledTable";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import PlaceholderTable from "../ui/PlaceholderTable";
import Modal from "../ui/Modal";

const Row = ({ template, setShowDelete }) => {
  const d = parseISO(template.timestamp);
  return (
    <TableRow>
      <TableCell scope="row" size="xlarge">
        <strong>{template.name}</strong>
      </TableCell>
      <TableCell scope="row" size="medium">
        {formatRelative(d, new Date())}
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
                  history.push(`/dashboard/templates/${template.name}/edit`);
                  break;
                case "Delete":
                  setShowDelete({ show: true, name: template.name });
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
  setShowDelete: PropTypes.func,
  template: PropTypes.shape({
    name: PropTypes.string,
    timestamp: PropTypes.string
  })
};

const Header = () => (
  <TableHeader>
    <TableRow>
      <TableCell scope="col" border="bottom" size="medium">
        <strong>Name</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="medium">
        <strong>Created At</strong>
      </TableCell>
      <TableCell align="end" scope="col" border="bottom" size="small">
        <strong>Action</strong>
      </TableCell>
    </TableRow>
  </TableHeader>
);

const TemplateTable = React.memo(({ list, setShowDelete }) => (
  <StyledTable>
    <Header />
    <TableBody>
      {list.map(t => (
        <Row template={t} key={t.name} setShowDelete={setShowDelete} />
      ))}
    </TableBody>
  </StyledTable>
));

TemplateTable.displayName = "TemplateTable";
TemplateTable.propTypes = {
  list: PropTypes.array,
  setShowDelete: PropTypes.func
};

const DeleteForm = ({ name, callApi, hideModal }) => {
  const deleteTemplate = async name => {
    await axios.delete(`/api/templates/${name}`);
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
            await deleteTemplate(name);
            await callApi({ url: "/api/templates" });
            setSubmitting(false);
            hideModal();
          }}
        />
      </Box>
    </Box>
  );
};

DeleteForm.propTypes = {
  name: PropTypes.string,
  callApi: PropTypes.func,
  hideModal: PropTypes.func
};

const List = () => {
  const [showDelete, setShowDelete] = useState({ show: false, name: "" });
  const [currentPage, setPage] = useState({ current: -1, tokens: [""] });
  const hideModal = () => setShowDelete({ show: false, name: "" });

  const [state, callApi] = useApi(
    {
      url: "/api/templates"
    },
    {
      next_token: "",
      collection: [],
      init: true
    }
  );

  let table = null;
  if (state.isLoading) {
    table = <PlaceholderTable header={Header} numCols={3} numRows={5} />;
  } else if (state.data.collection.length > 0) {
    table = (
      <TemplateTable
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
          title={`Delete template ${showDelete.name} ?`}
          hideModal={hideModal}
          form={
            <DeleteForm
              id={showDelete.name}
              callApi={callApi}
              hideModal={hideModal}
            />
          }
        />
      )}
      <Box gridArea="nav" direction="row">
        <Box>
          <Heading level="2" margin={{ bottom: "xsmall" }}>
            Templates
          </Heading>
        </Box>
        <Box margin={{ left: "medium", top: "medium" }}>
          <Button
            primary
            color="status-ok"
            label="Create new"
            icon={<Add />}
            reverse
            onClick={() => history.push("/dashboard/templates/new")}
          />
        </Box>
      </Box>
      <Box gridArea="main">
        <Box animation="fadeIn">
          {table}

          {!state.isLoading && state.data.collection.length === 0 ? (
            <Box align="center" margin={{ top: "large" }}>
              <Heading level="2">Create your first template.</Heading>
            </Box>
          ) : null}
        </Box>
        {!state.isLoading && state.data.collection.length > 0 ? (
          <Box direction="row" alignSelf="end" margin={{ top: "medium" }}>
            <Box margin={{ right: "small" }}>
              <Button
                icon={<FormPreviousLink />}
                label="Previous"
                onClick={() => {
                  const t = currentPage.tokens[currentPage.current];
                  callApi({
                    url: `/api/templates?next_token=${encodeURIComponent(t)}`
                  });
                  const removeNumOfTokens = currentPage.current > 0 ? 2 : 1;
                  currentPage.tokens.splice(-1, removeNumOfTokens);

                  setPage({
                    current: currentPage.current - 1,
                    tokens: currentPage.tokens
                  });
                }}
                disabled={currentPage.current === -1}
              />
            </Box>
            <Box>
              <Button
                icon={<FormNextLink />}
                reverse
                label="Next"
                onClick={() => {
                  const { next_token } = state.data;
                  callApi({
                    url: `/api/templates?next_token=${encodeURIComponent(
                      next_token
                    )}`
                  });
                  currentPage.tokens.push(next_token);

                  setPage({
                    current: currentPage.current + 1,
                    tokens: currentPage.tokens
                  });
                }}
                disabled={state.data.next_token === ""}
              />
            </Box>
          </Box>
        ) : null}
      </Box>
    </Grid>
  );
};

export default List;
