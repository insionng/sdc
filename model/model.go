package model

import (
	"errors"
	"fmt"
	//_ "github.com/lib/pq"
	"github.com/lunny/xorm"
	_ "github.com/mattn/go-sqlite3"
	"os"
	"sdc/helper"
	"strconv"
	"strings"
	"time"
)

const (
	//dbtype = "pgsql"
	dbtype = "sqlite"
)

var (
	Engine *xorm.Engine
)

type User struct {
	Id            int64
	Pid           int64  //用在归属地 归属学校 归属组织 等方面
	Email         string `xorm:"index"`
	Password      string `xorm:"index"`
	Username      string `xorm:"index"`
	Nickname      string `xorm:"index"`
	Realname      string `xorm:"index"`
	Content       string `xorm:"text index"`
	Avatar        string
	AvatarLarge   string
	AvatarMedium  string
	AvatarSmall   string
	Birth         time.Time
	Province      string
	City          string
	Company       string
	Address       string
	Postcode      string
	Mobile        string
	Website       string
	Sex           int64
	Qq            string
	Msn           string
	Weibo         string
	Ctype         int64     `xorm:"index"`
	Role          int64     `xorm:"index"`
	Created       time.Time `xorm:"created index"`
	Updated       time.Time `xorm:"updated index"`
	Hotness       float64   `xorm:"index"`
	Hotup         int64     `xorm:"index"`
	Hotdown       int64     `xorm:"index"`
	Hotscore      int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote       int64     `xorm:"index"` //Hotup  + 	Hotdown
	Views         int64     `xorm:"index"`
	LastLoginTime time.Time
	LastLoginIp   string
	LoginCount    int64 `xorm:"index"`
}

//category,Pid:root
type Category struct {
	Id             int64
	Pid            int64 `xorm:"index"`
	Uid            int64 `xorm:"index"`
	Order          int64
	Ctype          int64     `xorm:"index"`
	Title          string    `xorm:"index"`
	Content        string    `xorm:"text index"`
	Attachment     string    `xorm:"text"`
	Created        time.Time `xorm:"created index"`
	Updated        time.Time `xorm:"updated index"`
	Hotness        float64   `xorm:"index"`
	Hotup          int64     `xorm:"index"`
	Hotdown        int64     `xorm:"index"`
	Hotscore       int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote        int64     `xorm:"index"` //Hotup  + 	Hotdown
	Views          int64     `xorm:"index"`
	Author         string    `xorm:"index"` //这里指本分类创建者
	NodeTime       time.Time
	NodeCount      int64 `xorm:"index"`
	NodeLastUserId int64
}

//node,Pid:category
type Node struct {
	Id              int64
	Pid             int64 `xorm:"index"`
	Uid             int64 `xorm:"index"`
	Order           int64
	Ctype           int64     `xorm:"index"`
	Title           string    `xorm:"index"`
	Content         string    `xorm:"text index"`
	Attachment      string    `xorm:"text"`
	Created         time.Time `xorm:"created index"`
	Updated         time.Time `xorm:"updated index"`
	Hotness         float64   `xorm:"index"`
	Hotup           int64     `xorm:"index"`
	Hotdown         int64     `xorm:"index"`
	Hotscore        int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote         int64     `xorm:"index"` //Hotup  + 	Hotdown
	Views           int64     `xorm:"index"`
	Author          string    `xorm:"index"` //节点的创建者
	TopicTime       time.Time
	TopicCount      int64 `xorm:"index"`
	TopicLastUserId int64
}

//由于cid nid uid都可以是topic的上级所以默认不设置pid字段,这里默认Pid是nid
type Topic struct {
	Id                int64
	Cid               int64 `xorm:"index"`
	Nid               int64 `xorm:"index"`
	Uid               int64 `xorm:"index"`
	Order             int64
	Ctype             int64     `xorm:"index"`
	Title             string    `xorm:"index"`
	Content           string    `xorm:"text index"`
	Attachment        string    `xorm:"text index"`
	Thumbnails        string    `xorm:"index"` //Original remote file
	ThumbnailsLarge   string    `xorm:"index"` //200x300
	ThumbnailsMedium  string    `xorm:"index"` //200x150
	ThumbnailsSmall   string    `xorm:"index"` //70x70
	Tags              string    `xorm:"index"`
	Created           time.Time `xorm:"created index"`
	Updated           time.Time `xorm:"index"`
	Hotness           float64   `xorm:"index"`
	Hotup             int64     `xorm:"index"`
	Hotdown           int64     `xorm:"index"`
	Hotscore          int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote           int64     `xorm:"index"` //Hotup  + 	Hotdown
	Views             int64     `xorm:"index"`
	Author            string    `xorm:"index"`
	Category          string    `xorm:"index"`
	Node              string    `xorm:"index"`
	ReplyTime         time.Time
	ReplyCount        int64 `xorm:"index"`
	ReplyLastUserId   int64
	ReplyLastUsername string
	ReplyLastNickname string
}

//由于cid nid uid都可以是topic的上级所以默认不设置pid字段,这里默认Pid是nid
type Question struct {
	Id                int64
	Cid               int64 `xorm:"index"`
	Nid               int64 `xorm:"index"`
	Uid               int64 `xorm:"index"`
	Order             int64
	Ctype             int64     `xorm:"index"`
	Title             string    `xorm:"index"`
	Content           string    `xorm:"text index"`
	Attachment        string    `xorm:"text index"`
	Thumbnails        string    `xorm:"index"` //Original remote file
	ThumbnailsLarge   string    `xorm:"index"` //200x300
	ThumbnailsMedium  string    `xorm:"index"` //200x150
	ThumbnailsSmall   string    `xorm:"index"` //70x70
	Tags              string    `xorm:"index"`
	Created           time.Time `xorm:"created index"`
	Updated           time.Time `xorm:"index"`
	Hotness           float64   `xorm:"index"`
	Hotup             int64     `xorm:"index"`
	Hotdown           int64     `xorm:"index"`
	Hotscore          int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote           int64     `xorm:"index"` //Hotup  + 	Hotdown
	Views             int64     `xorm:"index"`
	Author            string    `xorm:"index"`
	Category          string    `xorm:"index"`
	Node              string    `xorm:"index"`
	ReplyTime         time.Time
	ReplyCount        int64 `xorm:"index"`
	ReplyHotscore     int64 `xorm:"index"` //Hotup  -	Hotdown
	ReplyLastUserId   int64
	ReplyLastUsername string
	ReplyLastNickname string
	FavoriteCount     int64 `xorm:"index"`
}

//reply,Pid:topic
type Reply struct {
	Id                int64
	Uid               int64 `xorm:"index"`
	Pid               int64 `xorm:"index"` //Topic id
	Order             int64
	Ctype             int64     `xorm:"index"`
	Content           string    `xorm:"text index"`
	Attachment        string    `xorm:"text"`
	Created           time.Time `xorm:"created index"`
	Updated           time.Time `xorm:"index"`
	Hotness           float64   `xorm:"index"`
	Hotup             int64     `xorm:"index"`
	Hotdown           int64     `xorm:"index"`
	Hotscore          int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote           int64     `xorm:"index"` //Hotup  + Hotdown
	Views             int64     `xorm:"index"`
	Author            string    `xorm:"index"`
	AuthorSignature   string    `xorm:"index"`
	Email             string    `xorm:"index"`
	Website           string    `xorm:"index"`
	ReplyTime         time.Time
	ReplyCount        int64 `xorm:"index"`
	ReplyLastUserId   int64
	ReplyLastUsername string
	ReplyLastNickname string
}

type File struct {
	Id              int64
	Uid             int64 `xorm:"index"`
	Cid             int64 `xorm:"index"`
	Nid             int64 `xorm:"index"`
	Pid             int64 `xorm:"index"`
	Order           int64
	Ctype           int64 `xorm:"index"`
	Filename        string
	Content         string `xorm:"text index"`
	Hash            string
	Location        string `xorm:"index"`
	Url             string `xorm:"index"`
	Size            int64
	Created         time.Time `xorm:"created index"`
	Updated         time.Time `xorm:"updated index"`
	Hotness         float64   `xorm:"index"`
	Hotup           int64     `xorm:"index"`
	Hotdown         int64     `xorm:"index"`
	Hotscore        int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote         int64     `xorm:"index"` //Hotup  + 	Hotdown
	Views           int64     `xorm:"index"`
	ReplyTime       time.Time
	ReplyCount      int64 `xorm:"index"`
	ReplyLastUserId int64
}

