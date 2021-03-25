package rbac.authz

# role-permissions assignments
role_permissions := {
  # test roles and permissions
  "engineering": [{"action": "read",  "object": "server123"}],
  "webdev":      [{"action": "read",  "object": "server123"},
                    {"action": "write", "object": "server123"}],
  "hr":          [{"action": "read",  "object": "database456"}]
}

default allow = false

# Allow admins to do anything.
allow {
	user_is_admin
}

user_is_admin {
	input.roles[_] == "admin"
}

allow {
    r := input.roles[_]
    permissions := role_permissions[r]
    p := permissions[_]
    p == {"action": input.action, "object": input.object}
}