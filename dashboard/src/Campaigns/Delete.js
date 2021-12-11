import React, { useState } from "react";
import PropTypes from "prop-types";
import { mainInstance as axios } from "../axios";
import { Box, Button } from "grommet";
import { ButtonWithLoader } from "../ui";
import { endpoints } from "../network/endpoints";

const DeleteCampaign = ({ id, onSuccess, hideModal }) => {
  const deleteCampaign = async (id) => {
    await axios.delete(endpoints.deleteCampaigns(id));
  };

  const [isSubmitting, setSubmitting] = useState(false);
  return (
    <Box direction="row" alignSelf="end" pad="small">
      <Box margin={{ right: "small" }}>
        <Button label="Cancel" onClick={() => hideModal()} />
      </Box>
      <Box>
        <ButtonWithLoader
          primary
          label="Delete"
          color="#FF4040"
          disabled={isSubmitting}
          onClick={async () => {
            setSubmitting(true);
            await deleteCampaign(id);
            setSubmitting(false);
            onSuccess();
            hideModal();
          }}
        />
      </Box>
    </Box>
  );
};

DeleteCampaign.propTypes = {
  id: PropTypes.number,
  onSuccess: PropTypes.func,
  hideModal: PropTypes.func,
};

export default DeleteCampaign;