type Image struct {
	Id              int64
	Uid             int64 `xorm:"index"`
	Cid             int64 `xorm:"index"`
	Nid             int64 `xorm:"index"`
	Pid             int64 `xorm:"index"`
	Order           int64
	Ctype           int64  `xorm:"index"`
	Fingerprint     string `xorm:"index"`
	Filename        string `xorm:"index"`
	Content         string `xorm:"text index"`
	Hash            string `xorm:"index"`
	Location        string `xorm:"index"`
	Url             string `xorm:"index"`
	Size            int64
	Width           int64
	Height          int64
	Created         time.Time `xorm:"created index"`
	Updated         time.Time `xorm:"updated index"`
	Hotness         float64   `xorm:"index"`
	Hotup           int64     `xorm:"index"`
	Hotdown         int64     `xorm:"index"`
	Hotscore        int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote         int64     `xorm:"index"` //Hotup  + 	Hotdown
	Views           int64     `xorm:"index"`
	Author          string    `xorm:"index"`
	ReplyTime       time.Time
	ReplyCount      int64 `xorm:"index"`
	ReplyLastUserId int64
}

type Timeline struct {
	Id                int64
	Cid               int64 `xorm:"index"`
	Nid               int64 `xorm:"index"`
	Uid               int64 `xorm:"index"`
	Order             int64
	Ctype             int64     `xorm:"index"`
	Title             string    `xorm:"index"`
	Content           string    `xorm:"text index"`
	Attachment        string    `xorm:"text index"`
	Created           time.Time `xorm:"created index"`
	Updated           time.Time `xorm:"updated index"`
	Hotness           float64   `xorm:"index"`
	Hotup             int64     `xorm:"index"`
	Hotdown           int64     `xorm:"index"`
	Hotscore          int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote           int64     `xorm:"index"` //Hotup  + 	Hotdown
	Views             int64     `xorm:"index"`
	Author            string    `xorm:"index"`
	AuthorSignature   string    `xorm:"index"`
	Category          string    `xorm:"index"`
	Node              string    `xorm:"index"`
	ReplyTime         time.Time
	ReplyCount        int64 `xorm:"index"`
	ReplyLastUserId   int64
	ReplyLastUsername string
	ReplyLastNickname string
}

type Ads struct {
	Id                int64
	Cid               int64 `xorm:"index"`
	Nid               int64 `xorm:"index"`
	Uid               int64 `xorm:"index"`
	Order             int64
	Ctype             int64  `xorm:"index"`
	Title             string `xorm:"index"`
	Content           string `xorm:"text index"`
	Attachment        string `xorm:"text"`
	Begintime         time.Time
	Endtime           time.Time
	Created           time.Time `xorm:"created index"`
	Updated           time.Time `xorm:"updated index"`
	Hotness           float64   `xorm:"index"`
	Hotup             int64     `xorm:"index"`
	Hotdown           int64     `xorm:"index"`
	Hotscore          int64     `xorm:"index"` //Hotup  -	Hotdown
	Hotvote           int64     `xorm:"index"` //Hotup  + 	Hotdown
	Views             int64     `xorm:"index"`
	Author            string    `xorm:"index"`
	Category          string    `xorm:"index"`
	Node              string    `xorm:"index"`
	ReplyTime         time.Time
	ReplyCount        int64 `xorm:"index"`
	ReplyLastUserId   int64
	ReplyLastUsername string
	ReplyLastNickname string
}

type QuestionMark struct {
	Id  int64
	Uid int64 `xorm:"index"`
	Qid int64 `xorm:"index"`
}

type AnswerMark struct {
	Id  int64
	Uid int64 `xorm:"index"`
	Aid int64 `xorm:"index"`
}

func init() {
	SetEngine()
	CreatTables()
	initData()
}

func ConDb() (*xorm.Engine, error) {
	switch {
	case dbtype == "sqlite":
		return xorm.NewEngine("sqlite3", "./data/sqlite.db")

	case dbtype == "mysql":
		return xorm.NewEngine("mysql", "user=mysql password=^%*&^%#@* dbname=mysql")

	case dbtype == "pgsql":
		// "user=postgres password=^%*&^%#@* dbname=pgsql sslmode=disable maxcons=10 persist=true"
		//return xorm.NewEngine("postgres", "host=110.76.39.205 user=postgres password=^%*&^%#@* dbname=pgsql sslmode=disable")
		return xorm.NewEngine("postgres", "user=postgres password=^%*&^%#@* dbname=sdc sslmode=disable")
		//return xorm.NewEngine("postgres", "host=127.0.0.1 port=6432 user=postgres password=^%*&^%#@* dbname=sdc sslmode=disable")
	}
	return nil, errors.New("尚未设定数据库连接")
}

func SetEngine() (*xorm.Engine, error) {

	var err error
	Engine, err = ConDb()
	//Engine.Mapper = xorm.SameMapper{}

	Engine.ShowDebug = true
	Engine.ShowErr = true
	Engine.ShowSQL = true
	//Engine.SetMaxConns(10)
	//Engine.Pool.SetMaxConns(5)
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 10000)
	Engine.SetDefaultCacher(cacher)

	f, err := os.Create("./logs/sql.log")
	if err != nil {
		println(err.Error())
		panic("OMG!")
	}
	Engine.Logger = f

	return Engine, err
}

func CreatTables() error {

	return Engine.Sync(&User{}, &Category{}, &Node{}, &Topic{}, &Reply{}, &Question{}, &QuestionMark{}, &AnswerMark{}, &File{}, &Image{}, &Timeline{}) //Engine.CreateTables(&User{}, &Category{}, &Node{}, &Topic{}, &Reply{}, &Kv{}, &File{}, &Image{}, &Timeline{})
}

func initData() {
	//用户等级划分：正数是普通用户，负数是管理员各种等级划分，为0则尚未注册
	if usr, err := GetUserByRole(-1000); usr == nil && err == nil {
		if id, err := AddUser("insion@sudochina.com", "root", "root", "root", helper.Encrypt_hash("rootpass", nil), -1000); err == nil {
			fmt.Println("Default Email:insion@sudochina.com ,Username:root ,Password:rootpass Userid:", id)
		} else {
			fmt.Print("create root got errors:", err)
		}

	}
	fmt.Println("The sdc system has started!")
}

func GetUserByRole(role int64) (*User, error) {
	user := new(User)
	if has, err := Engine.Where("role=?", role).Get(user); has {
		return user, err
	} else {
		return nil, err
	}
}

func GetUsersOnHotness(offset int, limit int, path string) (*[]User, error) {
	users := new([]User)
	err := Engine.Limit(limit, offset).Desc(path).Desc("hotness").Find(users)
	return users, err

}

func GetUserByUsername(username string) (*User, error) {
	user := new(User)
	if has, err := Engine.Where("username=?", username).Get(user); has {
		return user, err
	} else {
		return nil, err
	}
}

func GetUserByNickname(nickname string) (*User, error) {

	user := &User{Nickname: nickname}
	has, err := Engine.Get(user)
	if has {
		return user, err
	} else {

		return nil, err
	}
}

//返回值尽量返回指针 不然会出现诡异的问题
func GetUserByEmail(email string) (*User, error) {

	var user = &User{Email: email}
	has, err := Engine.Get(user)
	if has {
		return user, err
	} else {

		return nil, err
	}
}

func GetUser(id int64) (*User, error) {

	user := new(User)
	has, err := Engine.Id(id).Get(user)
	if has {
		return user, err
	} else {

		return nil, err
	}
}

func SearchTopic(content string, offset int, limit int, path string) (*[]Topic, error) {
	//排序首先是热值优先，然后是时间优先。
	if content != "" {

		keyword := "%" + content + "%"

		tps := new([]Topic)

		err := Engine.Where("title like ? or content like ?", keyword, keyword).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
		return tps, err
	}
	return nil, errors.New("搜索内容为空!")
}

func SearchQuestion(content string, offset int, limit int, path string) (*[]Question, error) {
	//排序首先是热值优先，然后是时间优先。
	if content != "" {

		keyword := "%" + content + "%"

		qst := new([]Question)

		err := Engine.Where("title like ? or content like ? or tags like ?", keyword, keyword, keyword).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(qst)
		return qst, err
	}
	return nil, errors.New("搜索内容为空!")
}

