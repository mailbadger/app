import React, { useState, useContext, useEffect, useRef, Fragment } from "react"
import { useApi, useInterval } from "../../hooks"
import { mainInstance as axios } from "../../network/axios"
import { NotificationsContext } from "../../Notifications/context"
import { StyledHeaderButton } from "./StyledSections"
import { endpoints } from "../../network/endpoints"

export const ExportSubscribers = () => {
    const linkEl = useRef(null)
    const { createNotification } = useContext(NotificationsContext)
    const [notification, setNotification] = useState()
    const [filename, setFilename] = useState("")
    const [retries, setRetries] = useState(-1)
    const [state, callApi] = useApi(
        {
            url: endpoints.getSubscribersExport,
        },
        null,
        true
    )

    useInterval(
        async () => {
            await callApi({
                url: `${endpoints.getSubscribersExportDownload}?filename=${filename}`,
            })
            setRetries(retries - 1)
        },
        retries > 0 ? 1000 : null
    )

    useEffect(() => {
        if (notification) {
            createNotification(notification.message, notification.status)
        }
    }, [notification])

    useEffect(() => {
        if (!state.isLoading && state.isError && state.data) {
            if (state.data.status === "failed" && retries > 0 && retries < 50) {
                setRetries(-1)
                setNotification({
                    message: state.data.message,
                    status: "status-error",
                })
            }
        }

        if (!state.isLoading && !state.isError && state.data) {
            if (retries > 0) {
                setRetries(-1)
                linkEl.current.click()
            }
        }
    }, [state])

    return (
        <Fragment>
            <StyledHeaderButton
                width="110"
                disabled={retries > 0 || state.isLoading}
                onClick={async () => {
                    try {
                        const res = await axios.post(
                            endpoints.postSubscribersExport
                        )
                        setFilename(res.data.file_name)
                        setRetries(50)
                    } catch (e) {
                        console.error("Unable to generate report", e)
                    }
                }}
                label="Export"
            />
            {!state.isLoading && !state.isError && state.data && (
                <a ref={linkEl} href={state.data.url} />
            )}
        </Fragment>
    )
}
