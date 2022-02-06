import React from "react"
import PropTypes from "prop-types"
import { ErrorMessage } from "formik"
import { Box, Button, FormField, TextInput } from "grommet"
import { ButtonWithLoader } from "../../../ui"

export const CreateForm = ({
    handleSubmit,
    handleChange,
    isSubmitting,
    hideModal,
}) => (
    <Box
        direction="column"
        fill
        margin={{ left: "medium", right: "medium", bottom: "medium" }}
    >
        <form onSubmit={handleSubmit}>
            <Box>
                <FormField htmlFor="name" label="Group Name">
                    <TextInput
                        name="name"
                        onChange={handleChange}
                        placeholder="My group"
                    />
                    <ErrorMessage name="name" />
                </FormField>
                <Box direction="row" alignSelf="end" margin={{ top: "medium" }}>
                    <Box margin={{ right: "small" }}>
                        <Button label="Cancel" onClick={() => hideModal()} />
                    </Box>
                    <Box>
                        <ButtonWithLoader
                            type="submit"
                            primary
                            disabled={isSubmitting}
                            label="Save Group"
                        />
                    </Box>
                </Box>
            </Box>
        </form>
    </Box>
)

CreateForm.propTypes = {
    hideModal: PropTypes.func,
    handleSubmit: PropTypes.func,
    handleChange: PropTypes.func,
    isSubmitting: PropTypes.bool,
}
