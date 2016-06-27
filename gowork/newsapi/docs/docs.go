package docs

import (
	"encoding/json"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/swagger"
)

const (
    Rootinfo string = `{"apiVersion":"1.0.0","swaggerVersion":"1.2","apis":[{"path":"/news","description":"Operations about news\n"}],"info":{"title":"beego Test API","description":"beego has a very cool tools to autogenerate documents for your API","contact":"astaxie@gmail.com","termsOfServiceUrl":"http://beego.me/","license":"Url http://www.apache.org/licenses/LICENSE-2.0.html"}}`
    Subapi string = `{"/news":{"apiVersion":"1.0.0","swaggerVersion":"1.2","basePath":"","resourcePath":"/news","produces":["application/json","application/xml","text/plain","text/html"],"apis":[{"path":"/getnews","description":"","operations":[{"httpMethod":"GET","nickname":"getnews","type":"","summary":"get news by gid","parameters":[{"paramType":"query","name":"gid","description":"\"The gid for news\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":false,"minimum":0,"maximum":0},{"paramType":"query","name":"limit","description":"\"The litit for news default 30\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":false,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"news success","responseModel":""},{"code":403,"message":"gid not exist","responseModel":""}]}]},{"path":"/getnewsinfo","description":"","operations":[{"httpMethod":"GET","nickname":"getnewsinfo","type":"","summary":"get newsinfo by type or time","parameters":[{"paramType":"query","name":"type","description":"\"type:1 国际,type:4 原油,type:5 贵金属,type:6 央行,type:7 外汇,type:13 独家\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":false,"minimum":0,"maximum":0},{"paramType":"query","name":"datetime","description":"\"根据时间来查找\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":false,"minimum":0,"maximum":0},{"paramType":"query","name":"limit","description":"\"The litit for news default 30\"","dataType":"string","type":"","format":"","allowMultiple":false,"required":false,"minimum":0,"maximum":0}],"responseMessages":[{"code":200,"message":"news success","responseModel":""},{"code":403,"message":"fail","responseModel":""}]}]}]}}`
    BasePath string= "/v1"
)

var rootapi swagger.ResourceListing
var apilist map[string]*swagger.APIDeclaration

func init() {
	err := json.Unmarshal([]byte(Rootinfo), &rootapi)
	if err != nil {
		beego.Error(err)
	}
	err = json.Unmarshal([]byte(Subapi), &apilist)
	if err != nil {
		beego.Error(err)
	}
	beego.GlobalDocAPI["Root"] = rootapi
	for k, v := range apilist {
		for i, a := range v.APIs {
			a.Path = urlReplace(k + a.Path)
			v.APIs[i] = a
		}
		v.BasePath = BasePath
		beego.GlobalDocAPI[strings.Trim(k, "/")] = v
	}
}


func urlReplace(src string) string {
	pt := strings.Split(src, "/")
	for i, p := range pt {
		if len(p) > 0 {
			if p[0] == ':' {
				pt[i] = "{" + p[1:] + "}"
			} else if p[0] == '?' && p[1] == ':' {
				pt[i] = "{" + p[2:] + "}"
			}
		}
	}
	return strings.Join(pt, "/")
}
