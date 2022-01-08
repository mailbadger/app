/* eslint-disable no-unused-vars */
import React, { useState } from "react"
import PropTypes from "prop-types"
import { Anchor, ThemeContext, Box } from "grommet"
import { StyledLink } from "./Sidebar/StyledSections"

const AnchorLink = (props) => {
    const [isHovered, setHoveredState] = useState(false)

    const { hover, color, active, fromSidebar, icon, label } = props
    let hfw = "bold"
    let hcolor = ""
    if (hover) {
        hfw = hover.fontWeight
        hcolor = hover.color
    }

    const { fontWeight } = props
    let fw = props.active ? "bold" : "500"
    if (fontWeight) {
        fw = fontWeight
    }

    let dark = active ? "white" : "light-1"
    let light = {}

    if (!fromSidebar) {
        light = active ? "white" : "dark-1"
    }

    const shouldDisplayTooltip = fromSidebar && isHovered

    let linkStyle = {}
    if (shouldDisplayTooltip) {
        linkStyle = {
            borderTopRightRadius: "0",
            borderBottomRightRadius: "0",
        }
    }

    const tooltipStyle = {
        display: "inline-flex",
        fontSize: "16px",
        color: "#000",
        fontWeight: "500",
        textAlign: "left",
        position: "fixed",
        height: "48px",
        backgroundColor: "#fadcff",
        padding: "13px 8.6px 13px 0",
        marginLeft: "143px",
        borderTopRightRadius: "4px",
        borderBottomRightRadius: "4px",
        zIndex: "1",
        minWidth: "100px",
    }

    return (
        <ThemeContext.Extend
            value={{
                anchor: {
                    textDecoration: "none",
                    fontWeight: fw,
                    color: {
                        dark,
                        light,
                    },
                    extend: {
                        boxShadow: "none",
                    },
                    hover: {
                        textDecoration: "none",
                        fontWeight: hfw,
                        extend: {
                            color: "accent-1",
                        },
                    },
                },
            }}
        >
            <Anchor
                onMouseEnter={() => setHoveredState(true)}
                onMouseLeave={() => setHoveredState(false)}
                as={({
                    active,
                    colorProp,
                    hasIcon,
                    hasLabel,
                    focus,
                    fromSidebar,
                    ...rest
                }) => (
                    <StyledLink
                        hovered={shouldDisplayTooltip ? 1 : 0}
                        {...rest}
                    >
                        {fromSidebar && <Box>{icon}</Box>}
                        {shouldDisplayTooltip && (
                            <Box style={tooltipStyle}>{label}</Box>
                        )}
                    </StyledLink>
                )}
                {...props}
            />
        </ThemeContext.Extend>
    )
}

AnchorLink.propTypes = {
    active: PropTypes.bool,
    fontWeight: PropTypes.string,
    hover: PropTypes.shape({
        fontWeight: PropTypes.string,
        color: PropTypes.string,
    }),
    color: PropTypes.shape({
        active: PropTypes.string,
        idle: PropTypes.string,
    }),
    fromSidebar: PropTypes.bool,
    icon: PropTypes.element,
    label: PropTypes.string,
}

export default AnchorLink
