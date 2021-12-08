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
import truncate from "../../../utils/truncate";
import { endpoints } from "../../../network/endpoints";

const Row = memo(({ complaint }) => {
  const d = parseISO(complaint.created_at);
  return (
    <TableRow>
      <TableCell scope="row" size="large">
        {complaint.recipient}
      </TableCell>
      <TableCell scope="row" size="small">
        {complaint.type}
      </TableCell>
      <TableCell scope="row" size="small">
        {truncate(complaint.user_agent, 50)}&hellip;
      </TableCell>
      <TableCell scope="row" size="large">
        {formatRelative(d, new Date())}
      </TableCell>
    </TableRow>
  );
});

Row.displayName = "Row";
Row.propTypes = {
  complaint: PropTypes.shape({
    id: PropTypes.number,
    campaign_id: PropTypes.number,
    recipient: PropTypes.string,
    type: PropTypes.string,
    user_agent: PropTypes.string,
    feedback_id: PropTypes.string,
    created_at: PropTypes.string,
  }),
};

const Header = () => (
  <TableHeader>
    <TableRow>
      <TableCell scope="col" border="bottom" size="small">
        <strong>Recipient</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="small">
        <strong>Type</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="small">
        <strong>User Agent</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="small">
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
      {list.map((c) => (
        <Row complaint={c} key={c.id} />
      ))}
    </TableBody>
  </StyledTable>
));

Table.displayName = "Table";
Table.propTypes = {
  list: PropTypes.array,
};

const Complaints = ({ campaignId }) => {
  const [state, callApi] = useApi(
    {
      url: endpoints.getCampaignComplaints(campaignId),
    },
    {
      collection: [],
    }
  );

  if (state.isLoading) {
    return <PlaceholderTable header={Header} numCols={6} numRows={3} />;
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

Complaints.propTypes = {
  campaignId: PropTypes.number,
};

export { Table, Row, Header, Complaints };

export default Complaints;
