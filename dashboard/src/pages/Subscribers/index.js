import React from "react"
import { Switch } from "react-router-dom"

import ProtectedRoute from "../../ProtectedRoute"
import List from "./List"
import { Row, Header, SubscriberTable } from "./SubscriberTable"
import Import from "./Import"
import BulkDelete from "./BulkDelete"

const Subscribers = () => {
    return (
        <Switch>
            <ProtectedRoute
                exact
                path="/dashboard/subscribers"
                component={() => <List />}
            />
            <ProtectedRoute
                exact
                path="/dashboard/subscribers/import"
                component={Import}
            />
            <ProtectedRoute
                exact
                path="/dashboard/subscribers/bulk-delete"
                component={BulkDelete}
            />
        </Switch>
    )
}

export { Row, Header, SubscriberTable as Table }

export default Subscribers
