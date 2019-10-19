import React, { Fragment } from "react";
import PropTypes from "prop-types";
import {
  FormClose,
  Logout,
  UserSettings,
  Group,
  List,
  Send,
  Template
} from "grommet-icons";
import { Box, Button, Collapsible, Layer } from "grommet";
import { NavLink } from "react-router-dom";

const StyledNavLink = props => (
  <NavLink
    style={{
      textDecoration: "none",
      marginLeft: "10px",
      color: "#AA9AD4",
      textTransform: "uppercase"
    }}
    {...props}
  >
    {props.children}
  </NavLink>
);

StyledNavLink.propTypes = {
  children: PropTypes.element.isRequired
};

const NavLinks = () => (
  <Fragment>
    <Box margin={{ top: "small" }}>
      <Box
        pad="xsmall"
        direction="row"
        style={{ borderBottom: "1px solid #53418B" }}
      >
        <Group size="medium" color="#AA9AD4" style={{ marginLeft: "10px" }} />
        <StyledNavLink to="/dashboard/subscribers">Subscribers</StyledNavLink>
      </Box>
      <Box
        pad="xsmall"
        direction="row"
        style={{
          borderBottom: "1px solid #53418B",
          borderTop: "1px solid #7058BA"
        }}
      >
        <List size="medium" color="#AA9AD4" style={{ marginLeft: "10px" }} />
        <StyledNavLink to="/dashboard/segments">Segments</StyledNavLink>
      </Box>
      <Box
        pad="xsmall"
        direction="row"
        style={{
          borderBottom: "1px solid #53418B",
          borderTop: "1px solid #7058BA"
        }}
      >
        <Send size="medium" color="#AA9AD4" style={{ marginLeft: "10px" }} />
        <StyledNavLink to="/dashboard/campaigns">Campaigns</StyledNavLink>
      </Box>
      <Box
        pad="xsmall"
        direction="row"
        style={{
          borderBottom: "1px solid #53418B",
          borderTop: "1px solid #7058BA"
        }}
      >
        <Template
          size="medium"
          color="#AA9AD4"
          style={{ marginLeft: "10px" }}
        />
        <StyledNavLink to="/dashboard/templates">Templates</StyledNavLink>
      </Box>
      <Box
        pad="xsmall"
        direction="row"
        style={{ borderTop: "1px solid #7058BA" }}
      />
    </Box>
    <Box>
      <Box
        pad="xsmall"
        direction="row"
        style={{ borderBottom: "1px solid #53418B" }}
      />
      <Box
        pad="xsmall"
        direction="row"
        style={{
          borderBottom: "1px solid #53418B",
          borderTop: "1px solid #7058BA"
        }}
      >
        <UserSettings
          size="medium"
          color="#AA9AD4"
          style={{ marginLeft: "10px" }}
        />
        <StyledNavLink to="/dashboard/settings">Settings</StyledNavLink>
      </Box>
      <Box
        pad="xsmall"
        direction="row"
        style={{
          borderBottom: "1px solid #53418B",
          borderTop: "1px solid #7058BA"
        }}
      >
        <Logout size="medium" color="#AA9AD4" style={{ marginLeft: "10px" }} />
        <StyledNavLink to="/logout">Logout</StyledNavLink>
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
            elevation="small"
            direction="column"
            justify="between"
            background="brand"
            style={{ position: "", minHeight: "100vh" }}
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

Sidebar.propTypes = {
  showSidebar: PropTypes.bool,
  size: PropTypes.string,
  closeSidebar: PropTypes.func
};

export default Sidebar;
