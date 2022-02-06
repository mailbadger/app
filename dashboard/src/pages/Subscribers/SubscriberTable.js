import React from "react"
import PropTypes from "prop-types"
import { parseISO, formatRelative } from "date-fns"
import { TableBody, TableRow, TableCell } from "grommet"
import { StyledTable } from "../../ui"
import { getColumnSize } from "../../ui/DashboardDataTable"
import { StyledTableHeader } from "../../ui/DashboardStyledTable"

export const Row = ({ subscriber, actions }) => {
    const ca = parseISO(subscriber.created_at)
    const ua = parseISO(subscriber.updated_at)
    return (
        <TableRow>
            <TableCell scope="row" size="medium">
                <strong>{subscriber.email}</strong>
            </TableCell>
            <TableCell scope="row" size="medium">
                {formatRelative(ca, new Date())}
            </TableCell>
            <TableCell scope="row" size="medium">
                {formatRelative(ua, new Date())}
            </TableCell>
            <TableCell scope="row" size="xsmall" align="end">
                {actions}
            </TableCell>
        </TableRow>
    )
}
Row.propTypes = {
    subscriber: PropTypes.shape({
        email: PropTypes.string,
        id: PropTypes.number,
        created_at: PropTypes.string,
        updated_at: PropTypes.string,
    }),
    actions: PropTypes.element,
}

export const Header = ({ size }) => (
    <StyledTableHeader>
        <TableRow>
            <TableCell scope="col" border="bottom" size={getColumnSize(size)}>
                <strong>Email</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="small">
                <strong>Created At</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="small">
                <strong>Updated At</strong>
            </TableCell>
            <TableCell scope="col" border="bottom" size="128px">
                <strong> {""}</strong>
            </TableCell>
            <TableCell align="center" scope="col" border="bottom" size="small">
                <strong>Actions</strong>
            </TableCell>
        </TableRow>
    </StyledTableHeader>
)
Header.propTypes = {
    size: PropTypes.string,
}
export const SubscriberTable = React.memo(({ list, actions }) => (
    <StyledTable>
        <Header />
        <TableBody>
            {list.map((s) => (
                <Row subscriber={s} key={s.id} actions={actions(s)} />
            ))}
        </TableBody>
    </StyledTable>
))
SubscriberTable.displayName = "SubscriberTable"
SubscriberTable.propTypes = {
    list: PropTypes.array,
    actions: PropTypes.func,
}