func SearchNode(content string, offset int, limit int, path string) (*[]Node, error) {
	//排序首先是热值优先，然后是时间优先。
	if content != "" {

		keyword := "%" + content + "%"

		nds := new([]Node)

		err := Engine.Where("title like ? or content like ?", keyword, keyword).Limit(limit, offset).Desc(path, "views", "topic_count", "created").Find(nds)
		return nds, err
	}
	return nil, errors.New("搜索内容为空!")
}

func SearchCategory(content string, offset int, limit int, path string) (*[]Category, error) {
	//排序首先是热值优先，然后是时间优先。
	if content != "" {

		keyword := "%" + content + "%"

		cats := new([]Category)

		err := Engine.Where("title like ? or content like ?", keyword, keyword).Limit(limit, offset).Desc(path, "views", "node_count", "created").Find(cats)
		return cats, err
	}
	return nil, errors.New("搜索内容为空!")
}

func AddUser(email string, username string, nickname string, realname string, password string, role int64) (int64, error) {

	usr := new(User)
	usr = &User{Email: email, Password: password, Username: username, Nickname: nickname, Realname: realname, Role: role, Created: time.Now()}

	if _, err := Engine.Insert(usr); err == nil {

		return usr.Id, err
	} else {

		return -1, err
	}

}

func AddNode(title string, content string, cid int64, uid int64) (int64, error) {
	nd := new(Node)
	nd = &Node{Pid: cid, Uid: uid, Title: title, Content: content, Created: time.Now()}
	if _, err := Engine.Insert(nd); err == nil {

		return nd.Id, err
	} else {

		return -1, err
	}

}

func AddTopic(title string, content string, cid int64, nid int64, uid int64) (int64, error) {
	cat, ea := GetCategory(cid)
	if ea != nil {
		return -1, ea
	}
	nd, eb := GetNode(nid)
	if eb != nil {
		return -1, eb
	}
	usr, ec := GetUser(uid)
	if ec != nil {
		return -1, ec
	}
	if ea == nil && eb == nil && ec == nil {
		tp := new(Topic)
		tp = &Topic{Cid: cid, Nid: nid, Uid: uid, Title: title, Content: content, Author: usr.Username,
			Category: cat.Title, Node: nd.Title, Created: time.Now()}

		if _, err := Engine.Insert(tp); err == nil {

			return tp.Id, err
		} else {

			return -1, err
		}
	}
	return -1, errors.New("AddTopic发生错误!")
}

func AddTimeline(title string, content string, cid int64, nid int64, uid int64, author string, author_signature string) (int64, error) {
	tml := new(Timeline)
	tml = &Timeline{Cid: cid, Nid: nid, Uid: uid, Title: title, Content: content, Author: author, AuthorSignature: author_signature, Created: time.Now()}

	if _, err := Engine.Insert(tml); err == nil {

		return tml.Id, err
	} else {

		return -1, err
	}
}

func DelTimeline(lid int64) error {
	if row, err := Engine.Id(lid).Delete(new(Timeline)); err != nil || row == 0 {
		fmt.Println("row:::", row)
		fmt.Println("删除时光记录错误:", err)  //此时居然为空!
		return errors.New("删除时光记录错误!") //错误还要我自己构造?!
	} else {
		return nil
	}

}

func GetTimeline(lid int64) (*Timeline, error) {
	tl := new(Timeline)
	_, err := Engine.Where("id=?", lid).Get(tl)

	return tl, err
}

func GetTimelines(offset int, limit int, path string, uid int64) (*[]Timeline, error) {
	tls := new([]Timeline)
	err := errors.New("")
	if uid == 0 {
		err = Engine.Limit(limit, offset).Desc(path).Find(tls)
	} else {
		if err = Engine.Where("uid=?", uid).Limit(limit, offset).Desc(path).Find(tls); err != nil {
			err = Engine.Where("uid=?", uid).NoCache().Limit(limit, offset).Desc(path).Find(tls)
		}
	}
	return tls, err
}

func GetTopicsByHotnessNodes(nodelimit int, topiclimit int) []*[]Topic {
	//找出最热的节点:views优先 然后按 hotness排序 大概找出5到10个节点
	//按上面找的节点读取下级话题
	nds, _ := GetNodes(0, nodelimit, "hotness")
	topics := make([]*[]Topic, 0)

	if len(*nds) > 0 {
		i := 0
		for _, v := range *nds {
			i = i + 1
			tps := GetTopicsByNid(v.Id, 0, topiclimit, 0, "views")
			if len(*tps) != 0 {
				topics = append(topics, tps)
			}
			if i == len(*nds)-1 {
				break
			}
		}
	}

	return topics

}

func GetTopicsByScoreNodes(nodelimit int, topiclimit int) []*[]Topic {
	//找出最热的节点:views优先 然后按 hotness排序 大概找出5到10个节点
	//按上面找的节点读取下级话题
	nds, _ := GetNodes(0, nodelimit, "hotscore")
	topics := make([]*[]Topic, 0)

	if len(*nds) > 0 {
		i := 0
		for _, v := range *nds {
			i = i + 1
			tps := GetTopicsByNid(v.Id, 0, topiclimit, 0, "views")
			if len(*tps) != 0 {
				topics = append(topics, tps)
			}
			if i == len(*nds) {
				break
			}
		}
	}

	return topics

}

func GetTopicsByHotnessCategory(catlimit int, topiclimit int) []*[]Topic {
	//找出最热的分类:views优先 然后按 hotness排序 大概找出5到10个节点
	//按上面找的节点读取下级话题
	cats, _ := GetCategorys(0, catlimit, "hotness")
	topics := make([]*[]Topic, 0)

	if len(cats) > 0 {
		i := 0
		for _, v := range cats {
			i = i + 1
			//(cid int64, offset int, limit int, ctype int64, path string)
			tps := GetTopicsByCid(v.Id, 0, topiclimit, 0, "views")
			if len(*tps) != 0 {
				topics = append(topics, tps)
			}
			if i == len(cats) {
				break
			}
		}
	}

	return topics

}

func GetImagesByHotnessFingerprint(usrlimit int, imagelimit int) []*[]Image {

	//获取最热的用户排名
	users, _ := GetUsersOnHotness(0, usrlimit, "views")

	images := make([]*[]Image, 0)
	if len(*users) > 0 {
		i := 0
		for _, v := range *users {
			i = i + 1
			imgs, _ := GetImagesOnViewsHotnessFingerprintByUidExcludeCid(v.Id, 10) //按用户id取图  并排除头像
			if len(*imgs) != 0 {
				images = append(images, imgs)
			}
			if i == len(*users) {
				break
			}
		}
	}

	return images

}

func SetTimelineContentByRid(lid int64, Content string) error {
	if row, err := Engine.Table(&Timeline{}).Where("id=?", lid).Update(map[string]interface{}{"content": Content}); err != nil || row == 0 {
		fmt.Println("SetTimelineContentByRid  row:::", row)
		fmt.Println("SetTimelineContentByRid出现错误:", err)
		return err
	} else {
		return nil
	}

}

func GetNode(id int64) (*Node, error) {

	nd := new(Node)
	has, err := Engine.Id(id).Get(nd)
	if has {
		return nd, err
	} else {

		return nil, err
	}

}

func GetNodeByTitle(title string) (*Node, error) {
	nd := new(Node)
	nd.Title = title
	has, err := Engine.Get(nd)
	if has {
		return nd, err
	} else {

		return nil, err
	}
}

func GetNodes(offset int, limit int, path string) (*[]Node, error) {
	nds := new([]Node)
	err := Engine.Limit(limit, offset).Desc(path).Find(nds)
	return nds, err
}

func AddCategory(title string, content string) (int64, error) {
	cat := new(Category)
	cat = &Category{Title: title, Content: content, Created: time.Now()}

	if _, err := Engine.Insert(cat); err == nil {

		return cat.Id, err
	} else {
		return -1, err
	}

}

func SetQuestionMark(uid int64, qid int64) (int64, error) {
	qstm := new(QuestionMark)
	qstm = &QuestionMark{Uid: uid, Qid: qid}
	rows, err := Engine.Insert(qstm)
	return rows, err
}

func SetAnswerMark(uid int64, aid int64) (int64, error) {

	ansm := new(AnswerMark)
	ansm = &AnswerMark{Uid: uid, Aid: aid}
	rows, err := Engine.Insert(ansm)
	return rows, err
}

