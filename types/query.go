package types

type (
	QueryApplicationOptions struct {
		Pagination *Pagination
	}
	QueryGroupOptions struct {
		Pagination *Pagination
	}
	QueryMenuOptions struct {
		Pagination *Pagination
	}
	QueryPermissionOptions struct {
		Pagination *Pagination
	}
	QueryRoleOptions struct {
		Pagination *Pagination
	}
	QueryUserOptions struct {
		Username   string
		Email      string
		Mobile     string
		States     []int8
		Pagination *Pagination
	}
)
