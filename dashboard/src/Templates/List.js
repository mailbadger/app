import React, { Fragment, useState } from "react";
import { parse } from "date-fns";
import { More } from "grommet-icons";
import axios from "axios";
import useApi from "../hooks/useApi";
import {
  TableHeader,
  TableBody,
  TableRow,
  TableCell,
  Box,
  Button,
  Layer,
  Heading,
  Select
} from "grommet";
import history from "../history";
import StyledTable from "../ui/StyledTable";
import StyledButton from "../ui/StyledButton";

const deleteTemplate = async name => {
  await axios.delete(`/api/templates/${name}`);
};

const Row = ({ template, setShowDelete }) => {
  const res = parse(template.timestamp);
  return (
    <TableRow>
      <TableCell scope="row" size="xlarge">
        {template.name}
      </TableCell>
      <TableCell scope="row" size="medium">
        {res.toUTCString()}
      </TableCell>
      <TableCell scope="row" size="xsmall">
        <Select
          alignSelf="center"
          plain
          icon={<More />}
          options={["Edit", "Preview", "Send Test", "Delete"]}
          onChange={({ option }) => {
            (function() {
              switch (option) {
                case "Edit":
                  history.push(`/dashboard/templates/${template.name}/edit`);
                  break;
                case "Preview":
                  setShowDelete({ show: true, name: template.name });
                  break;
                case "Send Test":
                  setShowDelete({ show: true, name: template.name });
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

const TemplateTable = React.memo(({ list, setShowDelete }) => (
  <StyledTable caption="Templates">
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
          size="xsmall"
        />
      </TableRow>
    </TableHeader>
    <TableBody>
      {list.map(t => (
        <Row template={t} key={t.name} setShowDelete={setShowDelete} />
      ))}
    </TableBody>
  </StyledTable>
));

const DeleteLayer = ({ setShowDelete, name, callApi }) => {
  const hideModal = () => setShowDelete({ show: false, name: "" });
  return (
    <Layer onEsc={() => hideModal()} onClickOutside={() => hideModal()}>
      <Heading margin="small" level="4">
        Delete template {name} ?
      </Heading>
      <Box direction="row" alignSelf="end" pad="small">
        <Box margin={{ right: "small" }}>
          <Button label="Cancel" onClick={() => hideModal()} />
        </Box>
        <Box>
          <Button
            label="Delete"
            onClick={() => {
              deleteTemplate(name);
              callApi({ url: "/api/templates" });
              hideModal();
            }}
          />
        </Box>
      </Box>
    </Layer>
  );
};

const List = () => {
  const [showDelete, setShowDelete] = useState({ show: false, name: "" });
  const [currentPage, setPage] = useState({ current: -1, tokens: [""] });

  const [state, callApi] = useApi(
    {
      url: "/api/templates"
    },
    {
      next_token: "",
      list: []
    }
  );

  if (state.isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <Fragment>
      {showDelete.show && (
        <DeleteLayer
          name={showDelete.name}
          setShowDelete={setShowDelete}
          callApi={callApi}
        />
      )}
      <TemplateTable list={state.data.list} setShowDelete={setShowDelete} />
      <Box direction="row" alignSelf="end" margin={{ top: "medium" }}>
        <Box margin={{ right: "small" }}>
          <StyledButton
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
          <StyledButton
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
    </Fragment>
  );
};

export default List;
