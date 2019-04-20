import React, { Fragment } from "react";
import { FormClose, Logout, UserSettings, Template } from "grommet-icons";
import { Box, Button, Collapsible, Layer } from "grommet";
import { NavLink } from "react-router-dom";

const NavLinks = () => (
  <Fragment>
    <Box margin={{ top: "small" }}>
      <Box pad="xsmall" direction="row" border="bottom">
        <NavLink
          style={{ textDecoration: "none", marginLeft: "10px" }}
          to="/dashboard/subscribers"
        >
          Subscribers
        </NavLink>
      </Box>
      <Box pad="xsmall" direction="row" border="bottom">
        <NavLink
          style={{ textDecoration: "none", marginLeft: "10px" }}
          to="/dashboard/lists"
        >
          Lists
        </NavLink>
      </Box>
      <Box pad="xsmall" direction="row" border="bottom">
        <NavLink
          style={{ textDecoration: "none", marginLeft: "10px" }}
          to="/dashboard/campaigns"
        >
          Campaigns
        </NavLink>
      </Box>
      <Box pad="xsmall" direction="row" border="bottom">
        <NavLink
          style={{ textDecoration: "none", marginLeft: "10px" }}
          to="/dashboard/templates"
        >
          <Template size="medium" />
          Templates
        </NavLink>
      </Box>
    </Box>
    <Box>
      <Box pad="xsmall" direction="row" border="bottom">
        <UserSettings size="medium" />
        <NavLink
          style={{ textDecoration: "none", marginLeft: "10px" }}
          to="/dashboard/settings"
        >
          Settings
        </NavLink>
      </Box>
      <Box pad="xsmall" direction="row">
        <Logout size="medium" />
        <NavLink
          style={{ textDecoration: "none", marginLeft: "10px" }}
          to="/logout"
        >
          Logout
        </NavLink>
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
