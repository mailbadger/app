import React, { useState, useContext, Fragment } from "react"
import PropTypes from "prop-types"
import { parseISO, formatRelative } from "date-fns"
import { More, Add } from "grommet-icons"
import { mainInstance as axios } from "../../network/axios"
import { Formik, ErrorMessage } from "formik"
import { string, object } from "yup"

import {
    Box,
    Button,
    Heading,
    Select,
    FormField,
    TextInput,
    ResponsiveContext,
} from "grommet"
import history from "../../utils/history"
import { ButtonWithLoader, Modal } from "../../ui"
import { NotificationsContext } from "../../Notifications/context"
import DeleteSegment from "./Delete"
import { endpoints } from "../../network/endpoints"
import { DashboardDataTable, getColumnSize } from "../../ui/DashboardDataTable"
import {
    StyledActions,
    StyledHeaderButtons,
    StyledHeaderTitle,
    StyledHeaderWrapper,
    StyledImportButton,
} from "../Subscribers/StyledSections"
import { useCallApiDataTable } from "../../hooks/useCallApiDataTable"
import DashboardPlaceholderTable from "../../ui/DashboardPlaceholderTable"
import { LinkWrapper } from "../../ui/LinkWrapper"
import { faPlus } from "@fortawesome/free-solid-svg-icons"

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
    const [searchInput, setSearchInput] = useState("")

    const contextSize = useContext(ResponsiveContext)

    const columns = [
        {
            property: "name",
            header: "Name",
            size: getColumnSize(contextSize),
            render: function Cell({ name, id }) {
                return (
                    <LinkWrapper to={`/dashboard/groups/${id}`}>
                        <span>{name}</span>
                    </LinkWrapper>
                )
            },
        },
        {
            property: "subscribers_in_segment",
            header: "Total subscribers",
            size: "small",
        },
        {
            property: "created",
            header: "Created At",
            size: "small",
            render: function Cell({ created_at }) {
                const dateCreatedAt = new Date(created_at)
                return <span>{dateCreatedAt.toLocaleDateString("en-US")}</span>
            },
        },
        {
            property: "updated",
            header: "Updated At",
            size: "small",
            render: function Cell(data) {
                const dateUpdatedAt = parseISO(data.updated_at)
                return <span>{formatRelative(dateUpdatedAt, new Date())}</span>
            },
        },
        {
            property: "actions",
            header: "Actions",
            size: "small",
            align: "center",
            render: function Cell({ id, name }) {
                return (
                    <StyledActions>
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
                                                `/dashboard/groups/${id}`
                                            )
                                            break
                                        case "Delete":
                                            setShowDelete({
                                                show: true,
                                                name: name,
                                                id: id,
                                            })
                                            break
                                        default:
                                            return null
                                    }
                                })()
                            }}
                        />
                    </StyledActions>
                )
            },
        },
    ]

    const [callApi, state, onClickPrev, onClickNext] = useCallApiDataTable(
        endpoints.getGroups
    )

    const handleChange = (e) => {
        setSearchInput(e.target.value)
    }
    let table = null
    if (state.isLoading) {
        table = (
            <DashboardPlaceholderTable
                columns={columns}
                numCols={columns.length}
                numRows={10}
            />
        )
    } else if (state.data && state.data.collection.length > 0) {
        table = (
            <DashboardDataTable
                columns={columns}
                data={state.data.collection}
                isLoading={state.isLoading}
                onClickNext={onClickNext}
                onClickPrev={onClickPrev}
                prevLinks={state.data.links.previous}
                nextLinks={state.data.links.next}
                searchInput={searchInput}
                handleChange={handleChange}
            />
        )
    }

    return (
        <Fragment>
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
            <StyledHeaderWrapper
                size={contextSize}
                gridArea="nav"
                margin={{
                    left: "40px",
                    right: "100px",
                    bottom: "22px",
                    top: "40px",
                }}
            >
                <StyledHeaderTitle size={contextSize}>Groups</StyledHeaderTitle>
                <StyledHeaderButtons
                    size={contextSize}
                    margin={{ left: "auto" }}
                >
                    <Fragment>
                        <StyledImportButton
                            width="256"
                            margin={{ right: "small" }}
                            icon={<Add icon={faPlus} />}
                            label="Create new"
                            onClick={() => openCreateModal(true)}
                        />
                    </Fragment>
                </StyledHeaderButtons>
            </StyledHeaderWrapper>
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
            </Box>
        </Fragment>
    )
}

export default List
