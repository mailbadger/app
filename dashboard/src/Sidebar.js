import React, { Fragment, useState } from "react"
import PropTypes from "prop-types"
import { FormClose } from "grommet-icons"
import { useLocation } from "react-router-dom"
import { Box, Button, Layer } from "grommet"
import { UserMenu } from "./ui"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import {
    faUsers,
    faSitemap,
    faBullhorn,
    faCopy,
    faTools,
} from "@fortawesome/free-solid-svg-icons"
import logo from "./images/logo-small.png"
import {
    StyledMenuItemContainer,
    StyledAnchorLink,
} from "./ui/Sidebar/StyledSections"

const links = [
    {
        to: "/dashboard/subscribers",
        label: "Subscribers",
        icon: <FontAwesomeIcon icon={faUsers} />,
    },
    {
        to: "/dashboard/groups",
        label: "Groups",
        icon: <FontAwesomeIcon icon={faSitemap} />,
    },
    {
        to: "/dashboard/campaigns",
        label: "Campaigns",
        icon: <FontAwesomeIcon icon={faBullhorn} />,
    },
    {
        to: "/dashboard/templates",
        label: "Templates",
        icon: <FontAwesomeIcon icon={faCopy} />,
    },
    {
        to: "/dashboard/settings",
        label: "Settings",
        icon: <FontAwesomeIcon icon={faTools} />,
    },
]

const NavLinks = () => {
    let location = useLocation()
    const [active, setActive] = useState()

    return (
        <Fragment>
            {links.map((link) => {
                const isActive =
                    active === link.label ||
                    location.pathname.startsWith(link.to)

                return (
                    <StyledMenuItemContainer active={isActive} key={link.label}>
                        <Box direction="row">
                            <StyledAnchorLink
                                to={link.to}
                                size="medium"
                                icon={link.icon}
                                active={isActive}
                                onClick={() => setActive(link.label)}
                                fromSidebar={true}
                                label={link.label}
                            />
                        </Box>
                    </StyledMenuItemContainer>
                )
            })}
        </Fragment>
    )
}

const Sidebar = (props) => {
    const { showSidebar, size, closeSidebar } = props

    return (
        <Fragment>
            {!showSidebar || size !== "small" ? (
                <Box
                    overflow="auto"
                    background="white"
                    margin={{ top: "0", bottom: "1px", horizontal: "0" }}
                    pad={{ top: "20px", bottom: "24px", horizontal: "0" }}
                    style={{ boxShadow: "0 3px 6px 0 #fadcff" }}
                >
                    <Box align="center" margin={{ bottom: "15px" }}>
                        <Box height="80px">
                            <img style={{ height: "100%" }} src={logo} />
                        </Box>
                    </Box>
                    <Box
                        align="center"
                        gap={size === "small" ? "medium" : "small"}
                    >
                        <NavLinks />
                    </Box>
                    <Box flex />
                    <Box>
                        <UserMenu alignSelf="center" />
                    </Box>
                </Box>
            ) : (
                <Layer>
                    <Box
                        background="light-2"
                        tag="header"
                        justify="end"
                        align="center"
                        direction="row"
                    >
                        <Button icon={<FormClose />} onClick={closeSidebar} />
                    </Box>
                    <Box
                        fill
                        background="light-2"
                        direction="column"
                        justify="between"
                    >
                        <NavLinks />
                    </Box>
                </Layer>
            )}
        </Fragment>
    )
}

Sidebar.propTypes = {
    showSidebar: PropTypes.bool,
    size: PropTypes.string,
    closeSidebar: PropTypes.func,
}

export default Sidebar
