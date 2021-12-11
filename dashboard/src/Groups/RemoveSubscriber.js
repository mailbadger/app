import React, { useState, useContext } from "react";
import PropTypes from "prop-types";
import { Box, Button } from "grommet";

import { mainInstance as axios } from "../axios";
import { ButtonWithLoader } from "../ui";
import { NotificationsContext } from "../Notifications/context";
import { endpoints } from "../network/endpoints";

const RemoveSubscriber = ({ id, subId, onSuccess, onCancel }) => {
  const { createNotification } = useContext(NotificationsContext);

  const removeSubscriber = async (id) => {
    await axios.delete(endpoints.deleteSubscriberFromGroup(id,subId));
  };

  const [isSubmitting, setSubmitting] = useState(false);
  return (
    <Box direction="row" alignSelf="end" pad="small">
      <Box margin={{ right: "small" }}>
        <Button label="Cancel" onClick={onCancel} />
      </Box>
      <Box>
        <ButtonWithLoader
          primary
          label="Remove"
          color="#FF4040"
          disabled={isSubmitting}
          onClick={async () => {
            setSubmitting(true);
            try {
              await removeSubscriber(id);
              onSuccess();
            } catch (e) {
              let msg = "Unable to remove subscriber from group. Try again.";
              const { response } = e;
              if (response) {
                msg = response.data.message;
              }
              createNotification(msg, "status-error");
            }
            setSubmitting(false);
          }}
        />
      </Box>
    </Box>
  );
};

RemoveSubscriber.propTypes = {
  id: PropTypes.number,
  subId: PropTypes.number,
  callApi: PropTypes.func,
  onSuccess: PropTypes.func,
  onCancel: PropTypes.func,
};

export default RemoveSubscriber;
