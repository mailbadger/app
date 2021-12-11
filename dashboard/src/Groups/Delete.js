import React, { useState } from "react";
import PropTypes from "prop-types";
import { mainInstance as axios } from "../axios";
import { Box, Button } from "grommet";
import { ButtonWithLoader } from "../ui";
import { endpoints } from "../network/endpoints";

const DeleteSegment = ({ id, onSuccess, onCancel }) => {
  const deleteSegment = async (id) => {
    await axios.delete(endpoints.deleteGroups(id));
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
          label="Delete"
          color="#FF4040"
          disabled={isSubmitting}
          onClick={async () => {
            setSubmitting(true);
            await deleteSegment(id);
            setSubmitting(false);
            onSuccess();
          }}
        />
      </Box>
    </Box>
  );
};

DeleteSegment.propTypes = {
  id: PropTypes.number,
  callApi: PropTypes.func,
  onSuccess: PropTypes.func,
  onCancel: PropTypes.func,
};

export default DeleteSegment;
