import React, { useState } from "react";
import { parseISO } from "date-fns";
import { More, Add } from "grommet-icons";
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
  Layer,
  Heading,
  Select
} from "grommet";
import history from "../history";
import StyledTable from "../ui/StyledTable";
import StyledButton from "../ui/StyledButton";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import PlaceholderRow from "../ui/PlaceholderRow";

const deleteSegment = async id => {
  await axios.delete(`/api/segments/${id}`);
};

const Row = ({ segment, setShowDelete }) => {
  const res = parseISO(segment.created_at);
  return (
    <TableRow>
      <TableCell scope="row" size="xlarge">
        {segment.name}
      </TableCell>
      <TableCell scope="row" size="medium">
        {res.toUTCString()}
      </TableCell>
      <TableCell scope="row" size="xsmall">
        <Select
          alignSelf="center"
          plain
          icon={<More />}
          options={["Edit", "Delete"]}
          onChange={({ option }) => {
            (function() {
              switch (option) {
                case "Edit":
                  history.push(`/dashboard/segments/${segment.id}/edit`);
                  break;
                case "Delete":
                  setShowDelete({
                    show: true,
                    name: segment.name,
                    id: segment.id
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

const Header = () => (
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
        size="small"
      />
    </TableRow>
  </TableHeader>
);
const PlaceholderTable = () => (
  <StyledTable caption="Segments">
    <Header />
    <TableBody>
      <PlaceholderRow columns={3} />
      <PlaceholderRow columns={3} />
      <PlaceholderRow columns={3} />
      <PlaceholderRow columns={3} />
      <PlaceholderRow columns={3} />
    </TableBody>
  </StyledTable>
);

const TemplateTable = React.memo(({ list, setShowDelete }) => (
  <StyledTable caption="Segments">
    <Header />
    <TableBody>
      {list.map(s => (
        <Row segment={s} key={s.id} setShowDelete={setShowDelete} />
      ))}
    </TableBody>
  </StyledTable>
));

const DeleteLayer = ({ setShowDelete, name, id, callApi }) => {
  const hideModal = () => setShowDelete({ show: false, name: "", id: "" });
  const [isSubmitting, setSubmitting] = useState(false);

  return (
    <Layer onEsc={() => hideModal()} onClickOutside={() => hideModal()}>
      <Heading margin="small" level="4">
        Delete segment {name} ?
      </Heading>
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
              await deleteSegment(id);
              await callApi({ url: "/api/segments" });
              setSubmitting(false);
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

  const [state, callApi] = useApi(
    {
      url: "/api/segments"
    },
    {
      collection: [],
      init: true
    }
  );

  let table = null;
  if (state.isLoading) {
    table = <PlaceholderTable />;
  } else if (state.data.collection.length > 0) {
    table = (
      <TemplateTable
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
        <DeleteLayer
          name={showDelete.name}
          setShowDelete={setShowDelete}
          callApi={callApi}
        />
      )}
      <Box gridArea="nav" direction="row">
        <Box>
          <Heading level="2" margin={{ bottom: "xsmall" }}>
            Segments
          </Heading>
        </Box>
        <Box margin={{ left: "medium", top: "medium" }}>
          <ButtonWithLoader
            label="Create new"
            icon={<Add color="#ffffff" />}
            reverse
            onClick={() => history.push("/dashboard/segments/new")}
          />
        </Box>
      </Box>
      <Box gridArea="main">
        <Box animation="fadeIn">
          {table}

          {!state.isLoading && state.data.collection.length === 0 ? (
            <Box align="center" margin={{ top: "large" }}>
              <Heading level="2">Create your first segment.</Heading>
            </Box>
          ) : null}
        </Box>
        {!state.isLoading && state.data.collection.length > 0 ? (
          <Box direction="row" alignSelf="end" margin={{ top: "medium" }}>
            <Box margin={{ right: "small" }}>
              <StyledButton
                label="Previous"
                onClick={() => console.log("Previous")}
              />
            </Box>
            <Box>
              <StyledButton label="Next" onClick={() => console.log("Next")} />
            </Box>
          </Box>
        ) : null}
      </Box>
    </Grid>
  );
};

export default List;
