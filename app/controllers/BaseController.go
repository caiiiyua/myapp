package controllers

import (
	"encoding/json"
	"fmt"
	"myapp/app/models"
	"myapp/app/models/entity"
	"net/http"
	"net/url"
	"strings"

	"github.com/revel/revel"
)

// Base of all other controllers
type BaseController struct {
	*revel.Controller
	XOrmTnController
}

func (c BaseController) IsLogined(id string) bool {
	user, ok := c.Session["user"]
	userId, ok2 := c.Session["id"]
	return ok && ok2 && user != "" && userId == id
}

func (c BaseController) userService() models.UserService {
	return models.DefaultUserService(c.XOrmSession)
}

func (c BaseController) GetUserId() string {
	if userid, ok := c.Session["UserId"]; ok {
		return userid
	}
	return ""
}

func (c BaseController) GetUsername() string {
	if username, ok := c.Session["Username"]; ok {
		return username
	}
	return ""
}

func (c BaseController) GetEmail() string {
	if email, ok := c.Session["Email"]; ok {
		return email
	}
	return ""
}

func (c BaseController) GetUserVerified() bool {
	if verified, ok := c.Session["Verified"]; ok {
		if verified == "1" {
			return true
		}
	}
	return false
}

func (c BaseController) GetSession(k string) string {
	if v, ok := c.Session[k]; ok {
		return v
	}
	return ""
}

func (c BaseController) SetSession(u entity.User) {
	if u.Id >= 0 {
		c.Session["UserId"] = string(u.Id)
		c.Session["Username"] = u.Username
		c.Session["Email"] = u.Email
		if u.Verified {
			c.Session["Verified"] = "1"
		} else {
			c.Session["Verified"] = "0"
		}
	}
}

func (c BaseController) SetLocale() string {
	locale := string(c.Request.Locale)
	lang := locale
	if strings.Contains(locale, "-") {
		pos := strings.Index(locale, "-")
		lang = locale[:pos]
	}
	if lang != "zh" {
		lang = "en"
	}
	return lang
}

func (c BaseController) checkReg() revel.Result {
	return c.check(Auth.Register)
}

func (c BaseController) checkAuth() revel.Result {
	return c.check(Auth.Login)
}

func (c BaseController) checkLogined(id string) revel.Result {
	if !c.IsLogined(id) {
		c.Validation.Error("Need logined").Key("email")
	}
	return c.check(Auth.Login)
}

func (c BaseController) check(f interface{}) revel.Result {
	if c.Validation.HasErrors() {
		c.Validation.Keep()
		c.FlashParams()
		return c.Redirect(f)
	}
	return nil
}

var WeChatOAuth = struct {
	AppId          string
	Secret         string
	CodeUrl        string
	AccessTokenUrl string
	UserInfoUrl    string
}{
	"wx6212752719ca7a9f",
	"secret", // need protect
	"https://open.weixin.qq.com/connect/oauth2/authorize?appid=APPID&redirect_uri=REDIRECT_URI&response_type=code&scope=snsapi_userinfo&state=STATE#wechat_redirect",
	"https://api.weixin.qq.com/sns/oauth2/access_token?appid=APPID&secret=SECRET&code=CODE&grant_type=authorization_code",
	"https://api.weixin.qq.com/sns/userinfo?access_token=ACCESS_TOKEN&openid=OPENID&lang=zh_CN",
}

func (c BaseController) GetCodeUrl(appId, redirectUri, state string) string {
	u, _ := url.Parse(WeChatOAuth.CodeUrl)
	q := u.Query()
	q.Set("appid", appId)
	q.Set("redirect_uri", "http://inaiping.wang/"+redirectUri)
	q.Set("state", state)
	u.RawQuery = q.Encode()
	fmt.Println(u.String())
	return u.String()
}

func (c BaseController) GetAccessTokenUrl(appId, secret, code, state string) string {
	u, _ := url.Parse(WeChatOAuth.AccessTokenUrl)
	q := u.Query()
	q.Set("appid", appId)
	q.Set("secret", secret)
	q.Set("code", code)
	u.RawQuery = q.Encode()
	fmt.Println(u.String())
	return u.String()
}

func (c BaseController) GetUserInfoUrl(accessToken, openId string) string {
	u, _ := url.Parse(WeChatOAuth.UserInfoUrl)
	q := u.Query()
	q.Set("access_token", accessToken)
	q.Set("openid", openId)
	u.RawQuery = q.Encode()
	fmt.Println(u.String())
	return u.String()
}

func (c BaseController) WeChatGetCode(url string) (code, state string) {
	codeUrl := c.GetCodeUrl(WeChatOAuth.AppId, url, url)
	resp, _ := http.Get(codeUrl)
	defer resp.Body.Close()
	me := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		revel.ERROR.Println(err)
	}
	fmt.Println(me)
	return me["code"].(string), me["state"].(string)
}

func (c BaseController) WeChatGetAccessToken(code, state string) (token, openId string) {
	tokenUrl := c.GetAccessTokenUrl(WeChatOAuth.AppId, WeChatOAuth.Secret, code, state)
	resp, _ := http.Get(tokenUrl)
	defer resp.Body.Close()
	me := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		revel.ERROR.Println(err)
	}
	fmt.Println(me)
	return me["access_token"].(string), me["openid"].(string)
}

func (c BaseController) WeChatGetUserInfo(accessToken, openId string) (nickName, sex, city string) {
	userUrl := c.GetUserInfoUrl(accessToken, openId)
	resp, _ := http.Get(userUrl)
	defer resp.Body.Close()
	me := map[string]interface{}{}
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		revel.ERROR.Println(err)
	}
	fmt.Println(me)
	return me["nickname"].(string), me["sex"].(string), me["city"].(string)
}
