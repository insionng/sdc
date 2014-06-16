package lib

import (
	//"fmt"
	"github.com/astaxie/beego"
	"html/template"
	"runtime"
	"sdc/helper"
	//"sdc/model"
)

var (
	sess_username string
	sess_uid      int64
	sess_role     int64
	sess_email    string
	sess_content  string
	Counter       map[string]string
)

type BaseHandler struct {
	beego.Controller
}

type AuthHandler struct {
	BaseHandler
}

type ApiHandler struct {
	BaseHandler
}

type RootHandler struct {
	BaseHandler
}

/*
func init() {
	Cache = cache.NewMemCache()
	Cache.Every = 300 //該單位為秒，0為不過期，259200 三天,604800 即一個星期清空一次緩存 ,300s即 5分鐘
	Cache.Start()
}
*/

/*
func Cached(k string, v string, t int) string {
	ivalue := []byte(v)

	if Cache.IsExist(k) {
		ivalue = Cache.Get(k).([]byte)
	} else {
		Cache.Put(k, ivalue, t)
	}
	return string(ivalue)
}
*/

//用户等级划分：正数是普通用户，负数是管理员各种等级划分，为0则尚未注册
func (self *BaseHandler) Prepare() {
	/*
		sess := self.StartSession()
		defer sess.SessionRelease()
		fmt.Println("SID::", sess.SessionID())
	*/
	//从session里读出登录信息
	sess_username, _ = self.GetSession("username").(string)
	sess_uid, _ = self.GetSession("userid").(int64)
	sess_role, _ = self.GetSession("userrole").(int64)
	sess_email, _ = self.GetSession("useremail").(string)
	sess_content, _ = self.GetSession("usercontent").(string)
	//把登录信息写入模板容器

	self.Data["userid"] = sess_uid
	self.Data["username"] = sess_username
	self.Data["userrole"] = sess_role
	self.Data["useremail"] = sess_email
	self.Data["usercontent"] = sess_content

	self.Data["xsrfdata"] = template.HTML(self.XsrfFormHtml())
	//self.Data["xsrfdata"] = template.HTML(`<input type="hidden" name="_xsrf" value="` + self.XsrfToken() + `"/>`)
	//self.Data["xsrftoken"] = self.XsrfToken()

	//for uploadify and &root.RMultipleUploaderHandler{}
	//self.Data["token"] = utils.Encrypt_hash("00d5930debddeb10eb172bba94eed053d12c42c0428f3e09", nil)

	//加载网站常用数据
	/*
		if cate, err := model.GetCategorys(0, 0, "id"); err == nil {
			self.Data["categorys"] = cate
		} else {
			fmt.Println("分类数据查询出错", err)
		}
		if nd, err := model.GetNodes(0, 0, "created"); err == nil {
			self.Data["nodes"] = *nd
		} else {
			fmt.Println("节点数据查询出错", err)
		}

		if nd, err := model.GetNodes(0, 4, "views"); err == nil {
			self.Data["nodes_views_4"] = *nd
		} else {
			fmt.Println("热门节点(4_VIEWS)数据查询出错", err)
		}
		if nd, err := model.GetNodes(0, 20, "views"); err == nil {
			self.Data["nodes_views_20"] = *nd
		} else {
			fmt.Println("热门节点(20_VIEWS)数据查询出错", err)
		}
	*/
}

//会员或管理员前台权限认证
func (self *AuthHandler) Prepare() {
	self.BaseHandler.Prepare()

	if sess_role == 0 {
		self.Redirect("/user/signin/", 302)
	}
}

//API权限认证
func (self *ApiHandler) Prepare() {

	self.BaseHandler.Prepare()
	//返回401未认证状态终止服务
	if sess_role == 0 {
		self.Abort("401")
	}
}

//管理员后台后台认证
func (self *RootHandler) Prepare() {
	self.BaseHandler.Prepare()

	if !helper.IsSpider(self.Ctx.Request.UserAgent()) {
		if sess_role != -1000 {
			self.Redirect("/user/signin/", 302)
		} else {
			self.Data["remoteproto"] = self.Ctx.Request.Proto
			self.Data["remotehost"] = self.Ctx.Request.Host
			self.Data["remoteos"] = runtime.GOOS
			self.Data["remotearch"] = runtime.GOARCH
			self.Data["remotecpus"] = runtime.NumCPU()
			self.Data["golangver"] = runtime.Version()
		}
	} else {
		self.Redirect("/", 302)
	}
}
