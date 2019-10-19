import React from "react";
import PropTypes from "prop-types";
import { Tab, Box, Tabs } from "grommet";

const CustomTabs = props => {
  return (
    <Tabs justify="start">
      {props.tabs.map((tab, index) => (
        <Tab margin="small" key={index} title={tab.title}>
          <Box pad="medium">{tab.children}</Box>
        </Tab>
      ))}
    </Tabs>
  );
};

CustomTabs.propTypes = {
  tabs: PropTypes.arrayOf(
    PropTypes.shape({
      title: PropTypes.string,
      children: PropTypes.element.isRequired
    })
  )
};

export default CustomTabs;
