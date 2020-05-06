import React, { useState, useContext, memo } from "react";
import PropTypes from "prop-types";
import { parseISO, formatRelative } from "date-fns";
import { More, Add, FormPreviousLink, FormNextLink } from "grommet-icons";
import { mainInstance as axios } from "../axios";
import { useApi } from "../hooks";
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
import history from "../history";
import {
  StyledTable,
  ButtonWithLoader,
  PlaceholderTable,
  Modal,
  Badge,
  Notice,
  BarLoader,
} from "../ui";
import CreateCampaign from "./Create";
import EditCampaign from "./Edit";
import { SesKeysContext } from "../Settings/SesKeysContext";

const Row = memo(({ campaign, setShowDelete, setShowEdit, hasSesKeys }) => {
  const d = parseISO(campaign.created_at);
  const statusColors = {
    draft: "#CCCCCC",
    sending: "#00739D",
    sent: "#00C781",
    scheduled: "#FFCA58",
  };

  let opts = ["Delete"];
  if (hasSesKeys) {
    opts.unshift("Edit");
  }

  return (
    <TableRow>
      <TableCell scope="row" size="large">
        <strong>{campaign.name}</strong>
      </TableCell>
      <TableCell scope="row" size="large">
        <Badge color={statusColors[campaign.status]}>{campaign.status}</Badge>
      </TableCell>
      <TableCell scope="row" size="large">
        {campaign.template_name}
      </TableCell>
      <TableCell scope="row" size="large">
        {formatRelative(d, new Date())}
      </TableCell>
      <TableCell scope="row" size="xsmall" align="end">
        <Select
          alignSelf="center"
          plain
          icon={<More />}
          options={opts}
          onChange={({ option }) => {
            (function () {
              switch (option) {
                case "Edit":
                  setShowEdit({
                    show: true,
                    id: campaign.id,
                  });
                  break;
                case "Delete":
                  setShowDelete({
                    show: true,
                    name: campaign.name,
                    id: campaign.id,
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
});

Row.propTypes = {
  campaign: PropTypes.shape({
    name: PropTypes.string,
    id: PropTypes.number,
    status: PropTypes.string,
    template_name: PropTypes.string,
    created_at: PropTypes.string,
  }),
  setShowDelete: PropTypes.func,
  setShowEdit: PropTypes.func,
  hasSesKeys: PropTypes.bool,
};

Row.displayName = "Row";

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

Header.displayName = "Header";

const CampaignsTable = memo(
  ({ list, setShowDelete, hasSesKeys, setShowEdit }) => (
    <StyledTable>
      <Header />
      <TableBody>
        {list.map((c) => (
          <Row
            campaign={c}
            key={c.id}
            setShowDelete={setShowDelete}
            setShowEdit={setShowEdit}
            hasSesKeys={hasSesKeys}
          />
        ))}
      </TableBody>
    </StyledTable>
  )
);

CampaignsTable.displayName = "CampaignsTable";
CampaignsTable.propTypes = {
  list: PropTypes.array,
  setShowDelete: PropTypes.func,
  setShowEdit: PropTypes.func,
  hasSesKeys: PropTypes.bool,
};

const DeleteForm = ({ id, callApi, hideModal }) => {
  const deleteCampaign = async (id) => {
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
  hideModal: PropTypes.func,
};

const List = () => {
  const [showDelete, setShowDelete] = useState({
    show: false,
    name: "",
    id: "",
  });
  const [showEdit, setShowEdit] = useState({ show: false, id: "" });
  const [showCreate, openCreateModal] = useState(false);
  const hideModal = () => setShowDelete({ show: false, name: "", id: "" });
  const hideEditModal = () => setShowEdit({ show: false, id: "" });
  const { keys, isLoading: keysLoading, error: keysError } = useContext(
    SesKeysContext
  );

  const hasSesKeys = !keysLoading && !keysError && keys !== null;

  const [state, callApi] = useApi(
    {
      url: "/api/campaigns",
    },
    {
      collection: [],
      init: true,
    }
  );

  const hasCampaigns =
    !state.isLoading && !state.isError && state.data.collection.length > 0;

  if (keysLoading) {
    return (
      <Box alignSelf="center" margin="20%">
        <BarLoader size={15} />
      </Box>
    );
  }

  if (!hasSesKeys && !hasCampaigns) {
    return (
      <Box>
        <Box align="center" margin={{ top: "large" }}>
          <Heading level="2">Please provide your AWS SES keys first.</Heading>
        </Box>
        <Box align="center" margin={{ top: "medium" }}>
          <Button
            primary
            color="status-ok"
            label="Add SES Keys"
            icon={<Add />}
            reverse
            onClick={() => history.push("/dashboard/settings")}
          />
        </Box>
      </Box>
    );
  }

  let table = null;
  if (state.isLoading || keysLoading) {
    table = <PlaceholderTable header={Header} numCols={3} numRows={3} />;
  } else if (hasCampaigns) {
    table = (
      <CampaignsTable
        isLoading={state.isLoading}
        list={state.data.collection}
        setShowDelete={setShowDelete}
        setShowEdit={setShowEdit}
        hasSesKeys={hasSesKeys}
      />
    );
  }

  return (
    <Grid
      rows={["fill", "fill"]}
      columns={["small", "large", "xsmall"]}
      gap="small"
      margin="medium"
      areas={[
        ["nav", "nav", "nav"],
        ["main", "main", "main"],
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
      {showCreate && (
        <Modal
          title={`Create campaign`}
          hideModal={() => openCreateModal(false)}
          form={
            <CreateCampaign
              callApi={callApi}
              hideModal={() => openCreateModal(false)}
            />
          }
        />
      )}
      {showEdit.show && (
        <Modal
          title={`Edit campaign`}
          hideModal={hideEditModal}
          form={
            <EditCampaign
              id={showEdit.id}
              callApi={callApi}
              hideModal={hideEditModal}
            />
          }
        />
      )}
      <Box gridArea="nav" direction="row">
        <Box alignSelf="center" margin={{ right: "small" }}>
          <Heading level="2">Campaigns</Heading>
        </Box>
        <Box alignSelf="center">
          <Button
            primary
            color="status-ok"
            label="Create new"
            icon={<Add />}
            reverse
            onClick={() => {
              if (hasSesKeys) {
                openCreateModal(true);
              }
            }}
            disabled={!hasSesKeys}
          />
        </Box>
        {!hasSesKeys && hasCampaigns && (
          <Box margin={{ left: "auto" }} alignSelf="center">
            <Notice
              message="Set your SES keys in order to send, create or edit campaigns."
              status="status-warning"
            />
          </Box>
        )}
      </Box>
      <Box gridArea="main">
        <Box animation="fadeIn">
          {table}

          {!state.isLoading &&
            !state.isError &&
            state.data.collection.length === 0 && (
              <Box align="center" margin={{ top: "large" }}>
                <Heading level="2">Create your first campaign.</Heading>
              </Box>
            )}
        </Box>
        {hasCampaigns ? (
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
