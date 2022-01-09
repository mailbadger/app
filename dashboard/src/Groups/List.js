import React, { useState, useContext } from "react"
import PropTypes from "prop-types"
import { parseISO, formatRelative } from "date-fns"
import { More, Add, FormPreviousLink, FormNextLink } from "grommet-icons"
import { mainInstance as axios } from "../axios"
import { Formik, ErrorMessage } from "formik"
import { string, object } from "yup"

import { useApi } from "../hooks"
import {
    TableHeader,
    TableBody,
    TableRow,
    TableCell,
    Box,
    Button,
    Heading,
    Select,
    FormField,
    TextInput,
} from "grommet"
import history from "../history"
import {
    StyledTable,
    ButtonWithLoader,
    PlaceholderTable,
    Modal,
    AnchorLink,
} from "../ui"
import { NotificationsContext } from "../Notifications/context"
import DeleteSegment from "./Delete"
import { endpoints } from "../network/endpoints"

const Row = ({ segment, setShowDelete }) => {
    const ca = parseISO(segment.created_at)
    const ua = parseISO(segment.updated_at)
    return (
        <TableRow>
            <TableCell scope="row" size="xlarge">
                <AnchorLink
                    size="medium"
                    fontWeight="bold"
                    to={`/dashboard/groups/${segment.id}`}
                    label={segment.name}
                />
            </TableCell>
            <TableCell scope="row" size="xlarge">
                <strong>{segment.subscribers_in_segment}</strong>
            </TableCell>
            <TableCell scope="row" size="medium">
                {formatRelative(ca, new Date())}
            </TableCell>
            <TableCell scope="row" size="medium">
                {formatRelative(ua, new Date())}
            </TableCell>
            <TableCell scope="row" size="xsmall" align="end">
                <Select
                    alignSelf="center"
                    plain
                    icon={<More />}
                    options={["View", "Delete"]}
                    onChange={({ option }) => {
                        ;(function () {
                            switch (option) {
                                case "View":
                                    history.push(
                                        `/dashboard/groups/${segment.id}`
                                    )
                                    break
                                case "Delete":
                                    setShowDelete({
                                        show: true,
                                        name: segment.name,
                                        id: segment.id,
                                    })
                                    break
                                default:
                                    return null
                            }
                        })()
                    }}
                />
            </TableCell>
        </TableRow>
    )
}

Row.propTypes = {
    segment: PropTypes.shape({
        name: PropTypes.string,
        id: PropTypes.number,
        subscribers_in_segment: PropTypes.number,
        created_at: PropTypes.string,
        updated_at: PropTypes.string,
    }),
    setShowDelete: PropTypes.func,
}

const Header = () => (
    <TableHeader>
        <TableRow>
            <TableCell scope="col" border="bottom" size="small">
                <strong>Name</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="small">
                <strong>Total Subscribers</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="small">
                <strong>Created At</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="small">
                <strong>Updated At</strong>
            </TableCell>
            <TableCell align="end" scope="col" border="bottom" size="small">
                <strong>Action</strong>
            </TableCell>
        </TableRow>
    </TableHeader>
)

const SegmentTable = React.memo(({ list, setShowDelete }) => (
    <StyledTable>
        <Header />
        <TableBody>
            {list.map((s) => (
                <Row segment={s} key={s.id} setShowDelete={setShowDelete} />
            ))}
        </TableBody>
    </StyledTable>
))

SegmentTable.displayName = "SegmentTable"
SegmentTable.propTypes = {
    list: PropTypes.array,
    setShowDelete: PropTypes.func,
}

const segmentValidation = object().shape({
    name: string()
        .required("Please enter a group name.")
        .max(191, "The name must not exceed 191 characters."),
})

const CreateForm = ({
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

const CreateSegment = ({ callApi, hideModal }) => {
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

const DeleteForm = ({ id, callApi, hideModal }) => {
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

DeleteForm.propTypes = {
    id: PropTypes.number,
    callApi: PropTypes.func,
    hideModal: PropTypes.func,
}

const List = () => {
    const [showDelete, setShowDelete] = useState({ show: false, name: "" })
    const [showCreate, openCreateModal] = useState(false)
    const hideModal = () => setShowDelete({ show: false, name: "", id: "" })

    const [state, callApi] = useApi(
        {
            url: endpoints.getGroups,
        },
        {
            collection: [],
            init: true,
        }
    )

    let table = null
    if (state.isLoading) {
        table = (
            <PlaceholderTable
                width="100%"
                header={Header}
                numCols={4}
                numRows={8}
            />
        )
    } else if (state.data && state.data.collection.length > 0) {
        table = (
            <SegmentTable
                isLoading={state.isLoading}
                list={state.data.collection}
                setShowDelete={setShowDelete}
            />
        )
    }

    return (
        <>
            {showDelete.show && (
                <Modal
                    title={`Delete group ${showDelete.name} ?`}
                    hideModal={hideModal}
                    form={
                        <DeleteSegment
                            id={showDelete.id}
                            onSuccess={async () => {
                                await callApi({ url: endpoints.getGroups })
                                hideModal()
                            }}
                            onCancel={hideModal}
                        />
                    }
                />
            )}
            {showCreate && (
                <Modal
                    title={`Create group`}
                    hideModal={() => openCreateModal(false)}
                    form={
                        <CreateSegment
                            callApi={callApi}
                            hideModal={() => openCreateModal(false)}
                        />
                    }
                />
            )}
            <Box
                gridArea="nav"
                direction="row"
                border={{ side: "bottom", color: "light-4" }}
            >
                <Box margin={{ right: "small" }} alignSelf="center">
                    <Heading level="2">Groups</Heading>
                </Box>
                <Box alignSelf="center">
                    <Button
                        primary
                        color="status-ok"
                        label="Create new"
                        icon={<Add />}
                        reverse
                        onClick={() => openCreateModal(true)}
                    />
                </Box>
            </Box>
            <Box gridArea="main">
                <Box animation="fadeIn">
                    {table}

                    {!state.isLoading &&
                    !state.isError &&
                    state.data.collection.length === 0 ? (
                        <Box align="center" margin={{ top: "small" }}>
                            <Heading level="2">
                                Create your first group.
                            </Heading>
                        </Box>
                    ) : null}
                </Box>
                {!state.isLoading && state.data.collection.length > 0 ? (
                    <Box
                        direction="row"
                        alignSelf="end"
                        margin={{ top: "medium" }}
                    >
                        <Box margin={{ right: "small" }}>
                            <Button
                                icon={<FormPreviousLink />}
                                label="Previous"
                                disabled={state.data.links.previous === null}
                                onClick={() => {
                                    callApi({
                                        url: state.data.links.previous,
                                    })
                                }}
                            />
                        </Box>
                        <Box>
                            <Button
                                icon={<FormNextLink />}
                                reverse
                                label="Next"
                                disabled={state.data.links.next === null}
                                onClick={() => {
                                    callApi({
                                        url: state.data.links.next,
                                    })
                                }}
                            />
                        </Box>
                    </Box>
                ) : null}
            </Box>
        </>
    )
}

export default List
