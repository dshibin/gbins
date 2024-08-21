package gcurl

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"gbins/gconf"
	"gbins/glog"
	"strings"
	"time"
)

func GGet(ctx *gin.Context, serverName string, params string, header ...map[string]string) (err error, r []byte) {
	cconf := gconf.GConfByClientName(serverName)
	err, r = request(ctx, http.MethodGet, fmt.Sprintf("%s?%s", cconf.Target, params), "", cconf.Timeout, header)
	return
}

func GPost(ctx *gin.Context, serverName string, params string, body interface{}, header ...map[string]string) (err error, r []byte) {
	cconf := gconf.GConfByClientName(serverName)
	bodyStr, _ := json.Marshal(body)
	err, r = request(ctx, http.MethodPost, fmt.Sprintf("%s?%s", cconf.Target, params), string(bodyStr), cconf.Timeout, append([]map[string]string{{"Content-Type": gin.MIMEJSON}}, header...))
	return
}

func request(ctx *gin.Context, method string, url string, body string, timeout int64, header []map[string]string) (err error, rsp []byte) {
	st := time.Now()
	req, _ := http.NewRequestWithContext(ctx, method, url, nil)
	if method != http.MethodGet {
		req, _ = http.NewRequestWithContext(ctx, method, url, strings.NewReader(body))
	}
	for _, h := range header {
		for k, v := range h {
			req.Header.Set(k, v)
		}
	}

	if timeout == 0 {
		timeout = 3000
	}

	resp, err := (&http.Client{
		Timeout: time.Millisecond * time.Duration(timeout),
	}).Do(req)

	t := time.Since(st)
	if err != nil {
		glog.Errorf(ctx, "err : %s , url :%s ,body : %s ", err.Error(), url, body)
		return
	}
	defer resp.Body.Close()
	rsp, _ = ioutil.ReadAll(resp.Body)
	if strings.Contains(resp.Header.Get("Content-Type"), "image") {
		glog.Infoz(ctx, fmt.Sprintf("url : %s ,body : %s ", url, body), zap.String(glog.ConstTotalTime, t.String()))
	} else {
		glog.Infoz(ctx, fmt.Sprintf("url : %s ,body : %s ,rsp : %s ", url, body, string(rsp)), zap.String(glog.ConstTotalTime, t.String()))
	}
	return
}
