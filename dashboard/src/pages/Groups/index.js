import React from "react"
import { Switch } from "react-router-dom"

import ProtectedRoute from "../../ProtectedRoute"
import List from "./List"
import Details from "./Details"

const Groups = () => (
    <Switch>
        <ProtectedRoute
            exact
            path="/dashboard/groups/:id"
            component={Details}
        />
        <ProtectedRoute
            exact
            path="/dashboard/groups"
            component={() => <List />}
        />
    </Switch>
)

export default Groups
