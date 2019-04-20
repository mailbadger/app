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
    <TableCell scope="row" size="xlarge">
      <strong>{template.name}</strong>
    </TableCell>
    <TableCell scope="row">{template.timestamp}</TableCell>
  </TableRow>
);

const TemplateTable = React.memo(({ list }) => (
  <Table caption="Templates">
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
      {list.map(t => (
        <Row template={t} key={t.name} />
      ))}
    </TableBody>
  </Table>
));

const List = () => {
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
      <TemplateTable list={state.data.list} />
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
