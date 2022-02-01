import React, { useState, useContext, Fragment } from "react"
import { parseISO, formatRelative } from "date-fns"
import { More, Add } from "grommet-icons"

import { Box, Heading, Select, ResponsiveContext } from "grommet"
import history from "../../utils/history"
import { Modal } from "../../ui"
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
import { CreateSegment } from "./CreateSegmet/CreateSegment"

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
