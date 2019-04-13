import React, { Fragment } from "react";
import { FormClose } from "grommet-icons";
import { Box, Button, Collapsible, Layer } from "grommet";

const Sidebar = props => {
  const { showSidebar, size, closeSidebar } = props;

  return (
    <Fragment>
      {!showSidebar || size !== "small" ? (
        <Collapsible direction="horizontal" open={showSidebar}>
          <Box
            flex
            width="18em"
            background="dark-2"
            elevation="small"
            align="center"
            justify="center"
          >
            sidebar
          </Box>
        </Collapsible>
      ) : (
        <Layer>
          <Box
            background="light-2"
            tag="header"
            justify="end"
            align="center"
            direction="row"
          >
            <Button icon={<FormClose />} onClick={closeSidebar} />
          </Box>
          <Box fill background="light-2" align="center" justify="center">
            wut
          </Box>
        </Layer>
      )}
    </Fragment>
  );
};

export default Sidebar;
