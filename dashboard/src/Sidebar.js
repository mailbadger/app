import React, { Fragment, useState } from "react";
import PropTypes from "prop-types";
import {
  FormClose,
  Configure,
  Group,
  List,
  Send,
  Template,
} from "grommet-icons";
import { useLocation } from "react-router-dom";
import { Box, Button, Collapsible, Layer } from "grommet";
import { AnchorLink } from "./ui";

const links = [
  {
    to: "/dashboard/subscribers",
    label: "Subscribers",
    icon: <Group size="medium" />,
  },
  {
    to: "/dashboard/segments",
    label: "Segments",
    icon: <List size="medium" />,
  },
  {
    to: "/dashboard/campaigns",
    label: "Campaigns",
    icon: <Send size="medium" />,
  },
  {
    to: "/dashboard/templates",
    label: "Templates",
    icon: <Template size="medium" />,
  },
  {
    to: "/dashboard/settings",
    label: "Settings",
    icon: <Configure size="medium" />,
  },
];

const NavLinks = () => {
  let location = useLocation();
  const [active, setActive] = useState();

  return (
    <Fragment>
      <Box margin={{ top: "small" }} pad="large">
        {links.map((link) => (
          <Box pad="xsmall" direction="row" key={link.label}>
            <AnchorLink
              to={link.to}
              size="medium"
              icon={link.icon}
              label={link.label}
              active={
                active === link.label || location.pathname.startsWith(link.to)
              }
              onClick={() => setActive(link.label)}
            />
          </Box>
        ))}
      </Box>
    </Fragment>
  );
};

const Sidebar = (props) => {
  const { showSidebar, size, closeSidebar } = props;

  return (
    <Fragment>
      {!showSidebar || size !== "small" ? (
        <Collapsible direction="horizontal" open={showSidebar}>
          <Box flex width="18em" direction="column" justify="between">
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
  closeSidebar: PropTypes.func,
};

export default Sidebar;
