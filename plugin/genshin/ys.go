// Package genshin 原神抽卡
package genshin

import (
	"archive/zip"
	"fmt"
	"github.com/FloatTech/AnimeAPI/wallet"
	fcext "github.com/FloatTech/floatbox/ctxext"
	"github.com/FloatTech/floatbox/process"
	"github.com/FloatTech/imgfactory"
	sql "github.com/FloatTech/sqlite"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	"github.com/golang/freetype"
	"github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type zipfilestructure map[string][]*zip.File

type dbs struct {
	db *sql.Sqlite
	sync.RWMutex
}

// 全局设置
type setting struct {
	Key   string
	Value string
}

// 全局设置
type settings struct {
	ThreeRate int
	FourRate  int
	Four2Rate int
	FiveRate  int
	//five2Rate   int
	ThreeReword int
	FourReword  int
	Four2Reword int
	FiveReword  int
	Five2Reword int
}

var (
	totl                   uint64 // 累计抽奖次数
	filetree               = make(zipfilestructure, 32)
	starN3, starN4, starN5 *zip.File
	namereg                = regexp.MustCompile(`_(.*)\.png`)
	min                    = 10
	max                    = 1000
	database               = &dbs{
		db: &sql.Sqlite{},
	}

	engine = control.Register("genshin", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "模拟抽卡",
		Help:             "- 十连\n- 切换卡池\n- 十连<数字>",
		PublicDataFolder: "Genshin",
	}).ApplySingle(ctxext.DefaultSingle)

	getDb = fcext.DoOnceOnSuccess(func(ctx *zero.Ctx) bool {
		database.db.DBPath = engine.DataFolder() + "ys.db"
		err := database.db.Open(time.Hour * 24)
		if err == nil {
			err = database.db.Create("setting", &setting{})
			if err != nil {
				ctx.SendChain(message.Text("[ERROR]:", err))
				return false
			}
			return true
		}
		ctx.SendChain(message.Text("[ERROR]:", err))
		return false
	})
)

