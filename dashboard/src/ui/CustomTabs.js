import React from "react";
import PropTypes from "prop-types";
import { Tab, Box, Tabs, ThemeContext } from "grommet";

const CustomTabs = (props) => {
  return (
    <ThemeContext.Extend
      value={{
        tab: {
          color: "text",
          active: {
            background: "light-3",
          },
          hover: {
            background: "background-back",
            color: "control",
          },
          pad: "small",
          margin: "none",
          extend: {
            fontWeight: "bold",
          },
        },
        tabs: {
          header: {
            extend: {
              borderTopLeftRadius: "12px",
              borderTopRightRadius: "12px",
            },
          },
        },
      }}
    >
      <Tabs flex>
        {props.tabs.map((tab, index) => (
          <Tab key={index} title={tab.title}>
            <Box round={{ corner: "bottom", size: "small" }} background="white">
              {tab.children}
            </Box>
          </Tab>
        ))}
      </Tabs>
    </ThemeContext.Extend>
  );
};

CustomTabs.propTypes = {
  tabs: PropTypes.arrayOf(
    PropTypes.shape({
      icon: PropTypes.element,
      title: PropTypes.string,
      children: PropTypes.element.isRequired,
    })
  ),
};

export default CustomTabs;
