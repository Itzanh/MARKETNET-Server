/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import "testing"

// ===== SETTINGS

func TestSettings(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := getSettingsRecordById(1)
	ok := s.updateSettingsRecord()
	if !ok {
		t.Error("Can't update settings")
		return
	}
}

func TestSettingsKey(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := getSettingsRecordByEnterprise("MARKETNET")
	ok := s.updateSettingsRecord()
	if !ok {
		t.Error("Can't update settings")
		return
	}
}

func TestSettingsAll(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	s := getSettingsRecords()
	if len(s) == 0 || s[0].Id <= 0 {
		t.Error("Can't scan")
		return
	}
}

// ===== USERS

func TestGetUser(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	users := getUser(1)
	if len(users) == 0 || users[0].Id <= 0 {
		t.Error("Can't scan users")
		return
	}
}

func TestGetUserByUsername(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	user := getUserByUsername(1, "marketnet")
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
	ok := ui.insertUser(1)
	if !ok {
		t.Error("Insert error, user not inserted")
		return
	}

	// update
	users := getUser(1)
	user := users[len(users)-1]

	user.Language = "es"
	user.EnterpriseId = 1
	ok = user.updateUser()
	if !ok {
		t.Error("Update error, user not updated")
		return
	}

	// attempts incorrect login
	ul := UserLogin{
		Username:   "gotestuser",
		Password:   "go.user",
		Enterprise: "MARKETNET",
	}
	result, _, _ := ul.login("127.0.0.1")
	if result.Ok {
		t.Error("Can login with incorrect password!!!")
		return
	}

	// attempt correct login
	ul = UserLogin{
		Username:   "gotestuser",
		Password:   "go.test.user",
		Enterprise: "MARKETNET",
	}
	result, _, _ = ul.login("127.0.0.1")
	if !result.Ok {
		t.Error("Can't login!!!")
		return
	}

	// change password
	up := UserPassword{
		Id:       user.Id,
		Password: "go1234testuser",
	}
	ok = up.userPassword(1)
	if !ok {
		t.Error("Can't update the user's password")
		return
	}

	// attempts incorrect login
	ul = UserLogin{
		Username:   "gotestuser",
		Password:   "go.user",
		Enterprise: "MARKETNET",
	}
	result, _, _ = ul.login("127.0.0.1")
	if result.Ok {
		t.Error("Can login with incorrect password!!!")
		return
	}

	// attempt correct login
	ul = UserLogin{
		Username:   "gotestuser",
		Password:   "go1234testuser",
		Enterprise: "MARKETNET",
	}
	result, _, _ = ul.login("127.0.0.1")
	if !result.Ok {
		t.Error("Can't login!!!")
		return
	}

	// deactivate user
	user.EnterpriseId = 1
	ok = user.offUser()
	if !ok {
		t.Error("Can't deactivate user")
		return
	}

	// attempts incorrect login
	ul = UserLogin{
		Username:   "gotestuser",
		Password:   "go1234testuser",
		Enterprise: "MARKETNET",
	}
	result, _, _ = ul.login("127.0.0.1")
	if result.Ok {
		t.Error("Can login with deactivated user!!!")
		return
	}

	// reactivate user
	ok = user.offUser()
	user.EnterpriseId = 1
	if !ok {
		t.Error("Can't reactivate user")
		return
	}

	// attempt correct login
	ul = UserLogin{
		Username:   "gotestuser",
		Password:   "go1234testuser",
		Enterprise: "MARKETNET",
	}
	result, _, _ = ul.login("127.0.0.1")
	if !result.Ok {
		t.Error("Can't login!!!")
		return
	}

	// delete
	user.EnterpriseId = 1
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

	groups := getGroup(1)
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
		Name:         "Test",
		EnterpriseId: 1,
	}
	ok := g.insertGroup()
	if !ok {
		t.Error("Insert error, group not inserted")
		return
	}

	groups := getGroup(1)
	g = groups[len(groups)-1]
	g.Sales = true
	g.EnterpriseId = 1
	ok = g.updateGroup()
	if !ok {
		t.Error("Update error, group not updated")
		return
	}

	groups = getGroup(1)
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

	users := getUser(1)
	for i := 0; i < len(users); i++ {
		groups := getUserGroups(users[i].Id, 1)
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
	ui.insertUser(1)
	users := getUser(1)
	user := users[len(users)-1]

	// create group
	g := Group{
		Name:         "Test",
		EnterpriseId: 1,
	}
	g.insertGroup()
	groups := getGroup(1)
	g = groups[len(groups)-1]

	// create user group
	ug := UserGroup{
		UserId:  user.Id,
		GroupId: g.Id,
	}
	ok := ug.insertUserGroup()
	if !ok {
		t.Error("Insert error, user group not inserted")
		return
	}

	// delete user group
	userGroups := getUserGroups(user.Id, 1)
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
	ui.insertUser(1)
	users := getUser(1)
	user := users[len(users)-1]

	// create groups
	gSales := Group{
		Name:         "Test sales",
		Sales:        true,
		EnterpriseId: 1,
	}
	gSales.insertGroup()
	gPurchases := Group{
		Name:         "Test purchases",
		Purchases:    true,
		EnterpriseId: 1,
	}
	gPurchases.insertGroup()

	// create user group
	ug := UserGroup{
		UserId:  user.Id,
		GroupId: gSales.Id,
	}
	ok := ug.insertUserGroup()
	if !ok {
		t.Error("Insert error, user group not inserted")
		return
	}
	ug = UserGroup{
		UserId:  user.Id,
		GroupId: gPurchases.Id,
	}
	ok = ug.insertUserGroup()
	if !ok {
		t.Error("Insert error, user group not inserted")
		return
	}

	// check permissions
	permissions := getUserPermissions(user.Id, 1)
	if !permissions.Sales || !permissions.Purchases || permissions.Masters || permissions.Warehouse || permissions.Manufacturing || permissions.Preparation || permissions.Admin || permissions.PrestaShop || permissions.Accounting {
		t.Error("Permissions not set correctly")
		return
	}

	// delete user group
	ug = UserGroup{
		UserId:  user.Id,
		GroupId: gSales.Id,
	}
	ug.deleteUserGroup()
	if !ok {
		t.Error("Delete error, can't delete user group")
		return
	}
	ug = UserGroup{
		UserId:  user.Id,
		GroupId: gPurchases.Id,
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

	apiKeys := getApiKeys(1)
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
		Name:          "Test key",
		UserCreatedId: 1,
		UserId:        1,
		EnterpriseId:  1,
		Auth:          "P",
	}
	ok := key.insertApiKey()
	if !ok {
		t.Error("Insert error, can't insert api key")
		return
	}

	ok, _, _, _ = checkApiKeyByTokenAuthType(*key.Token, "P")
	if !ok {
		t.Error("The API key can't get authenticated")
		return
	}

	apiKeys := getApiKeys(1)
	key = apiKeys[len(apiKeys)-1]

	ok = key.offApiKey()
	if !ok {
		t.Error("Can't deactivate API key")
		return
	}

	ok, _, _, _ = checkApiKeyByTokenAuthType(*key.Token, "P")
	if ok {
		t.Error("The API key can be accessed after deactivating")
		return
	}

	ok = key.offApiKey()
	if !ok {
		t.Error("Can't deactivate API key")
		return
	}

	ok, _, _, _ = checkApiKeyByTokenAuthType(*key.Token, "P")
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

func TestEvaluatePasswordSecureCloud(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	// test complexity
	res := evaluatePasswordSecureCloud(1, "AAAAAAAA")
	if res.PasswordComplexity == true {
		t.Error("Password complexity OK in incorrect password")
		return
	}

	res = evaluatePasswordSecureCloud(1, "12345678")
	if res.PasswordComplexity == true {
		t.Error("Password complexity OK in incorrect password")
		return
	}

	res = evaluatePasswordSecureCloud(1, "ABCD1234")
	if res.PasswordComplexity == false {
		t.Error("Password complexity ERR in correct password")
		return
	}

	// test blacklist
	sqlStatement := `SELECT pwd FROM public.pwd_blacklist LIMIT 1`
	row := db.QueryRow(sqlStatement)
	var passwdInBlacklist string
	row.Scan(passwdInBlacklist)

	res = evaluatePasswordSecureCloud(1, passwdInBlacklist)
	if res.PasswordInBlacklist == false {
		t.Error("Password blacklist OK in incorrect password")
		return
	}

	res = evaluatePasswordSecureCloud(1, "passwdInBlacklist")
	if res.PasswordInBlacklist == true {
		t.Error("Password blacklist ERR in correct password")
		return
	}

	// test hash blacklist
	insertSinglePwdBlacklistHash("miblacklist")
	res = evaluatePasswordSecureCloud(1, "miblacklist")
	if res.PasswordHashInBlacklist == false {
		t.Error("Password hash blacklist OK in incorrect password")
		return
	}
}

// ===== PERMISSION DICTIONARY

func TestGetPermissionDictionary(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	permissions := getPermissionDictionary(1)
	if len(permissions) == 0 || permissions[0].Key == "" {
		t.Error("Can't scan permisssion dictionary")
		return
	}
}

func TestGetGroupPermissionDictionary(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	permissions := getGroupPermissionDictionary(1, 1)
	if len(permissions.In) > 0 && permissions.In[0].PermissionKeyId == "" {
		t.Error("Can't scan permisssion dictionary")
		return
	}
	if len(permissions.Out) > 0 && permissions.Out[0].Key == "" {
		t.Error("Can't scan permisssion dictionary")
		return
	}
}

func TestInsertDeletePermissionDictionaryGroup(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	p := PermissionDictionaryGroup{
		GroupId:         1,
		PermissionKeyId: "CANT_CREATE_PRODUCT",
		EnterpriseId:    1,
	}
	if !p.insertPermissionDictionaryGroup() {
		t.Error("Can't insert permisssion dictionary group")
		return
	}
	if !p.deletePermissionDictionaryGroup() {
		t.Error("Can't delete permisssion dictionary group")
		return
	}
}

func TestGetPermissionDictionaryUserGroupInForWebClient(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	p := getPermissionDictionaryUserGroupInForWebClient(1)
	if len(p) > 0 && p[0] == "" {
		t.Error("Can't scan permisssion dictionary")
		return
	}
}

func TestGetUserPermission(t *testing.T) {
	if db == nil {
		ConnectTestWithDB(t)
	}

	permissions := getGroupPermissionDictionary(1, 2)
	if len(permissions.In) > 0 && !getUserPermission(permissions.In[0].PermissionKeyId, 1, 2) {
		t.Error("Permission error")
		return
	}
	if len(permissions.Out) > 0 && getUserPermission(permissions.Out[0].Key, 1, 2) {
		t.Error("Permission error")
		return
	}
}
