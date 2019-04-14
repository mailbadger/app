import React, { Fragment } from "react";
import { FormClose } from "grommet-icons";
import { Box, Button, Collapsible, Layer } from "grommet";
import { Link } from "react-router-dom";

const NavLinks = () => (
  <Fragment>
    <Box margin={{ top: "small" }}>
      <Box pad="xsmall" border="bottom">
        <Link style={{ textDecoration: "none" }} to="/dashboard/subscribers">
          Subscribers
        </Link>
      </Box>
      <Box pad="xsmall" border="bottom">
        <Link style={{ textDecoration: "none" }} to="/dashboard/lists">
          Lists
        </Link>
      </Box>
      <Box pad="xsmall" border="bottom">
        <Link style={{ textDecoration: "none" }} to="/dashboard/campaigns">
          Campaigns
        </Link>
      </Box>
      <Box pad="xsmall" border="bottom">
        <Link style={{ textDecoration: "none" }} to="/dashboard/templates">
          Templates
        </Link>
      </Box>
    </Box>
    <Box>
      <Box pad="xsmall" border="bottom">
        <Link style={{ textDecoration: "none" }} to="/dashboard/settings">
          Settings
        </Link>
      </Box>
      <Box pad="xsmall">
        <Link style={{ textDecoration: "none" }} to="/logout">
          Logout
        </Link>
      </Box>
    </Box>
  </Fragment>
);

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
            direction="column"
            justify="between"
          >
            <NavLinks />
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
          <Box fill background="light-2" direction="column" justify="between">
            <NavLinks />
          </Box>
        </Layer>
      )}
    </Fragment>
  );
};

export default Sidebar;
