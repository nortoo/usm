package usm

import "errors"

// addPolicy adds a new policy for a given role, tenant, resource, and action.
func (c *Client) addPolicy(role, tenant, resource, action string) error {
	if c.casbinCli == nil {
		return errors.New("casbin client is not initialized")
	}
	if _, err := c.casbinCli.AddPolicy(role, tenant, resource, action); err != nil {
		return err
	}
	if err := c.casbinCli.SavePolicy(); err != nil {
		return err
	}

	return c.casbinCli.LoadPolicy()
}

// clearPolicy removes all policies for a given role.
func (c *Client) clearPolicy(role string) error {
	if c.casbinCli == nil {
		return errors.New("casbin client is not initialized")
	}
	_, err := c.casbinCli.RemoveFilteredNamedPolicy("p", 0, role)
	if err != nil {
		return err
	}
	return c.casbinCli.LoadPolicy()
}

// Authorize checks if a role has permission to perform an action on a resource within a tenant.
func (c *Client) Authorize(role, tenant, resource, action string) (bool, error) {
	return c.casbinCli.Enforce(role, tenant, resource, action)
}
