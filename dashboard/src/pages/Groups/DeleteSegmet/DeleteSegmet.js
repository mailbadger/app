import React, { useState } from "react"
import PropTypes from "prop-types"
import { mainInstance as axios } from "../../../network/axios"
import { Box, Button } from "grommet"
import { ButtonWithLoader } from "../../../ui"
import { endpoints } from "../../../network/endpoints"

// This component is not used
const DeleteSegmet = ({ id, callApi, hideModal }) => {
    const deleteSegment = async (id) => {
        await axios.delete(endpoints.deleteGroups(id))
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
                        await deleteSegment(id)
                        await callApi({ url: endpoints.getGroups })
                        setSubmitting(false)
                        hideModal()
                    }}
                />
            </Box>
        </Box>
    )
}
DeleteSegmet.propTypes = {
    id: PropTypes.number,
    callApi: PropTypes.func,
    hideModal: PropTypes.func,
}
