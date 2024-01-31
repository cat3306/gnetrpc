package service

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/smtp"
	"strconv"
	"time"

	"github.com/cat3306/gnetrpc"
	"github.com/cat3306/gnetrpc/example/gameserver/models"
	"github.com/cat3306/gnetrpc/example/gameserver/thirdmodule"
	"github.com/cat3306/gnetrpc/protocol"
	"github.com/cat3306/gnetrpc/rpclog"
	"github.com/google/uuid"
)

type Account struct {
	salt   string
	secret string
}

const (
	CodeTimeOut   = 3 * time.Minute
	GenSignUpCode = "GenSignUpCode"
	Login         = ""
)

func (a *Account) Init(v ...interface{}) gnetrpc.IService {
	return a
}
func (a *Account) Alias() string {
	return ""
}
func (a *Account) HashPwd(pwd string) string {
	h := sha256.New()
	h.Write([]byte(a.salt))
	h.Write([]byte("$"))
	h.Write([]byte(pwd))
	pre := "sha256"
	return fmt.Sprintf("%s:%s", pre, hex.EncodeToString(h.Sum(nil)))
}
func (a *Account) ComparePwd(reqPwd string, hashPwd string) bool {
	s := a.HashPwd(reqPwd)
	return s == hashPwd
}

type LoginReq struct {
	Email string `json:"Email"`
	Pwd   string `json:"Pwd"`
}
type LoginRsp struct {
	UserId string `json:"UserId"`
	Nick   string `json:"Nick"`
}

// {"Email":"1273014435@qq.com","Pwd":"123"}
func (a *Account) Login(ctx *protocol.Context, req *LoginReq, rsp *ApiRsp, tag struct{}) *gnetrpc.CallMode {
	data, err := a.login(req, ctx)
	if err != nil {
		rpclog.Errorf("Login err:%s,req:%+v", err.Error(), req)
		rsp.Err(err.Error())
	} else {
		ctx.Conn.SetProperty("user_id", data.UserId)
		rsp.Ok(data)
	}
	return gnetrpc.CallSelf()
}

func (a *Account) login(req *LoginReq, ctx *protocol.Context) (*LoginRsp, error) {
	if req.Email == "" || req.Pwd == "" {
		return nil, errors.New("invalid args")
	}
	userProfile := models.UserProfile{}
	err := userProfile.FindByEmail(thirdmodule.MysqlDb, req.Email)
	if err != nil {
		return nil, err
	}
	if userProfile.Id == 0 {
		return nil, errors.New("你还没注册，亲亲")
	}
	if !a.ComparePwd(req.Pwd, userProfile.Pwd) {
		return nil, errors.New("密码错误")
	}
	var okInt int64
	okInt, err = thirdmodule.CacheSelect(0).Exists(context.Background(), userProfile.UserId).Result()
	if err != nil {
		return nil, err
	}
	if okInt == 1 {
		return nil, errors.New("已登录")
	}
	//	ctx.DelUserId()
	rsp := &LoginRsp{
		UserId: userProfile.UserId,
		Nick:   userProfile.Nick,
	}
	data := map[string]string{
		"email": userProfile.Email,
		"nick":  userProfile.Nick,
	}
	info, _ := json.Marshal(data)
	err = thirdmodule.CacheSelect(0).Set(context.Background(), userProfile.UserId, string(info), -1).Err()
	if err != nil {
		return nil, err
	}
	ctx.Conn.SetProperty(UserInfoKey, &userProfile)
	return rsp, nil
}
func (a *Account) Logout(ctx *protocol.Context, req *struct{}, rsp *ApiRsp, tag struct{}) *gnetrpc.CallMode {

	uv, _ := ctx.Conn.GetProperty("user_id")
	userId := uv.(string)
	err := a.logout(userId)
	if err != nil {
		rpclog.Errorf("Logout err:%s,req:%+v", err.Error(), req)
		rsp.Err(err.Error())
	} else {
		ctx.Conn.DelProperty("user_id")
		rsp.Ok(nil)
	}
	return gnetrpc.CallSelf()
}
func (a *Account) logout(userId string) error {

	return thirdmodule.CacheSelect(0).Del(context.Background(), userId).Err()
}

type RegisterReq struct {
	Email string `json:"Email"`
	Pwd   string `json:"Pwd"`
	Nick  string `json:"Nick"`
	Code  string `json:"Code"`
}

