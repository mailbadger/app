import React, { Fragment } from "react";
import { FormField, Button, TextInput } from "grommet";
import { Formik, ErrorMessage } from "formik";
import "codemirror/lib/codemirror.css";
import "codemirror/theme/material.css";

import "codemirror/mode/xml/xml";
import "codemirror/mode/javascript/javascript";
import { UnControlled as CodeMirror } from "react-codemirror2";

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
          <TextInput name="name" onChange={handleChange} />
          <ErrorMessage name="name" />
        </FormField>
        <FormField label="Subject" htmlFor="subject">
          <TextInput name="subject" onChange={handleChange} />
          <ErrorMessage name="subject" />
        </FormField>

        <CodeMirror
          value={values.html_part}
          options={{
            mode: "xml",
            theme: "material",
            lineNumbers: true
          }}
          onChange={editor => {
            console.log(editor.getValue());
            setFieldValue("html_part", editor.getValue(), true);
          }}
        />

        <Button type="submit" primary label="Submit" />
      </form>
    </Fragment>
  );
};

const CreateTemplateForm = () => {
  return (
    <Formik
      onSubmit={values => console.log(values)}
      // validationSchema={loginValidation}
      render={Form}
    />
  );
};

export default CreateTemplateForm;
