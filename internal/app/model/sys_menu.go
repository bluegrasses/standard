package model

type SysMenu struct {
	Model
	Name       string    `gorm:"comment:'菜单名称(英文名, 可用于国际化)'" json:"name"`
	Title      string    `gorm:"comment:'菜单标题(无法国际化时使用)'" json:"title"`
	Icon       string    `gorm:"comment:'菜单图标'" json:"icon"`
	Path       string    `gorm:"comment:'菜单访问路径'" json:"path"`
	Redirect   string    `gorm:"comment:'重定向路径'" json:"redirect"`
	Component  string    `gorm:"comment:'前端组件路径'" json:"component"`
	Permission string    `gorm:"comment:'权限标识'" json:"permission"`
	Sort       int       `gorm:"type:int(3);comment:'菜单顺序(同级菜单, 从0开始, 越小显示越靠前)'" json:"sort"`
	Status     *bool     `gorm:"type:tinyint(1);default:1;comment:'菜单状态(正常/禁用, 默认正常)'" json:"status"` // 由于设置了默认值, 这里使用ptr, 可避免赋值失败
	Visible    *bool     `gorm:"type:tinyint(1);default:1;comment:'菜单可见性(可见/隐藏, 默认可见)'" json:"visible"`
	Breadcrumb *bool     `gorm:"type:tinyint(1);default:1;comment:'面包屑可见性(可见/隐藏, 默认可见)'" json:"breadcrumb"`
	ParentId   uint      `gorm:"default:0;comment:'父菜单编号(编号为0时表示根菜单)'" json:"parentId"`
	Creator    string    `gorm:"comment:'创建人'" json:"creator"`
	Children   []SysMenu `gorm:"-" json:"children"`
	Roles      []SysRole `json:"roles" gorm:"many2many:relation_role_menu;"`
}

func (m SysMenu) TableName() string {
	return "sys_menu"
}
