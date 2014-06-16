package main

import (
	"github.com/astaxie/beego"
	"sdc/helper"
	//blogs_handler "sdc/logic/blogs/handlers"
	"sdc/core"
	//"sdc/logic/mzr/handler"
	"sdc/handler"
)

func main() {

	//未登录用户的数据 不用缓存，而是直接 用静态文件即可
	//主要存储静态化后的话题页面
	beego.SetStaticPath("/file", "file")
	beego.SetStaticPath("/static", "static")

	//URL定义规范:必须以/结尾

	//首页
	beego.Router("/", &handler.MainHandler{})
	beego.Router("/question/", &handler.MainHandler{})
	//首页 ?page
	beego.Router("/page-:page([1-9]\\d*)/", &handler.MainHandler{})

	//首页 hotness类
	beego.Router("/:tab([A-Za-z]+)/", &handler.MainHandler{})
	//http://localhost/lastest/page-2/
	beego.Router("/:tab([A-Za-z]+)/page-:page([1-9]\\d*)/", &handler.MainHandler{})

	//详情页面
	beego.Router("/:qid([1-9]\\d*)/", &handler.QuestionHandler{})
	beego.Router("/question/:qid([1-9]\\d*)/", &handler.QuestionHandler{})

	//搜索话题
	beego.Router("/search/", &handler.SearchHandler{})
	//同时支持page和keyword参数
	beego.Router("/search/:keyword([\\+\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)/page-:page([1-9]\\d*)/", &handler.SearchHandler{})
	//支持keyword参数 为了兼容所有搜索条件 这里需要用:all
	beego.Router("/search/:keyword([\\+\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)/", &handler.SearchHandler{})
	//搜索标签
	beego.Router("/tag/:keyword([\\+\\x{4e00}-\\x{9fa5}A-Z0-9a-z_-]+)/", &handler.SearchHandler{})

	//创建问题
	beego.Router("/new/question/", &handler.NewQuestionHandler{})

	//创建答案
	beego.Router("/new/answer/:qid:int/", &handler.NewAnswerHandler{})

	//采纳答案
	beego.Router("/accept/answer/:aid:int/:qid:int/", &handler.AcceptAnswerHandler{})

	//忽略答案
	beego.Router("/ignore/answer/:aid:int/:qid:int/", &handler.IgnoreAnswerHandler{})

	//关闭问题
	beego.Router("/close/question/:qid:int/", &handler.CloseQuestionHandler{})

	//开放问题
	beego.Router("/open/question/:qid:int/", &handler.OpenQuestionHandler{})

	//删除答案
	beego.Router("/delete/answer/:aid:int/:qid:int/", &handler.DeleteAnswerHandler{})

	//编辑问题
	beego.Router("/edit/question/:qid:int/", &handler.EditQuestionHandler{})

	//编辑答案
	beego.Router("/edit/answer/:aid:int/", &handler.EditAnswerHandler{})

	//删除问题
	beego.Router("/delete/question/:qid:int/", &handler.DeleteQuestionHandler{})
	//删除话题
	//beego.Router("/delete/topic/:tid:int/", &handler.DeleteTopicHandler{})

	//访问次数
	beego.Router("/view/question/:name([A-Za-z]+)/:id:int/", &handler.ViewQuestionHandler{})

	//hotness
	beego.Router("/like/:name([A-Za-z]+)/:id:int/", &handler.LikeHandler{})
	beego.Router("/hate/:name([A-Za-z]+)/:id:int/", &handler.HateHandler{})

	//登录
	beego.Router("/user/signin/", &handler.SigninHandler{})
	//退出
	beego.Router("/user/signout/", &handler.SignoutHandler{})
	//注册
	beego.Router("/user/signup/", &handler.SignupHandler{})

	//核心接口 话题接口
	beego.RESTRouter("/core/topic/", &core.TopicHandler{})
	//beego.Router("/core/node", &core.NodeHandler{})

	//blogs start
	/*
		beego.Router("/blogs/", &blogs_handler.MainHandler{})
		beego.Router("/blogs/category/:cid:int", &blogs_handler.MainHandler{})
		beego.Router("/blogs/search", &blogs_handler.SearchHandler{})

		beego.Router("/blogs/node/:nid:int", &blogs_handler.NodeHandler{})
		beego.Router("/blogs/view/:tid:int", &blogs_handler.ViewHandler{})

		beego.Router("/blogs/like/:name:string/:id:int", &blogs_handler.LikeHandler{})
		beego.Router("/blogs/hate/:name:string/:id:int", &blogs_handler.HateHandler{})

		beego.Router("/blogs/new/category", &blogs_handler.NewCategoryHandler{})
		beego.Router("/blogs/new/node", &blogs_handler.NewNodeHandler{})
		beego.Router("/blogs/new/topic", &blogs_handler.NewTopicHandler{})
		beego.Router("/blogs/new/reply/:tid:int", &blogs_handler.NewReplyHandler{})

		beego.Router("/blogs/modify/category", &blogs_handler.ModifyCategoryHandler{})
		beego.Router("/blogs/modify/node", &blogs_handler.ModifyNodeHandler{})

		beego.Router("/blogs/topic/delete/:tid:int", &blogs_handler.TopicDeleteHandler{})
		beego.Router("/blogs/topic/edit/:tid:int", &blogs_handler.TopicEditHandler{})

		beego.Router("/blogs/node/delete/:nid:int", &blogs_handler.NodeDeleteHandler{})
		beego.Router("/blogs/node/edit/:nid:int", &blogs_handler.NodeEditHandler{})

		beego.Router("/blogs/delete/reply/:rid:int", &blogs_handler.DeleteReplyHandler{})
	*/
	//blogs end

	/*
		beego.Router("/category/:cid:int/", &handler.CategoryHandler{})
		beego.Router("/category/:cid:int/page-:page([1-9]\\d*)/", &handler.CategoryHandler{})

		beego.Router("/category/:tab([A-Za-z]+)/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/", &handler.CategoryHandler{})
		//http://localhost/category/lastest/page-2/
		beego.Router("/category/:tab([A-Za-z]+)/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/page-:page([1-9]\\d*)/", &handler.CategoryHandler{})

		beego.Router("/category/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/", &handler.CategoryHandler{})
		beego.Router("/category/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/page-:page([1-9]\\d*)/", &handler.CategoryHandler{})
	*/

	//浏览节点 "/node/:tab([A-Za-z]+)/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/"优先级必须高于"/node/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/"
	/*
		beego.Router("/node/:nid:int/", &handler.NodeHandler{})
		beego.Router("/node/:nid:int/page-:page([1-9]\\d*)/", &handler.NodeHandler{})

		beego.Router("/node/:tab([A-Za-z]+)/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/", &handler.NodeHandler{})
		beego.Router("/node/:tab([A-Za-z]+)/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/page-:page([1-9]\\d*)/", &handler.NodeHandler{})

		beego.Router("/node/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/", &handler.NodeHandler{})
		beego.Router("/node/:name([\\x{4e00}-\\x{9fa5}A-Za-z]+)/page-:page([1-9]\\d*)/", &handler.NodeHandler{})
	*/

	//捕抓话题  (获取图片 获取文本 等等..)
	//beego.Router("/catch/topic/", &handler.NewTopicHandler{})

	/*
		beego.Router("/timeline/", &handler.TimelineHandler{})
		beego.Router("/user/:username([A-Za-z]+)/", &handler.TimelineHandler{})
		beego.Router("/userid/:userid:int/", &handler.TimelineHandler{})
		//发布时光记录
		beego.Router("/new/timeline/", &handler.TimelineHandler{})
		//删除时光
		beego.Router("/delete/timeline/:lid:int/", &handler.DeleteTimelineHandler{})
	*/

	//发现话题 以汇总资讯为方向
	//beego.Router("/discover/topic/", &handler.DiscoverHandler{})

	//浏览单图
	//beego.Router("/image/:mid:int/", &handler.ImageHandler{})

	//创建分类
	//beego.Router("/new/category/", &handler.NewCategoryHandler{})

	//创建节点
	//beego.Router("/new/node/", &handler.NewNodeHandler{})

	//beego.Router("/delete/reply/:rid([0-9]+)", &handler.DeleteReplyHandler{})

	//编辑分类
	//beego.Router("/edit/category/", &handler.EditCategoryHandler{})
	//编辑节点
	//beego.Router("/edit/node/", &handler.EditNodeHandler{})
	//编辑话题
	//beego.Router("/edit/topic/:tid:int/", &handler.EditTopicHandler{})

	/*
		beego.Router("/delete/node/:nid([0-9]+)", &handler.NodeDeleteHandler{})

	*/

	//个人设定
	//beego.Router("/user/settings/", &handler.Settings{})
	//beego.AutoRouter(&handler.Settings{})

	//beego.Router("/avatar/:username([A-Za-z]+)/:filename([A-Za-z]+)/", &handler.AvatarHandler{})

	//外部URL路由
	//beego.Router("/url/", &handler.UrlHandler{})

	//上传文件
	//beego.Router("/upload/", &handler.UploaderHandler{})

	//模板函数
	beego.AddFuncMap("timesince", helper.TimeSince)
	beego.AddFuncMap("tags", helper.Tags)
	beego.AddFuncMap("metric", helper.Metric)
	beego.AddFuncMap("gravatar", helper.Gravatar)
	beego.AddFuncMap("markdown", helper.Markdown)
	beego.AddFuncMap("markdown2text", helper.Markdown2Text)

	beego.SessionOn = true
	beego.SessionName = "sdc"
	//beego.SessionProvider = "file"
	//beego.SessionSavePath = "./session"
	beego.AutoRender = true
	beego.CopyRequestBody = true //必须开启,不然core api部分会无法正常工作

	//runtime.GOMAXPROCS(2)
	beego.Run()
}
