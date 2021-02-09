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
import { Box, Button, Layer } from "grommet";
import { AnchorLink, GradientBadger, UserMenu } from "./ui";

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
        {links.map((link) => (
          <Box pad="xsmall" direction="row" key={link.label}>
            <AnchorLink
              to={link.to}
              size="medium"
              icon={link.icon}
              active={
                active === link.label || location.pathname.startsWith(link.to)
              }
              onClick={() => setActive(link.label)}
            />
          </Box>
        ))}
    </Fragment>
  );
};

const Sidebar = (props) => {
  const { showSidebar, size, closeSidebar } = props;

  return (
    <Fragment>
      {!showSidebar || size !== "small" ? (
          <Box
            overflow="auto"
            background="brand"
          >
            <Box align="center" pad={{ vertical: 'small' }}>
              <GradientBadger />
            </Box>
            <Box align="center" gap={size === 'small' ? 'medium' : 'small'}>
              <NavLinks />
            </Box>
            <Box flex />
            <Box pad={{ vertical: 'small' }}>
              <UserMenu alignSelf="center" />
            </Box>
          </Box>
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
