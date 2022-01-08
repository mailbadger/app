import React, { useContext } from "react"
import { Box, DataTable, Button, ResponsiveContext } from "grommet"
import PropTypes from "prop-types"
import styled, { css } from "styled-components"
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome"
import {
    faChevronRight,
    faChevronLeft,
    faSearch,
} from "@fortawesome/free-solid-svg-icons"

const StyledWrapper = styled(Box)`
    min-height: auto;
    position: relative;
    border-radius: 12px 12px 0px 0px;
`

export const tableHeading = css`
    height: 50px;
    margin: 10px 0 0;
    background-color: rgba(245, 245, 250, 0.4);
    font-weight: bold;
`

const StyledSearchIcon = styled(FontAwesomeIcon)`
	position: absolute;
	margin-top: 26px;
	margin-left: 40px;
	width: 20px;
	height: 20px;
}
`
const StyledInput = styled.input`
    ${tableHeading};
    border: none;
    padding-left: 65px;
    outline: none;
    font-size: 17px;
    font-family: "Poppins Bold";

    &::-webkit-input-placeholder,
    &::-moz-placeholder,
    &:-ms-input-placeholder,
    &:-moz-placeholder {
        font-size: 13px;
        line-height: 1.38;
        color: #737d8b;
    }
`

const StyledDataTable = styled(DataTable)`
	background:white;
	padding: 10px 0 0;
	display: table;
  overflow: auto;
  
  @media (width: 1440px) {
    width: auto;
   }

	${this} thead {
		tr {
		width: 100%;
		font-size: 13px;
		${tableHeading};
		}
		
		th {
			border:none;

			&:first-of-type{
				padding-left: 30px;
			}
		}
	}

	${this} tbody {
		th, td {
			padding-top: 2px;
			padding-bottom: 1px;
		}

		tr {
			input {
				width: 85px;
				text-align:center;

				&:hover {
					color:#541388;
				}
			}
			
			&:hover {
				background-color: rgba(84, 19, 136, 0.1);
			}

			th {
				font-size: 16px;
				font-weight: bold;
				color: #541388;

				&:first-of-type{
					padding-left: 30px;
				}
			}

			span {
				white-space: nowrap;
				text-overflow: ellipsis;
				overflow: hidden;
			}

			td {
				button{
					&:focus {
						box-shadow: none;
					}
				}

				&:nth-of-type(3) {	
					div {
						background-color: #fadcff;
						border-radius: 8px;
						padding: 8.5px 24px 8.5px;
						font-size: 14px;
						font-weight: bold;
						color: #541388;
				}
			}
		}	
	}
`

const StyledLinkIcon = styled(Box)`
    width: 40px;
    height: 40px;
    margin: 0;
    border-radius: 6px;
    background-color: rgba(84, 19, 136, 0.1);
    justify-content: center;
    align-items: center;

    ${this} svg {
        width: 12.1px !important;
        height: 20px;
    }
`

const StyledButton = styled(Button)`
    border: none;
    font-size: 14px;
    text-align: left;
    color: #000000;
    font-weight: normal;

    &:focus {
        box-shadow: 0 0 6px 2px #390099;
    }

    &:first-of-type {
        padding-right: 5px;
    }
`

export const DashboardWrapper = styled(Box)`
    padding: ${(props) =>
        props.contextSize === "large" ? "0 100px 15px" : "20px"};
    display: table;
`

export const getColumnSize = (size) => (size === "large" ? "medium" : "small")

export const DashboardSearchPlaceholder = ({
    searchInput,
    handleChange,
    pad,
}) => {
    const onChange = handleChange ? handleChange : () => {}
    return (
        <StyledWrapper background=" white" fill pad={pad ? pad : {}}>
            <Box
                background="white"
                style={{
                    borderRadius: "12px 12px 0px 0px",
                    position: "relative",
                }}
                pad={{ top: "30px" }}
            >
                <StyledInput
                    placeholder="Search in this table ..."
                    name="searchInput"
                    value={searchInput || ""}
                    onChange={onChange}
                    label="Search"
                />
                <StyledSearchIcon icon={faSearch} />
            </Box>
        </StyledWrapper>
    )
}

DashboardSearchPlaceholder.propTypes = {
    searchInput: PropTypes.string,
    handleChange: PropTypes.func,
    pad: PropTypes.object,
}

export const DashboardDataTable = ({
    columns,
    data,
    isLoading,
    onClickPrev,
    onClickNext,
    prevLinks,
    nextLinks,
    searchInput,
    handleChange,
}) => {
    const size = useContext(ResponsiveContext)

    return (
        <DashboardWrapper fill="horizontal" contextSize={size} overflow="auto">
            <DashboardSearchPlaceholder
                searchInput={searchInput}
                handleChange={handleChange}
            />
            <StyledDataTable
                contextSize={size}
                rowProps={{ email: { pad: "large" } }}
                columns={columns}
                data={data}
            />
            {!isLoading && data.length > 0 ? (
                <Box
                    width="100%"
                    pad={{
                        top: "10px",
                        right: "30px",
                        bottom: "15px",
                        left: "0",
                    }}
                    background="white"
                    style={{
                        minHeight: "auto",
                        borderRadius: "0px 0px 12px 12px",
                    }}
                >
                    <Box direction="row" alignSelf="end">
                        <Box margin={{ right: "small" }} justify="center">
                            <StyledButton
                                icon={
                                    <StyledLinkIcon>
                                        <FontAwesomeIcon icon={faChevronLeft} />
                                    </StyledLinkIcon>
                                }
                                label="Prev"
                                disabled={prevLinks === null}
                                onClick={onClickPrev}
                            />
                        </Box>
                        <Box justify="center" margin={{ right: "0" }}>
                            <StyledButton
                                icon={
                                    <StyledLinkIcon>
                                        <FontAwesomeIcon
                                            icon={faChevronRight}
                                        />
                                    </StyledLinkIcon>
                                }
                                reverse
                                label="Next"
                                disabled={nextLinks === null}
                                onClick={onClickNext}
                            />
                        </Box>
                    </Box>
                </Box>
            ) : null}
        </DashboardWrapper>
    )
}

DashboardDataTable.propTypes = {
    columns: PropTypes.array,
    data: PropTypes.array,
    isLoading: PropTypes.bool,
    onClickPrev: PropTypes.func,
    onClickNext: PropTypes.func,
    prevLinks: PropTypes.string || undefined,
    nextLinks: PropTypes.string || undefined,
    searchInput: PropTypes.string,
    handleChange: PropTypes.func,
}
