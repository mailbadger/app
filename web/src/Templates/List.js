import React, { Fragment, useState } from "react";
import { parse } from "date-fns";
import { Edit, Trash } from "grommet-icons";
import axios from "axios";
import useApi from "../hooks/useApi";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableCell,
  Box,
  Button,
  Layer,
  Heading
} from "grommet";
import history from "../history";

const deleteTemplate = async name => {
  try {
    await axios.delete(`/api/templates/${name}`);
  } catch (error) {
    console.log(error.response.data);
  }
};

const Row = ({ template, setShowDelete }) => {
  const res = parse(template.timestamp);
  return (
    <TableRow>
      <TableCell scope="row" size="large">
        <strong>{template.name}</strong>
      </TableCell>
      <TableCell scope="row" size="medium">
        {res.toUTCString()}
      </TableCell>
      <TableCell scope="row">
        <Button
          plain
          icon={<Edit />}
          onClick={() =>
            history.push(`/dashboard/templates/${template.name}/edit`)
          }
        />
      </TableCell>
      <TableCell scope="row">
        <Button
          plain
          icon={<Trash />}
          onClick={() => setShowDelete({ show: true, name: template.name })}
        />
      </TableCell>
    </TableRow>
  );
};

const TemplateTable = React.memo(({ list, setShowDelete }) => (
  <Table caption="Templates">
    <TableHeader>
      <TableRow>
        <TableCell scope="col" border="bottom" size="xlarge">
          Name
        </TableCell>
        <TableCell scope="col" border="bottom">
          Created At
        </TableCell>
        <TableCell scope="col" border="bottom">
          Edit
        </TableCell>
        <TableCell scope="col" border="bottom">
          Delete
        </TableCell>
      </TableRow>
    </TableHeader>
    <TableBody>
      {list.map(t => (
        <Row template={t} key={t.name} setShowDelete={setShowDelete} />
      ))}
    </TableBody>
  </Table>
));

const DeleteLayer = ({ setShowDelete, name, callApi }) => {
  return (
    <Layer
      onEsc={() => setShowDelete({ show: false, name: "" })}
      onClickOutside={() => setShowDelete({ show: false, name: "" })}
    >
      <Heading margin="small" level="4">
        Delete template {name} ?
      </Heading>
      <Box direction="row" alignSelf="end" pad="small">
        <Box margin={{ right: "small" }}>
          <Button
            label="Cancel"
            onClick={() => setShowDelete({ show: false, name: "" })}
          />
        </Box>
        <Box>
          <Button
            label="Delete"
            onClick={() => {
              deleteTemplate(name);
              callApi({ url: "/api/templates" });
              setShowDelete({ show: false, name: "" });
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
          <Button
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
