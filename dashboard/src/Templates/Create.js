import React, { Fragment, useState } from "react";
import { Controlled as CodeMirror } from "react-codemirror2";
import { FormField } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import axios from "axios";
import qs from "qs";

import "codemirror/lib/codemirror.css";
import "codemirror/theme/material.css";
import "codemirror/mode/xml/xml";
import "codemirror/mode/javascript/javascript";

import history from "../history";

import ButtonWithLoader from "../ui/ButtonWithLoader";
import StyledTextInput from "../ui/StyledTextInput";

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
  htmlPart: string().required("Please enter a valid HTML")
});

const Form = ({
  setHtml,
  html,
  handleSubmit,
  handleChange,
  setFieldValue,
  isSubmitting,
  errors
}) => {
  return (
    <Fragment>
      {errors && errors.message && <div>{errors.message}</div>}

      <form onSubmit={handleSubmit}>
        <FormField htmlFor="name">
          Template Name
          <StyledTextInput
            name="name"
            onChange={handleChange}
            placeholder="MyTemplate"
          />
          <ErrorMessage name="name" />
        </FormField>
        <br />
        Template Subject
        <FormField htmlFor="subject">
          <StyledTextInput
            name="subject"
            onChange={handleChange}
            placeholder="Greetings, {{name}}"
          />
          <ErrorMessage name="subject" />
        </FormField>
        <br />
        HTML Content
        <FormField htmlFor="htmlPart">
          <CodeMirror
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
        <br />
        <ButtonWithLoader
          type="submit"
          primary
          disabled={isSubmitting}
          label="Save Template"
        />
      </form>
    </Fragment>
  );
};

const CreateTemplateForm = () => {
  const [html, setHtml] = useState(initialHtml);

  const handleSubmit = async (values, { setSubmitting, setErrors }) => {
    const callApi = async () => {
      try {
        await axios.post(
          "/api/templates",
          qs.stringify({
            name: values.name,
            content: values.htmlPart,
            subject: values.subject
          })
        );

        history.push(`/dashboard/templates`);
      } catch (error) {
        setErrors(error.response.data);
      }
    };

    await callApi();

    //done submitting, set submitting to false
    setSubmitting(false);

    return;
  };

  return (
    <Formik
      onSubmit={handleSubmit}
      validationSchema={templateValidation}
      initialValues={{
        htmlPart: html
      }}
      render={props => <Form setHtml={setHtml} html={html} {...props} />}
    />
  );
};

export default CreateTemplateForm;
