import React, {
  useState,
  useContext,
  useEffect,
  useRef,
  Fragment,
} from "react";
import PropTypes from "prop-types";
import { parseISO, formatRelative } from "date-fns";
import { More } from "grommet-icons";
import {
  TableBody,
  TableRow,
  TableCell,
  Box,
  Heading,
  ResponsiveContext,
  Select,
} from "grommet";

import history from "../history";
import { useApi, useInterval } from "../hooks";
import { mainInstance as axios } from "../axios";
import { StyledTable, Modal } from "../ui";
import { NotificationsContext } from "../Notifications/context";
import CreateSubscriber from "./Create";
import DeleteSubscriber from "./Delete";
import EditSubscriber from "./Edit";
import { DashboardDataTable, getColumnSize } from "../ui/DashboardDataTable";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faPlus } from "@fortawesome/free-solid-svg-icons";
import {
  StyledHeaderWrapper,
  StyledHeaderButtons,
  StyledHeaderTitle,
  StyledHeaderButton,
  StyledImportButton,
  StyledActions,
} from "./StyledSections";
import { StyledTableHeader } from "../ui/DashboardStyledTable";
import DashboardPlaceholderTable from "../ui/DashboardPlaceholderTable";
import { endpoints } from "../network/endpoints";

export const Row = ({ subscriber, actions }) => {
  const ca = parseISO(subscriber.created_at);
  const ua = parseISO(subscriber.updated_at);
  return (
    <TableRow>
      <TableCell scope="row" size="medium">
        <strong>{subscriber.email}</strong>
      </TableCell>
      <TableCell scope="row" size="medium">
        {formatRelative(ca, new Date())}
      </TableCell>
      <TableCell scope="row" size="medium">
        {formatRelative(ua, new Date())}
      </TableCell>
      <TableCell scope="row" size="xsmall" align="end">
        {actions}
      </TableCell>
    </TableRow>
  );
};

Row.propTypes = {
  subscriber: PropTypes.shape({
    email: PropTypes.string,
    id: PropTypes.number,
    created_at: PropTypes.string,
    updated_at: PropTypes.string,
  }),
  actions: PropTypes.element,
};

export const Header = ({ size }) => (
  <StyledTableHeader>
    <TableRow>
      <TableCell scope="col" border="bottom" size={getColumnSize(size)}>
        <strong>Email</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="small">
        <strong>Created At</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="small">
        <strong>Updated At</strong>
      </TableCell>
      <TableCell scope="col" border="bottom" size="128px">
        <strong> {""}</strong>
      </TableCell>
      <TableCell align="center" scope="col" border="bottom" size="small">
        <strong>Actions</strong>
      </TableCell>
    </TableRow>
  </StyledTableHeader>
);

Header.propTypes = {
  size: PropTypes.string,
};
export const SubscriberTable = React.memo(({ list, actions }) => (
  <StyledTable>
    <Header />
    <TableBody>
      {list.map((s) => (
        <Row subscriber={s} key={s.id} actions={actions(s)} />
      ))}
    </TableBody>
  </StyledTable>
));

SubscriberTable.displayName = "SubscriberTable";
SubscriberTable.propTypes = {
  list: PropTypes.array,
  actions: PropTypes.func,
};

const ExportSubscribers = () => {
  const linkEl = useRef(null);
  const { createNotification } = useContext(NotificationsContext);
  const [notification, setNotification] = useState();
  const [filename, setFilename] = useState("");
  const [retries, setRetries] = useState(-1);
  const [state, callApi] = useApi(
    {
      url: endpoints.getSubscribersExport,
    },
    null,
    true
  );

  useInterval(
    async () => {
      await callApi({
        url: `${endpoints.getSubscribersExportDownload}?filename=${filename}`,
      });
      setRetries(retries - 1);
    },
    retries > 0 ? 1000 : null
  );

  useEffect(() => {
    if (notification) {
      createNotification(notification.message, notification.status);
    }
  }, [notification]);

  useEffect(() => {
    if (!state.isLoading && state.isError && state.data) {
      if (state.data.status === "failed" && retries > 0 && retries < 50) {
        setRetries(-1);
        setNotification({
          message: state.data.message,
          status: "status-error",
        });
      }
    }

    if (!state.isLoading && !state.isError && state.data) {
      if (retries > 0) {
        setRetries(-1);
        linkEl.current.click();
      }
    }
  }, [state]);

  return (
    <Fragment>
      <StyledHeaderButton
        width="110"
        disabled={retries > 0 || state.isLoading}
        onClick={async () => {
          try {
            const res = await axios.post(endpoints.postSubscribersExport);
            setFilename(res.data.file_name);
            setRetries(50);
          } catch (e) {
            console.error("Unable to generate report", e);
          }
        }}
        label="Export"
      />
      {!state.isLoading && !state.isError && state.data && (
        <a ref={linkEl} href={state.data.url} />
      )}
    </Fragment>
  );
};

