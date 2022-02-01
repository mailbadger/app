import React, { useState, useContext } from "react"
import PropTypes from "prop-types"
import { mainInstance as axios } from "../../../network/axios"
import { Box, Button } from "grommet"
import { ButtonWithLoader } from "../../../ui"
import { NotificationsContext } from "../../../Notifications/context"
import { endpoints } from "../../../network/endpoints"

export const DeleteForm = ({ id, callApi, hideModal }) => {
    const { createNotification } = useContext(NotificationsContext)

    const deleteTemplate = async (id) => {
        await axios.delete(endpoints.deleleteTemplates(id))
    }

    const [isSubmitting, setSubmitting] = useState(false)
    return (
        <Box direction="row" alignSelf="end" pad="small">
            <Box margin={{ right: "small" }}>
                <Button label="Cancel" onClick={() => hideModal()} />
            </Box>
            <Box>
                <ButtonWithLoader
                    primary
                    label="Delete"
                    color="#FF4040"
                    disabled={isSubmitting}
                    onClick={async () => {
                        setSubmitting(true)
                        try {
                            await deleteTemplate(id)
                            await callApi({ url: endpoints.getTemplates })
                            hideModal()
                        } catch (e) {
                            console.error(e)
                            createNotification(
                                "Unable to delete template, please try again.",
                                "status-error"
                            )
                        }

                        setSubmitting(false)
                    }}
                />
            </Box>
        </Box>
    )
}
DeleteForm.propTypes = {
    id: PropTypes.number,
    callApi: PropTypes.func,
    hideModal: PropTypes.func,
}