func init() {

	engine.OnFullMatch("切换卡池").SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
			if !ok {
				ctx.SendChain(message.Text("找不到服务!"))
				return
			}
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			store := (storage)(c.GetData(gid))
			if store.setmode(!store.is5starsmode()) {
				process.SleepAbout1sTo2s()
				ctx.SendChain(message.Text("切换到五星卡池~"))
			} else {
				process.SleepAbout1sTo2s()
				ctx.SendChain(message.Text("切换到普通卡池~"))
			}
			err := c.SetData(gid, int64(store))
			if err != nil {
				process.SleepAbout1sTo2s()
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
			}
		})

	engine.OnFullMatch("十连", fcext.DoOnceOnSuccess(
		func(ctx *zero.Ctx) bool {
			zipfile := engine.DataFolder() + "Genshin.zip"
			_, err := engine.GetLazyData("Genshin.zip", false)
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return false
			}
			err = parsezip(zipfile)
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return false
			}
			return true
		},
	), getDb).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
			if !ok {
				ctx.SendChain(message.Text("找不到服务!"))
				return
			}
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			store := (storage)(c.GetData(gid))
			sts := database.getSettings()
			img, str, mode, _, err := randnums(10, store, sts)
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			b, err := imgfactory.ToBytes(img)
			if err != nil {
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			if mode {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID,
					message.Text("恭喜你抽到了: \n", str), message.ImageBytes(b)))
			} else {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID,
					message.Text("十连成功~"), message.ImageBytes(b)))
			}
		})

	engine.OnRegex(`^十连\s*(\d+)`, fcext.DoOnceOnSuccess(
		func(ctx *zero.Ctx) bool {
			zipfile := engine.DataFolder() + "Genshin.zip"
			_, err := engine.GetLazyData("Genshin.zip", false)
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return false
			}
			err = parsezip(zipfile)
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return false
			}
			return true
		},
	), getDb).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
			if !ok {
				ctx.SendChain(message.Text("找不到服务!"))
				return
			}
			invest, err := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			uid := ctx.Event.UserID
			money := wallet.GetWalletOf(uid)
			if invest < min {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("单次最低投入", min, "个糖果哦~")))
				invest = min
			}
			if invest > max {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("单次最高投入", max, "个糖果哦~")))
				invest = max
			}
			if money < invest {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("啊这,你的糖果不够~")))
				return
			}
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			store := (storage)(c.GetData(gid))
			sts := database.getSettings()
			img, str, mode, reword, err := randnums(10, store, sts)
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			b, cl := writer.ToBytes(img)
			rst := ""
			if !store.is5starsmode() {
				realReword := reword * invest / 1000
				insert := realReword - invest
				_ = wallet.InsertWalletOf(zero.BotConfig.SuperUsers[1], invest)
				_ = wallet.InsertWalletOf(zero.BotConfig.SuperUsers[1], -realReword)
				err := wallet.InsertWalletOf(ctx.Event.UserID, insert)
				fmt.Println("[ys] 投入糖果", invest)
				fmt.Println("[ys] 产出数值", reword, "‰")
				fmt.Println("[ys] 产出糖果", realReword)
				fmt.Println("[ys] 真实收获", insert)
				if err != nil {
					fmt.Println("[ys] 钱包坏掉力", err)
					ctx.SendChain(message.Text("[ys] 钱包坏掉力"))
					return
				}
				if realReword < 0 {
					rst = fmt.Sprintf("你投入了%d块糖果,并被扣除了%d块糖果~", invest, -realReword)
				} else if realReword == 0 {
					rst = fmt.Sprintf("你投入了%d块糖果,然后无事发生~", invest)
				} else {
					rst = fmt.Sprintf("你投入了%d块糖果,并收获了%d块糖果~", invest, realReword)
				}
			} else {
				rst = "当前是五星模式!"
			}
			if mode {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID,
					message.Text("恭喜你抽到了: \n", str), message.ImageBytes(b), message.Text("\n", rst)))
			} else {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("十连成功~"), message.ImageBytes(b), message.Text("\n", rst)))
			}

			cl()
		})

	engine.OnRegex(`^五十连\s*(\d+)`, fcext.DoOnceOnSuccess(
		func(ctx *zero.Ctx) bool {
			zipfile := engine.DataFolder() + "Genshin.zip"
			_, err := engine.GetLazyData("Genshin.zip", false)
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return false
			}
			err = parsezip(zipfile)
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return false
			}
			return true
		},
	), getDb).SetBlock(true).Limit(ctxext.LimitByUser).
		Handle(func(ctx *zero.Ctx) {
			c, ok := ctx.State["manager"].(*ctrl.Control[*zero.Ctx])
			if !ok {
				ctx.SendChain(message.Text("找不到服务!"))
				return
			}
			invest, err := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			uid := ctx.Event.UserID
			money := wallet.GetWalletOf(uid)
			if invest < min {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("单次最低投入", min, "个糖果哦~")))
				invest = min
			}
			if invest > max {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("单次最高投入", max, "个糖果哦~")))
				invest = max
			}
			if money < (invest * 5) {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("啊这,你的糖果不够~")))
				return
			}
			gid := ctx.Event.GroupID
			if gid == 0 {
				gid = -ctx.Event.UserID
			}
			store := (storage)(c.GetData(gid))
			sts := database.getSettings()
			img1, str1, _, reword1, err := randnums(10, store, sts)
			img2, str2, _, reword2, err2 := randnums(10, store, sts)
			img3, str3, _, reword3, err3 := randnums(10, store, sts)
			img4, str4, _, reword4, err4 := randnums(10, store, sts)
			img5, str5, _, reword5, err5 := randnums(10, store, sts)
			if err != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			if err2 != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			if err3 != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			if err4 != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}
			if err5 != nil {
				fmt.Println(err)
				ctx.SendChain(message.Text("ERROR: ", err))
				return
			}

			if str1 == "" {
				str1 = "什么都没有~"
			}
			if str2 == "" {
				str2 = "什么都没有~"
			}
			if str3 == "" {
				str3 = "什么都没有~"
			}
			if str4 == "" {
				str4 = "什么都没有~"
			}
			if str5 == "" {
				str5 = "什么都没有~"
			}

			b1, cl1 := writer.ToBytes(img1)
			b2, cl2 := writer.ToBytes(img2)
			b3, cl3 := writer.ToBytes(img3)
			b4, cl4 := writer.ToBytes(img4)
			b5, cl5 := writer.ToBytes(img5)
			rst := ""
			if !store.is5starsmode() {
				realReword1 := reword1 * invest / 1000
				realReword2 := reword2 * invest / 1000
				realReword3 := reword3 * invest / 1000
				realReword4 := reword4 * invest / 1000
				realReword5 := reword5 * invest / 1000
				realReword := realReword1 + realReword2 + realReword3 + realReword4 + realReword5
				reword := reword1 + reword2 + reword3 + reword4 + reword5
				insert := realReword - (invest * 5)
				_ = wallet.InsertWalletOf(zero.BotConfig.SuperUsers[1], invest*5)
				_ = wallet.InsertWalletOf(zero.BotConfig.SuperUsers[1], -realReword)
				err := wallet.InsertWalletOf(ctx.Event.UserID, insert)
				fmt.Println("[ys] 投入糖果", invest*5)
				fmt.Println("[ys] 产出数值", reword, "‰")
				fmt.Println("[ys] 产出糖果", realReword)
				fmt.Println("[ys] 真实收获", insert)
				if err != nil {
					fmt.Println("[ys] 钱包坏掉力", err)
					ctx.SendChain(message.Text("[ys] 钱包坏掉力"))
					return
				}
				if realReword < 0 {
					rst = fmt.Sprintf("你投入了%d块糖果,并被扣除了%d块糖果~", invest*5, -realReword)
				} else if realReword == 0 {
					rst = fmt.Sprintf("你投入了%d块糖果,然后无事发生~", invest*5)
				} else {
					rst = fmt.Sprintf("你投入了%d块糖果,并收获了%d块糖果~", invest*5, realReword)
				}
			} else {
				rst = "当前是五星模式!"
			}
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID,
				message.Text("恭喜你抽到了: \n", str1), message.ImageBytes(b1),
				message.Text("\n", str2), message.ImageBytes(b2),
				message.Text("\n", str3), message.ImageBytes(b3),
				message.Text("\n", str4), message.ImageBytes(b4),
				message.Text("\n", str5), message.ImageBytes(b5),
				message.Text("\n", rst)))

			cl1()
			cl2()
			cl3()
			cl4()
			cl5()
		})
}

