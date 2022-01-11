import React, { useState } from "react"
import PropTypes from "prop-types"
import { Box, Button } from "grommet"

import { mainInstance as axios } from "../../network/axios"
import { ButtonWithLoader } from "../../ui"
import { endpoints } from "../../network/endpoints"

const DeleteSubscriber = ({ id, callApi, hideModal }) => {
    const deleteSubscriber = async (id) => {
        await axios.delete(endpoints.deleteSubscribers(id))
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
                        await deleteSubscriber(id)
                        await callApi({ url: endpoints.getGroups })
                        setSubmitting(false)
                        hideModal()
                    }}
                />
            </Box>
        </Box>
    )
}

DeleteSubscriber.propTypes = {
    id: PropTypes.number,
    callApi: PropTypes.func,
    hideModal: PropTypes.func,
}

export default DeleteSubscriber
