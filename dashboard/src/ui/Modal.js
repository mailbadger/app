import React from "react";
import PropTypes from "prop-types";
import { Layer, Box, Heading } from "grommet";

const Modal = ({ hideModal, title, form, width = "30em" }) => {
  return (
    <Layer onEsc={() => hideModal()} onClickOutside={() => hideModal()}>
      <Box width={width}>
        <Heading
          margin={{ top: "20px", bottom: "20px", left: "20px" }}
          level="3"
        >
          {title}
        </Heading>
        {form}
      </Box>
    </Layer>
  );
};

Modal.propTypes = {
  hideModal: PropTypes.func,
  title: PropTypes.string || PropTypes.element,
  form: PropTypes.element.isRequired,
  width: PropTypes.string,
};

export default Modal;
