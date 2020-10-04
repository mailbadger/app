import React, { memo } from "react";
import PropTypes from "prop-types";
import {
  TableHeader,
  TableBody,
  TableRow,
  TableCell,
  Box,
  Heading,
} from "grommet";

import { StyledTable, PlaceholderTable } from "../../../ui";
import { useApi } from "../../../hooks";

const Row = memo(({ click }) => {
  return (
    <TableRow>
      <TableCell scope="row" size="small">
        {click.link}
      </TableCell>
      <TableCell scope="row" size="small">
        {click.unique}
      </TableCell>
      <TableCell scope="row" size="large">
        {click.total}
      </TableCell>
    </TableRow>
  );
});

Row.displayName = "Row";
Row.propTypes = {
  click: PropTypes.shape({
    link: PropTypes.string,
    unique: PropTypes.number,
    total: PropTypes.number,
  }),
};

const Header = () => (
  <TableHeader>
    <TableRow>
      <TableCell scope="col" border="bottom" size="large">
        <strong>Link</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xxsmall">
        <strong>Unique</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xxsmall">
        <strong>Total</strong>
      </TableCell>
    </TableRow>
  </TableHeader>
);

Header.displayName = "Header";

const Table = memo(({ list }) => (
  <StyledTable>
    <Header />
    <TableBody>
      {list.map((c, i) => (
        <Row click={c} key={i} />
      ))}
    </TableBody>
  </StyledTable>
));

Table.displayName = "Table";
Table.propTypes = {
  list: PropTypes.array,
};

const Clicks = ({ campaignId }) => {
  const [state] = useApi(
    {
      url: `/api/campaigns/${campaignId}/clicks`,
    },
    {
      collection: [],
    }
  );

  if (state.isLoading) {
    return <PlaceholderTable header={Header} numCols={6} numRows={6} />;
  }
  if (!state.isLoading && !state.isError) {
    return (
      <>
        <Table list={state.data.collection} />
        {state.data.collection.length === 0 && (
            <Box align="center">
              <Heading level="3">Clicks list is currently empty.</Heading>
            </Box>
        )}
      </>
    );
  }
  return null;
};

Clicks.propTypes = {
  campaignId: PropTypes.number,
};

export { Table, Row, Header, Clicks };

export default Clicks;
