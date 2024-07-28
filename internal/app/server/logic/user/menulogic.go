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
				Icon:         "icon-dashboard",
				Id:           8,
				RequiresAuth: true,
				Locale:       "项目概览",
				IgnoreCache:  true,
			},
			Component: "LAYOUT",
			Path:      "/home",
			Redirect:  "/home/workplace",
			Children: []*RespMenuListItem{
				{
					Name: "home",
					Meta: &MenuMeta{
						Id:           8,
						RequiresAuth: true,
						Locale:       "工作台",
						IgnoreCache:  true,
					},
					Component: "/dashboard/workplace/index",
					Path:      "workplace",
				},
			},
		},
		{
			Name: "system",
			Meta: &MenuMeta{
				Id:          61,
				Icon:        "icon-settings",
				Locale:      "系统设置",
				IgnoreCache: true,
			},
			Component: "LAYOUT",
			Path:      "/system",
			Redirect:  "/system/node",
			Children: []*RespMenuListItem{
				{
					Component: "/system/node/index",
					Meta: &MenuMeta{
						Id:           12,
						Locale:       "节点管理",
						RequiresAuth: true,
						IgnoreCache:  true,
					},
					Name: "node",
					Path: "node",
				},
				{
					Component: "/system/account/index",
					Meta: &MenuMeta{
						Id:           121,
						Locale:       "用户管理",
						RequiresAuth: true,
						IgnoreCache:  true,
					},
					Name: "account",
					Path: "account",
				},
			},
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
			Redirect:  "/project/list",
			Children: []*RespMenuListItem{
				{
					Name:      "list",
					Component: "/project/list/index",
					Meta: &MenuMeta{
						Id:           671,
						Locale:       "项目列表",
						IgnoreCache:  true,
						RequiresAuth: true,
					},
					Path: "list",
				},
			},
		},
		{
			Name: "user",
			Meta: &MenuMeta{
				Id:          69,
				HideInMenu:  true,
				Locale:      "用户管理",
				IgnoreCache: true,
			},
			Component: "LAYOUT",
			Path:      "/user",
			Redirect:  "/user/info",
			Children: []*RespMenuListItem{
				{
					Name:      "info",
					Component: "/user/info/index",
					Meta: &MenuMeta{
						Id:          670,
						HideInMenu:  true,
						Locale:      "用户信息",
						IgnoreCache: true,
					},
					Path: "info",
				},
				{
					Name:      "setting",
					Component: "/user/setting/index",
					Meta: &MenuMeta{
						Id:          61,
						HideInMenu:  true,
						Locale:      "用户信息",
						IgnoreCache: true,
					},
					Path: "setting",
				},
			},
		},
	}

	//{
	//	"component": "/dashboard/workplace/index",
	//	"meta": {
	//	"icon": "icon-dashboard",
	//		"id": 8,
	//		"ignoreCache": true,
	//		"locale": "概况",
	//		"requiresAuth": true
	//},
	//	"name": "home",
	//	"path": "/home"
	//}
	return res, nil
}
