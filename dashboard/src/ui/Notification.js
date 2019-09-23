import React from "react";
import { Box, Button, Layer, Text } from "grommet";
import { FormClose, StatusGood } from "grommet-icons";

const Notification = ({ onClose, message, status = "status-ok" }) => (
  <Layer
    position="bottom"
    modal={false}
    margin={{ vertical: "medium", horizontal: "small" }}
    onEsc={onClose}
    responsive={false}
    plain
  >
    <Box
      align="center"
      direction="row"
      gap="small"
      justify="between"
      round="medium"
      elevation="medium"
      pad={{ vertical: "xsmall", horizontal: "small" }}
      background={status}
    >
      <Box align="center" direction="row" gap="xsmall">
        <StatusGood />
        <Text>{message}</Text>
      </Box>
      <Button icon={<FormClose />} onClick={onClose} plain />
    </Box>
  </Layer>
);

export default Notification;
