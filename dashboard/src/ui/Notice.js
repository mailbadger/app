import React, { useState } from "react";
import PropTypes from "prop-types";
import { Box, Button, Text } from "grommet";
import { FormClose } from "grommet-icons";

import StatusIcons from "./StatusIcons";

const Notice = ({ message, status }) => {
  const [closed, setClosed] = useState(false);
  if (closed) {
    return null;
  }

  return (
    <Box
      border={{
        side: "all",
        size: "small",
        color: "dark-2",
      }}
      align="center"
      direction="row"
      gap="small"
      justify="between"
      elevation="medium"
      pad={{ vertical: "xsmall", horizontal: "small" }}
      background={status}
    >
      <Box align="center" direction="row" gap="xsmall">
        {StatusIcons[status]}
        <Text>{message}</Text>
      </Box>
      <Button icon={<FormClose />} onClick={() => setClosed(true)} plain />
    </Box>
  );
};

Notice.propTypes = {
  message: PropTypes.string,
  status: PropTypes.string,
};

export default Notice;
