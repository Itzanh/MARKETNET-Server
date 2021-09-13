package main

import "testing"

// ===== SETTINGS

func TestSettings(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := getSettingsRecord()
	ok := s.updateSettingsRecord()
	if !ok {
		t.Error("Can't update settings")
		return
	}
}

// ===== USERS

func TestGetUser(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	users := getUser()
	if len(users) == 0 || users[0].Id <= 0 {
		t.Error("Can't scan users")
		return
	}
}

func TestGetUserByUsername(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	user := getUserByUsername("marketnet")
	if user.Id <= 0 {
		t.Error("Can't scan user")
		return
	}
}

func TestGetUserRow(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	user := getUserRow(1)
	if user.Id <= 0 {
		t.Error("Can't scan user")
		return
	}
}

func TestUserInsertUpdateDeleteLoginPassword(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// insert
	ui := UserInsert{
		Username: "gotestuser",
		FullName: "Go test user",
		Password: "go.test.user",
		Language: "en",
	}
	ok := ui.insertUser()
	if !ok {
		t.Error("Insert error, user not inserted")
		return
	}

	// update
	users := getUser()
	user := users[len(users)-1]

	user.Language = "es"
	ok = user.updateUser()
	if !ok {
		t.Error("Update error, user not updated")
		return
	}

	// attempts incorrect login
	ul := UserLogin{
		Username: "gotestuser",
		Password: "go.user",
	}
	result, _ := ul.login("127.0.0.1")
	if result.Ok {
		t.Error("Can login with incorrect password!!!")
		return
	}

	// attempt correct login
	ul = UserLogin{
		Username: "gotestuser",
		Password: "go.test.user",
	}
	result, _ = ul.login("127.0.0.1")
	if !result.Ok {
		t.Error("Can't login!!!")
		return
	}

	// change password
	up := UserPassword{
		Id:       user.Id,
		Password: "go1234testuser",
	}
	ok = up.userPassword()
	if !ok {
		t.Error("Can't update the user's password")
		return
	}

	// attempts incorrect login
	ul = UserLogin{
		Username: "gotestuser",
		Password: "go.user",
	}
	result, _ = ul.login("127.0.0.1")
	if result.Ok {
		t.Error("Can login with incorrect password!!!")
		return
	}

	// attempt correct login
	ul = UserLogin{
		Username: "gotestuser",
		Password: "go1234testuser",
	}
	result, _ = ul.login("127.0.0.1")
	if !result.Ok {
		t.Error("Can't login!!!")
		return
	}

	// deactivate user
	ok = user.offUser()
	if !ok {
		t.Error("Can't deactivate user")
		return
	}

	// attempts incorrect login
	ul = UserLogin{
		Username: "gotestuser",
		Password: "go1234testuser",
	}
	result, _ = ul.login("127.0.0.1")
	if result.Ok {
		t.Error("Can login with deactivated user!!!")
		return
	}

	// reactivate user
	ok = user.offUser()
	if !ok {
		t.Error("Can't reactivate user")
		return
	}

	// attempt correct login
	ul = UserLogin{
		Username: "gotestuser",
		Password: "go1234testuser",
	}
	result, _ = ul.login("127.0.0.1")
	if !result.Ok {
		t.Error("Can't login!!!")
		return
	}

	// delete
	ok = user.deleteUser()
	if !ok {
		t.Error("Delete error, user not deleted")
		return
	}
}

// ===== GROUP

func TestGetGroup(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	groups := getGroup()
	if len(groups) == 0 || groups[0].Id <= 0 {
		t.Error("Can't scan groups")
		return
	}
}

func TestGroupInsertUpdateDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	g := Group{
		Name: "Test",
	}
	ok := g.insertGroup()
	if !ok {
		t.Error("Insert error, group not inserted")
		return
	}

	groups := getGroup()
	g = groups[len(groups)-1]
	g.Sales = true
	ok = g.updateGroup()
	if !ok {
		t.Error("Update error, group not updated")
		return
	}

	groups = getGroup()
	g = groups[len(groups)-1]
	if !g.Sales {
		t.Error("Update not successful")
		return
	}

	ok = g.deleteGroup()
	if !ok {
		t.Error("Delete error, group not deleted")
		return
	}
}

// ===== USER GROUP

