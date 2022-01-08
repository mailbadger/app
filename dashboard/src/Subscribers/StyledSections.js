import { SecondaryButton } from "../ui"

const { default: styled } = require("styled-components")
const { Box } = require("grommet")

export const StyledHeaderButtons = styled(Box)`
    ${(props) =>
        props.size === "large"
            ? `
    align-self: center;
    align-items: center;
    flex-direction: row;
    margin-left: auto;
    `
            : `
    align-self: start;
    align-items: start;
    flex-direction: column;
    margin-left: 0;

    button {
        margin-top: 10px;
    }
`}
`

export const StyledHeaderWrapper = styled(Box)`
    min-height: auto;
    flex-direction: ${(props) => (props.size === "large" ? "row" : "column")};
`

export const StyledHeaderTitle = styled(Box)`
    font-size: 50px;
    font-weight: bold;
    ${(props) =>
        props.size === "large"
            ? `
	align-self: center;
	 `
            : `align-self: start; 
	 margin-bottom: 30px;`}
`

export const StyledHeaderButton = styled(SecondaryButton)`
	width: ${(props) => (props.width ? `${props.width}px` : "auto")}
	height: 39px;
	border-radius: 20px;
	border: solid 1px #1c1c24;
	font-weight: normal;
	color: #000000;
	font-size: 16px;

	:hover {
	background-color: #000000;
	color: white;
	}

	&:focus {
		box-shadow:none;
	}
`

export const StyledImportButton = styled(StyledHeaderButton)`
    background-color: #fadcff;
    border: none;
    font-size: 20px;
    height: 50px;

    &:hover {
        background-color: #541388;
        color: white;
    }
`

export const StyledActions = styled(Box)`
    flex-direction: row;
    justify-content: center;

    button {
        div {
            div:nth-child(1) {
                display: none;
            }
        }

        svg {
            &:hover,
            &:focus {
                stroke: #541388;
            }
        }
    }
`
