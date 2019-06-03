import React, { Fragment, useState, useEffect } from "react";
import { Controlled as CodeMirror } from "react-codemirror2";
import { FormField, Button, TextInput } from "grommet";
import { Formik, ErrorMessage } from "formik";
import { string, object } from "yup";
import axios from "axios";
import qs from "qs";

import history from "../history";
import useApi from "../hooks/useApi";

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
  values,
  errors
}) => {
  return (
    <Fragment>
      {errors && errors.message && <div>{errors.message}</div>}

      <form onSubmit={handleSubmit}>
        <FormField label="Subject" htmlFor="subject">
          <TextInput
            value={values.subject}
            name="subject"
            onChange={handleChange}
            placeholder="Greetings, {{name}}"
          />
          <ErrorMessage name="subject" />
        </FormField>

        <FormField label="HTML Template" htmlFor="htmlPart">
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
        <Button type="submit" primary disabled={isSubmitting} label="Submit" />
      </form>
    </Fragment>
  );
};

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

      history.push(`/dashboard/templates`);
    } catch (error) {
      setErrors(error.response.data);
    }
  };

  await callApi();

  //done submitting, set submitting to false
  setSubmitting(false);
};

const EditTemplateForm = ({ match }) => {
  const [html, setHtml] = useState();
  const [state] = useApi({
    url: `/api/templates/${match.params.id}`
  });

  useEffect(() => {
    if (!state.isLoading && state.data) {
      setHtml(state.data.html_part);
    }
  }, [state]);

  if (state.isLoading) {
    return <div>Loading...</div>;
  }

  if (state.isError) {
    return <div>Template not found.</div>;
  }

  return (
    <Fragment>
      {!state.isLoading && state.data && (
        <Formik
          onSubmit={handleSubmit(match.params.id)}
          validationSchema={templateValidation}
          initialValues={{
            subject: state.data.subject_part
          }}
          render={props => <Form setHtml={setHtml} html={html} {...props} />}
        />
      )}
    </Fragment>
  );
};

export default EditTemplateForm;
