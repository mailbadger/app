import React from "react";
import PropTypes from "prop-types";
import { Grid, Box, Heading, Button, ThemeContext } from "grommet";
import { UserAdd, SubtractCircle, Download, Edit } from "grommet-icons";

import useApi from "../hooks/useApi";
import LoadingOverlay from "../ui/LoadingOverlay";

const ActionButtons = () => (
  <ThemeContext.Extend
    value={{
      button: {
        border: {
          radius: "18px",
        },
        padding: {
          vertical: "2px",
          horizontal: "12px",
        },
      },
      text: {
        medium: {
          size: "14px",
        },
      },
    }}
  >
    <Box margin={{ right: "small" }}>
      <Button
        size="small"
        gap="xsmall"
        label="Import subscribers"
        icon={<UserAdd size="20px" />}
      />
    </Box>
    <Box margin={{ right: "small" }}>
      <Button
        gap="xsmall"
        label="Remove subscribers"
        icon={<SubtractCircle size="20px" />}
      />
    </Box>
    <Box margin={{ right: "small" }}>
      <Button
        gap="xsmall"
        label="Export segment"
        icon={<Download size="20px" />}
      />
    </Box>
  </ThemeContext.Extend>
);

const Details = ({ match }) => {
  const [state] = useApi({
    url: `/api/segments/${match.params.id}`,
  });

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
      rows={["1fr", "1fr", "1fr"]}
      columns={["fill"]}
      margin="medium"
      areas={[
        { name: "title", start: [0, 0], end: [0, 2] },
        { name: "actions", start: [0, 1], end: [0, 1] },
        { name: "main", start: [0, 2], end: [0, 2] },
      ]}
    >
      {!state.isLoading && state.data && (
        <>
          <Box gridArea="title" direction="row">
            <Heading level="2" alignSelf="center">
              {state.data.name}
            </Heading>
            <Button
              a11yTitle="edit segment name"
              alignSelf="center"
              icon={<Edit a11yTitle="edit segment name" color="dark-1" />}
            />
          </Box>
          <Box gridArea="actions" direction="row" align="start">
            <ActionButtons />
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
