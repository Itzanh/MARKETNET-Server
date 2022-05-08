package main

import "gorm.io/gorm"

// When adding fields to this table/struct, add them also to:
// The Permissions struct, and the getUserPermissions function at the bottom of this file.
// The getUserGroupsIn function at the beginning of the user_group.go file.
// The initialGroup function in the initial_data.go file, the admin always has all the permissions set.
type Group struct {
	Id            int32    `json:"id" gorm:"index:group_id_enterprise,unique:true,priority:1"`
	Name          string   `json:"name" gorm:"type:character varying(50);not null:true"`
	Sales         bool     `json:"sales" gorm:"not null:true"`
	Purchases     bool     `json:"purchases" gorm:"not null:true"`
	Masters       bool     `json:"masters" gorm:"not null:true"`
	Warehouse     bool     `json:"warehouse" gorm:"not null:true"`
	Manufacturing bool     `json:"manufacturing" gorm:"not null:true"`
	Preparation   bool     `json:"preparation" gorm:"not null:true"`
	Admin         bool     `json:"admin" gorm:"not null:true"`
	PrestaShop    bool     `json:"prestashop" gorm:"column:prestashop;not null:true"`
	Accounting    bool     `json:"accounting" gorm:"not null:true"`
	EnterpriseId  int32    `json:"-" gorm:"column:enterprise;not null:true;index:group_id_enterprise,unique:true,priority:2"`
	Enterprise    Settings `json:"-" gorm:"foreignKey:EnterpriseId;references:Id"`
	PointOfSale   bool     `json:"pointOfSale" gorm:"not null:true"`
}

func (g *Group) TableName() string {
	return "group"
}

func getGroup(enterpriseId int32) []Group {
	var groups []Group = make([]Group, 0)
	dbOrm.Model(&Group{}).Where("enterprise = ?", enterpriseId).Order("id ASC").Find(&groups)
	return groups
}

func getGroupsPermissionDictionary(enterpriseId int32, permission string) []Group {
	var groups []Group = make([]Group, 0)
	var permissionDictionaryGroups []PermissionDictionaryGroup = make([]PermissionDictionaryGroup, 0)
	dbOrm.Model(&PermissionDictionaryGroup{}).Where(`permission_key = ? AND enterprise = ?`, permission, enterpriseId).Order(`"group" ASC`).Preload("Group").Find(&permissionDictionaryGroups)

	for i := 0; i < len(permissionDictionaryGroups); i++ {
		groups = append(groups, permissionDictionaryGroups[i].Group)
	}

	return groups
}

func (g *Group) isValid() bool {
	return !(len(g.Name) > 50)
}

func (g *Group) BeforeCreate(tx *gorm.DB) (err error) {
	var group Group
	tx.Model(&Group{}).Last(&group)
	g.Id = group.Id + 1
	return nil
}

func (g *Group) insertGroup() bool {
	if !g.isValid() {
		return false
	}

	result := dbOrm.Create(&g)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (g *Group) updateGroup() bool {
	if g.Id <= 0 || !g.isValid() {
		return false
	}

	var group Group
	result := dbOrm.Where("id = ? AND enterprise = ?", g.Id, g.EnterpriseId).First(&group)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	group.Name = g.Name
	group.Sales = g.Sales
	group.Purchases = g.Purchases
	group.Masters = g.Masters
	group.Warehouse = g.Warehouse
	group.Manufacturing = g.Manufacturing
	group.Preparation = g.Preparation
	group.Admin = g.Admin
	group.PrestaShop = g.PrestaShop
	group.Accounting = g.Accounting
	group.PointOfSale = g.PointOfSale

	result = dbOrm.Save(&group)
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

func (g *Group) deleteGroup() bool {
	if g.Id <= 0 {
		return false
	}

	result := dbOrm.Where("id = ? AND enterprise = ?", g.Id, g.EnterpriseId).Delete(&Group{})
	if result.Error != nil {
		log("DB", result.Error.Error())
		return false
	}

	return true
}

type Permissions struct {
	Sales                bool     `json:"sales"`
	Purchases            bool     `json:"purchases"`
	Masters              bool     `json:"masters"`
	Warehouse            bool     `json:"warehouse"`
	Manufacturing        bool     `json:"manufacturing"`
	Preparation          bool     `json:"preparation"`
	Admin                bool     `json:"admin"`
	PrestaShop           bool     `json:"prestashop"`
	Accounting           bool     `json:"accounting"`
	PointOfSale          bool     `json:"pointOfSale"`
	PermissionDictionary []string `json:"permissionDictionary"`
}

func getUserPermissions(userId int32, enterpriseId int32) Permissions {
	ug := getUserGroups(userId, enterpriseId)
	p := Permissions{}

	for i := 0; i < len(ug.GroupsIn); i++ {
		if ug.GroupsIn[i].Sales {
			p.Sales = true
		}
		if ug.GroupsIn[i].Purchases {
			p.Purchases = true
		}
		if ug.GroupsIn[i].Masters {
			p.Masters = true
		}
		if ug.GroupsIn[i].Warehouse {
			p.Warehouse = true
		}
		if ug.GroupsIn[i].Manufacturing {
			p.Manufacturing = true
		}
		if ug.GroupsIn[i].Preparation {
			p.Preparation = true
		}
		if ug.GroupsIn[i].Admin {
			p.Admin = true
		}
		if ug.GroupsIn[i].PrestaShop {
			p.PrestaShop = true
		}
		if ug.GroupsIn[i].Accounting {
			p.Accounting = true
		}
		if ug.GroupsIn[i].PointOfSale {
			p.PointOfSale = true
		}
	}

	p.PermissionDictionary = getPermissionDictionaryUserGroupInForWebClient(userId)
	return p
}
