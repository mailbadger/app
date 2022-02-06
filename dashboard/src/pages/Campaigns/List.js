import React, { useState, useContext, Fragment } from "react"
import PropTypes from "prop-types"
import { parseISO, formatRelative } from "date-fns"
import { More, Add } from "grommet-icons"
import { Box, Heading, ResponsiveContext, Select, Text } from "grommet"
// import history from "../../utils/history"
import { Modal, Badge, Notice, BarLoader } from "../../ui"
import CreateCampaign from "./Create"
import EditCampaign from "./Edit"
import DeleteCampaign from "./Delete"
import { SesKeysContext } from "../Settings/SesKeysContext"
import { endpoints } from "../../network/endpoints"
import DashboardPlaceholderTable from "../../ui/DashboardPlaceholderTable"
import { DashboardDataTable } from "../../ui/DashboardDataTable"
import { useCallApiDataTable } from "../../hooks/useCallApiDataTable"
import {
    StyledHeaderButtons,
    StyledHeaderTitle,
    StyledHeaderWrapper,
    StyledImportButton,
} from "../Subscribers/StyledSections"
import { faPlus } from "@fortawesome/free-solid-svg-icons"
import { LinkWrapper } from "../../ui/LinkWrapper"

const NameLink = ({ id, name, status, hasSesKeys }) => {
    let to = `/dashboard/campaigns/send/${id}`
    if (status === "sent") {
        to = `/dashboard/campaigns/${id}/report`
    }
    if (status === "draft" && !hasSesKeys) {
        return (
            <Text weight="bold" size="small">
                {name}
            </Text>
        )
    }
    return <LinkWrapper to={to}>{name}</LinkWrapper>
}

NameLink.displayName = "NameLink"
NameLink.propTypes = {
    name: PropTypes.string,
    id: PropTypes.number,
    status: PropTypes.string,
    hasSesKeys: PropTypes.bool,
}

const StatusBadge = ({ status }) => {
    const statusColors = {
        draft: "#CCCCCC",
        sending: "#00739D",
        sent: "#00C781",
        scheduled: "#FFCA58",
    }

    return <Badge color={statusColors[status]}>{status}</Badge>
}

StatusBadge.propTypes = {
    status: PropTypes.string,
}
StatusBadge.displayName = "StatusBadge"

