import React, { useState, useContext } from "react"
import PropTypes from "prop-types"
import { parseISO, formatRelative } from "date-fns"
import { More, Add, FormPreviousLink, FormNextLink } from "grommet-icons"
import { mainInstance as axios } from "../../network/axios"
import { useApi } from "../../hooks"
import {
    TableHeader,
    TableBody,
    TableRow,
    TableCell,
    Box,
    Button,
    Heading,
    Select,
} from "grommet"
import history from "../../utils/history"
import { StyledTable, ButtonWithLoader, PlaceholderTable, Modal } from "../../ui"
import { NotificationsContext } from "../../Notifications/context"
import { endpoints } from "../../network/endpoints"

const Row = ({ template, setShowDelete }) => {
    const createdAt = parseISO(template.created_at)
    const updatedAt = parseISO(template.updated_at)
    return (
        <TableRow>
            <TableCell scope="row" size="xlarge">
                <strong>{template.name}</strong>
            </TableCell>
            <TableCell scope="row" size="medium">
                <strong>{template.subject_part}</strong>
            </TableCell>
            <TableCell scope="row" size="medium">
                {formatRelative(createdAt, new Date())}
            </TableCell>
            <TableCell scope="row" size="medium">
                {formatRelative(updatedAt, new Date())}
            </TableCell>
            <TableCell scope="row" size="xsmall">
                <Select
                    alignSelf="center"
                    plain
                    icon={<More />}
                    options={["Edit", "Delete"]}
                    onChange={({ option }) => {
                        ;(function () {
                            switch (option) {
                                case "Edit":
                                    history.push(
                                        `/dashboard/templates/${template.id}/edit`
                                    )
                                    break
                                case "Delete":
                                    setShowDelete({
                                        show: true,
                                        name: template.name,
                                        id: template.id,
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
    setShowDelete: PropTypes.func,
    template: PropTypes.shape({
        id: PropTypes.number,
        name: PropTypes.string,
        subject_part: PropTypes.string,
        created_at: PropTypes.string,
        updated_at: PropTypes.string,
    }),
}

const Header = () => (
    <TableHeader>
        <TableRow>
            <TableCell scope="col" border="bottom" size="medium">
                <strong>Name</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="medium">
                <strong>Subject</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="medium">
                <strong>Created At</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="medium">
                <strong>Updated At</strong>
            </TableCell>
            <TableCell align="end" scope="col" border="bottom" size="small">
                <strong>Action</strong>
            </TableCell>
        </TableRow>
    </TableHeader>
)

const TemplateTable = React.memo(({ list, setShowDelete }) => (
    <StyledTable>
        <Header />
        <TableBody>
            {list.map((t) => (
                <Row template={t} key={t.name} setShowDelete={setShowDelete} />
            ))}
        </TableBody>
    </StyledTable>
))

TemplateTable.displayName = "TemplateTable"
TemplateTable.propTypes = {
    list: PropTypes.array,
    setShowDelete: PropTypes.func,
}

const DeleteForm = ({ id, callApi, hideModal }) => {
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

const List = () => {
    const [showDelete, setShowDelete] = useState({
        show: false,
        name: "",
        id: null,
    })
    const hideModal = () => setShowDelete({ show: false, name: "", id: null })

    const [state, callApi] = useApi(
        {
            url: endpoints.getTemplates,
        },
        {
            next_token: "",
            collection: [],
        }
    )

    let table = null
    if (state.isLoading) {
        table = <PlaceholderTable header={Header} numCols={4} numRows={5} />
    } else if (!state.isError && state.data.collection.length > 0) {
        table = (
            <TemplateTable
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
                    title={`Delete template ${showDelete.name} ?`}
                    hideModal={hideModal}
                    form={
                        <DeleteForm
                            id={showDelete.id}
                            callApi={callApi}
                            hideModal={hideModal}
                        />
                    }
                />
            )}
            <Box
                gridArea="nav"
                direction="row"
                border={{ side: "bottom", color: "light-4" }}
            >
                <Box alignSelf="center" margin={{ right: "small" }}>
                    <Heading level="2">Templates</Heading>
                </Box>
                <Box alignSelf="center">
                    <Button
                        primary
                        color="status-ok"
                        label="Create new"
                        icon={<Add />}
                        reverse
                        onClick={() => history.push("/dashboard/templates/new")}
                    />
                </Box>
            </Box>
            <Box gridArea="main">
                <Box animation="fadeIn">
                    {table}

                    {!state.isLoading &&
                    !state.isError &&
                    state.data.collection.length === 0 ? (
                        <Box align="center" margin={{ top: "large" }}>
                            <Heading level="2">
                                Create your first template.
                            </Heading>
                        </Box>
                    ) : null}
                </Box>
                {!state.isLoading &&
                !state.isError &&
                state.data.collection.length > 0 ? (
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