func randnums(nums int, store storage, sts settings) (rgba *image.RGBA, str string, replyMode bool, reword int, err error) {
	var (
		fours, fives                  = make([]*zip.File, 0, 10), make([]*zip.File, 0, 10)                           // 抽到 四, 五星角色
		threeArms, fourArms, fiveArms = make([]*zip.File, 0, 10), make([]*zip.File, 0, 10), make([]*zip.File, 0, 10) // 抽到 三 , 四, 五星武器
		fourN, fiveN                  = 0, 0                                                                         // 抽到 四, 五星角色的数量
		bgs                           = make([]*zip.File, 0, 10)                                                     // 背景图片名
		threeN2, fourN2, fiveN2       = 0, 0, 0                                                                      // 抽到 三 , 四, 五星武器的数量
		hero, stars                   = make([]*zip.File, 0, 10), make([]*zip.File, 0, 10)                           // 角色武器名, 储存星级图标

		cicon                   = make([]*zip.File, 0, 10)                                                            // 元素图标
		fivebg, fourbg, threebg = filetree["five_bg.jpg"][0], filetree["four_bg.jpg"][0], filetree["three_bg.jpg"][0] // 背景图片名
		fivelen                 = len(filetree["five"])
		five2len                = len(filetree["five2"])
		threelen                = len(filetree["Three"])
		fourlen                 = len(filetree["four"])
		four2len                = len(filetree["four2"])

		//threeReword, fourReword, fiveReword = -200, 350, 7000
		luckyMode = 0
	)
	fmt.Println(sts)
	threeRate := sts.ThreeRate // 800
	fourRate := sts.FourRate   // 890
	four2Rate := sts.Four2Rate // 980
	fiveRate := sts.FiveRate   // 990
	//five2Rate := settings.threeRate
	threeReword := sts.ThreeReword
	fourReword := sts.FourReword
	four2Reword := sts.Four2Reword
	fiveReword := sts.FiveReword
	five2Reword := sts.Five2Reword

	if totl%90 == 0 { // 累计9次加入一个五星
		switch rand.Intn(2) {
		case 0:
			fiveN++
			fmt.Println("[ys] 保底5x +", fiveReword, ",now ", reword)
			reword += fiveReword
			fives = append(fives, filetree["five"][rand.Intn(fivelen)])
		case 1:
			fiveN2++
			fmt.Println("[ys] 保底5x +", five2Reword, ",now ", reword)
			reword += five2Reword
			fiveArms = append(fiveArms, filetree["five2"][rand.Intn(five2len)])
		}
		nums--
	}

	if store.is5starsmode() { // 5星模式
		for i := 0; i < nums; i++ {
			switch rand.Intn(2) {
			case 0:
				fiveN++
				fives = append(fives, filetree["five"][rand.Intn(fivelen)])
			case 1:
				fiveN2++
				fiveArms = append(fiveArms, filetree["five2"][rand.Intn(five2len)])
			}
		}
	} else { // 默认模式
		lucky := rand.Intn(100000)
		fmt.Println("[ys] lucky ", lucky)
		switch {
		case lucky >= 0 && lucky <= 99900:
			for i := 0; i < nums; i++ {
				a := rand.Intn(1000) // 抽卡几率 三星80% 四星19% 五星1%
				switch {
				case a >= 0 && a <= threeRate:
					threeN2++
					reword += threeReword
					fmt.Println("[ys] 3x +", threeReword, ",now ", reword, ",a ", a)
					threeArms = append(threeArms, filetree["Three"][rand.Intn(threelen)])
				case a > threeRate && a <= fourRate:
					fourN++
					reword += fourReword
					fmt.Println("[ys] 4x +", fourReword, ",now ", reword, ",a ", a)
					fours = append(fours, filetree["four"][rand.Intn(fourlen)]) // 随机角色
				case a > fourRate && a <= four2Rate:
					fourN2++
					reword += four2Reword
					fmt.Println("[ys] 4x +", four2Reword, ",now ", reword, ",a ", a)
					fourArms = append(fourArms, filetree["four2"][rand.Intn(four2len)]) // 随机武器
				case a > four2Rate && a <= fiveRate:
					fiveN++
					reword += fiveReword
					fmt.Println("[ys] 5x +", fiveReword, ",now ", reword, ",a ", a)
					fives = append(fives, filetree["five"][rand.Intn(fivelen)])
				default:
					fiveN2++
					reword += five2Reword
					fmt.Println("[ys] 5x +", five2Reword, ",now ", reword, ",a ", a)
					fiveArms = append(fiveArms, filetree["five2"][rand.Intn(five2len)])
				}
			}
		case lucky >= 99901 && lucky <= 99990:
			luckyMode = 1
			fmt.Println("[ys] 幸运模式")
			for i := 0; i < nums; i++ {
				switch rand.Intn(2) {
				case 0:
					fiveN++
					reword += fiveReword
					fives = append(fives, filetree["five"][rand.Intn(fivelen)])
				case 1:
					fiveN2++
					reword += five2Reword
					fiveArms = append(fiveArms, filetree["five2"][rand.Intn(five2len)])
				}
			}
		default:
			luckyMode = 2
			fmt.Println("[ys] 超级幸运模式")
			for i := 0; i < nums; i++ {
				switch rand.Intn(2) {
				case 0:
					fiveN++
					reword += fiveReword * 10
					fives = append(fives, filetree["five"][rand.Intn(fivelen)])
				case 1:
					fiveN2++
					reword += five2Reword * 10
					fiveArms = append(fiveArms, filetree["five2"][rand.Intn(five2len)])
				}
			}
		}

		if fourN+fourN2 == 0 && threeN2 > 0 { // 没有四星时自动加入
			reword -= threeReword
			threeN2--
			threeArms = threeArms[:len(threeArms)-1]
			switch rand.Intn(2) {
			case 0:
				fourN++
				reword += fourReword
				fmt.Println("[ys] 保底4x -", threeReword, " ,+", fourReword, ",now ", reword)
				fours = append(fours, filetree["four"][rand.Intn(fourlen)]) // 随机角色
			case 1:
				fourN2++
				reword += four2Reword
				fmt.Println("[ys] 保底4x -", threeReword, " ,+", four2Reword, ",now ", reword)
				fourArms = append(fourArms, filetree["four2"][rand.Intn(four2len)]) // 随机武器
			}
		}
		_ = atomic.AddUint64(&totl, 1)
	}

	icon := func(f *zip.File) *zip.File {
		name := f.Name
		name = name[strings.LastIndex(name, "/")+1:strings.Index(name, "_")] + ".png"
		logrus.Debugln("[genshin]get named file", name)
		return filetree[name][0]
	}

	he := func(cnt int, id int, f *zip.File, bg *zip.File) {
		var hen *[]*zip.File
		for i := 0; i < cnt; i++ {
			switch id {
			case 1:
				hen = &threeArms
			case 2:
				hen = &fourArms
			case 3:
				hen = &fours
			case 4:
				hen = &fiveArms
			case 5:
				hen = &fives
			}
			bgs = append(bgs, bg) // 加入颜色背景
			hero = append(hero, (*hen)[i])
			stars = append(stars, f)               // 加入星级图标
			cicon = append(cicon, icon((*hen)[i])) // 加入元素图标
		}
	}

	if fiveN > 0 { // 按顺序加入
		he(fiveN, 5, starN5, fivebg) // 五星角色
		str += reply(fives, 1, str)
		replyMode = true
	}
	if fourN > 0 {
		he(fourN, 3, starN4, fourbg) // 四星角色
	}
	if fiveN2 > 0 {
		he(fiveN2, 4, starN5, fivebg) // 五星武器
		str += reply(fiveArms, 2, str)
		replyMode = true
	}
	if fourN2 > 0 {
		he(fourN2, 2, starN4, fourbg) // 四星武器
	}
	if threeN2 > 0 {
		he(threeN2, 1, starN3, threebg) // 三星武器
	}
	if luckyMode == 1 {
		str += "\n你触发了幸运模式,全五星激活!"
	}
	if luckyMode == 2 {
		str += "\n你触发了超级幸运模式,超高收益全五星激活!"
	}

	var c1, c2, c3 uint8 = 50, 50, 50 // 背景颜色

	img00, err := filetree["bg0.jpg"][0].Open() // 打开背景图片
	if err != nil {
		return
	}

	rectangle := image.Rect(0, 0, 1920, 1080) // 图片宽度, 图片高度
	rgba = image.NewRGBA(rectangle)
	draw.Draw(rgba, rgba.Bounds(), image.NewUniform(color.RGBA{c1, c2, c3, 255}), image.Point{}, draw.Over)
	context := freetype.NewContext() // 创建一个新的上下文
	context.SetDPI(72)               // 每英寸 dpi
	context.SetClip(rgba.Bounds())
	context.SetDst(rgba)

	defer img00.Close()
	img0, err := jpeg.Decode(img00) // 读取一个本地图像
	if err != nil {
		return
	}

	offset := image.Pt(0, 0) // 图片在背景上的位置
	draw.Draw(rgba, img0.Bounds().Add(offset), img0, image.Point{}, draw.Over)

	w1, h1 := 230, 0
	for i := 0; i < len(hero); i++ {
		if i > 0 {
			w1 += 146 // 图片宽度
		}

		imgs, err := bgs[i].Open() // 取出背景图片
		if err != nil {
			return nil, "", false, reword, err
		}
		defer imgs.Close()

		img, _ := jpeg.Decode(imgs)
		offset := image.Pt(w1, h1)
		draw.Draw(rgba, img.Bounds().Add(offset), img, image.Point{}, draw.Over)

		imgs1, err := hero[i].Open() // 取出图片名
		if err != nil {
			return nil, "", false, reword, err
		}
		defer imgs1.Close()

		img1, _ := png.Decode(imgs1)
		offset1 := image.Pt(w1, h1)
		draw.Draw(rgba, img1.Bounds().Add(offset1), img1, image.Point{}, draw.Over)

		imgs2, err := stars[i].Open() // 取出星级图标
		if err != nil {
			return nil, "", false, reword, err
		}
		defer imgs2.Close()

		img2, _ := png.Decode(imgs2)
		offset2 := image.Pt(w1, h1)
		draw.Draw(rgba, img2.Bounds().Add(offset2), img2, image.Point{}, draw.Over)

		imgs3, err := cicon[i].Open() // 取出类型图标
		if err != nil {
			return nil, "", false, reword, err
		}
		defer imgs3.Close()

		img3, _ := png.Decode(imgs3)
		offset3 := image.Pt(w1, h1)
		draw.Draw(rgba, img3.Bounds().Add(offset3), img3, image.Point{}, draw.Over)
	}
	imgs4, err := filetree["Reply.png"][0].Open() // "分享" 图标
	if err != nil {
		return nil, "", false, reword, err
	}
	defer imgs4.Close()
	img4, err := png.Decode(imgs4)
	if err != nil {
		return nil, "", false, reword, err
	}
	offset4 := image.Pt(1270, 945) // 宽, 高
	draw.Draw(rgba, img4.Bounds().Add(offset4), img4, image.Point{}, draw.Over)
	return
}

