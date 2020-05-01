import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { Grid, Box, Heading, Button } from "grommet";
import { Edit } from "grommet-icons";

import { useApi } from "../hooks";
import { LoadingOverlay, Modal } from "../ui";
import EditSegment from "./Edit";

const Details = ({ match }) => {
  const [segment, setSegment] = useState();
  const [showEdit, setShowEdit] = useState(false);

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
      columns={["fill"]}
      margin="medium"
      areas={[
        { name: "title", start: [0, 0], end: [0, 1] },
        { name: "main", start: [0, 1], end: [0, 1] },
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
          <Box gridArea="title" direction="row">
            <Heading level="2" alignSelf="center">
              {segment.name}
            </Heading>
            <Button
              a11yTitle="edit segment name"
              alignSelf="center"
              icon={<Edit a11yTitle="edit segment name" color="dark-1" />}
              onClick={() => setShowEdit(true)}
            />
          </Box>

          <Box gridArea="main">
            <Heading level="2">Main</Heading>
          </Box>
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