func IsQuestionMark(uid int64, qid int64) bool {

	qstm := new(QuestionMark)

	if has, err := Engine.Where("uid=? and qid=?", uid, qid).Get(qstm); err != nil {
		fmt.Println(err)
		return false
	} else {
		if has {
			if qstm.Uid == uid {
				return true
			} else {
				return false
			}

		} else {
			return false
		}
	}

}

func IsAnswerMark(uid int64, aid int64) bool {

	ansm := new(AnswerMark)

	if has, err := Engine.Where("uid=? and aid=?", uid, aid).Get(ansm); err != nil {
		return false
	} else {
		if has {
			if ansm.Uid == uid {
				return true
			} else {
				return false
			}

		} else {
			return false
		}
	}

}

func GetCategorys(offset int, limit int, path string) (cate []*Category, err error) {
	err = Engine.Limit(limit, offset).Desc(path).Find(&cate)
	return cate, err
}

func GetCategory(id int64) (*Category, error) {

	cat := new(Category)
	has, err := Engine.Id(id).Get(cat)

	if has {
		return cat, err
	} else {

		return nil, err
	}
}

func GetTopic(id int64) (*Topic, error) {

	tp := new(Topic)

	has, err := Engine.Id(id).Get(tp)
	if has {
		return tp, err
	} else {

		return nil, err
	}
}

func GetQuestion(id int64) (*Question, error) {

	qs := new(Question)
	has, err := Engine.Id(id).Get(qs)
	if has {
		return qs, err
	} else {

		return nil, err
	}
}

func GetAnswer(id int64) (*Reply, error) {

	ans := new(Reply)
	has, err := Engine.Id(id).Get(ans)
	if has {
		return ans, err
	} else {

		return nil, err
	}
}

func GetTopics(offset int, limit int, path string) (*[]Topic, error) {
	tps := new([]Topic)
	err := Engine.Limit(limit, offset).Desc(path).Find(tps)
	return tps, err
}

func GetQuestions(offset int, limit int, path string) (*[]Question, error) {

	tps := new([]Question)
	if path == "unanswered" {

		err := Engine.Where("reply_count=? or ctype=?", 0, 0).Limit(limit, offset).Desc("id").Find(tps)
		return tps, err
	} else {

		err := Engine.Limit(limit, offset).Desc(path).Find(tps)
		return tps, err
	}
}

func GetTopicsCount(offset int, limit int, path string) (int64, error) {
	total, err := Engine.Limit(limit, offset).Count(&Topic{})
	return total, err
}

func GetQuestionsCount(offset int, limit int, path string) (int64, error) {

	if path == "unanswered" {
		total, err := Engine.Where("reply_count=? or ctype=?", 0, 0).Limit(limit, offset).Count(&Question{})
		return total, err
	} else {
		total, err := Engine.Limit(limit, offset).Count(&Question{})
		return total, err

	}
}

func GetTopicsByCategoryCount(category string, offset int, limit int, path string) (int64, error) {
	total, err := Engine.Where("category=?", category).Limit(limit, offset).Count(&Topic{})
	return total, err
}

//大数据下会出现极其严重的性能问题 亟待改善
func GetTopicsByCid(cid int64, offset int, limit int, ctype int64, path string) *[]Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]Topic)
	switch {
	case path == "asc":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Limit(limit, offset).Asc("id").Find(tps)
		} else {
			Engine.Where("cid=?", cid).Limit(limit, offset).Asc("id").Find(tps)

		}
	case path == "views" || path == "reply_count":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Desc(path).Limit(limit, offset).Find(tps)

		} else {
			if cid == 0 {
				Engine.Desc(path).Limit(limit, offset).Find(tps)
			} else {
				Engine.Where("cid=?", cid).Desc(path).Limit(limit, offset).Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
		} else {
			if cid == 0 {
				Engine.Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("cid=?", cid).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			}
		}

	}
	return tps
}

func GetQuestionsByCid(cid int64, offset int, limit int, ctype int64, path string) *[]Question {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]Question)
	switch {
	case path == "asc":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Limit(limit, offset).Asc("id").Find(tps)
		} else {
			Engine.Where("cid=?", cid).Limit(limit, offset).Asc("id").Find(tps)

		}
	case path == "views" || path == "reply_count":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Desc(path).Limit(limit, offset).Find(tps)

		} else {
			if cid == 0 {
				Engine.Desc(path).Limit(limit, offset).Find(tps)
			} else {
				Engine.Where("cid=?", cid).Desc(path).Limit(limit, offset).Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Where("cid=? and ctype=?", cid, ctype).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
		} else {
			if cid == 0 {
				Engine.Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("cid=?", cid).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			}
		}

	}
	return tps
}

func GetTopicsByCidOnBetween(cid int64, startid int64, endid int64, offset int, limit int, ctype int64, path string) (tps []*Topic) {
	switch {
	case path == "asc":
		if ctype != 0 {
			Engine.Where("cid=? and ctype=? and id>? and id<?", cid, ctype, startid-1, endid+1).Limit(limit, offset).Asc("id").Find(&tps)
		} else {
			if cid == 0 {
				Engine.Where("id>? and id<?", startid-1, endid+1).Limit(limit, offset).Asc("id").Find(&tps)
			} else {
				Engine.Where("cid=? and id>? and id<?", cid, startid-1, endid+1).Limit(limit, offset).Asc("id").Find(&tps)
			}
		}
	default: //Desc
		if ctype != 0 {
			Engine.Where("cid=? and ctype=? and id>? and id<?", cid, ctype, startid-1, endid+1).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(&tps)
		} else {
			if cid == 0 {
				Engine.Where("id>? and id<?", startid-1, endid+1).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(&tps)
			} else {
				Engine.Where("cid=? and id>? and id<?", cid, startid-1, endid+1).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(&tps)
			}
		}
	}

	return tps
}

func GetTopicsByCategory(category string, offset int, limit int, ctype int64, path string) *[]Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]Topic)
	switch {
	case path == "asc":
		if ctype != 0 {
			Engine.Where("category=? and ctype=?", category, ctype).Limit(limit, offset).Asc("id").Find(tps)
		} else {
			Engine.Where("category=?", category).Limit(limit, offset).Asc("id").Find(tps)

		}
	case path == "views" || path == "reply_count":
		if ctype != 0 {
			Engine.Where("category=? and ctype=?", category, ctype).Desc(path).Limit(limit, offset).Find(tps)

		} else {
			if category == "" {
				Engine.Desc(path).Limit(limit, offset).Find(tps)
			} else {
				Engine.Where("category=?", category).Desc(path).Limit(limit, offset).Find(tps)
			}

		}
	default:
		if ctype != 0 {
			Engine.Where("category=? and ctype=?", category, ctype).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
		} else {
			if category == "" {
				Engine.Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("category=?", category).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			}
		}

	}
	return tps
}

func GetTopicsByUid(uid int64, offset int, limit int, ctype int64, path string) *[]Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]Topic)

	switch {
	case path == "asc":
		if uid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("uid=? and ctype=?", uid, ctype).Limit(limit, offset).Asc("id").Find(tps)
			} else {
				Engine.Where("uid=?", uid).Limit(limit, offset).Asc("id").Find(tps)
			}
		}
	default:
		if uid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("uid=? and ctype=?", uid, ctype).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("uid=?", uid).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			}
		}
	}
	return tps
}

func GetTopicsByNid(nodeid int64, offset int, limit int, ctype int64, path string) *[]Topic {
	//排序首先是热值优先，然后是时间优先。
	tps := new([]Topic)

	switch {
	case path == "asc":
		if nodeid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("nid=? and ctype=?", nodeid, ctype).Limit(limit, offset).Asc("id").Find(tps)
			} else {
				Engine.Where("nid=?", nodeid).Limit(limit, offset).Asc("id").Find(tps)
			}
		}
	default:
		if nodeid == 0 {
			//q.Offset(offset).Limit(limit).OrderByDesc(path).OrderByDesc("views").OrderByDesc("reply_count").OrderByDesc("created").FindAll(&allt)
			return nil
		} else {
			if ctype != 0 {
				Engine.Where("nid=? and ctype=?", nodeid, ctype).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			} else {
				Engine.Where("nid=?", nodeid).Limit(limit, offset).Desc(path, "views", "reply_count", "created").Find(tps)
			}
		}
	}
	return tps
}