func parsezip(zipFile string) error {
	zipReader, err := zip.OpenReader(zipFile) // will not close
	if err != nil {
		return err
	}
	for _, f := range zipReader.File {
		if f.FileInfo().IsDir() {
			filetree[f.Name] = make([]*zip.File, 0, 32)
			continue
		}
		f.Name = f.Name[8:]
		i := strings.LastIndex(f.Name, "/")
		if i < 0 {
			filetree[f.Name] = []*zip.File{f}
			logrus.Debugln("[genshin]insert file", f.Name)
			continue
		}
		folder := f.Name[:i]
		if folder != "" {
			filetree[folder] = append(filetree[folder], f)
			logrus.Debugln("[genshin]insert file into", folder)
			if folder == "gacha" {
				switch f.Name[i+1:] {
				case "ThreeStar.png":
					starN3 = f
				case "FourStar.png":
					starN4 = f
				case "FiveStar.png":
					starN5 = f
				}
			}
		}
	}
	return nil
}

// 取出角色武器名
func reply(z []*zip.File, num int, nameStr string) string {
	var tmp strings.Builder
	tmp.Grow(128)
	switch {
	case num == 1:
		tmp.WriteString("★五星角色★\n")
	case num == 2 && len(nameStr) > 0:
		tmp.WriteString("\n★五星武器★\n")
	default:
		tmp.WriteString("★五星武器★\n")
	}
	for i := range z {
		tmp.WriteString(namereg.FindStringSubmatch(z[i].Name)[1] + " * ")
	}
	return tmp.String()
}