const getData = (subscribersData, setShowEdit, setShowDelete) => {
  const data = [];

  for (let i = 0; i < subscribersData.length; i += 1) {
    const { email, created_at, updated_at, id } = subscribersData[i];

    const dateCreatedAt = new Date(created_at);
    const dateUpdatedAt = parseISO(updated_at);

    data.push({
      email,
      created: dateCreatedAt.toLocaleDateString("en-US"),
      updated: formatRelative(dateUpdatedAt, new Date()),
      tags: "Subscribers",
      actions: (
        <StyledActions>
          <Select
            alignSelf="center"
            plain
            defaultValue="View"
            icon={<More />}
            options={["View", "Edit", "Delete"]}
            onChange={({ option }) => {
              (() => {
                switch (option) {
                  case "Edit":
                    setShowEdit({
                      show: true,
                      id,
                    });
                    break;
                  case "View":
                    history.push(`/dashboard/groups/${id}`);
                    break;
                  case "Delete":
                    setShowDelete({
                      show: true,
                      email,
                      id,
                    });
                    break;
                  default:
                    return "null";
                }
              })();
            }}
          />
        </StyledActions>
      ),
    });
  }

  return data;
};

const search = (setFilteredData, searchInput, data) => {
  let filteredData = [];
  if (searchInput !== "") {
    filteredData = data.filter((entry) => {
      const foundMatch = Object.values(entry).some((entryValue) =>
        entryValue.toString().toLowerCase().includes(searchInput.toLowerCase())
      );

      if (foundMatch) {
        return entry;
      }
    });
  }

  setFilteredData(filteredData);
};

const List = () => {
  const [showDelete, setShowDelete] = useState({
    show: false,
    email: "",
    id: "",
  });
  const [showEdit, setShowEdit] = useState({ show: false, id: "" });
  const [showCreate, openCreateModal] = useState(false);
  const [searchInput, setSearchInput] = useState("");
  const [filteredData, setFilteredData] = useState([]);

  const hideDeleteModal = () =>
    setShowDelete({ show: false, email: "", id: "" });
  const hideEditModal = () => setShowEdit({ show: false, id: "" });

  const contextSize = useContext(ResponsiveContext);
  const columns = [
    { property: "email", header: "Email", size: getColumnSize(contextSize) },
    { property: "created", header: "Created At", size: "small" },
    { property: "updated", header: "Updated At", size: "small" },
    { property: "tags", header: "", size: "128px", align: "center" },
    { property: "actions", header: "Actions", size: "small", align: "center" },
  ];

  const [state, callApi] = useApi(
    {
      url: endpoints.getGroups,
    },
    {
      collection: [],
      init: true,
    }
  );

  const onClickPrev = () => {
    callApi({
      url: state.data.links.previous,
    });
  };

  const onClickNext = () => {
    callApi({
      url: state.data.links.next,
    });
  };

  const dataFromApi = getData(
    state.data.collection,
    setShowEdit,
    setShowDelete
  );

  const handleChange = (e) => {
    setSearchInput(e.target.value);
  };

  useEffect(() => {
    search(setFilteredData, searchInput, dataFromApi);
  }, [searchInput]);

  let table = null;
  const data = filteredData && filteredData.length ? filteredData : dataFromApi;

  if (state.isLoading) {
    table = (
      <DashboardPlaceholderTable
        columns={columns}
        numCols={columns.length}
        numRows={10}
      />
    );
  } else if (data.length > 0) {
    table = (
      <DashboardDataTable
        columns={columns}
        data={data}
        isLoading={state.isLoading}
        onClickNext={onClickNext}
        onClickPrev={onClickPrev}
        prevLinks={state.data.links.previous}
        nextLinks={state.data.links.next}
        searchInput={searchInput}
        handleChange={handleChange}
      />
    );
  }

  return (
    <Fragment>
      {showDelete.show && (
        <Modal
          title={`Delete subscriber ${showDelete.email} ?`}
          hideModal={hideDeleteModal}
          form={
            <DeleteSubscriber
              id={showDelete.id}
              callApi={callApi}
              hideModal={hideDeleteModal}
            />
          }
        />
      )}
      {showCreate && (
        <Modal
          title={`Create subscriber`}
          hideModal={() => openCreateModal(false)}
          form={
            <CreateSubscriber
              callApi={callApi}
              hideModal={() => openCreateModal(false)}
            />
          }
        />
      )}
      {showEdit.show && (
        <Modal
          title={`Edit subscriber`}
          hideModal={hideEditModal}
          form={
            <EditSubscriber
              id={showEdit.id}
              callApi={callApi}
              hideModal={hideEditModal}
            />
          }
        />
      )}
      <StyledHeaderWrapper
        size={contextSize}
        gridArea="nav"
        margin={{ left: "40px", right: "100px", bottom: "22px", top: "40px" }}
      >
        <StyledHeaderTitle size={contextSize}>Subscribers</StyledHeaderTitle>
        <StyledHeaderButtons size={contextSize} margin={{ left: "auto" }}>
          <Fragment>
            <StyledImportButton
              width="256"
              margin={{ right: "small" }}
              icon={<FontAwesomeIcon icon={faPlus} />}
              label="Import from file"
              onClick={() => history.push("/dashboard/subscribers/import")}
            />
            <StyledHeaderButton
              width="154"
              margin={{ right: "small" }}
              label="Create New"
              onClick={() => openCreateModal(true)}
            />
            <StyledHeaderButton
              width="184"
              margin={{ right: "small" }}
              label="Delete from file"
              onClick={() => history.push("/dashboard/subscribers/bulk-delete")}
            />
            <ExportSubscribers />
          </Fragment>
        </StyledHeaderButtons>
      </StyledHeaderWrapper>
      <Box gridArea="main">
        <Box animation="fadeIn">
          {table}
          {!state.isLoading && data.length === 0 ? (
            <Box align="center" margin={{ top: "large" }}>
              <Heading level="2">Create your first subscriber.</Heading>
            </Box>
          ) : null}
        </Box>
      </Box>
    </Fragment>
  );
};

export default List;
