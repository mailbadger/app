import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { Grid, Box, Heading } from "grommet";
import { Edit, Trash } from "grommet-icons";

import { useApi } from "../hooks";
import { LoadingOverlay, Modal, SecondaryButton } from "../ui";
import EditSegment from "./Edit";
import DeleteSegment from "./Delete";
import history from "../history";

const Details = ({ match }) => {
  const [segment, setSegment] = useState();
  const [showEdit, setShowEdit] = useState(false);
  const [showDelete, setShowDelete] = useState(false);

  const [state] = useApi({
    url: `/api/segments/${match.params.id}`,
  });

  useEffect(() => {
    if (state.isLoading || state.isError) {
      return;
    }

    setSegment(state.data);
  }, [state]);

  if (state.isLoading) {
    return <LoadingOverlay />;
  }

  if (state.isError) {
    return (
      <Box margin="15%" alignSelf="center">
        <Heading>Segment not found.</Heading>
      </Box>
    );
  }

  return (
    <Grid
      rows={["1fr", "1fr"]}
      columns={["small", "large", "xsmall"]}
      margin="medium"
      gap="small"
      areas={[
        ["title", "title", "gap"],
        ["main", "main", "main"],
      ]}
    >
      {segment && (
        <>
          {showEdit && (
            <Modal
              title={`Edit segment`}
              hideModal={() => setShowEdit(false)}
              form={
                <EditSegment
                  segment={segment}
                  setSegment={setSegment}
                  hideModal={() => setShowEdit(false)}
                />
              }
            />
          )}
          {showDelete && (
            <Modal
              title={`Delete segment ${segment.name} ?`}
              hideModal={() => setShowDelete(false)}
              form={
                <DeleteSegment
                  id={segment.id}
                  onSuccess={() => history.replace("/dashboard/segments")}
                  onCancel={() => setShowDelete(false)}
                />
              }
            />
          )}
          <Box gridArea="title" direction="row">
            <Heading level="2" alignSelf="center">
              {segment.name}
            </Heading>
            <Box direction="row" margin={{ left: "auto" }}>
              <SecondaryButton
                margin={{ right: "small" }}
                a11yTitle="edit segment name"
                alignSelf="center"
                icon={<Edit a11yTitle="edit segment name" color="dark-1" />}
                label="Edit"
                onClick={() => setShowEdit(true)}
              />
              <SecondaryButton
                a11yTitle="delete segment"
                alignSelf="center"
                icon={<Trash a11yTitle="delete segment" color="dark-1" />}
                label="Delete"
                onClick={() => setShowDelete(true)}
              />
            </Box>
          </Box>

          <Box gridArea="main"></Box>
        </>
      )}
    </Grid>
  );
};

Details.propTypes = {
  match: PropTypes.shape({
    params: PropTypes.shape({
      id: PropTypes.string,
    }),
  }),
};

export default Details;
