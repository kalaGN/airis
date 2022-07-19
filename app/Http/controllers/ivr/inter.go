package ivr

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/kataras/iris/v12"
)

type CommonRes struct {
	Code string `json:"code"`
	Id   string `json:"id"`
}

type CommonResg struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

func Create(ctx iris.Context) {
	rand.Seed(time.Now().UnixNano())

	d1 := CommonRes{"0", "cret0000" + strconv.Itoa(rand.Intn(10000))}
	ctx.JSON(d1)
}

func Getflowname(ctx iris.Context) {
	rand.Seed(time.Now().UnixNano())
	d1 := CommonResg{"0", "name0000" + strconv.Itoa(rand.Intn(10000))}
	ctx.JSON(d1)
}
func Copy(ctx iris.Context) {
	rand.Seed(time.Now().UnixNano())
	d1 := CommonRes{"0", "copy0000" + strconv.Itoa(rand.Intn(10000))}
	ctx.JSON(d1)
}
func Del(ctx iris.Context) {
	d1 := CommonRes{"200", "ok"}
	ctx.JSON(d1)
}
