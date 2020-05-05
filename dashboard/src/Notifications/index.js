import React from "react";
import { Box, Button, Layer, Text } from "grommet";
import { FormClose } from "grommet-icons";

import { NotificationConsumer } from "./context";
import { StatusIcons } from "../ui";

const Notification = () => (
  <NotificationConsumer>
    {({ close, message, status, isOpen }) =>
      isOpen ? (
        <Layer
          position="bottom"
          modal={false}
          margin={{ vertical: "medium", horizontal: "small" }}
          onEsc={close}
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
              {StatusIcons[status]}
              <Text>{message}</Text>
            </Box>
            <Button icon={<FormClose />} onClick={close} plain />
          </Box>
        </Layer>
      ) : null
    }
  </NotificationConsumer>
);

export default Notification;
