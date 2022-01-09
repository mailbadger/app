import React from "react"
import { ThemeContext, Button } from "grommet"

const SecondaryButton = (props) => (
    <ThemeContext.Extend
        value={{
            button: {
                border: {
                    radius: 0,
                    color: "dark-1",
                    width: "2px",
                },
                padding: {
                    vertical: "2px",
                    horizontal: "5px",
                },
            },
            text: {
                medium: {
                    size: "14px",
                },
            },
        }}
    >
        <Button hoverIndicator="light-3" {...props} />
    </ThemeContext.Extend>
)

export default SecondaryButton
