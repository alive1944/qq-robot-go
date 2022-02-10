package msg

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/buger/jsonparser"
)

type RecvNormalMsg struct {
	Anonymous   string `json:"anonymous"` // 匿名，群属性
	GroupId     int64  `json:"group_id"`  // 群ID
	Font        int64  `json:"font"`
	Message     string `json:"message"`
	MessageId   int64  `json:"message_id"`
	MessageType string `json:"message_type"`
	PostType    string `json:"PostType"`
	RowMessage  string `json:"row_message"`
	SelfId      int64  `json:"self_id"`
	TargetId    int64  `json:"target_id"` // 发送目标的user_id 私聊属性
	SubType     string `json:"sub_type"`
	Time        int64  `json:"time"`
	UserId      int64  `json:"user_id"`
	Sender      struct {
		Age      int64  `json:"age"`
		Area     string `json:"area"`  // 地区，群属性
		Card     string `json:"card"`  // 卡片？，群属性
		Level    string `json:"level"` // 等级，群属性
		Role     string `json:"admin"` // 角色，群属性
		Nickname string `json:"nickname"`
		Title    string `json:"title"` // 角色title，群属性（名字前面的称谓）
		Sex      string `json:"sex"`
		UserId   int64  `json:"user_id"`
	}
}

func (m *RecvNormalMsg) GetEventType() {

}

// 删除所有CQ码后的纯文本，并且会删除首尾的空格
func (m *RecvNormalMsg) GetPureMsg() string {
	reg := regexp.MustCompile(`\[.*?\]`)
	pu := reg.ReplaceAllString(m.Message, "")
	pu = strings.Trim(pu, " ")
	return pu
}

func (m RecvNormalMsg) IsGroup() bool {
	return m.MessageType == MSG_TYPE_GROUP
}

func (m RecvNormalMsg) IsPrivate() bool {
	return m.MessageType == MSG_TYPE_PRIVATE
}

func (m RecvNormalMsg) IsAtMe() bool {
	atList := CQAtDecode(m.Message)
	for _, qq := range atList {
		if qq == CurLoginQQ {
			return true
		}
	}
	return false
}

func (m RecvNormalMsg) IsAtAll() bool {
	atList := CQAtDecode(m.Message)
	if len(atList) > 0 && atList[0] == 0 {
		return true
	}

	return false
}

const (
	MSG_TYPE_PRIVATE string = "private"
	MSG_TYPE_GROUP   string = "group"
	// ...
)

const (
	SUB_TYPE_NORMAL string = "normal"
	SUB_TYPE_FRIEND string = "friend"
)

func NewRecvMsgObj(recv []byte) *RecvNormalMsg {
	postType, err := jsonparser.GetString(recv, "post_type")
	if err != nil {
		// 获取不到信息类型，直接return掉
		return nil
	}

	if postType == "message" {
		var recvMeg *RecvNormalMsg
		err2 := json.Unmarshal(recv, &recvMeg)
		if err2 != nil {
			return nil
		}

		return recvMeg
	}

	return nil
}
