import React from "react";
import { Grid, Box, Heading, Button } from "grommet";

import List from "./List";

const Templates = () => {
  return (
    <Grid
      rows={["small", "medium"]}
      columns={["fit"]}
      areas={[
        { name: "header", start: [0, 0], end: [0, 0] },
        { name: "main", start: [0, 1], end: [0, 1] }
      ]}
    >
      <Box gridArea="header" margin="medium" background="light-4">
        <Heading>Templates</Heading>
      </Box>
      <Box gridArea="main" margin="medium" background="light-2">
        <List />
      </Box>
    </Grid>
  );
};

export default Templates;
