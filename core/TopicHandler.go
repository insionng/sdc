package core

import (
	"encoding/json"
	"fmt"
	"sdc/helper"
	"sdc/lib"
	"sdc/model"
	"strconv"
)

var (
	aesPublicKey  = helper.AesPublicKey
	rsaPublicKey  = helper.RsaPublicKey
	rsaPrivateKey = helper.RsaPrivateKey
)

type TopicHandler struct {
	lib.BaseHandler
}

//创建帖子
func (self *TopicHandler) Post() {

	if hash := self.GetString("hash"); hash != "" {
		//(decrypt bool, hash string, status string, content []byte, aesPublicKey string, rsaPublicKey []byte, rsaPrivateKey []byte)
		if rsa_decrypt_content, err := helper.ReceivingPackets(true, hash, "POST", self.Ctx.Input.RequestBody, aesPublicKey, rsaPublicKey, rsaPrivateKey); err == nil {
			tp := new(model.Topic)
			json.Unmarshal(rsa_decrypt_content, &tp)
			if tid, err := model.PostTopic(tp); err != nil {

				self.Data["json"] = "Post failed!"
			} else {
				self.Data["json"] = `{"TopicId:"` + strconv.Itoa(int(tid)) + `}`
			}

			self.ServeJson()
		} else {

			fmt.Println("401 Unauthorized!")
			self.Abort("401")
		}
	} else {

		fmt.Println("401 Unauthorized!")
		self.Abort("401")
	}
}

//获取帖子
func (self *TopicHandler) Get() {
	tid, _ := self.GetInt(":objectId") //beego api模式下，提交的参数名总是唤作objectId

	if tid > 0 {
		tp, _ := model.GetTopic(tid)

		self.Data["json"] = *tp

	} else {
		tps, _ := model.GetTopics(0, 0, "id")
		self.Data["json"] = *tps
	}
	self.ServeJson()
}

//更新帖子
func (self *TopicHandler) Put() {

	if hash := self.GetString("hash"); hash != "" {

		if rsa_decrypt_content, err := helper.ReceivingPackets(true, hash, "PUT", self.Ctx.Input.RequestBody, aesPublicKey, rsaPublicKey, rsaPrivateKey); err == nil {

			tid, _ := self.GetInt(":objectId")
			tp := new(model.Topic)
			json.Unmarshal(rsa_decrypt_content, &tp)

			if _, err := model.PutTopic(tid, tp); err != nil {
				self.Data["json"] = "Update failed!"
			} else {
				self.Data["json"] = "Update success!"
			}
			self.ServeJson()
		} else {

			fmt.Println("401 Unauthorized!")
			self.Abort("401")
		}
	} else {

		fmt.Println("401 Unauthorized!")
		self.Abort("401")
	}
}

//删除帖子
func (self *TopicHandler) Delete() {
	if hash := self.GetString("hash"); hash != "" {

		if rsa_decrypt_content, err := helper.ReceivingPackets(true, hash, "DELETE", self.Ctx.Input.RequestBody, aesPublicKey, rsaPublicKey, rsaPrivateKey); err == nil {

			tid, _ := self.GetInt(":objectId")
			var tp *model.Topic
			json.Unmarshal(rsa_decrypt_content, &tp)
			if tid == tp.Id && tid > 0 {
				if e := model.DelTopic(tid, 1, -100000); e != nil {
					self.Data["json"] = "Delete failed!"

				} else {

					self.Data["json"] = "Delete success!"
				}
				//self.ServeJson()
				self.Ctx.WriteString(self.Data["json"].(string))

			} else {

				fmt.Println("401 Unauthorized!")
				self.Abort("401")
			}

		} else {

			fmt.Println("401 Unauthorized!")
			self.Abort("401")
		}
	} else {

		fmt.Println("401 Unauthorized!")
		self.Abort("401")
	}
}
