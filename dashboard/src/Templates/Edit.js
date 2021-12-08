import React, { useState, useEffect, useContext } from "react";
import PropTypes from "prop-types";
import { Controlled as CodeMirror } from "react-codemirror2";
import { FormField, Box, Heading, ResponsiveContext } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import { mainInstance as axios } from "../axios";
import qs from "qs";
import { ButtonWithLoader, StyledTextInput, StyledSpinner } from "../ui";
import history from "../history";
import { useApi } from "../hooks";
import { NotificationsContext } from "../Notifications/context";
import { FormPropTypes } from "../PropTypes";
import { endpoints } from "../network/endpoints";

const templateValidation = object().shape({
  name: string().required("Please enter a template name."),
  subject: string().required("Please enter a subject for the email."),
  html_part: string().required("Please enter a valid HTML"),
});

const Form = ({
  html,
  setHtml,
  handleSubmit,
  handleChange,
  setFieldValue,
  isSubmitting,
  values,
}) => {
  return (
    <Box direction="column">
      <form onSubmit={handleSubmit}>
        <Box>
          <FormField htmlFor="name" label="Template Name">
            <StyledTextInput
              value={values.name}
              name="name"
              onChange={handleChange}
              placeholder="HelloWorld"
            />
            <ErrorMessage name="name" />
          </FormField>
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
          <FormField htmlFor="html_part" label="HTML Content">
            <CodeMirror
              style={{ height: "100%" }}
              value={html}
              options={{
                mode: "xml",
                theme: "material",
                lineNumbers: true,
              }}
              onBeforeChange={(editor, data, value) => {
                setHtml(value);
              }}
              onChange={(editor) => {
                setFieldValue("html_part", editor.getValue(), true);
              }}
            />
            <ErrorMessage name="html_part" />
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
  const [success, setSuccess] = useState(false);

  const [state] = useApi({
    url: endpoints.putTemplates(match.params.id),
  });

  const { createNotification } = useContext(NotificationsContext);
  const size = useContext(ResponsiveContext);

  let width = "100%";
  if (size === "large") {
    width = "60%";
  }

  const handleSubmit = (id) => async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.put(
          endpoints.putTemplates(id),
          qs.stringify({
            name: values.name,
            html_part: values.html_part,
            text_part: values.html_part,
            subject: values.subject,
          })
        );

        createNotification("Template has been updated successfully.");
        setSuccess(true);
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
    if (success) {
      history.push("/dashboard/templates");
    }
  }, [success]);

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
    <Box direction="column" margin="medium" animation="fadeIn">
      {!state.isLoading && state.data && (
        <>
          <Box pad={{ left: "medium" }} margin={{ bottom: "small" }}>
            <Heading level="2">Edit Template</Heading>
          </Box>
          <Box background="white" pad="medium" width={width} alignSelf="start">
            <Formik
              onSubmit={handleSubmit(match.params.id)}
              validationSchema={templateValidation}
              initialValues={{
                subject: state.data.subject_part,
                name: state.data.name,
              }}
            >
              {(props) => <Form setHtml={setHtml} html={html} {...props} />}
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
      id: PropTypes.string,
    }),
  }),
};

export default EditTemplateForm;