func GetTopicsByNode(node string, offset int, limit int, path string) (*[]Topic, error) {
	tps := new([]Topic)
	err := Engine.Where("node=?", node).Limit(limit, offset).Desc(path).Find(tps)
	return tps, err
}

func PostTopic(tp *Topic) (int64, error) {

	if _, err := Engine.Insert(tp); err == nil && tp != nil {
		n, _ := Engine.Where("nid=?", tp.Nid).Count(&Topic{})
		if row, err := Engine.Table(&Node{}).Where("id=?", tp.Nid).Update(map[string]interface{}{"author": tp.Author, "topic_time": time.Now(), "topic_count": n, "topic_last_user_id": tp.Uid}); err != nil || row == 0 {
			fmt.Println("PostTopic更新node表话题相关信息 row:::", row)
			fmt.Println("PostTopic更新node表话题相关信息时出现错误:", err)

		}
		return tp.Id, err
	} else {
		return -1, err
	}

}

func PostQuestion(qs *Question) (*Question, int64, error) {

	rows, err := Engine.Insert(qs)
	return qs, rows, err
}

func PutTimeline(lid int64, tl *Timeline) (int64, error) {
	//覆盖式更新
	row, err := Engine.Update(tl, &Timeline{Id: lid}) //该方法目前返回的row为执行SQL所影响的行数
	return row, err
}

func PutTopic(tid int64, tp *Topic) (int64, error) {
	//覆盖式更新
	row, err := Engine.Update(tp, &Topic{Id: tid}) //该方法目前返回的row为执行SQL所影响的行数
	return row, err
}

func PutQuestion(qid int64, qs *Question) (int64, error) {

	row, err := Engine.Update(qs, &Question{Id: qid}) //该方法目前返回的row为执行SQL所影响的行数
	return row, err
}

func PutReply(rid int64, rp *Reply) (int64, error) {

	row, err := Engine.Update(rp, &Reply{Id: rid}) //该方法目前返回的row为执行SQL所影响的行数
	return row, err
}

func PutAnswer(aid int64, ans *Reply) (int64, error) {

	row, err := Engine.Update(ans, &Reply{Id: aid}) //该方法目前返回的row为执行SQL所影响的行数
	return row, err
}

func PutUser(uid int64, usr *User) (int64, error) {

	row, err := Engine.Update(usr, &User{Id: uid})
	return row, err
}

func PutNode(nid int64, nd *Node) (int64, error) {
	//覆盖式更新
	row, err := Engine.Update(nd, &Node{Id: nid})
	return row, err
}

func PutCategory(cid int64, cat *Category) (int64, error) {
	//覆盖式更新
	row, err := Engine.Update(cat, &Category{Id: cid})
	return row, err

}
func PutImage(mid int64, img *Image) (int64, error) {
	//覆盖式更新
	img.Updated = time.Now()
	row, err := Engine.Update(img, &Image{Id: mid})

	return row, err

}

//map[string]interface{}{"ctype": ctype}
func UpdateCategory(cid int64, catmap *map[string]interface{}) error {
	cat := new(Category)
	if row, err := Engine.Table(cat).Where("id=?", cid).Update(catmap); err != nil || row == 0 {
		fmt.Println("UpdateCategory  row:::", row)
		fmt.Println("UpdateCategory出现错误:", err)
		return err
	} else {
		return nil
	}

}

//map[string]interface{}{"ctype": ctype}
func UpdateNode(nid int64, nodemap *map[string]interface{}) error {
	nd := new(Node)
	if row, err := Engine.Table(nd).Where("id=?", nid).Update(nodemap); err != nil || row == 0 {
		fmt.Println("UpdateNode  row:::", row)
		fmt.Println("UpdateNode出现错误:", err)
		return err
	} else {
		return nil
	}

}

//map[string]interface{}{"ctype": ctype}
func UpdateAnswer(aid int64, answermap map[string]interface{}) error {
	if row, err := Engine.Table(&Reply{}).Where("id=?", aid).Update(answermap); err != nil || row == 0 {
		fmt.Println("UpdateAnswer  row:::", row)
		fmt.Println("UpdateAnswer出现错误:", err)
		return err
	} else {
		return nil
	}

}

