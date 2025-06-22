package types

type (
	QueryApplicationOptions struct {
		Pagination *Pagination
		WithTotal  bool
	}
	QueryGroupOptions struct {
		IsDefault  []bool
		Pagination *Pagination
		WithTotal  bool
	}
	QueryMenuOptions struct {
		Pagination *Pagination
		WithTotal  bool
	}
	QueryPermissionOptions struct {
		Pagination *Pagination
		WithTotal  bool
	}
	QueryRoleOptions struct {
		IsDefault  []bool
		Pagination *Pagination
		WithTotal  bool
	}
	QueryUserOptions struct {
		Username   string
		Email      string
		Mobile     string
		States     []int8
		RoleID     uint
		GroupID    uint
		Pagination *Pagination
		WithTotal  bool
	}
)
