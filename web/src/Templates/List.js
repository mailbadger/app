import React, { Fragment, useState } from "react";
import useApi from "../hooks/useApi";
import {
  Table,
  TableHeader,
  TableBody,
  TableRow,
  TableCell,
  Box,
  Button
} from "grommet";

const Row = ({ template }) => (
  <TableRow>
    <TableCell scope="row">
      <strong>{template.name}</strong>
    </TableCell>
    <TableCell scope="row">{template.timestamp}</TableCell>
  </TableRow>
);

const List = () => {
  const [prevToken, setPrevius] = useState("");

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
      <Table>
        <TableHeader>
          <TableRow>
            <TableCell scope="col" border="bottom">
              Name
            </TableCell>
            <TableCell scope="col" border="bottom">
              Created At
            </TableCell>
          </TableRow>
        </TableHeader>
        <TableBody>
          {state.data.list.map(t => (
            <Row template={t} key={t.name} />
          ))}
        </TableBody>
      </Table>
      <Box direction="row" margin={{ top: "medium" }}>
        <Box margin={{ right: "small" }}>
          <Button
            label="Previous"
            onClick={() => {
              if (prevToken) {
                callApi({ url: `/api/templates?next_token=${prevToken}` });
              }
            }}
          />
        </Box>
        <Box>
          <Button
            label="Next"
            onClick={() => {
              if (state.data.next_token) {
                callApi({
                  url: `/api/templates?next_token=${state.data.next_token}`
                });
                setPrevius(state.data.next_token);
              }
            }}
          />
        </Box>
      </Box>
    </Fragment>
  );
};

export default List;