func TestGetUserGroups(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	users := getUser()
	for i := 0; i < len(users); i++ {
		groups := getUserGroups(users[i].Id)
		for j := 0; j < len(groups.GroupsIn); j++ {
			if groups.GroupsIn[j].Id <= 0 {
				t.Error("Can't scan user groups in")
			}
		}
		for j := 0; j < len(groups.GroupsOut); j++ {
			if groups.GroupsOut[j].Id <= 0 {
				t.Error("Can't scan user groups in")
			}
		}
	}
}

func TestUserGroupInsertDelete(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// create user
	ui := UserInsert{
		Username: "gotestuser",
		FullName: "Go test user",
		Password: "go.test.user",
		Language: "en",
	}
	ui.insertUser()
	users := getUser()
	user := users[len(users)-1]

	// create group
	g := Group{
		Name: "Test",
	}
	g.insertGroup()
	groups := getGroup()
	g = groups[len(groups)-1]

	// create user group
	ug := UserGroup{
		User:  user.Id,
		Group: g.Id,
	}
	ok := ug.insertUserGroup()
	if !ok {
		t.Error("Insert error, user group not inserted")
		return
	}

	// delete user group
	userGroups := getUserGroups(user.Id)
	if len(userGroups.GroupsIn) == 0 {
		t.Error("Can't scan user groups")
		return
	}

	ug.deleteUserGroup()
	if !ok {
		t.Error("Delete error, can't delete user group")
		return
	}

	// delete group
	g.deleteGroup()

	// delete user
	user.deleteUser()
}

func TestPermissions(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// create user
	ui := UserInsert{
		Username: "gotestuser",
		FullName: "Go test user",
		Password: "go.test.user",
		Language: "en",
	}
	ui.insertUser()
	users := getUser()
	user := users[len(users)-1]

	// create groups
	gSales := Group{
		Name:  "Test sales",
		Sales: true,
	}
	gSales.insertGroup()
	gPurchases := Group{
		Name:      "Test purchases",
		Purchases: true,
	}
	gPurchases.insertGroup()

	// create user group
	ug := UserGroup{
		User:  user.Id,
		Group: gSales.Id,
	}
	ok := ug.insertUserGroup()
	if !ok {
		t.Error("Insert error, user group not inserted")
		return
	}
	ug = UserGroup{
		User:  user.Id,
		Group: gPurchases.Id,
	}
	ok = ug.insertUserGroup()
	if !ok {
		t.Error("Insert error, user group not inserted")
		return
	}

	// check permissions
	permissions := getUserPermissions(user.Id)
	if !permissions.Sales || !permissions.Purchases || permissions.Masters || permissions.Warehouse || permissions.Manufacturing || permissions.Preparation || permissions.Admin || permissions.PrestaShop || permissions.Accounting {
		t.Error("Permissions not set correctly")
		return
	}

	// delete user group
	ug = UserGroup{
		User:  user.Id,
		Group: gSales.Id,
	}
	ug.deleteUserGroup()
	if !ok {
		t.Error("Delete error, can't delete user group")
		return
	}
	ug = UserGroup{
		User:  user.Id,
		Group: gPurchases.Id,
	}
	ug.deleteUserGroup()
	if !ok {
		t.Error("Delete error, can't delete user group")
		return
	}

	// delete group
	gSales.deleteGroup()
	gPurchases.deleteGroup()

	// delete user
	user.deleteUser()
}

// ===== API KEYS

func TestGetApiKeys(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	apiKeys := getApiKeys()
	if len(apiKeys) > 0 && apiKeys[0].Id <= 0 {
		t.Error("Can't scan API keys")
		return
	}
}

func TestApiKeys(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	key := ApiKey{
		Name:        "Test key",
		UserCreated: 1,
		User:        1,
	}
	ok := key.insertApiKey()
	if !ok {
		t.Error("Insert error, can't insert api key")
		return
	}

	ok, _ = checkApiKey(key.Token)
	if !ok {
		t.Error("The API key can't get authenticated")
		return
	}

	apiKeys := getApiKeys()
	key = apiKeys[len(apiKeys)-1]

	ok = key.offApiKey()
	if !ok {
		t.Error("Can't deactivate API key")
		return
	}

	ok, _ = checkApiKey(key.Token)
	if !ok {
		t.Error("The API key can be accessed after deactivating")
		return
	}

	ok = key.offApiKey()
	if !ok {
		t.Error("Can't deactivate API key")
		return
	}

	ok, _ = checkApiKey(key.Token)
	if !ok {
		t.Error("The API key can't get authenticated")
		return
	}

	ok = key.deleteApiKey()
	if !ok {
		t.Error("Delete error, can't delete api key")
		return
	}
}
