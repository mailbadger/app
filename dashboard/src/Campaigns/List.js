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
import Badge from "../ui/Badge";

const Row = ({ campaign, setShowDelete }) => {
  const d = parseISO(campaign.created_at);
  const statusColors = {
    draft: "#CCCCCC",
    sending: "#00739D",
    sent: "#00C781",
    scheduled: "#FFCA58"
  };
  return (
    <TableRow>
      <TableCell scope="row" size="xxsmall">
        <strong>{campaign.name}</strong>
      </TableCell>
      <TableCell scope="row" size="xxsmall">
        <Badge color={statusColors[campaign.status]}>{campaign.status}</Badge>
      </TableCell>
      <TableCell scope="row" size="xxsmall">
        {campaign.template_name}
      </TableCell>
      <TableCell scope="row" size="xxsmall">
        {formatRelative(d, new Date())}
      </TableCell>
      <TableCell scope="row" size="xxsmall" align="end">
        <Select
          alignSelf="center"
          plain
          icon={<More />}
          options={["Edit", "Delete"]}
          onChange={({ option }) => {
            (function() {
              switch (option) {
                case "Edit":
                  history.push(`/dashboard/campaigns/${campaign.id}/edit`);
                  break;
                case "Delete":
                  setShowDelete({
                    show: true,
                    name: campaign.name,
                    id: campaign.id
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
  campaign: PropTypes.shape({
    name: PropTypes.string,
    id: PropTypes.number,
    status: PropTypes.string,
    template_name: PropTypes.string,
    created_at: PropTypes.string
  }),
  setShowDelete: PropTypes.func
};

const Header = () => (
  <TableHeader>
    <TableRow>
      <TableCell scope="col" border="bottom" size="xxsmall">
        <strong>Name</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xxsmall">
        <strong>Status</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xxsmall">
        <strong>Template</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="xxsmall">
        <strong>Created At</strong>
      </TableCell>
      <TableCell
        style={{ textAlign: "right" }}
        align="end"
        scope="col"
        border="bottom"
        size="xxsmall"
      >
        <strong>Action</strong>
      </TableCell>
    </TableRow>
  </TableHeader>
);

const CampaignsTable = React.memo(({ list, setShowDelete }) => (
  <StyledTable>
    <Header />
    <TableBody>
      {list.map(c => (
        <Row campaign={c} key={c.id} setShowDelete={setShowDelete} />
      ))}
    </TableBody>
  </StyledTable>
));

CampaignsTable.displayName = "CampaignsTable";
CampaignsTable.propTypes = {
  list: PropTypes.array,
  setShowDelete: PropTypes.func
};

const DeleteForm = ({ id, callApi, hideModal }) => {
  const deleteCampaign = async id => {
    await axios.delete(`/api/campaigns/${id}`);
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
            await deleteCampaign(id);
            await callApi({ url: "/api/campaigns" });
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
  hideModal: PropTypes.func
};

const List = () => {
  const [showDelete, setShowDelete] = useState({ show: false, name: "" });
  const hideModal = () => setShowDelete({ show: false, name: "", id: "" });

  const [state, callApi] = useApi(
    {
      url: "/api/campaigns"
    },
    {
      collection: [],
      init: true
    }
  );

  let table = null;
  if (state.isLoading) {
    table = <PlaceholderTable header={Header} numCols={3} numRows={3} />;
  } else if (state.data.collection.length > 0) {
    table = (
      <CampaignsTable
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
          title={`Delete campaign ${showDelete.name} ?`}
          hideModal={hideModal}
          form={
            <DeleteForm
              id={showDelete.id}
              callApi={callApi}
              hideModal={hideModal}
            />
          }
        />
      )}
      <Box gridArea="nav" direction="row">
        <Box>
          <Heading level="2" margin={{ bottom: "xsmall" }}>
            Campaigns
          </Heading>
        </Box>
        <Box margin={{ left: "medium", top: "medium" }}>
          <Button
            primary
            color="status-ok"
            label="Create new"
            icon={<Add />}
            reverse
            onClick={() => console.log("new campaign!")}
          />
        </Box>
      </Box>
      <Box gridArea="main">
        <Box animation="fadeIn">
          {table}

          {!state.isLoading && state.data.collection.length === 0 ? (
            <Box align="center" margin={{ top: "large" }}>
              <Heading level="2">Create your first campaign.</Heading>
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
                    url: state.data.links.previous
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
                    url: state.data.links.next
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
