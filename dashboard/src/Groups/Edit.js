import React, { useContext } from "react";
import PropTypes from "prop-types";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import qs from "qs";
import { Box, Button, FormField, TextInput } from "grommet";

import { mainInstance as axios } from "../axios";
import { ButtonWithLoader } from "../ui";
import { NotificationsContext } from "../Notifications/context";
import { endpoints } from "../network/endpoints";

const segmentValidation = object().shape({
  name: string()
    .required("Please enter a group name.")
    .max(191, "The name must not exceed 191 characters."),
});

const EditForm = ({
  handleSubmit,
  handleChange,
  isSubmitting,
  hideModal,
  values,
}) => (
  <Box
    direction="column"
    fill
    margin={{ left: "medium", right: "medium", bottom: "medium" }}
  >
    <form onSubmit={handleSubmit}>
      <Box>
        <FormField htmlFor="name" label="Group Name">
          <TextInput
            name="name"
            value={values.name}
            onChange={handleChange}
            placeholder="My group"
          />
          <ErrorMessage name="name" />
        </FormField>
        <Box direction="row" alignSelf="end" margin={{ top: "medium" }}>
          <Box margin={{ right: "small" }}>
            <Button label="Cancel" onClick={() => hideModal()} />
          </Box>
          <Box>
            <ButtonWithLoader
              type="submit"
              primary
              disabled={isSubmitting}
              label="Save Group"
            />
          </Box>
        </Box>
      </Box>
    </form>
  </Box>
);

EditForm.propTypes = {
  hideModal: PropTypes.func,
  handleSubmit: PropTypes.func,
  handleChange: PropTypes.func,
  isSubmitting: PropTypes.bool,
  values: PropTypes.shape({
    name: PropTypes.string,
  }),
};

const EditSegment = ({ setSegment, hideModal, segment }) => {
  const { createNotification } = useContext(NotificationsContext);

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const postForm = async () => {
      try {
        await axios.put(
          endpoints.putGroups(segment.id),
          qs.stringify({
            name: values.name,
          })
        );
        createNotification("Group has been edited successfully.");

        segment.name = values.name;
        setSegment(segment);

        //done submitting, set submitting to false
        setSubmitting(false);

        hideModal();
      } catch (error) {
        if (error.response) {
          const { message, errors } = error.response.data;

          setErrors(errors);

          const msg = message
            ? message
            : "Unable to edit group. Please try again.";

          createNotification(msg, "status-error");

          //done submitting, set submitting to false
          setSubmitting(false);
        }
      }
    };

    await postForm();
  };

  return (
    <Box direction="row">
      <Formik
        initialValues={{ name: segment.name }}
        onSubmit={handleSubmit}
        validationSchema={segmentValidation}
      >
        {(props) => <EditForm {...props} hideModal={hideModal} />}
      </Formik>
    </Box>
  );
};

EditSegment.propTypes = {
  setSegment: PropTypes.func,
  hideModal: PropTypes.func,
  segment: PropTypes.shape({
    id: PropTypes.number,
    name: PropTypes.string,
  }),
};

export default EditSegment;
