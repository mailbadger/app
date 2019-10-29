import React, { useState, useEffect, useContext } from "react";
import PropTypes from "prop-types";
import { Controlled as CodeMirror } from "react-codemirror2";
import { FormField, Box, Heading } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import axios from "axios";
import qs from "qs";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import StyledTextInput from "../ui/StyledTextInput";
import StyledSpinner from "../ui/StyledSpinner";
import history from "../history";
import useApi from "../hooks/useApi";
import { NotificationsContext } from "../Notifications/context";
import { FormPropTypes } from "../PropTypes";

const templateValidation = object().shape({
  subject: string().required("Please enter a subject for the email."),
  htmlPart: string().required("Please enter a valid HTML")
});

const Form = ({
  html,
  setHtml,
  handleSubmit,
  handleChange,
  setFieldValue,
  isSubmitting,
  values
}) => {
  return (
    <Box direction="column">
      <form onSubmit={handleSubmit}>
        <Box>
          <FormField htmlFor="subject" label="Template Subject">
            <StyledTextInput
              value={values.subject}
              name="subject"
              onChange={handleChange}
              placeholder="Greetings, {{name}}"
            />
            <ErrorMessage name="subject" />
          </FormField>
        </Box>
        <Box margin={{ top: "small" }}>
          <FormField htmlFor="htmlPart" label="HTML Content">
            <CodeMirror
              style={{ height: "100%" }}
              value={html}
              options={{
                mode: "xml",
                theme: "material",
                lineNumbers: true
              }}
              onBeforeChange={(editor, data, value) => {
                setHtml(value);
              }}
              onChange={editor => {
                setFieldValue("htmlPart", editor.getValue(), true);
              }}
            />
            <ErrorMessage name="htmlPart" />
          </FormField>
        </Box>
        <Box margin={{ top: "small" }} align="start">
          <ButtonWithLoader
            type="submit"
            primary
            disabled={isSubmitting}
            label="Save Template"
          />
        </Box>
      </form>
    </Box>
  );
};

Form.propTypes = FormPropTypes;

const EditTemplateForm = ({ match }) => {
  const [html, setHtml] = useState();
  const [state] = useApi({
    url: `/api/templates/${match.params.id}`
  });

  const { createNotification } = useContext(NotificationsContext);

  const handleSubmit = id => async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.put(
          `/api/templates/${id}`,
          qs.stringify({
            content: values.htmlPart,
            subject: values.subject
          })
        );

        createNotification("Template has been updated successfully.");

        history.push(`/dashboard/templates`);
      } catch (error) {
        if (error.response) {
          setErrors(error.response.data);

          const { message } = error.response.data;
          const msg = message
            ? message
            : "Unable to create template. Please try again.";

          createNotification(msg, "status-error");
        }
      }
    };

    await callApi();

    //done submitting, set submitting to false
    setSubmitting(false);
  };

  useEffect(() => {
    if (!state.isLoading && state.data) {
      setHtml(state.data.html_part);
    }
  }, [state]);

  if (state.isLoading) {
    return (
      <Box margin="15%" alignSelf="center">
        <StyledSpinner color="#2e2e2e" size={8} />
      </Box>
    );
  }

  if (state.isError) {
    return (
      <Box margin="15%" alignSelf="center">
        <Heading>Template not found.</Heading>
      </Box>
    );
  }

  return (
    <Box
      direction="column"
      margin="medium"
      background="#ffffff"
      elevation="medium"
      animation="fadeIn"
    >
      {!state.isLoading && state.data && (
        <>
          <Box pad={{ left: "medium" }} margin={{ bottom: "small" }}>
            <Heading size={3}>Edit Template</Heading>
          </Box>
          <Box pad={{ left: "medium", right: "medium", bottom: "medium" }} fill>
            <Formik
              onSubmit={handleSubmit(match.params.id)}
              validationSchema={templateValidation}
              initialValues={{
                subject: state.data.subject_part
              }}
            >
              {props => <Form setHtml={setHtml} html={html} {...props} />}
            </Formik>
          </Box>
        </>
      )}
    </Box>
  );
};

EditTemplateForm.propTypes = {
  match: PropTypes.shape({
    params: PropTypes.shape({
      id: PropTypes.string
    })
  })
};

export default EditTemplateForm;
