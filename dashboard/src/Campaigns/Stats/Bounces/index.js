import React, { memo } from "react";
import PropTypes from "prop-types";
import { parseISO, formatRelative } from "date-fns";
import {
  TableHeader,
  TableBody,
  TableRow,
  TableCell,
  Box,
  Button,
} from "grommet";
import { FormPreviousLink, FormNextLink } from "grommet-icons";

import { StyledTable, PlaceholderTable } from "../../../ui";
import { useApi } from "../../../hooks";
import { endpoints } from "../../../network/endpoints";

const Row = memo(({ bounce }) => {
  const d = parseISO(bounce.created_at);
  return (
    <TableRow>
      <TableCell scope="row" size="large">
        {bounce.recipient}
      </TableCell>
      <TableCell scope="row" size="small">
        {bounce.type}
      </TableCell>
      <TableCell scope="row" size="small">
        {bounce.sub_type}
      </TableCell>
      <TableCell scope="row" size="xsmall">
        {bounce.status}
      </TableCell>
      <TableCell scope="row" size="large">
        {bounce.diagnostic_code}
      </TableCell>
      <TableCell scope="row" size="large">
        {formatRelative(d, new Date())}
      </TableCell>
    </TableRow>
  );
});

Row.displayName = "Row";
Row.propTypes = {
  bounce: PropTypes.shape({
    id: PropTypes.number,
    campaign_id: PropTypes.number,
    recipient: PropTypes.string,
    type: PropTypes.string,
    sub_type: PropTypes.string,
    action: PropTypes.string,
    status: PropTypes.string,
    diagnostic_code: PropTypes.string,
    feedback_id: PropTypes.string,
    created_at: PropTypes.string,
  }),
};

const Header = () => (
  <TableHeader>
    <TableRow>
      <TableCell scope="col" border="bottom" size="xsmall">
        <strong>Recipient</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xxsmall">
        <strong>Type</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xxsmall">
        <strong>Sub Type</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xxsmall">
        <strong>Status</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xsmall">
        <strong>Diagnostic Code</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xsmall">
        <strong>Created At</strong>
      </TableCell>
    </TableRow>
  </TableHeader>
);

Header.displayName = "Header";

const Table = memo(({ list }) => (
  <StyledTable>
    <Header />
    <TableBody>
      {list.map((b) => (
        <Row bounce={b} key={b.id} />
      ))}
    </TableBody>
  </StyledTable>
));

Table.displayName = "Table";
Table.propTypes = {
  list: PropTypes.array,
};

const Bounces = ({ campaignId }) => {
  const [state, callApi] = useApi(
    {
      url: endpoints.getCampaignBounces(campaignId),
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
        {state.data.collection.length > 0 ? (
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
      </>
    );
  }
  return null;
};

Bounces.propTypes = {
  campaignId: PropTypes.number,
};

export { Table, Row, Header, Bounces };

export default Bounces;
