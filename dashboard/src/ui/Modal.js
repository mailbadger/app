import React from "react"
import PropTypes from "prop-types"
import { Layer, Box, Heading } from "grommet"

const Modal = ({ hideModal, title, form }) => {
    return (
        <Layer onEsc={() => hideModal()} onClickOutside={() => hideModal()}>
            <Box width="30em">
                <Heading margin="small" level="3">
                    {title}
                </Heading>
                {form}
            </Box>
        </Layer>
    )
}

Modal.propTypes = {
    hideModal: PropTypes.func,
    title: PropTypes.string,
    form: PropTypes.element.isRequired,
}

export default Modal
