import React, { useState, useContext, memo } from "react"
import PropTypes from "prop-types"
import { parseISO, formatRelative } from "date-fns"
import { More, Add, FormPreviousLink, FormNextLink } from "grommet-icons"
import { useApi } from "../../hooks"
import {
    TableHeader,
    TableRow,
    TableCell,
    Box,
    Button,
    Heading,
    Select,
    Text,
} from "grommet"
import history from "../../utils/history"
import {
    PlaceholderTable,
    Modal,
    Badge,
    Notice,
    BarLoader,
    ListGrid,
    AnchorLink,
    StyledDataTable,
} from "../../ui"
import CreateCampaign from "./Create"
import EditCampaign from "./Edit"
import DeleteCampaign from "./Delete"
import { SesKeysContext } from "../Settings/SesKeysContext"
import { endpoints } from "../../network/endpoints"

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
    return <AnchorLink size="small" fontWeight="bold" to={to} label={name} />
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
    let opts = ["Delete"]
    if (hasSesKeys) {
        opts.unshift("Edit")
    }
    return (
        <Select
            alignSelf="center"
            plain
            icon={<More />}
            options={opts}
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

const columns = [
    {
        property: "name",
        primary: true,
        render: NameLink,
        header: "Name",
        search: true,
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
        align: "end",
        render: ActionDropdown,
    },
]

const PlaceholderHeader = () => (
    <TableHeader>
        <TableRow>
            <TableCell scope="col" border="bottom" size="xxsmall">
                <strong>Name</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="xxsmall">
                <strong>Status</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="xxsmall">
                <strong>Template</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="xxsmall">
                <strong>Subject</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="xxsmall">
                <strong>Created At</strong>
            </TableCell>
            <TableCell
                style={{ textAlign: "right" }}
                align="end"
                scope="col"
                border="bottom"
                size="xxsmall"
            >
                <strong>Action</strong>
            </TableCell>
        </TableRow>
    </TableHeader>
)

PlaceholderHeader.displayName = "PlaceholderHeader"

const CampaignsTable = memo(
    ({ list, setShowDelete, hasSesKeys, setShowEdit }) => {
        return (
            <StyledDataTable
                columns={columns}
                data={list.map((c) => ({
                    ...c,
                    setShowDelete,
                    hasSesKeys,
                    setShowEdit,
                }))}
                background={{
                    header: "white",
                    body: ["light-1", "white"],
                }}
                size="medium"
            />
        )
    }
)

CampaignsTable.displayName = "CampaignsTable"
CampaignsTable.propTypes = {
    list: PropTypes.array,
    setShowDelete: PropTypes.func,
    setShowEdit: PropTypes.func,
    hasSesKeys: PropTypes.bool,
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

    const hasSesKeys = !keysLoading && !keysError && keys !== null

    const [state, callApi] = useApi(
        {
            url: endpoints.getCampaigns,
        },
        {
            collection: [],
            init: true,
        }
    )

    const hasCampaigns =
        !state.isLoading && !state.isError && state.data.collection.length > 0

    if (keysLoading) {
        return (
            <Box gridArea="nav" alignContent="center" margin="20%">
                <BarLoader size={15} />
            </Box>
        )
    }

    if (!hasSesKeys && !hasCampaigns) {
        return (
            <Box>
                <Box align="center" margin={{ top: "large" }}>
                    <Heading level="2">
                        Please provide your AWS SES keys first.
                    </Heading>
                </Box>
                <Box align="center" margin={{ top: "medium" }}>
                    <Button
                        primary
                        color="status-ok"
                        label="Add SES Keys"
                        icon={<Add />}
                        reverse
                        onClick={() => history.push("/dashboard/settings")}
                    />
                </Box>
            </Box>
        )
    }

    let table = null
    if (state.isLoading || keysLoading) {
        table = (
            <PlaceholderTable
                header={PlaceholderHeader}
                numCols={5}
                numRows={3}
            />
        )
    } else if (hasCampaigns) {
        table = (
            <CampaignsTable
                isLoading={state.isLoading}
                list={state.data.collection}
                setShowDelete={setShowDelete}
                setShowEdit={setShowEdit}
                hasSesKeys={hasSesKeys}
            />
        )
    }

    return (
        <ListGrid>
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
                <Box gridArea="nav" direction="row">
                    <Box alignSelf="center" margin={{ right: "small" }}>
                        <Heading level="2">Campaigns</Heading>
                    </Box>
                    <Box alignSelf="center">
                        <Button
                            primary
                            color="status-ok"
                            label="Create new"
                            icon={<Add />}
                            reverse
                            onClick={() => {
                                if (hasSesKeys) {
                                    openCreateModal(true)
                                }
                            }}
                            disabled={!hasSesKeys}
                        />
                    </Box>
                    {!hasSesKeys && hasCampaigns && (
                        <Box margin={{ left: "auto" }} alignSelf="center">
                            <Notice
                                message="Set your SES keys in order to send, create or edit campaigns."
                                status="status-warning"
                                color="white"
                                borderColor="status-warning"
                            />
                        </Box>
                    )}
                </Box>
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
                    {hasCampaigns ? (
                        <Box
                            direction="row"
                            alignSelf="end"
                            margin={{ top: "medium" }}
                        >
                            <Box margin={{ right: "small" }}>
                                <Button
                                    icon={<FormPreviousLink />}
                                    label="Previous"
                                    disabled={
                                        state.data.links.previous === null
                                    }
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
        </ListGrid>
    )
}

export default List
