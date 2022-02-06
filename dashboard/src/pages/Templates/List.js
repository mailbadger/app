import React, { useState, useContext, Fragment } from "react"
import { parseISO, formatRelative } from "date-fns"
import { More, Add } from "grommet-icons"
import { Box, Heading, Select, ResponsiveContext } from "grommet"
import history from "../../utils/history"
import { Modal } from "../../ui"
import { endpoints } from "../../network/endpoints"
import { DashboardDataTable, getColumnSize } from "../../ui/DashboardDataTable"
import DashboardPlaceholderTable from "../../ui/DashboardPlaceholderTable"
import { useCallApiDataTable } from "../../hooks/useCallApiDataTable"
import {
    StyledActions,
    StyledHeaderButtons,
    StyledHeaderTitle,
    StyledHeaderWrapper,
    StyledImportButton,
} from "../Subscribers/StyledSections"
import { faPlus } from "@fortawesome/free-solid-svg-icons"
import { DeleteForm } from "./DeleteTemplate/DeleteForm"

const List = () => {
    const [showDelete, setShowDelete] = useState({
        show: false,
        name: "",
        id: null,
    })
    const hideModal = () => setShowDelete({ show: false, name: "", id: null })
    const [searchInput, setSearchInput] = useState("")

    const contextSize = useContext(ResponsiveContext)

    const columns = [
        {
            property: "name",
            header: "Name",
            size: getColumnSize(contextSize),
        },
        {
            property: "subject_part",
            header: "Subject",
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
                            options={["Edit", "Delete"]}
                            onChange={({ option }) => {
                                ;(function () {
                                    switch (option) {
                                        case "Edit":
                                            history.push(
                                                `/dashboard/templates/${id}/edit`
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
        endpoints.getTemplates
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
    } else if (!state.isError && state.data.collection.length > 0) {
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
                    Templates
                </StyledHeaderTitle>
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
                            onClick={() =>
                                history.push("/dashboard/templates/new")
                            }
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
                        <Box align="center" margin={{ top: "large" }}>
                            <Heading level="2">
                                Create your first template.
                            </Heading>
                        </Box>
                    ) : null}
                </Box>
            </Box>
        </>
    )
}

export default List
