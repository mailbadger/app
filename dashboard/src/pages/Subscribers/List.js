import React, { useState, useContext, useEffect, Fragment } from "react"
import { parseISO, formatRelative } from "date-fns"
import { More } from "grommet-icons"
import { Box, Heading, ResponsiveContext, Select } from "grommet"

import history from "../../utils/history"
import { Modal } from "../../ui"
import CreateSubscriber from "./Create"
import DeleteSubscriber from "./Delete"
import EditSubscriber from "./Edit"
import { DashboardDataTable, getColumnSize } from "../../ui/DashboardDataTable"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import { faPlus } from "@fortawesome/free-solid-svg-icons"
import {
    StyledHeaderWrapper,
    StyledHeaderButtons,
    StyledHeaderTitle,
    StyledHeaderButton,
    StyledImportButton,
    StyledActions,
} from "./StyledSections"
import DashboardPlaceholderTable from "../../ui/DashboardPlaceholderTable"
import { endpoints } from "../../network/endpoints"
import { useCallApiDataTable } from "../../hooks/useCallApiDataTable"
import { ExportSubscribers } from "./ExportSubscribers"

const search = (setFilteredData, searchInput, data) => {
    let filteredData = []
    if (searchInput !== "") {
        filteredData = data.filter((entry) => {
            const foundMatch = Object.values(entry).some((entryValue) => {
                return (
                    entryValue &&
                    entryValue
                        .toString()
                        .toLowerCase()
                        .includes(searchInput.toLowerCase())
                )
            })
            if (foundMatch) {
                return entry
            }
        })
    }
    setFilteredData(filteredData)
}

const List = () => {
    const [showDelete, setShowDelete] = useState({
        show: false,
        email: "",
        id: "",
    })
    const [showEdit, setShowEdit] = useState({ show: false, id: "" })
    const [showCreate, openCreateModal] = useState(false)
    const [searchInput, setSearchInput] = useState("")
    const [filteredData, setFilteredData] = useState([])

    const hideDeleteModal = () =>
        setShowDelete({ show: false, email: "", id: "" })
    const hideEditModal = () => setShowEdit({ show: false, id: "" })

    const contextSize = useContext(ResponsiveContext)
    const columns = [
        {
            property: "email",
            header: "Email",
            size: getColumnSize(contextSize),
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
            render: function Cell({ updated_at }) {
                const dateUpdatedAt = parseISO(updated_at)
                return <span>{formatRelative(dateUpdatedAt, new Date())}</span>
            },
        },
        {
            property: "actions",
            header: "Actions",
            size: "small",
            align: "center",
            render: function Cell({ id, email }) {
                return (
                    <StyledActions>
                        <Select
                            alignSelf="center"
                            plain
                            defaultValue="View"
                            icon={<More />}
                            options={["View", "Edit", "Delete"]}
                            onChange={({ option }) => {
                                ;(() => {
                                    switch (option) {
                                        case "Edit":
                                            setShowEdit({
                                                show: true,
                                                id,
                                            })
                                            break
                                        case "View":
                                            history.push(
                                                `/dashboard/groups/${id}`
                                            )
                                            break
                                        case "Delete":
                                            setShowDelete({
                                                show: true,
                                                email,
                                                id,
                                            })
                                            break
                                        default:
                                            return "null"
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
        endpoints.getSubscribers
    )

    const handleChange = (e) => {
        setSearchInput(e.target.value)
    }

    useEffect(() => {
        search(setFilteredData, searchInput, state.data.collection)
    }, [searchInput])

    let table = null
    const data =
        filteredData && filteredData.length
            ? filteredData
            : state.data.collection

    if (state.isLoading) {
        table = (
            <DashboardPlaceholderTable
                columns={columns}
                numCols={columns.length}
                numRows={10}
            />
        )
    } else if (data.length > 0) {
        table = (
            <DashboardDataTable
                columns={columns}
                data={data}
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
                    title={`Delete subscriber ${showDelete.email} ?`}
                    hideModal={hideDeleteModal}
                    form={
                        <DeleteSubscriber
                            id={showDelete.id}
                            callApi={callApi}
                            hideModal={hideDeleteModal}
                        />
                    }
                />
            )}
            {showCreate && (
                <Modal
                    title={`Create subscriber`}
                    hideModal={() => openCreateModal(false)}
                    form={
                        <CreateSubscriber
                            callApi={callApi}
                            hideModal={() => openCreateModal(false)}
                        />
                    }
                />
            )}
            {showEdit.show && (
                <Modal
                    title={`Edit subscriber`}
                    hideModal={hideEditModal}
                    form={
                        <EditSubscriber
                            id={showEdit.id}
                            callApi={callApi}
                            hideModal={hideEditModal}
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
                <StyledHeaderTitle size={contextSize}>
                    Subscribers
                </StyledHeaderTitle>
                <StyledHeaderButtons
                    size={contextSize}
                    margin={{ left: "auto" }}
                >
                    <Fragment>
                        <StyledImportButton
                            width="256"
                            margin={{ right: "small" }}
                            icon={<FontAwesomeIcon icon={faPlus} />}
                            label="Import from file"
                            onClick={() =>
                                history.push("/dashboard/subscribers/import")
                            }
                        />
                        <StyledHeaderButton
                            width="154"
                            margin={{ right: "small" }}
                            label="Create New"
                            onClick={() => openCreateModal(true)}
                        />
                        <StyledHeaderButton
                            width="184"
                            margin={{ right: "small" }}
                            label="Delete from file"
                            onClick={() =>
                                history.push(
                                    "/dashboard/subscribers/bulk-delete"
                                )
                            }
                        />
                        <ExportSubscribers />
                    </Fragment>
                </StyledHeaderButtons>
            </StyledHeaderWrapper>
            <Box gridArea="main">
                <Box animation="fadeIn">
                    {table}
                    {!state.isLoading && data.length === 0 ? (
                        <Box align="center" margin={{ top: "large" }}>
                            <Heading level="2">
                                Create your first subscriber.
                            </Heading>
                        </Box>
                    ) : null}
                </Box>
            </Box>
        </Fragment>
    )
}

export default List
