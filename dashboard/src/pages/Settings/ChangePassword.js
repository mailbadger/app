import React, { useContext } from "react"
import { FormField, TextInput, Heading, Box } from "grommet"
import { Formik, ErrorMessage } from "formik"
import { string, object, ref, addMethod } from "yup"
import { mainInstance as axios } from "../../network/axios"

import { NotificationsContext } from "../../Notifications/context"
import equalTo from "../../utils/equalTo"
import { ButtonWithLoader } from "../../ui"
import { FormPropTypes } from "../../utils/PropTypes"
import { endpoints } from "../../network/endpoints"

addMethod(string, "equalTo", equalTo)

const changePassValidation = object().shape({
    password: string().required("Please enter your password."),
    new_password: string()
        .min(8, "Your password must be atleast 8 characters.")
        .required("Password must not be empty."),
    new_password_confirm: string()
        .equalTo(ref("new_password"), "Passwords don't match")
        .required("Confirm Password is required"),
})

const Form = ({ handleSubmit, handleChange, isSubmitting }) => (
    <Box
        round={{
            corner: "bottom",
        }}
        pad="medium"
        alignSelf="stretch"
        background="white"
        animation="fadeIn"
        margin={{ bottom: "medium" }}
    >
        <Box alignSelf="center">
            <Heading level="4" color="#564392">
                Change password
            </Heading>
            <Box width="medium">
                <form onSubmit={handleSubmit}>
                    <FormField label="Old password" htmlFor="password">
                        <TextInput
                            name="password"
                            type="password"
                            onChange={handleChange}
                        />
                        <ErrorMessage name="password" />
                    </FormField>
                    <FormField label="New password" htmlFor="new_password">
                        <TextInput
                            name="new_password"
                            type="password"
                            onChange={handleChange}
                        />
                        <ErrorMessage name="new_password" />
                    </FormField>
                    <FormField
                        label="Confirm new password"
                        htmlFor="new_password_confirm"
                    >
                        <TextInput
                            name="new_password_confirm"
                            type="password"
                            onChange={handleChange}
                        />
                        <ErrorMessage name="new_password_confirm" />
                    </FormField>

                    <Box margin={{ top: "medium" }}>
                        <ButtonWithLoader
                            type="submit"
                            primary
                            disabled={isSubmitting}
                            label="Update password"
                        />
                    </Box>
                </form>
            </Box>
        </Box>
    </Box>
)

Form.propTypes = FormPropTypes

const ChangePasswordForm = () => {
    const { createNotification } = useContext(NotificationsContext)

    const handleSubmit = async (values, { setSubmitting, setErrors }) => {
        const callApi = async () => {
            const params = {
                password: values.password,
                new_password: values.new_password,
            }
            try {
                const res = await axios.post(
                    endpoints.postChangePassword,
                    params
                )

                createNotification(res.data.message)
            } catch (error) {
                setErrors(error.response.data)

                const { message } = error.response.data
                const msg = message
                    ? message
                    : "Unable to update your password. Please try again."

                createNotification(msg, "status-error")
            }
        }

        await callApi()

        //done submitting, set submitting to false
        setSubmitting(false)

        return
    }

    return (
        <Formik
            initialValues={{
                password: "",
                new_password: "",
                new_password_confirm: "",
            }}
            onSubmit={handleSubmit}
            validationSchema={changePassValidation}
        >
            {Form}
        </Formik>
    )
}

export default ChangePasswordForm