const ActionDropdown = ({
    id,
    name,
    hasSesKeys,
    setShowEdit,
    setShowDelete,
}) => {
    let options = ["Delete"]
    if (hasSesKeys) {
        options.unshift("Edit")
    }
    return (
        <Select
            alignSelf="center"
            plain
            icon={<More />}
            options={options}
            onChange={({ option }) => {
                ;(function () {
                    switch (option) {
                        case "Edit":
                            setShowEdit({
                                show: true,
                                id: id,
                            })
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
    )
}

ActionDropdown.displayName = "ActionDropdown"
ActionDropdown.propTypes = {
    id: PropTypes.number,
    name: PropTypes.string,
    hasSesKeys: PropTypes.bool,
    setShowEdit: PropTypes.func,
    setShowDelete: PropTypes.func,
}

const List = () => {
    const [showDelete, setShowDelete] = useState({
        show: false,
        name: "",
        id: "",
    })
    const [showEdit, setShowEdit] = useState({ show: false, id: "" })
    const [showCreate, openCreateModal] = useState(false)
    const hideModal = () => setShowDelete({ show: false, name: "", id: "" })
    const hideEditModal = () => setShowEdit({ show: false, id: "" })
    const {
        keys,
        isLoading: keysLoading,
        error: keysError,
    } = useContext(SesKeysContext)
    const [searchInput, setSearchInput] = useState("")

    const hasSesKeys = !keysLoading && !keysError && keys !== null

    const columns = [
        {
            property: "name",
            render: NameLink,
            header: "Name",
        },
        {
            property: "status",
            header: "Status",
            render: StatusBadge,
        },
        {
            property: "template_name",
            header: "Template Name",
            render: (campaign) =>
                campaign.template ? campaign.template.name : "/",
        },
        {
            property: "subject",
            header: "Subject",
            render: (campaign) =>
                campaign.template ? campaign.template.subject_part : "/",
        },
        {
            property: "created_at",
            header: "Created At",
            render: (campaign) => {
                const d = parseISO(campaign.created_at)
                return formatRelative(d, new Date())
            },
        },
        {
            property: "action",
            header: "Action",
            render: ActionDropdown,
        },
    ]

    const contextSize = useContext(ResponsiveContext)

    const [callApi, state, onClickPrev, onClickNext] = useCallApiDataTable(
        endpoints.getCampaigns
    )

    const handleChange = (e) => {
        setSearchInput(e.target.value)
    }

    const hasCampaigns =
        !state.isLoading && !state.isError && state.data.collection.length > 0

    if (keysLoading) {
        return (
            <Box gridArea="nav" alignContent="center" margin="20%">
                <BarLoader size={15} />
            </Box>
        )
    }

    // if (!hasSesKeys && !hasCampaigns) {
    //     return (
    //         <Box>
    //             <Box align="center" margin={{ top: "large" }}>
    //                 <Heading level="2">
    //                     Please provide your AWS SES keys first.
    //                 </Heading>
    //             </Box>
    //             <Box align="center" margin={{ top: "medium" }}>
    //                 <Button
    //                     primary
    //                     color="status-ok"
    //                     label="Add SES Keys"
    //                     icon={<Add />}
    //                     reverse
    //                     onClick={() => history.push("/dashboard/settings")}
    //                 />
    //             </Box>
    //         </Box>
    //     )
    // }

    let table = null
    if (state.isLoading || keysLoading) {
        table = (
            <DashboardPlaceholderTable
                columns={columns}
                numCols={columns.length}
                numRows={10}
            />
        )
    } else if (hasCampaigns) {
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
                    title={`Delete campaign ${showDelete.name} ?`}
                    hideModal={hideModal}
                    form={
                        <DeleteCampaign
                            id={showDelete.id}
                            onSuccess={() =>
                                callApi({ url: endpoints.getCampaigns })
                            }
                            hideModal={hideModal}
                        />
                    }
                />
            )}
            {showCreate && (
                <Modal
                    title={`Create campaign`}
                    hideModal={() => openCreateModal(false)}
                    form={
                        <CreateCampaign
                            callApi={callApi}
                            hideModal={() => openCreateModal(false)}
                        />
                    }
                />
            )}
            {showEdit.show && (
                <Modal
                    title={`Edit campaign`}
                    hideModal={hideEditModal}
                    form={
                        <EditCampaign
                            id={showEdit.id}
                            onSuccess={() =>
                                callApi({ url: endpoints.getCampaigns })
                            }
                            hideModal={hideEditModal}
                        />
                    }
                />
            )}
            <>
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
                        Campaings
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
                                onClick={() => {
                                    if (hasSesKeys) {
                                        openCreateModal(true)
                                    }
                                }}
                                disabled={!hasSesKeys}
                            />
                        </Fragment>
                    </StyledHeaderButtons>
                </StyledHeaderWrapper>
                {!hasSesKeys && hasCampaigns && (
                    <Box alignSelf="center" pad="medium">
                        <Notice
                            message="Set your SES keys in order to send, create or edit campaigns."
                            status="status-warning"
                            color="white"
                            borderColor="status-warning"
                        />
                    </Box>
                )}
            </>
            <Box gridArea="main">
                <Box animation="fadeIn">
                    {table}

                    {!state.isLoading &&
                        !state.isError &&
                        state.data.collection.length === 0 && (
                            <Box align="center" margin={{ top: "large" }}>
                                <Heading level="2">
                                    Create your first campaign.
                                </Heading>
                            </Box>
                        )}
                </Box>
            </Box>
        </>
    )
}

export default List