func (sql *dbs) getSettings() (sts settings) {
	sql.Lock()
	defer sql.Unlock()
	err := sql.db.Create("setting", &setting{})
	if err != nil {
		return
	}
	var st setting

	err = sql.db.Find("setting", &st, "where key = 'ThreeRate'")
	if err != nil {
		sta := setting{Key: "ThreeRate", Value: "800"}
		_ = sql.db.Insert("setting", &sta)
	}
	sts.ThreeRate, _ = strconv.Atoi(st.Value)

	err = sql.db.Find("setting", &st, "where key = 'FourRate'")
	if err != nil {
		sta := setting{Key: "FourRate", Value: "890"}
		_ = sql.db.Insert("setting", &sta)
	}
	sts.FourRate, _ = strconv.Atoi(st.Value)

	err = sql.db.Find("setting", &st, "where key = 'Four2Rate'")
	if err != nil {
		sta := setting{Key: "Four2Rate", Value: "980"}
		_ = sql.db.Insert("setting", &sta)
	}
	sts.Four2Rate, _ = strconv.Atoi(st.Value)

	err = sql.db.Find("setting", &st, "where key = 'FiveRate'")
	if err != nil {
		sta := setting{Key: "FiveRate", Value: "990"}
		_ = sql.db.Insert("setting", &sta)
	}
	sts.FiveRate, _ = strconv.Atoi(st.Value)

	err = sql.db.Find("setting", &st, "where key = 'ThreeReword'")
	if err != nil {
		sta := setting{Key: "ThreeReword", Value: "-200"}
		_ = sql.db.Insert("setting", &sta)
	}
	sts.ThreeReword, _ = strconv.Atoi(st.Value)

	err = sql.db.Find("setting", &st, "where key = 'FourReword'")
	if err != nil {
		sta := setting{Key: "FourReword", Value: "350"}
		_ = sql.db.Insert("setting", &sta)
	}
	sts.FourReword, _ = strconv.Atoi(st.Value)

	err = sql.db.Find("setting", &st, "where key = 'Four2Reword'")
	if err != nil {
		sta := setting{Key: "Four2Reword", Value: "350"}
		_ = sql.db.Insert("setting", &sta)
	}
	sts.Four2Reword, _ = strconv.Atoi(st.Value)

	err = sql.db.Find("setting", &st, "where key = 'FiveReword'")
	if err != nil {
		sta := setting{Key: "FiveReword", Value: "350"}
		_ = sql.db.Insert("setting", &sta)
	}
	sts.FiveReword, _ = strconv.Atoi(st.Value)

	err = sql.db.Find("setting", &st, "where key = 'Five2Reword'")
	if err != nil {
		sta := setting{Key: "Five2Reword", Value: "7000"}
		_ = sql.db.Insert("setting", &sta)
	}
	sts.Five2Reword, _ = strconv.Atoi(st.Value)

	return
}
