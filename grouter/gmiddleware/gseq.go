package gmiddleware

import (
	"fmt"
	"github.com/dshibin/gbins/gconf"
	"github.com/dshibin/gbins/gret"
	"github.com/dshibin/gbins/gutil"
	"github.com/gin-gonic/gin"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

var seqArr = map[string]string{}
var seqmu sync.RWMutex
var smu sync.Mutex

const SeqKey = "gMiddlewareSeq"

func init() {
	go func() {
		limit := gconf.GConfig().Server.Limit
		for true {
			for true {
				if len(seqArr) >= limit {
					break
				}
				setSeqArr(createSeqMd5())
			}
			time.Sleep(1 * time.Second)
		}
	}()
}

func GSeq() gin.HandlerFunc {
	return func(c *gin.Context) {
		smu.Lock()
		defer smu.Unlock()
		if c.GetHeader(SeqKey) != "" {
			c.Set(SeqKey, c.GetHeader(SeqKey))
			return
		}
		if c.GetString(SeqKey) == "" {
			seq := CreateSeq()
			if seq == "" {
				gret.Ret(c, gret.ReqTooMore)
				return
			}
			c.Set(SeqKey, seq)
		}
		c.Next()
	}
}

func CreateSeq() string {
	t := time.Now().UnixNano() / 1e6
	seqMd5 := getSeqMd5()
	if seqMd5 == "" {
		return ""
	}
	return fmt.Sprintf("%s-%d-%s", gconf.GConfig().Server.App, t, seqMd5)
}

func getSeqMd5() (ret string) {
	seqmu.Lock()
	defer seqmu.Unlock()
	if len(seqArr) == 0 {
		return
	}
	for k, v := range seqArr {
		ret = v
		delete(seqArr, k)
		break
	}
	return ret
}

func createSeqMd5() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Int31()
	l := rand.Intn(10) + 5
	return gutil.Md5(strconv.FormatInt(int64(n), 10))[:l]
}

func setSeqArr(s string) bool {
	seqmu.Lock()
	defer seqmu.Unlock()
	if _, ok := seqArr[s]; ok {
		return false
	}
	seqArr[s] = s
	return true
}