func (a *Account) Register(ctx *protocol.Context, req *RegisterReq, rsp *ApiRsp, tag struct{}) *gnetrpc.CallMode {
	err := a.register(req)
	if err != nil {
		rpclog.Errorf("Register err:%s,req:%+v", err.Error(), req)
		rsp.Err(err.Error())
	} else {
		rsp.Ok("注册成功")
	}
	return gnetrpc.CallSelf()

}
func (a *Account) register(req *RegisterReq) error {
	stashCode, err := thirdmodule.CacheSelect(0).Get(context.Background(), GenSignUpCode+req.Email).Result()
	err = thirdmodule.IgnoreRedisNil(err)
	if err != nil {
		return err
	}
	if stashCode != req.Code {
		return fmt.Errorf("code wrong !,email:%s", req.Email)
	}
	user := models.UserProfile{}
	err = user.FindByEmail(thirdmodule.MysqlDb, req.Email)
	if err != nil {
		return err
	}
	if user.Id != 0 {
		rpclog.Errorf("已注册,email:%s", req.Email)

		return errors.New("已注册")
	}
	uid := uuid.New()
	user = models.UserProfile{
		Nick:       req.Nick,
		Pwd:        a.HashPwd(req.Pwd),
		Email:      req.Email,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
		UserId:     uid.String(),
	}
	if err = user.Create(thirdmodule.MysqlDb); err != nil {
		rpclog.Errorf(" user.Create failed to,err:%s", err.Error())
		return err
	}
	return nil
}

type EmailCodeReq struct {
	Email string `json:"email"`
}

func (a *Account) EmailCode(ctx *protocol.Context, req *EmailCodeReq, rsp *ApiRsp, tag struct{}) *gnetrpc.CallMode {
	err := a.emailCode(req)
	if err != nil {
		rpclog.Errorf("EmailCode err:%s,req:%+v", err.Error(), req)
		rsp.Err(err.Error())
	} else {
		rsp.Ok(nil)
	}
	return gnetrpc.CallSelf()
}
func GenEmailCode(digit int) string { //生成几位验证码

	AuthCode := ""
	seed := time.Now().Unix()
	rand.Seed(seed)
	for i := 0; i < digit; i++ {

		AuthCode += strconv.Itoa(rand.Intn(10))
	}
	return AuthCode

}
func (a *Account) emailCode(req *EmailCodeReq) error {
	user := models.UserProfile{}
	if err := user.FindByEmail(thirdmodule.MysqlDb, req.Email); err != nil {
		return err
	}
	if user.Id != 0 {
		rpclog.Errorf("email:%s,already register", req.Email)
		return errors.New("邮箱已注册")
	}
	code := GenEmailCode(6)
	err := thirdmodule.CacheSelect(0).Set(context.Background(), GenSignUpCode+req.Email, code, CodeTimeOut).Err()
	if err != nil {
		return err
	}
	en := base64.StdEncoding.EncodeToString([]byte("恶魔射手官方"))
	to := []string{fmt.Sprintf("%v", req.Email)}
	title := "【恶魔射手】"
	content := fmt.Sprintf("亲爱的用户,您的验证码为%v,有效时间%.0f分钟,祝您游戏愉快！", code, CodeTimeOut.Minutes())
	from := fmt.Sprintf("=?UTF-8?B?%s?= <1273014435@qq.com>", en)
	err = EmailSend(to, title, content, from)
	if err != nil {
		return err
	}
	return nil
}

func EmailSend(to []string, title string, context string, from string) error {
	userEmail := "1273014435@qq.com"
	mailSmtpPort := ":587"
	mailPassword := "tfmbksrpxxfvhjig"
	mailHost := "smtp.qq.com"
	auth := smtp.PlainAuth("", userEmail, mailPassword, mailHost)
	for _, v := range to {
		if v != "" {
			header := make(map[string]string)
			header["From"] = from
			header["To"] = v
			header["Subject"] = title
			header["Content-Type"] = "text/html;charset=UTF-8"
			body := context
			to := []string{v}
			messageStr := ""
			for k, v := range header {
				messageStr += fmt.Sprintf("%s: %s\r\n", k, v)
			}
			messageStr += "\r\n" + body
			msg := []byte(messageStr)
			err := smtp.SendMail(mailHost+mailSmtpPort, auth, userEmail, to, msg)
			if err != nil {
				return err
			}
		}
	}

	return nil

}
