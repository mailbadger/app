import PropTypes from "prop-types";

export let FormPropTypes = {
  handleSubmit: PropTypes.func.isRequired,
  handleChange: PropTypes.func.isRequired,
  isSubmitting: PropTypes.bool.isRequired,
  errors: PropTypes.shape({
    message: PropTypes.string
  })
};
