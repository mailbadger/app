import React, { Fragment, useState, useEffect } from "react"
import PropTypes from "prop-types"
import { NavLink } from "react-router-dom"
import { Paragraph } from "grommet"
import { mainInstance as axios } from "../../network/axios"
import { endpoints } from "../../network/endpoints"

const VerifyEmail = (props) => {
    const [data, setData] = useState({ message: "" })
    const {
        match: { params },
    } = props

    useEffect(() => {
        const callApi = async () => {
            try {
                const res = await axios.put(
                    `${endpoints.verifyEmail}/${params.token}`
                )
                setData(res.data)
            } catch (error) {
                setData(error.response.data)
            }
        }

        callApi()
    }, [params])

    if (data.message === "") {
        return <div>Loading...</div>
    }

    return (
        <Fragment>
            <Paragraph>{data.message}</Paragraph>
            <NavLink to="/dashboard">Go to app &gt;</NavLink>
        </Fragment>
    )
}

VerifyEmail.propTypes = {
    match: PropTypes.shape({
        params: PropTypes.shape({
            token: PropTypes.string,
        }),
    }),
}

export default VerifyEmail
