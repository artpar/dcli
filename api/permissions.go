// api/permission.go

package api

import (
	"fmt"
	"strconv"
	"strings"
)

const (
	None AuthPermission = 0

	GuestPeek AuthPermission = 1 << iota
	GuestRead
	GuestCreate
	GuestUpdate
	GuestDelete
	GuestExecute
	GuestRefer
	UserPeek
	UserRead
	UserCreate
	UserUpdate
	UserDelete
	UserExecute
	UserRefer
	GroupPeek
	GroupRead
	GroupCreate
	GroupUpdate
	GroupDelete
	GroupExecute
	GroupRefer
)

var permissionNameMap = map[string]AuthPermission{
	"GuestPeek":    GuestPeek,
	"GuestRead":    GuestRead,
	"GuestCreate":  GuestCreate,
	"GuestUpdate":  GuestUpdate,
	"GuestDelete":  GuestDelete,
	"GuestExecute": GuestExecute,
	"GuestRefer":   GuestRefer,
	"UserPeek":     UserPeek,
	"UserRead":     UserRead,
	"UserCreate":   UserCreate,
	"UserUpdate":   UserUpdate,
	"UserDelete":   UserDelete,
	"UserExecute":  UserExecute,
	"UserRefer":    UserRefer,
	"GroupPeek":    GroupPeek,
	"GroupRead":    GroupRead,
	"GroupCreate":  GroupCreate,
	"GroupUpdate":  GroupUpdate,
	"GroupDelete":  GroupDelete,
	"GroupExecute": GroupExecute,
	"GroupRefer":   GroupRefer,
}

func StringsToAuthPermission(names []string) (AuthPermission, error) {
	var perm AuthPermission
	for _, name := range names {
		name = strings.TrimSpace(name)
		if p, ok := permissionNameMap[name]; ok {
			perm |= p
		} else {
			return 0, fmt.Errorf("unknown permission name: %s", name)
		}
	}
	return perm, nil
}

func AuthPermissionToStrings(perm AuthPermission) []string {
	var names []string
	for name, p := range permissionNameMap {
		if perm&p != 0 {
			names = append(names, name)
		}
	}
	return names
}

func (c *Client) GetPermissions(entityType, objectID string) (AuthPermission, error) {
	path := fmt.Sprintf("%s/%s", entityType, objectID)
	respData, err := c.GetResource(path)
	if err != nil {
		return 0, err
	}

	// Assuming the permission field is named "permission" in the JSON response
	data, ok := respData["data"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response format")
	}

	attributes, ok := data["attributes"].(map[string]interface{})
	if !ok {
		return 0, fmt.Errorf("invalid response format")
	}

	permissionValue, ok := attributes["permission"]
	if !ok {
		return 0, fmt.Errorf("permission field not found")
	}

	var permValue int64
	switch v := permissionValue.(type) {
	case string:
		permValue, err = strconv.ParseInt(v, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid permission value: %v", err)
		}
	case float64:
		permValue = int64(v)
	default:
		return 0, fmt.Errorf("invalid permission value type")
	}

	return AuthPermission(permValue), nil
}

func (c *Client) SetPermissions(entityType, objectID string, perm AuthPermission) error {
	path := fmt.Sprintf("%s/%s", entityType, objectID)
	data := map[string]interface{}{
		"data": map[string]interface{}{
			"type": entityType,
			"id":   objectID,
			"attributes": map[string]interface{}{
				"permission": fmt.Sprintf("%d", perm),
			},
		},
	}

	_, err := c.PatchResource(path, data)
	return err
}
