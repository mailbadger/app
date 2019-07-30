import React from "react";
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

export default CustomTabs;
