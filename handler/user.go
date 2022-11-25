package handler

import (
	mydb "filestore-server/db"
	"filestore-server/util"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	// 用于加密的盐值(自定义)
	pwdSalt = "*#890"
)

// SignUpHandler : 注册接口
func SignUpHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := os.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		//http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	if len(username) < 3 || len(password) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	encpassword := util.Sha1([]byte(password + pwdSalt))
	flag := mydb.UserSignup(username, encpassword)
	if flag {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

// SignInHandler :用户登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := os.ReadFile("./static/view/signip.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		//http.Redirect(w, r, "/static/view/signin.html", http.StatusFound)
		return
	}
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encpassword := util.Sha1([]byte(password + pwdSalt))
	//1.检查登陆的密码是否正确
	flag := mydb.UserSignin(username, encpassword)
	if !flag {
		w.Write([]byte("FAILED"))
		return
	}
	//2.生成后续访问需要的token
	token := GenToken(username)
	upRes := mydb.UpdateToken(username, token)
	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}
	//3.登陆成功重定向到首页
	// 3. 登录成功后重定向到首页
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	w.Write(resp.JSONBytes())
}

// GenToken : 生成token
func GenToken(username string) string {
	// 40位字符:md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}

// UserInfoHandler ： 查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1.解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	// token := r.Form.Get("token")
	// //2.检查token是否有效
	// flag := IsTokenValid(token)
	// if !flag {
	// 	w.WriteHeader(http.StatusForbidden)
	// 	return
	// }
	//3.查询用户信息
	user, err := mydb.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	//4.组装并返回用户信息
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

// IsTokenValid : token是否有效
func IsTokenValid(token string) bool {
	// if len(token) != 40 {
	// 	return false
	// }
	// TODO: 判断token的时效性，是否过期
	// TODO: 从数据库表tbl_user_token查询username对应的token信息
	// TODO: 对比两个token是否一致
	return len(token) == 40
}