func DelTopic(id int64, uid int64, role int64) error {
	allow := false
	if role < 0 {
		allow = true
	}

	topic := new(Topic)

	if has, err := Engine.Id(id).Get(topic); has == true && err == nil {

		if topic.Uid == uid || allow {
			//检查附件字段并尝试删除文件
			if topic.Attachment != "" {

				if p := helper.Url2local(topic.Attachment); helper.Exist(p) {
					//验证是否管理员权限
					if allow {
						if err := os.Remove(p); err != nil {
							//可以输出错误，但不要反回错误，以免陷入死循环无法删掉
							fmt.Println("ROOT DEL TOPIC Attachment, TOPIC ID:", id, ",ERR:", err)
						}
					} else { //检查用户对文件的所有权
						if helper.VerifyUserfile(p, strconv.Itoa(int(uid))) {
							if err := os.Remove(p); err != nil {
								fmt.Println("DEL TOPIC Attachment, TOPIC ID:", id, ",ERR:", err)
							}
						}
					}

				}
			}

			//检查内容字段并尝试删除文件
			if topic.Content != "" {
				//若内容中存在图片则开始尝试删除图片
				delfiles_local := []string{}

				if m, n := helper.GetImages(topic.Content); n > 0 {

					for _, v := range m {
						if helper.IsLocal(v) {
							delfiles_local = append(delfiles_local, v)
							//如果本地同时也存在banner缓存文件,则加入旧图集合中,等待后面一次性删除
							if p := helper.Url2local(helper.SetSuffix(v, "_banner.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
							if p := helper.Url2local(helper.SetSuffix(v, "_large.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
							if p := helper.Url2local(helper.SetSuffix(v, "_medium.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
							if p := helper.Url2local(helper.SetSuffix(v, "_small.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
						}
					}
					for k, v := range delfiles_local {
						if p := helper.Url2local(v); helper.Exist(p) { //如若文件存在,则处理,否则忽略
							//先行判断是否缩略图  如果不是则执行删除image表记录的操作 因为缩略图是没有存到image表记录里面的
							isThumbnails := bool(true) //false代表不是缩略图 true代表是缩略图
							if (!strings.HasSuffix(v, "_large.jpg")) &&
								(!strings.HasSuffix(v, "_medium.jpg")) &&
								(!strings.HasSuffix(v, "_small.jpg")) {
								isThumbnails = false

							}
							//验证是否管理员权限
							if allow {
								if err := os.Remove(p); err != nil {
									fmt.Println("#", k, ",ROOT DEL FILE ERROR:", err)
								}

								//删除image表中已经被删除文件的记录
								if !isThumbnails {
									if e := DelImageByLocation(v); e != nil {
										fmt.Println("DelImageByLocation删除未使用文件", v, "的数据记录时候出现错误:", e)
									}
								}
							} else { //检查用户对文件的所有权
								if helper.VerifyUserfile(p, strconv.Itoa(int(uid))) {
									if err := os.Remove(p); err != nil {
										fmt.Println("#", k, ",DEL FILE ERROR:", err)
									}

									//删除image表中已经被删除文件的记录
									if !isThumbnails {
										if e := DelImageByLocation(v); e != nil {
											fmt.Println("v:", v)
											fmt.Println("DelImageByLocation删除未使用文件", v, "的数据记录时候出现错误:", e)
										}
									}
								}
							}

						}
					}
				}

			}
			//不管实际路径中是否存在文件均删除该数据库记录，以免数据库记录陷入死循环无法删掉
			if topic.Id == id {

				if row, err := Engine.Id(id).Delete(new(Topic)); err != nil || row == 0 {
					fmt.Println("row:::", row)
					fmt.Println("删除话题错误:", err)  //此时居然为空!
					return errors.New("删除话题错误!") //错误还要我自己构造?!
				} else {
					return nil
				}

			}
		}
		return errors.New("你无权删除此话题:" + strconv.Itoa(int(id)))
	}
	return errors.New("无法删除不存在的TOPIC ID:" + strconv.Itoa(int(id)))
}

func DelQuestion(id int64, uid int64, role int64) error {
	allow := bool(false)
	if role < 0 {
		allow = true
	}

	question := new(Question)

	if has, err := Engine.Id(id).Get(question); has == true && err == nil {

		if question.Uid == uid || allow {
			//检查附件字段并尝试删除文件
			if question.Attachment != "" {

				if p := helper.Url2local(question.Attachment); helper.Exist(p) {
					//验证是否管理员权限
					if allow {
						if err := os.Remove(p); err != nil {
							//可以输出错误，但不要反回错误，以免陷入死循环无法删掉
							fmt.Println("ROOT DEL Question Attachment, Question ID:", id, ",ERR:", err)
						}
					} else { //检查用户对文件的所有权
						if helper.VerifyUserfile(p, strconv.Itoa(int(uid))) {
							if err := os.Remove(p); err != nil {
								fmt.Println("DEL Question Attachment, Question ID:", id, ",ERR:", err)
							}
						}
					}

				}
			}

			//检查内容字段并尝试删除文件
			if question.Content != "" {
				//若内容中存在图片则开始尝试删除图片
				delfiles_local := []string{}

				if m, n := helper.GetImages(question.Content); n > 0 {

					for _, v := range m {
						if helper.IsLocal(v) {
							delfiles_local = append(delfiles_local, v)
							//如果本地同时也存在banner缓存文件,则加入旧图集合中,等待后面一次性删除
							if p := helper.Url2local(helper.SetSuffix(v, "_banner.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
							if p := helper.Url2local(helper.SetSuffix(v, "_large.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
							if p := helper.Url2local(helper.SetSuffix(v, "_medium.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
							if p := helper.Url2local(helper.SetSuffix(v, "_small.jpg")); helper.Exist(p) {
								delfiles_local = append(delfiles_local, p)
							}
						}
					}
					for k, v := range delfiles_local {
						if p := helper.Url2local(v); helper.Exist(p) { //如若文件存在,则处理,否则忽略
							//先行判断是否缩略图  如果不是则执行删除image表记录的操作 因为缩略图是没有存到image表记录里面的
							isThumbnails := bool(true) //false代表不是缩略图 true代表是缩略图
							if (!strings.HasSuffix(v, "_large.jpg")) &&
								(!strings.HasSuffix(v, "_medium.jpg")) &&
								(!strings.HasSuffix(v, "_small.jpg")) {
								isThumbnails = false

							}
							//验证是否管理员权限
							if allow {
								if err := os.Remove(p); err != nil {
									fmt.Println("#", k, ",ROOT DEL FILE ERROR:", err)
								}

								//删除image表中已经被删除文件的记录
								if !isThumbnails {
									if e := DelImageByLocation(v); e != nil {
										fmt.Println("DelImageByLocation删除未使用文件", v, "的数据记录时候出现错误:", e)
									}
								}
							} else { //检查用户对文件的所有权
								if helper.VerifyUserfile(p, strconv.Itoa(int(uid))) {
									if err := os.Remove(p); err != nil {
										fmt.Println("#", k, ",DEL FILE ERROR:", err)
									}

									//删除image表中已经被删除文件的记录
									if !isThumbnails {
										if e := DelImageByLocation(v); e != nil {
											fmt.Println("v:", v)
											fmt.Println("DelImageByLocation删除未使用文件", v, "的数据记录时候出现错误:", e)
										}
									}
								}
							}

						}
					}
				}
			}

			//不管实际路径中是否存在文件均删除该数据库记录，以免数据库记录陷入死循环无法删掉
			if question.Id == id {

				if row, err := Engine.Id(id).Delete(new(Question)); err != nil {
					fmt.Println("row:::", row)
					fmt.Println("删除话题错误:", err) //此时居然为空!
				}
			}

			//删除该问题下所有答案
			if err := DelReplysByPid(id); err != nil {
				fmt.Println("执行DelReplysByPid出错:", err)
				return err
			}
			return nil
		} else {
			return errors.New("你无权删除此问题:" + strconv.Itoa(int(id)))
		}
	}

	return errors.New("无法删除不存在的Question ID:" + strconv.Itoa(int(id)))
}

func GetImage(id int64) (*Image, error) {

	img := new(Image)
	has, err := Engine.Id(id).Get(img)

	if has {
		return img, err
	} else {

		return nil, err
	}
}

func GetImagesByCtype(ctype int64) (*[]Image, error) {
	img := new([]Image)
	err := Engine.Where("ctype=?", ctype).Find(img)
	return img, err
}

func GetImagesByCtypeWithUid(ctype int64, uid int64) (*[]Image, error) {
	img := new([]Image)
	err := Engine.Where("ctype=? and uid=?", ctype, uid).Find(img)
	return img, err
}

func GetImagesByCtypeWithUidAndTid(ctype int64, uid int64, tid int64) (*[]Image, error) {
	img := new([]Image)
	err := Engine.Where("ctype=? and uid=? and pid=?", ctype, uid, tid).Find(img)
	return img, err
}

func GetImagesByUidAndTid(uid int64, tid int64) (*[]Image, error) {
	img := new([]Image)
	err := Engine.Where("uid=? and pid=?", uid, tid).Find(img)
	return img, err
}

func GetImagesByCtypeWidthNid(ctype int64, nid int64) (*[]Image, error) {
	img := new([]Image)
	err := Engine.Where("ctype=? and nid=?", ctype, nid).Desc("hotness").Find(img)
	return img, err
}

func GetImagesByCtypeWidthCid(ctype int64, cid int64) (*[]Image, error) {
	img := new([]Image)
	err := Engine.Where("ctype=? and cid=?", ctype, cid).Desc("hotness").Find(img)
	return img, err
}

func GetImagesOnViewsHotnessFingerprintByUidExcludeCid(uid int64, cid int64) (*[]Image, error) {
	img := new([]Image)
	err := Engine.Where("uid=? and ctype<>?", uid, cid).Desc("views", "hotness", "fingerprint").Find(img)
	return img, err
}

func GetImagesOnViewsHotnessFingerprintByUid(uid int64) (*[]Image, error) {
	img := new([]Image)
	err := Engine.Where("uid=?", uid).Desc("views", "hotness", "fingerprint").Find(img)
	return img, err
}

func DelImageByLocation(location string) error {

	if row, err := Engine.Where("location=?", helper.Local2url(location)).Delete(new(Image)); err != nil || row == 0 {
		fmt.Println("row:", row)
		fmt.Println("DelImageBylocation删除话题错误:", err)
		return err
	} else {
		return nil
	}

}

func DelImageByMid(mid int64) error {

	if row, err := Engine.Where("id=?", mid).Delete(new(Image)); err != nil || row == 0 {
		fmt.Println("row:", row)
		fmt.Println("DelImageByMid删除话题错误:", err)
		return err
	} else {
		return nil
	}

}

func SetImageCtypeByMid(mid int64, ctype int64) error {
	if row, err := Engine.Table(&Image{}).Where("id=?", mid).Update(map[string]interface{}{"ctype": ctype}); err != nil || row == 0 {
		fmt.Println("SetImageCtypeByMId  row:::", row)
		fmt.Println("SetImageCtypeByMId出现错误:", err)
		return err
	} else {
		return nil
	}

}

func SetImageCtypeByLocation(location string, ctype int64) error {
	if row, err := Engine.Table(&Image{}).Where("location=?", location).Update(map[string]interface{}{"ctype": ctype}); err != nil || row == 0 {
		fmt.Println("SetImageCtypeBylocation  row:::", row)
		fmt.Println("SetImageCtypeBylocation出现错误:", err)
		return err
	} else {
		return nil
	}

}

func SetImageByLocationWithUid(location string, uid int64, tid int64, ctype int64) error {
	if row, err := Engine.Table(&Image{}).Where("location=? and uid=?", location, uid).Update(map[string]interface{}{"pid": tid, "ctype": ctype}); err != nil || row == 0 {
		fmt.Println("SetImageByLocationWithUid  row:::", row)
		fmt.Println("SetImageByLocationWithUid出现错误:", err)
		return err
	} else {
		return nil
	}

}

func SetImageCtypeByLocationWithUid(location string, uid int64, ctype int64) error {
	if row, err := Engine.Table(&Image{}).Where("location=? and uid=?", location, uid).Update(map[string]interface{}{"ctype": ctype}); err != nil || row == 0 {
		fmt.Println("SetImageCtypeByLocationWithUid  row:::", row)
		fmt.Println("SetImageCtypeByLocationWithUid出现错误:", err)
		return err
	} else {
		return nil
	}

}

func SetImageCtypeByLocationWithUidAndTid(location string, uid int64, tid int64, ctype int64) error {
	if row, err := Engine.Table(&Image{}).Where("location=? and uid=? and pid=?", location, uid, tid).Update(map[string]interface{}{"ctype": ctype}); err != nil || row == 0 {
		fmt.Println("SetImageCtypeByLocationWithUidAndTid  row:::", row)
		fmt.Println("SetImageCtypeByLocationWithUidAndTid出现错误:", err)
		return err
	} else {
		return nil
	}

}

func SetRecordforImageOnPost(tid int64, uid int64) {
	if tc, err := GetTopic(tid); err == nil && tc != nil {

		tpfiles, imgslist, delfiles := []string{}, []string{}, []string{}
		//获取成功发布后话题的本地图片集合:tpfiles
		if imgs, num := helper.GetImages(tc.Content); num > 0 {
			//记录已使用的图片集
			for _, v := range imgs {
				if helper.IsLocal(v) {
					tpfiles = append(tpfiles, v) //已使用的本地图片集合
					//再把已使用的图片到image表中进行对比,把表中ctype改为1标识为已使用
					if e := SetImageByLocationWithUid(v, uid, tid, 1); e != nil {
						fmt.Println("model.SetImageByLocationWithUid出现错误:", e)
					}
				}
			}
		}

		//获取image表中同一用户ctype为-1的图片集合:imgslist
		if imgs2, err := GetImagesByCtypeWithUid(-1, uid); err == nil {
			if len(*imgs2) > 0 {
				for _, v := range *imgs2 {
					imgslist = append(imgslist, v.Location)
				}
			}
		}

		delfiles = helper.DifferenceSets(imgslist, tpfiles) //用标识为-1的图集减去正式使用的图集等于没使用的图集
		for _, v := range delfiles {
			if helper.Exist(helper.Url2local(v)) {
				if e := os.Remove(helper.Url2local(v)); e != nil {
					fmt.Println("删除未使用文件", v, "时候出现错误:", e)
				}
			}
			//如果本地同时也存在banner缓存文件,则一同删除
			if p := helper.Url2local(helper.SetSuffix(v, "_banner.jpg")); helper.Exist(p) {
				if e := os.Remove(p); e != nil {
					fmt.Println("SetRecordforImageOnEdit删除未使用文件之banner", p, "时出现错误:", e)
				}
			}
			//删除image表中已经被删除文件的记录
			if e := DelImageByLocation(v); e != nil {
				fmt.Println("删除未使用文件", v, "的数据记录时候出现错误:", e)
			}

		}
	}
}

func SetRecordforImageOnEdit(tid int64, uid int64) {
	if tc, err := GetTopic(tid); tc != nil && err == nil {

		tpfiles, imgslist, delfiles := []string{}, []string{}, []string{}

		//获取image表中同一用户同一话题的本地图片集合 然后根据记录id对其进行ctype重置为-1
		if imgs, err := GetImagesByUidAndTid(uid, tid); err == nil {
			if len(*imgs) > 0 {
				for _, v := range *imgs {
					SetImageCtypeByMid(v.Id, -1) //重置为-1
				}
			}
		}

		//获取成功更新后话题的已使用本地图片集合:tpfiles
		if imgs, num := helper.GetImages(tc.Content); num > 0 {
			//记录已使用的图片集
			for _, v := range imgs {
				if helper.IsLocal(v) {
					tpfiles = append(tpfiles, v) //已使用的本地图片集合
					//再把已使用的图片到image表中进行对比,把表中ctype改为1标识为已使用
					if e := SetImageByLocationWithUid(v, uid, tid, 1); e != nil {
						fmt.Println("SetImageByLocationWithUid出现错误:", e)
					}
				}
			}
		}

		//获取image表中同一用户所有话题 ctype为-1的图片集合:imgslist
		if imgs2, err := GetImagesByCtypeWithUid(-1, uid); err == nil {
			if len(*imgs2) > 0 {
				for _, v := range *imgs2 {
					imgslist = append(imgslist, v.Location)
				}
			}
		}

		delfiles = helper.DifferenceSets(imgslist, tpfiles) //用标识为-1的图集减去正式使用的图集等于没使用的图集
		for _, v := range delfiles {

			if helper.Exist(helper.Url2local(v)) {

				if e := os.Remove(helper.Url2local(v)); e != nil {
					fmt.Println("SetRecordforImageOnEdit删除未使用文件", v, "的时候出现错误:", e)
				}

			}

			//如果本地同时也存在banner缓存文件,则一同删除
			if p := helper.Url2local(helper.SetSuffix(v, "_banner.jpg")); helper.Exist(p) {
				if e := os.Remove(p); e != nil {
					fmt.Println("SetRecordforImageOnEdit删除未使用文件之banner", p, "时出现错误:", e)
				}
			}

			//删除image表中已经被删除文件的记录
			if e := DelImageByLocation(v); e != nil {
				fmt.Println("DelImageByLocation删除未使用文件", v, "的数据记录时候出现错误:", e)
			}
		}
	}
}

func AddImage(path string, pid int64, ctype int64, uid int64) (int64, error) {
	fg, err := helper.GetImagePha(helper.Url2local(path))
	if err != nil {
		return -1, err
	}

	img := new(Image)
	img.Ctype = ctype
	img.Uid = uid
	img.Created = time.Now()
	img.Location = helper.Local2url(path)
	img.Fingerprint = fg
	img.Pid = pid
	return Engine.Insert(img)

}

func PostImage(img *Image) (int64, error) {
	img.Created = time.Now()
	if _, err := Engine.Insert(img); err == nil {

		return img.Id, err
	} else {

		return -1, err
	}
}

func GetAllReply() (allr *[]Reply) {
	Engine.Desc("id").Find(&allr)
	return allr
}

func GetReply(rid int64) (reply *Reply) {
	Engine.Where("id=?", rid).Get(&reply)
	return reply
}

func GetReplysByPid(id int64, ctype int64, offset int, limit int, path string) *[]Reply {
	rp := new([]Reply)

	//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
	//ctype为10 则是image的回应

	if id == 0 {
		Engine.Where("ctype=?", ctype).Limit(limit, offset).Desc(path).Find(rp)
	} else {

		if ctype == 0 {
			Engine.Where("(ctype=-1 or ctype=1) and pid=?", id).Limit(limit, offset).Desc(path).Find(rp)

		} else {

			Engine.Where("ctype=? and pid=?", ctype, id).Limit(limit, offset).Desc(path).Find(rp)
		}
	}
	return rp
}

func GetAnswersByPid(id int64, ctype int64, offset int, limit int, path string) *[]Reply {
	ans := new([]Reply)

	//ctype等于-1为游客  ctype等于1为正常会员 这里ctype等于0的情况则返回两者
	//ctype为10 则是image的回应

	if id == 0 {
		Engine.Where("ctype=?", ctype).Limit(limit, offset).Desc(path).Find(ans)
	} else {

		if ctype == 0 {
			Engine.Where("(ctype<>0) and pid=?", id).Limit(limit, offset).Desc(path).Find(ans)

		} else {

			Engine.Where("ctype=? and pid=?", ctype, id).Limit(limit, offset).Desc(path).Find(ans)
		}
	}
	return ans
}

func SetReplyCountByPid(qid int64) (int64, error) {
	n, _ := Engine.Where("pid=?", qid).Count(&Reply{Pid: qid})

	qs := new(Question)
	qs.ReplyCount = n
	affected, err := Engine.Where("id=?", qid).Cols("reply_count").Update(qs)
	return affected, err
}

func SetAcceptAnswer(qid int64, aid int64, uid int64, role int64) error {
	allow := bool(false)

	if ans, err := GetAnswer(aid); err == nil && ans != nil {

		if qs, err := GetQuestion(qid); err == nil && qs != nil {
			if (qs.Uid == uid) && (qs.Id == ans.Pid) {
				allow = true
			} else if role < 0 {
				allow = true
			}

			if allow {
				ans.Ctype = 2
				if _, err := PutAnswer(aid, ans); err != nil {
					fmt.Println("PutAnswer执行错误:", err)
					return err
				}
				qs.Ctype = 2
				if _, err := PutQuestion(qid, qs); err != nil {
					fmt.Println("PutQuestion执行错误:", err)
					return err
				} else {
					return err
				}
			} else {
				return errors.New("你没有权限采纳答案!")
			}

		} else {
			return err
		}
	} else {
		return err
	}
}

func SetIgnoreAnswer(qid int64, aid int64, uid int64, role int64) error {
	allow := bool(false)

	if anz, err := GetAnswer(aid); err == nil && anz != nil {

		if qs, err := GetQuestion(qid); err == nil && qs != nil {
			if (qs.Uid == uid) && (qs.Id == anz.Pid) {
				allow = true
			} else if role < 0 {
				allow = true
			}

			if allow {
				anz.Ctype = 1
				if _, err := PutAnswer(aid, anz); err != nil {
					fmt.Println("PutAnswer执行错误:", err)
					return err
				} else {
					//ctype==2 是已经采纳的答案,这里判断当前忽略操作之后是否仍有已采纳的答案,
					//如果有就不要修改问题的ctype,如果没有就把问题的ctype改为1,即表示当前问题尚未解决

					//查找剩余的已解决的答案集合
					if ans := GetAnswersByPid(qid, 2, 0, 0, "id"); ans != nil {
						if len(*ans) <= 0 {
							qs.Ctype = 1
						}
					}

					if qs.Ctype == 1 {

						if _, err := PutQuestion(qid, qs); err != nil {
							fmt.Println("PutQuestion执行错误:", err)
							return err
						} else {
							return err
						}
					}
				}
			} else {
				return errors.New("你没有权限忽略答案!")
			}

		} else {
			return err
		}
	} else {
		return err
	}
	return nil
}

func SetCtypeforQuestion(qid int64, uid int64, role int64, ctype int64) error {
	allow := bool(false)
	if qs, err := GetQuestion(qid); err == nil && qs != nil {
		if qs.Uid == uid {
			allow = true
		} else if role < 0 {
			allow = true
		}

		if allow {
			qs.Ctype = ctype
			if _, err := PutQuestion(qid, qs); err != nil {
				fmt.Println("PutQuestion执行错误:", err)
				return err
			} else {
				return err
			}
		} else {
			return errors.New("你没有权限执行SetCtypeforQuestion!")
		}

	} else {
		return err
	}
}

func GetReplyCountByPid(tid int64) int64 {
	n, _ := Engine.Where("pid=?", tid).Count(&Reply{Pid: tid})
	return n
}

func GetTopicCountByNid(nid int64) int64 {
	n, _ := Engine.Where("nid=?", nid).Count(&Topic{Nid: nid})
	return n
}

func GetNodeCountByPid(cid int64) int64 {
	n, _ := Engine.Where("pid=?", cid).Count(&Node{Pid: cid})
	return n
}

func GetCategoryCountByPid(pid int64) int64 {
	n, _ := Engine.Where("pid=?", pid).Count(&Category{Pid: pid})
	return n
}

func AddReply(tid int64, uid int64, ctype int64, content string, author string, author_signature string, email string, website string) (int64, error) {
	rp := new(Reply)
	rp.Pid = tid
	rp.Uid = uid
	rp.Ctype = ctype
	rp.Content = content
	rp.Author = author
	rp.AuthorSignature = author_signature
	rp.Email = email
	rp.Website = website
	rp.Created = time.Now()

	if _, err := Engine.Insert(rp); err == nil {
		//更新话题中的回应相关记录
		if row, err := Engine.Table(&Topic{}).Where("id=?", tid).Update(map[string]interface{}{"reply_time": time.Now(), "reply_count": GetReplyCountByPid(tid), "reply_last_user_id": uid}); err != nil || row == 0 {
			fmt.Println("AddReply #", rp.Id)
			fmt.Println("AddReply  row:::", row)
			fmt.Println("AddReply出现错误:", err)
		}
		return rp.Id, err
	}
	return -1, errors.New("AddReply无法插入数据!")
}

func AddAnswer(qid int64, uid int64, ctype int64, content string, author string, author_signature string, email string, website string) (int64, error) {
	ans := new(Reply)
	ans.Pid = qid
	ans.Uid = uid
	ans.Ctype = ctype
	ans.Content = content
	ans.Author = author
	ans.AuthorSignature = author_signature
	ans.Email = email
	ans.Website = website
	ans.Created = time.Now()
	ans.Updated = ans.Created

	if _, err := Engine.Insert(ans); err == nil {
		//更新问题中的回应相关记录
		qs := new(Question)
		qs.ReplyCount = GetReplyCountByPid(qid)
		qs.ReplyTime = ans.Updated
		qs.ReplyLastNickname = author
		qs.ReplyLastUsername = author
		qs.ReplyLastUserId = uid
		if affected, err := Engine.Where("id=?", qid).Cols("reply_count").Update(qs); err != nil || affected == 0 {
			fmt.Println("AddAnswer #", ans.Id)
			fmt.Println("AddAnswer affected:::", affected)
			fmt.Println("AddAnswer出现错误:", err)
		}
		/*
			if row, err := Engine.Table(&Question{}).Where("id=?", qid).Update(map[string]interface{}{"reply_time": time.Now(), "reply_count": GetReplyCountByPid(qid), "reply_last_user_id": uid}); err != nil || row == 0 {
				fmt.Println("AddAnswer #", ans.Id)
				fmt.Println("AddAnswer  row:::", row)
				fmt.Println("AddAnswer出现错误:", err)
			}
		*/
		return ans.Id, err
	}
	return -1, errors.New("AddAnswer无法插入数据!")
}

func SetReplyContentByRid(rid int64, Content string) error {
	if row, err := Engine.Table(&Reply{}).Where("id=?", rid).Update(map[string]interface{}{"content": Content}); err != nil || row == 0 {
		fmt.Println("SetReplyContentByRid  row:::", row)
		fmt.Println("SetReplyContentByRid出现错误:", err)
		return err
	} else {
		return nil
	}

}

func DelReply(rid int64) error {
	if row, err := Engine.Id(rid).Delete(new(Reply)); err != nil || row == 0 {
		fmt.Println("row:::", row)
		fmt.Println("删除回应错误:", err)  //此时居然为空!
		return errors.New("删除回应错误!") //错误还要我自己构造?!
	} else {
		return nil
	}

}

func DelReplysByPid(pid int64) error {
	ans := new([]Reply)
	if err := Engine.Where("pid=?", pid).Find(ans); err == nil && ans != nil {
		for _, v := range *ans {
			if err := DelAnswer(v.Id, v.Uid, -1000); err != nil {
				fmt.Println("DelAnswer:", err)
			}
		}
		return nil
	} else {
		return err
	}

}

func DelAnswer(aid int64, uid int64, role int64) error {
	allow := bool(false)
	if anz, err := GetAnswer(aid); err == nil && anz != nil {
		if anz.Uid == uid {
			allow = true
		} else if role < 0 {
			allow = true
		}
		if allow {
			if row, err := Engine.Id(aid).Delete(new(Reply)); err != nil || row == 0 {
				fmt.Println("row:::", row)
				fmt.Println("删除答案错误:", err)  //此时居然为空!
				return errors.New("删除答案错误!") //错误还要我自己构造?!
			} else {
				return nil
			}
		} else {
			return errors.New("你没有权限删除答案!")
		}

	} else {
		return errors.New("没有办法删除根本不存在的答案!")
	}

}

//最终发布到topic
func AtLinksPostImagesOnTopic(content string) (int64, string, error) {

	output := ""
	imgz, content := helper.AtPagesGetImages(content) //这里的content是微博文字,这提取@url,然后each url for get images
	if len(imgz) > 0 {
		for _, v := range imgz {
			output = output + "<p><img style='width:100%;' src='" + v + "'/></p>"
		}

		layout := "2006/1月2号 3点04分首发 "
		title := time.Now().UTC().Format(layout) + strconv.Itoa(len(imgz)) + "P美图!"
		id, err := AddTopic(title, output, 1, 1, 1)
		return id, content, err
	}

	return -1, content, errors.New("没有获得图片")
}

func GetAScoresByPid(pid int64) int64 {
	ascores := int64(0)
	if ans := GetAnswersByPid(pid, 0, 0, 0, "id"); ans != nil {
		if len(*ans) > 0 {
			for _, v := range *ans {
				ascores = ascores + v.Hotscore
			}
		}
	}
	return ascores
}

/*
func main() {

	//Engine.ShowSQL = true

	var tp = Topic{Title: " haha!"}
	PostTopic(tp)
	for i := 0; i < 3; i++ {
		fmt.Println(GetTopic(int64(i)))
	}
}
*/
