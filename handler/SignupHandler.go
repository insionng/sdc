package handler

import (
	//"fmt"
	"github.com/astaxie/beego"
	//"html/template"
	"sdc/helper"
	"sdc/lib"
	"sdc/model"
	//"strconv"
	"strings"
)

type SignupHandler struct {
	lib.BaseHandler
}

func (self *SignupHandler) Get() {
	self.TplNames = "sdc/signup.html"

	//signbar的值为 0则关闭提示栏  1则显示提示栏
	self.Ctx.SetCookie("signbar", "0", 31536000, "/")

	//侧栏九宫格推荐榜单
	//先行取出最热门的9个节点 然后根据节点获取该节点下最热门的话题
	/*
		if nd, err := model.GetNodes(0, 9, "hotness"); err == nil {
			if len(*nd) > 0 {
				for _, v := range *nd {

					i := 0
					output_start := `<ul class="widgets-popular widgets-similar clx">`
					output := ""
					if tps := model.GetTopicsByNid(v.Id, 0, 1, 0, "hotness"); err == nil {

						if len(*tps) > 0 {
							for _, v := range *tps {

								i += 1
								if i == 3 {
									output = output + `<li class="similar similar-third">`
									i = 0
								} else {
									output = output + `<li class="similar">`
								}
								output = output + `<a target="_blank" href="/` + strconv.Itoa(int(v.Id)) + `/" title="` + v.Title + `">
													<img src="` + v.ThumbnailsSmall + `" wdith="70" height="70" />
												</a>
											</li>`
							}
						}
					}
					output_end := "</ul>"
					if len(output) > 0 {
						output = output_start + output + output_end
						self.Data["topic_hotness_9_module"] = template.HTML(output)
					} else {
						self.Data["topic_hotness_9_module"] = nil
					}

				}
			}
		} else {
			fmt.Println("节点数据查询出错", err)
		}
	*/
}

func (self *SignupHandler) Post() {

	self.TplNames = "sdc/signup.html"

	flash := beego.NewFlash()
	email := strings.TrimSpace(strings.ToLower(self.GetString("email")))
	username := strings.ToLower(self.GetString("username"))
	password := self.GetString("password")
	repassword := self.GetString("repassword")

	if password == "" {
		flash.Error("密码为空~")
		flash.Store(&self.Controller)

		return

	}

	if password != repassword {
		flash.Error("两次密码不匹配~")
		flash.Store(&self.Controller)

		return

	}

	if helper.CheckPassword(password) == false {
		flash.Error("密码含有非法字符或密码过短(至少4~30位密码)!")
		flash.Store(&self.Controller)

		return

	}

	if username == "" {
		flash.Error("用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空~")
		flash.Store(&self.Controller)

		return

	}

	if helper.CheckUsername(username) == false {
		flash.Error("用户名是为永久性设定,不能少于4个字或多于30个字,请慎重考虑,不能为空~")
		flash.Store(&self.Controller)

		return
	}

	if helper.CheckEmail(email) == false {
		flash.Error("Email格式不合符规格~")
		flash.Store(&self.Controller)

		return

	}

	if usrinfo, err := model.GetUserByEmail(email); usrinfo != nil {

		flash.Error("此账号不能使用~")
		flash.Store(&self.Controller)

		return

	} else if err != nil {

		flash.Error("检索账号期间出错~")
		flash.Store(&self.Controller)

		return
	}

	if usrid, err := model.AddUser(email, username, "", "", helper.Encrypt_hash(password, nil), 1); err != nil {
		flash.Error("用户注册信息写入数据库时发生错误~")
		flash.Store(&self.Controller)

		return

	} else {

		if usrinfo, err := model.GetUser(usrid); err == nil {

			//注册账号成功,以下自动登录并设置session
			self.SetSession("userid", usrid)
			self.SetSession("username", usrinfo.Username)
			self.SetSession("userrole", usrinfo.Role)
			self.SetSession("useremail", usrinfo.Email)
			self.SetSession("usercontent", usrinfo.Content)

			//signbar的值为 0则关闭提示栏  1则显示提示栏
			self.Ctx.SetCookie("signbar", "1", 31536000, "/")

			//设置cookie

			//设置提示栏cookie标记
			//signbar的值为 0则关闭提示栏  1则显示提示栏
			self.Ctx.SetCookie("signbar", "0", 31536000, "/")

			flash.Notice("账号登录成功~")
			flash.Store(&self.Controller)

			//session 写入后直接跳到首页
			self.Redirect("/", 302)

		} else {

			flash.Notice("注册账号成功,请手动登录~")
			flash.Store(&self.Controller)

			//注册成功后直接跳转到登录页
			self.Redirect("/signin/", 302)

		}

	}

}
