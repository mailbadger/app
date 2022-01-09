import React, { useState, useContext, useEffect } from "react"
import { Controlled as CodeMirror } from "react-codemirror2"
import { FormField, Box, TextInput, Heading, ResponsiveContext } from "grommet"
import { Formik, ErrorMessage } from "formik"
import { string, object } from "yup"
import { mainInstance as axios } from "../axios"

import "codemirror/lib/codemirror.css"
import "codemirror/theme/material.css"
import "codemirror/mode/xml/xml"
import "codemirror/mode/javascript/javascript"

import { NotificationsContext } from "../Notifications/context"
import history from "../history"
import { ButtonWithLoader } from "../ui"
import { FormPropTypes } from "../PropTypes"
import { endpoints } from "../network/endpoints"

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
</html>`

const templateValidation = object().shape({
    name: string().required("Please enter a template name."),
    subject: string().required("Please enter a subject for the email."),
    html_part: string().required("Please enter a valid HTML"),
})

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
                    <FormField htmlFor="html_part" label="HTML Content">
                        <CodeMirror
                            value={html}
                            options={{
                                mode: "xml",
                                theme: "material",
                                lineNumbers: true,
                            }}
                            onBeforeChange={(editor, data, value) => {
                                setHtml(value)
                            }}
                            onChange={(editor) => {
                                setFieldValue(
                                    "html_part",
                                    editor.getValue(),
                                    true
                                )
                            }}
                        />
                        <ErrorMessage name="html_part" />
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
    )
}

Form.propTypes = FormPropTypes

const CreateTemplateForm = () => {
    const [html, setHtml] = useState(initialHtml)
    const [success, setSuccess] = useState(false)
    const { createNotification } = useContext(NotificationsContext)
    const size = useContext(ResponsiveContext)

    let width = "100%"
    if (size === "large") {
        width = "60%"
    }

    const handleSubmit = async (values, { setSubmitting, setErrors }) => {
        const callApi = async () => {
            const params = {
                name: values.name,
                html_part: values.html_part,
                text_part: values.html_part,
                subject: values.subject,
            }
            try {
                await axios.post(endpoints.postTemplates, params)
                createNotification("Template has been created successfully.")

                setSuccess(true)
            } catch (error) {
                if (error.response) {
                    const { message, errors } = error.response.data

                    setErrors(errors)
                    const msg = message
                        ? message
                        : "Unable to create template. Please try again."

                    createNotification(msg, "status-error")
                }
            }
        }

        await callApi()

        //done submitting, set submitting to false
        setSubmitting(false)

        return
    }

    useEffect(() => {
        if (success) {
            history.push("/dashboard/templates")
        }
    }, [success])

    return (
        <Box direction="column" margin="medium" animation="fadeIn">
            <Box pad={{ left: "medium" }} margin={{ bottom: "small" }}>
                <Heading level="2">Create Template</Heading>
            </Box>
            <Box
                round
                background="white"
                pad="medium"
                width={width}
                alignSelf="start"
            >
                <Formik
                    onSubmit={handleSubmit}
                    validationSchema={templateValidation}
                    initialValues={{
                        html_part: html,
                    }}
                >
                    {(props) => (
                        <Form setHtml={setHtml} html={html} {...props} />
                    )}
                </Formik>
            </Box>
        </Box>
    )
}

export default CreateTemplateForm
