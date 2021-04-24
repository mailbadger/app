package rbac.authz

# role-permissions assignments
role_permissions := {
  # test roles and permissions
  "engineering": [{"method": "GET",  "path": "/api/campaigns"},
                  {"method": "GET",  "path": "/api/campaigns/:id"}],
  "webdev":      [{"method": "GET",  "path": "/api/campaigns"},
                  {"method": "PUT",  "path": "/api/campaigns/:id"}],
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
    p == {"method": input.method, "path": input.path}
}