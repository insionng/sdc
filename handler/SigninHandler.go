package handler

import (
	//"fmt"
	"github.com/astaxie/beego"
	//"html/template"
	"sdc/helper"
	"sdc/lib"
	"sdc/model"
	//"strconv"
)

type SigninHandler struct {
	lib.BaseHandler
}

func (self *SigninHandler) Get() {

	remember, ckerr := self.Ctx.Request.Cookie("remember")
	sess_username, _ := self.GetSession("username").(string)

	//signbar的值为 0则关闭提示栏  1则显示提示栏
	self.Ctx.SetCookie("signbar", "0", 31536000, "/")
	self.Ctx.SetCookie("remember", "on", 31536000, "/")

	//如果未登录
	if sess_username == "" {
		if ckerr == nil {
			if remember.Value == "on" {
				self.Data["remember"] = "on"
			} else {
				self.Data["remember"] = nil
			}
		}
		self.TplNames = "sdc/signin.html"

	} else { //如果已登录
		self.Redirect("/", 302)
	}

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

func (self *SigninHandler) Post() {
	self.TplNames = "sdc/signin.html"

	flash := beego.NewFlash()
	email := self.GetString("email")
	password := self.GetString("password")
	remember := self.GetString("remember")

	if email == "" {
		flash.Error("EMAIL为空~")
		flash.Store(&self.Controller)

		return

	}

	if password == "" {
		flash.Error("密码为空~")
		flash.Store(&self.Controller)

		return

	}

	if helper.CheckEmail(email) == false {
		flash.Error("Email格式不合符规格~")
		flash.Store(&self.Controller)

		return

	}

	if helper.CheckPassword(password) == false {
		flash.Error("密码含有非法字符或密码过短(至少4~30位密码)!")
		flash.Store(&self.Controller)

		return

	}

	if usrinfo, err := model.GetUserByEmail(email); usrinfo != nil && err == nil {

		if helper.Validate_hash(usrinfo.Password, password) {

			//登录成功设置session
			self.SetSession("userid", usrinfo.Id)
			self.SetSession("username", usrinfo.Username)
			self.SetSession("userrole", usrinfo.Role)
			self.SetSession("useremail", usrinfo.Email)
			self.SetSession("usercontent", usrinfo.Content)

			//设置cookie

			//设置提示栏cookie标记
			//signbar的值为 0则关闭提示栏  1则显示提示栏
			self.Ctx.SetCookie("signbar", "0", 31536000, "/")
			if remember == "on" {

				self.Ctx.SetCookie("remember", "on", 31536000, "/")
			} else {

				self.Ctx.SetCookie("remember", "off", 31536000, "/")
			}
			self.Redirect("/", 302)
		} else {

			flash.Error("密码无法通过校验~")
			flash.Store(&self.Controller)
			return
		}
	} else {

		flash.Error("该账号不存在~")
		flash.Store(&self.Controller)
		return
	}
}
