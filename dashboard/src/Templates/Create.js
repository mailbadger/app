import React, { useState, useContext } from "react";
import { Controlled as CodeMirror } from "react-codemirror2";
import { FormField, Box, TextInput, Heading } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import { mainInstance as axios } from "../axios";
import qs from "qs";

import "codemirror/lib/codemirror.css";
import "codemirror/theme/material.css";
import "codemirror/mode/xml/xml";
import "codemirror/mode/javascript/javascript";

import { NotificationsContext } from "../Notifications/context";
import history from "../history";
import ButtonWithLoader from "../ui/ButtonWithLoader";
import { FormPropTypes } from "../PropTypes";

const initialHtml = `<!DOCTYPE html>
<html>
  <head>
    <meta charset="utf-8" />
    <title>Html example</title>
    <meta name="description" content="A simple HTML template" />
  </head>
  <body>
    <h1>Hello {{name}}</h1>
    <p>Your favorite animal is {{favoriteanimal}}.</p>
  </body>
</html>`;

const templateValidation = object().shape({
  name: string().required("Please enter a template name."),
  subject: string().required("Please enter a subject for the email."),
  htmlPart: string().required("Please enter a valid HTML"),
});

const Form = ({
  setHtml,
  html,
  handleSubmit,
  handleChange,
  setFieldValue,
  isSubmitting,
}) => {
  return (
    <Box direction="column">
      <form onSubmit={handleSubmit}>
        <Box>
          <FormField htmlFor="name" label="Template Name">
            <TextInput
              name="name"
              onChange={handleChange}
              placeholder="MyTemplate"
            />
            <ErrorMessage name="name" />
          </FormField>
        </Box>
        <Box margin={{ top: "small" }}>
          <FormField htmlFor="subject" label="Template Subject">
            <TextInput
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
                setFieldValue("htmlPart", editor.getValue(), true);
              }}
            />
            <ErrorMessage name="htmlPart" />
          </FormField>
        </Box>
        <Box margin={{ top: "medium" }} align="start">
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

const CreateTemplateForm = () => {
  const [html, setHtml] = useState(initialHtml);
  const { createNotification } = useContext(NotificationsContext);

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.post(
          "/api/templates",
          qs.stringify({
            name: values.name,
            content: values.htmlPart,
            subject: values.subject,
          })
        );
        createNotification("Template has been created successfully.");

        history.push(`/dashboard/templates`);
      } catch (error) {
        if (error.response) {
          const { message, errors } = error.response.data;

          setErrors(errors);
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

    return;
  };

  return (
    <Box
      direction="column"
      margin="medium"
      background="#ffffff"
      elevation="medium"
      animation="fadeIn"
    >
      <Box pad={{ left: "medium" }} margin={{ bottom: "small" }}>
        <Heading size={3}>Create Template</Heading>
      </Box>
      <Box pad={{ left: "medium", right: "medium", bottom: "medium" }} fill>
        <Formik
          onSubmit={handleSubmit}
          validationSchema={templateValidation}
          initialValues={{
            htmlPart: html,
          }}
        >
          {(props) => <Form setHtml={setHtml} html={html} {...props} />}
        </Formik>
      </Box>
    </Box>
  );
};

export default CreateTemplateForm;
