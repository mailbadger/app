import React, { useContext } from "react"
import PropTypes from "prop-types"
import { mainInstance as axios } from "../../../network/axios"
import { Formik } from "formik"
import { Box } from "grommet"
import { NotificationsContext } from "../../../Notifications/context"
import { endpoints } from "../../../network/endpoints"
import { CreateForm } from "./CreateForm"
import { segmentValidation } from "./utils"

export const CreateSegment = ({ callApi, hideModal }) => {
    const { createNotification } = useContext(NotificationsContext)

    const handleSubmit = async (values, { setSubmitting, setErrors }) => {
        const postForm = async () => {
            const params = {
                name: values.name,
            }
            try {
                await axios.post(endpoints.postGroups, params)
                createNotification("Group has been created successfully.")

                await callApi({ url: endpoints.getGroups })

                //done submitting, set submitting to false
                setSubmitting(false)

                hideModal()
            } catch (error) {
                if (error.response) {
                    const { message, errors } = error.response.data

                    setErrors(errors)

                    const msg = message
                        ? message
                        : "Unable to create group. Please try again."

                    createNotification(msg, "status-error")

                    //done submitting, set submitting to false
                    setSubmitting(false)
                }
            }
        }

        await postForm()

        return
    }

    return (
        <Box direction="row">
            <Formik
                initialValues={{ name: "" }}
                onSubmit={handleSubmit}
                validationSchema={segmentValidation}
            >
                {(props) => <CreateForm {...props} hideModal={hideModal} />}
            </Formik>
        </Box>
    )
}
CreateSegment.propTypes = {
    callApi: PropTypes.func,
    hideModal: PropTypes.func,
}
