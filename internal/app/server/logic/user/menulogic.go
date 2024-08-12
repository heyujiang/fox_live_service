package user

var MenuLogic = newMenuLogic()

type (
	menuLogic struct{}

	RespMenuListItem struct {
		Name      string              `json:"name,omitempty"`
		Meta      *MenuMeta           `json:"meta,omitempty"`
		Component string              `json:"component,omitempty"`
		Children  []*RespMenuListItem `json:"children,omitempty"`
		Path      string              `json:"path,omitempty"`
		Icon      string              `json:"icon,omitempty"`
		Redirect  string              `json:"redirect,omitempty"`
	}

	MenuMeta struct {
		Id                 int      `json:"id,omitempty"`
		Roles              []string `json:"roles,omitempty"`        //roles?: string[]; // Controls roles that have access to the page
		RequiresAuth       bool     `json:"requiresAuth,omitempty"` //requiresAuth: boolean; // Whether login is required to access the current page (every route must declare)
		Icon               string   `json:"icon,omitempty"`         //icon?: string; // The icon show in the side menu
		Locale             string   `json:"locale,omitempty"`       //locale?: string; // The locale name show in side menu and breadcrumb
		HideInMenu         bool     `json:"hideInMenu,omitempty"`   //hideInMenu?: boolean; // If true, it is not displayed in the side menu
		HideChildrenInMenu bool     `json:"hideChildrenInMenu"`     //hideChildrenInMenu?: boolean; // if set true, the children are not displayed in the side menu
		ActiveMenu         string   `json:"activeMenu,omitempty"`   //activeMenu?: string; // if set name, the menu will be highlighted according to the name you set
		Order              int      `json:"order,omitempty"`        //order?: number; // Sort routing menu items. If set key, the higher the value, the more forward it is
		NoAffix            bool     `json:"noAffix,omitempty"`      //noAffix?: boolean; // if set true, the tag will not affix in the tab-bar
		IgnoreCache        bool     `json:"ignoreCache,omitempty"`  //ignoreCache?: boolean; // if set true, the page will not be cached
		Title              string   `json:"title,omitempty"`
	}
)

func newMenuLogic() *menuLogic {
	return &menuLogic{}
}

func (m *menuLogic) GetMenus(uid int) ([]*RespMenuListItem, error) {
	res := []*RespMenuListItem{
		{
			Name: "home",
			Meta: &MenuMeta{
				Id:           8,
				Icon:         "icon-dashboard",
				RequiresAuth: true,
				Locale:       "项目概览",
				IgnoreCache:  true,
			},
			Component: "/dashboard/workplace/index",
			Path:      "/workplace",
		},
		{
			Name: "project",
			Meta: &MenuMeta{
				Id:          69,
				HideInMenu:  true,
				Locale:      "项目管理",
				IgnoreCache: true,
				Icon:        "icon-menu-unfold",
			},
			Component: "LAYOUT",
			Path:      "/project",
			Redirect:  "/project/project",
			Children: []*RespMenuListItem{
				{
					Name:      "project",
					Component: "/project/project/index",
					Meta: &MenuMeta{
						Id:           671,
						Locale:       "项目列表",
						IgnoreCache:  true,
						RequiresAuth: true,
					},
					Path: "project",
				},
				{
					Name:      "addproject",
					Component: "/project/addproject/index",
					Meta: &MenuMeta{
						Id:           672,
						Locale:       "新建项目",
						IgnoreCache:  true,
						RequiresAuth: true,
					},
					Path: "addproject",
				},
				{
					Name:      "detail",
					Component: "/project/detail/index",
					Meta: &MenuMeta{
						Id:           677,
						Locale:       "项目详情",
						IgnoreCache:  true,
						RequiresAuth: true,
						HideInMenu:   true,
					},
					Path: "detail",
				},
			},
		},
		{
			Name: "project_ing",
			Meta: &MenuMeta{
				Id:          70,
				HideInMenu:  true,
				Locale:      "项目进度",
				IgnoreCache: true,
				Icon:        "icon-schedule",
			},
			Component: "LAYOUT",
			Path:      "/project",
			Redirect:  "/project/project_ing",
			Children: []*RespMenuListItem{
				{
					Name:      "project_ing",
					Component: "/project/project_ing/index",
					Meta: &MenuMeta{
						Id:           673,
						Locale:       "项目进度列表",
						IgnoreCache:  true,
						RequiresAuth: true,
					},
					Path: "project_ing",
				},
			},
		}, {
			Name: "system",
			Meta: &MenuMeta{
				Id:          70,
				HideInMenu:  true,
				Locale:      "系统管理",
				IgnoreCache: true,
				Icon:        "icon-settings",
			},
			Component: "LAYOUT",
			Path:      "/system",
			Redirect:  "/system/account",
			Children: []*RespMenuListItem{
				{
					Component: "/system/account/index",
					Meta: &MenuMeta{
						Id:           121,
						Icon:         "icon-user",
						Locale:       "用户管理",
						RequiresAuth: true,
						IgnoreCache:  true,
					},
					Name: "account",
					Path: "account",
				},
				{
					Component: "/system/dept/index",
					Meta: &MenuMeta{
						Id:           121,
						Icon:         "icon-user",
						Locale:       "部门管理",
						RequiresAuth: true,
						IgnoreCache:  true,
					},
					Name: "dept",
					Path: "dept",
				},
				{
					Component: "/system/role/index",
					Meta: &MenuMeta{
						Id:           121,
						Icon:         "icon-user",
						Locale:       "角色管理",
						RequiresAuth: true,
						IgnoreCache:  true,
					},
					Name: "role",
					Path: "role",
				},
				{
					Component: "/system/rule/index",
					Meta: &MenuMeta{
						Id:           121,
						Icon:         "icon-user",
						Locale:       "权限管理",
						RequiresAuth: true,
						IgnoreCache:  true,
					},
					Name: "rule",
					Path: "rule",
				},
			},
		},
		{
			Component: "/system/node/index",
			Meta: &MenuMeta{
				Id:           12,
				Icon:         "icon-branch",
				Locale:       "节点管理",
				RequiresAuth: true,
				IgnoreCache:  true,
			},
			Name: "node",
			Path: "node",
		},
		{
			Name: "user_center",
			Meta: &MenuMeta{
				Id:          73,
				HideInMenu:  true,
				Locale:      "个人中心",
				IgnoreCache: true,
				Icon:        "icon-user",
			},
			Component: "LAYOUT",
			Path:      "/user",
			Redirect:  "/user/info",
			Children: []*RespMenuListItem{
				{
					Name:      "user_info",
					Component: "/user/info/index",
					Meta: &MenuMeta{
						Id:           678,
						Locale:       "我的项目",
						IgnoreCache:  true,
						RequiresAuth: true,
					},
					Path: "info",
				},
				{
					Name:      "user_weekday",
					Component: "/user/weekday/index",
					Meta: &MenuMeta{
						Id:           680,
						Locale:       "周报日报",
						IgnoreCache:  true,
						RequiresAuth: true,
					},
					Path: "weekday",
				},
				{
					Name:      "user_setting",
					Component: "/user/setting/index",
					Meta: &MenuMeta{
						Id:          679,
						Locale:      "个人设置",
						IgnoreCache: true,
					},
					Path: "setting",
				},
			},
		},
	}

	return res, nil
}
