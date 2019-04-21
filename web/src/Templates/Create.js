import React, { Fragment } from "react";
import { FormField, Button, TextInput } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/material.css";

import "codemirror/mode/xml/xml";
import "codemirror/mode/javascript/javascript";
import { UnControlled as CodeMirror } from "react-codemirror2";

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
  handleSubmit,
  handleChange,
  setFieldValue,
  values,
  errors
}) => {
  return (
    <Fragment>
      {errors && errors.message && <div>{errors.message}</div>}

      <form onSubmit={handleSubmit}>
        <FormField label="Template Name" htmlFor="name">
          <TextInput
            name="name"
            onChange={handleChange}
            placeholder="MyTemplate"
          />
          <ErrorMessage name="name" />
        </FormField>
        <FormField label="Subject" htmlFor="subject">
          <TextInput
            name="subject"
            onChange={handleChange}
            placeholder="Greetings, {{name}}"
          />
          <ErrorMessage name="subject" />
        </FormField>

        <FormField label="HTML Template" htmlFor="htmlPart">
          <CodeMirror
            value={values.htmlPart}
            options={{
              mode: "xml",
              theme: "material",
              lineNumbers: true
            }}
            onChange={editor => {
              console.log(editor.getValue());
              setFieldValue("htmlPart", editor.getValue(), true);
            }}
          />
          <ErrorMessage name="htmlPart" />
        </FormField>
        <Button type="submit" primary label="Submit" />
      </form>
    </Fragment>
  );
};

const CreateTemplateForm = () => {
  return (
    <Formik
      onSubmit={values => console.log(values)}
      validationSchema={templateValidation}
      initialValues={{
        htmlPart: initialHtml
      }}
      render={Form}
    />
  );
};

export default CreateTemplateForm;
