/*
This file is part of MARKETNET.

MARKETNET is free software: you can redistribute it and/or modify it under the terms of the GNU Affero General Public License as published by the Free Software Foundation, version 3 of the License.

MARKETNET is distributed in the hope that it will be useful, but WITHOUT ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License along with MARKETNET. If not, see <https://www.gnu.org/licenses/>.
*/

package main

import "gorm.io/gorm/clause"

type UserGroup struct {
	UserId  int32 `json:"userId" gorm:"primaryKey;column:user;not null:true"`
	User    User  `json:"user" gorm:"foreignKey:UserId;references:Id"`
	GroupId int32 `json:"groupId" gorm:"primaryKey;column:group;not null:true"`
	Group   Group `json:"group" gorm:"foreignKey:GroupId;references:Id"`
}

func (ug *UserGroup) TableName() string {
	return "user_group"
}

type UserGroups struct {
	GroupsIn  []Group `json:"groupsIn"`
	GroupsOut []Group `json:"groupsOut"`
}

func getUserGroups(userId int32, enterpriseId int32) UserGroups {
	groupsIn := getUserGroupsIn(userId)
	return UserGroups{
		GroupsIn:  groupsIn,
		GroupsOut: getUserGroupsOut(userId, groupsIn, enterpriseId),
	}
}

func getGroupUsers(groupId int32, enterpriseId int32) []User {
	// get group row
	var group Group
	result := dbOrm.Where("\"id\" = ?", groupId).First(&group)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return nil
	}
	if group.EnterpriseId != enterpriseId {
		return nil
	}

	var users []User = make([]User, 0)
	var userGroups []UserGroup = make([]UserGroup, 0)
	result = dbOrm.Model(&UserGroup{}).Where("\"group\" = ?", groupId).Order("\"user\" ASC").Preload(clause.Associations).Find(&userGroups)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return users
	}

	for i := 0; i < len(userGroups); i++ {
		users = append(users, userGroups[i].User)
	}

	return users
}

func getUserGroupsIn(userId int32) []Group {
	var userGroups []UserGroup = make([]UserGroup, 0)
	var groups []Group = make([]Group, 0)
	result := dbOrm.Model(&UserGroup{}).Where("\"user\" = ?", userId).Order("\"group\" ASC").Preload(clause.Associations).Find(&userGroups)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return groups
	}

	for i := 0; i < len(userGroups); i++ {
		groups = append(groups, userGroups[i].Group)
	}

	return groups
}

func getUserGroupsOut(userId int32, groupsIn []Group, enterpriseId int32) []Group {
	groups := getGroup(enterpriseId)

	for i := 0; i < len(groupsIn); i++ {
		for j := len(groups) - 1; j >= 0; j-- {
			if groupsIn[i].Id == groups[j].Id {
				groups = append(groups[0:j], groups[j+1:]...)
				break
			}
		}
	}

	return groups
}

func (u *UserGroup) isValid() bool {
	return !(u.UserId <= 0 || u.GroupId <= 0)
}

func (u *UserGroup) insertUserGroup() bool {
	if !u.isValid() {
		return false
	}

	result := dbOrm.Create(&u)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (u *UserGroup) deleteUserGroup() bool {
	if !u.isValid() {
		return false
	}

	result := dbOrm.Where("\"user\" = ? AND \"group\" = ?", u.UserId, u.GroupId).Delete(&UserGroup{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}
