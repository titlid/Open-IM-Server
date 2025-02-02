package callbackstruct

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

type CommonCallbackReq struct {
	SendID           string   `json:"sendID"`
	CallbackCommand  string   `json:"callbackCommand"`
	ServerMsgID      string   `json:"serverMsgID"`
	ClientMsgID      string   `json:"clientMsgID"`
	OperationID      string   `json:"operationID"`
	SenderPlatformID int32    `json:"senderPlatformID"`
	SenderNickname   string   `json:"senderNickname"`
	SessionType      int32    `json:"sessionType"`
	MsgFrom          int32    `json:"msgFrom"`
	ContentType      int32    `json:"contentType"`
	Status           int32    `json:"status"`
	CreateTime       int64    `json:"createTime"`
	Content          string   `json:"content"`
	Seq              uint32   `json:"seq"`
	AtUserIDList     []string `json:"atUserList"`
	SenderFaceURL    string   `json:"faceURL"`
	Ex               string   `json:"ex"`
}

func (c *CommonCallbackReq) GetCallbackCommand() string {
	return c.CallbackCommand
}

type CallbackReq interface {
	GetCallbackCommand() string
}

type CallbackResp interface {
	Parse() (err error)
}

type CommonCallbackResp struct {
	ActionCode int    `json:"actionCode"`
	ErrCode    int32  `json:"errCode"`
	ErrMsg     string `json:"errMsg"`
	ErrDlt     string `json:"errDlt"`
}

func (c CommonCallbackResp) Parse() error {
	if c.ActionCode != errs.NoError || c.ErrCode != errs.NoError {
		return errs.NewCodeError(int(c.ErrCode), c.ErrMsg).WithDetail(c.ErrDlt)
	}
	return nil
}

type UserStatusBaseCallback struct {
	CallbackCommand string `json:"callbackCommand"`
	OperationID     string `json:"operationID"`
	PlatformID      int    `json:"platformID"`
	Platform        string `json:"platform"`
}

func (c UserStatusBaseCallback) GetCallbackCommand() string {
	return c.CallbackCommand
}

type UserStatusCallbackReq struct {
	UserStatusBaseCallback
	UserID string `json:"userID"`
}

type UserStatusBatchCallbackReq struct {
	UserStatusBaseCallback
	UserIDList []string `json:"userIDList"`
}
