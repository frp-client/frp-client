package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/frp-client/frp-client/model"
	"github.com/frp-client/frp-client/utils"
	"log"
)

func (a *App) getTempSessionFile() string {
	return utils.AppTempFile(".session")
}

// 检查用户session数据
func (a *App) checkUserSession() (session model.UserSession, err error) {
	buf, err := utils.AesDecrypt(a.aesKey(), utils.ReadFileAsString(a.getTempSessionFile()))
	if err != nil {
		return session, err
	}
	if err = json.Unmarshal(buf, &session); err != nil {
		return session, err
	}
	return session, nil
}

func (a *App) apiClientLogin() (model.UserSession, error) {
	// 注册返回结构
	var registerResp = struct {
		Code uint   `json:"code"`
		Msg  string `json:"msg"`
		Data any    `json:"data"`
	}{}
	// 登录返回结构
	var loginResp = struct {
		Code uint   `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Id          int64  `json:"id"`
			Username    string `json:"username"`
			AccessToken string `json:"access_token"`
			JwtToken    string `json:"jwt_token"`
			CreatedAt   int64  `json:"created_at"`
		} `json:"data"`
	}{}
	// 登录成功检测返回结构
	var checkResp = struct {
		Code uint   `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			Id       int64  `json:"id"`
			Username string `json:"username"`
			ClientId string `json:"client_id"`
		} `json:"data"`
	}{}

	var clientId = utils.ClientId()
	var ts = utils.UnixTimestamp()
	session, err := a.checkUserSession()
	if err != nil {
		// 模拟注册登录
		_, err = utils.HttpJsonPostUnmarshal(
			utils.FormatUrl(baseURL, "/api/user/register"),
			utils.ToJsonByte(model.Map{
				"username":  clientId,
				"password":  utils.Sha256(fmt.Sprintf("%s,%s", clientId, clientId)),
				"sign":      utils.Sha256(fmt.Sprintf("%s,%d", clientId, ts)),
				"timestamp": ts,
			}), map[string]string{"X_CLIENT_ID": clientId},
			&registerResp,
		)
		if err != nil {
			return session, err
		}
		if registerResp.Code != 200 {
			return session, errors.New(fmt.Sprintf("服务器错误：%s", registerResp.Msg))
		}
		_, err = utils.HttpJsonPostUnmarshal(
			utils.FormatUrl(baseURL, "/api/user/login"),
			utils.ToJsonByte(model.Map{
				"username": clientId,
				"password": utils.Sha256(fmt.Sprintf("%s,%s", clientId, clientId)),
			}),
			map[string]string{"X_CLIENT_ID": clientId},
			&loginResp,
		)
		if err != nil {
			return session, err
		}
		if loginResp.Code != 200 {
			return session, errors.New(fmt.Sprintf("服务器错误：%s", loginResp.Msg))
		}

		session = model.UserSession{
			UserId:      loginResp.Data.Id,
			Username:    loginResp.Data.Username,
			MachineId:   clientId,
			JwtToken:    loginResp.Data.JwtToken,
			AccessToken: loginResp.Data.AccessToken,
			UpdatedAt:   utils.UnixTimestamp(),
		}
	}

	// 登录成功检测
	_, err = utils.HttpJsonGetUnmarshal(
		utils.FormatUrl(baseURL, "/api/user/check"),
		map[string]string{"X_CLIENT_ID": clientId, "Authorization": fmt.Sprintf("Bearer %s", session.JwtToken)},
		&checkResp,
	)
	if err != nil {
		return session, err
	}
	if checkResp.Code != 200 {
		return session, errors.New(fmt.Sprintf("服务器错误：%s", checkResp.Msg))
	}
	log.Println("[checkResp]", utils.ToJsonString(checkResp))
	// 以上为机器码注册用户

	//utils.AesDecrypt
	encrypt, err := utils.AesEncrypt(a.aesKey(), utils.ToJsonString(session))
	if err != nil {
		return session, err
	}

	return session, utils.SaveFileAsString(a.getTempSessionFile(), encrypt)
}

func (a *App) LoginSuccess(data model.LoginResp) bool {
	log.Println("[LoginSuccess] ", utils.ToJsonString(data))

	session, _ := a.checkUserSession()
	session.UserId = data.Id
	session.Username = data.Username
	session.JwtToken = data.JwtToken
	session.AccessToken = data.AccessToken
	session.UpdatedAt = utils.UnixTimestamp()

	encrypt, err := utils.AesEncrypt(a.aesKey(), utils.ToJsonString(session))
	if err != nil {
		return false
	}
	if err := utils.SaveFileAsString(a.getTempSessionFile(), encrypt); err != nil {
		return false
	}
	return true
}

func (a *App) aesKey() []byte {
	return []byte(a.ClientId()[:32])
}
