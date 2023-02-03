package aamy

import (
	"fmt"
	"github.com/FloatTech/AnimeAPI/wallet"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
	"github.com/FloatTech/zbputils/ctxext"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"
	"math/rand"
	"strconv"
	"time"
)

func init() { // 插件主体
	engine := control.Register("aamy", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault: false,
		Brief:            "",
		Help:             "神秘插件",
	})

	//engine.OnRegex(`\[CQ:reply,id=(\-?[0-9]+)\].*撤回.*`).SetBlock(true).
	//	Handle(func(ctx *zero.Ctx) {
	//		//[CQ:reply,id=543569241]撤回
	//		messageID := ctx.State["regex_matched"].([]string)[0]
	//		messageID = strings.Split(messageID, "=")[1]
	//		messageID = strings.Split(messageID, "]")[0]
	//
	//		ctx.CallAction("delete_msg", zero.Params{
	//			"message_id": messageID,
	//		})
	//	})

	//engine.OnFullMatch("我的权限", zero.OnlyToMe).SetBlock(true).
	//	Handle(func(ctx *zero.Ctx) {
	//		if zero.SuperUserPermission(ctx) {
	//			ctx.Send("超管")
	//		}
	//		if ctx.Event.Sender.Role == "owner" {
	//			ctx.Send("群主")
	//		}
	//		if ctx.Event.Sender.Role == "admin" {
	//			ctx.Send("群管")
	//		}
	//		ctx.Send("群员")
	//	})

	zero.On("notice/group_ban/ban", zero.OnlyToMe, zero.OnlyGroup).SetBlock(false).
		Handle(func(ctx *zero.Ctx) {
			oid := ctx.Event.OperatorID
			if !isSuperUserPermission(oid) {
				ctx.SetGroupLeave(ctx.Event.GroupID, false)
			}
		})

	engine.OnRegex(`转账\[CQ:at,qq=(\d+)\]\s*(\d+)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			target, _ := strconv.ParseInt(ctx.State["regex_matched"].([]string)[1], 10, 64)
			money, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
			selfMoney := wallet.GetWalletOf(uid)
			if selfMoney < money {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你的糖果不够!")))
				return
			}
			err := wallet.InsertWalletOf(uid, -money)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("转出失败,请检查!")))
				return
			}
			tap1 := 100
			tap2 := 200
			tap3 := 500
			tap4 := 1000
			tap5 := 2000
			tap6 := 5000
			tap7 := 10000
			tap8 := 20000
			tap9 := 50000
			tap10 := 100000
			tap11 := 200000
			tap12 := 500000
			tap1Tax := 50
			tap2Tax := 20
			tap3Tax := 19
			tap4Tax := 17
			tap5Tax := 15
			tap6Tax := 13
			tap7Tax := 11
			tap8Tax := 9
			tap9Tax := 7
			tap10Tax := 5
			tap11Tax := 3
			tap12Tax := 2

			tax := 0
			if money > tap1 {
				if money < tap2 {
					tax += (money - tap1) / tap1Tax
				} else {
					tax += (tap2 - tap1) / tap1Tax
				}
			}
			if money > tap2 {
				if money < tap3 {
					tax += (money - tap2) / tap2Tax
				} else {
					tax += (tap3 - tap2) / tap2Tax
				}
			}
			if money > tap3 {
				if money < tap4 {
					tax += (money - tap3) / tap3Tax
				} else {
					tax += (tap4 - tap3) / tap3Tax
				}
			}
			if money > tap4 {
				if money < tap5 {
					tax += (money - tap4) / tap4Tax
				} else {
					tax += (tap5 - tap4) / tap4Tax
				}
			}
			if money > tap5 {
				if money < tap6 {
					tax += (money - tap5) / tap5Tax
				} else {
					tax += (tap6 - tap5) / tap5Tax
				}
			}
			if money > tap6 {
				if money < tap7 {
					tax += (money - tap6) / tap6Tax
				} else {
					tax += (tap7 - tap6) / tap6Tax
				}
			}
			if money > tap7 {
				if money < tap8 {
					tax += (money - tap7) / tap7Tax
				} else {
					tax += (tap8 - tap7) / tap7Tax
				}
			}
			if money > tap8 {
				if money < tap9 {
					tax += (money - tap8) / tap8Tax
				} else {
					tax += (tap9 - tap8) / tap8Tax
				}
			}
			if money > tap9 {
				if money < tap10 {
					tax += (money - tap9) / tap9Tax
				} else {
					tax += (tap10 - tap9) / tap9Tax
				}
			}
			if money > tap10 {
				if money < tap11 {
					tax += (money - tap10) / tap10Tax
				} else {
					tax += (tap11 - tap10) / tap10Tax
				}
			}
			if money > tap11 {
				if money < tap12 {
					tax += (money - tap11) / tap11Tax
				} else {
					tax += (tap12 - tap11) / tap11Tax
				}
			}
			if money > tap12 {
				tax += (money - tap12) / tap12Tax
			}

			err = wallet.InsertWalletOf(target, money-tax)
			err = wallet.InsertWalletOf(zero.BotConfig.SuperUsers[1], tax)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("转入失败,请检查!")))
				return
			}
			if tax != 0 {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("糖果转账成功!\n转账途中掉了", tax, "个糖果!")))
			} else {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("糖果转账成功!")))
			}
		})

	engine.OnRegex(`无中生有\s*(\d+)`, zero.SuperUserPermission).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			money, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
			err := wallet.InsertWalletOf(uid, money)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("无中生有成功!")))
		})

	engine.OnFullMatch("查银行").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			money := wallet.GetWalletOf(zero.BotConfig.SuperUsers[1])
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("银行有", money, "块糖果!")))
		})

	engine.OnFullMatch("抢银行").Limit(ctxext.LimitByGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			bankMoney := wallet.GetWalletOf(zero.BotConfig.SuperUsers[1])
			thiefMoney := wallet.GetWalletOf(ctx.Event.UserID)
			if thiefMoney > bankMoney {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("别闹!你才是最大的银行!")))
				return
			}
			if thiefMoney < 51 {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("保安把你拦住了!")))
				return
			}
			lucky := rand.Intn(100000)
			if lucky == 1 && thiefMoney > 100 {
				err := wallet.InsertWalletOf(zero.BotConfig.SuperUsers[1], -bankMoney)
				if err != nil {
					fmt.Println(err)
					ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
					return
				}
				err = wallet.InsertWalletOf(ctx.Event.UserID, bankMoney)
				if err != nil {
					fmt.Println(err)
					ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
					return
				}
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你从银行手上抢到了", bankMoney, "块糖果!")))
				return
			}
			dropMoney := (thiefMoney / 10) + 200
			if dropMoney < 50 {
				dropMoney = thiefMoney
			}
			err := wallet.InsertWalletOf(zero.BotConfig.SuperUsers[1], dropMoney)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			err = wallet.InsertWalletOf(ctx.Event.UserID, -dropMoney)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("银行用源源不绝的糖果砸死了你!你在落荒而逃的途中掉了", dropMoney, "块糖果!")))
			return

		})

	engine.OnRegex(`捐银行\s*(\d+)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			donatMoney, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
			selfMoney := wallet.GetWalletOf(uid)
			if selfMoney < donatMoney {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你的糖果不够!")))
				return
			}
			err := wallet.InsertWalletOf(uid, -donatMoney)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			err = wallet.InsertWalletOf(zero.BotConfig.SuperUsers[1], donatMoney)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			pool := wallet.GetWalletOf(zero.BotConfig.SuperUsers[1])
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你捐了", donatMoney, "个糖果给银行!\n银行现在有", pool, "个糖果!")))
		})

	engine.OnFullMatch("恰低保").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			money := wallet.GetWalletOf(uid)
			if money >= 50 {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你不需要恰低保!")))
				return
			}
			bankMoney := wallet.GetWalletOf(zero.BotConfig.SuperUsers[1])
			userMoney := wallet.GetWalletOf(uid)
			if bankMoney < 50 {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("银行破产了!")))
				return
			}
			tranMoney := -userMoney + 50
			err := wallet.InsertWalletOf(zero.BotConfig.SuperUsers[1], -tranMoney)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			err = wallet.InsertWalletOf(uid, tranMoney)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("V你", tranMoney, "糖果,继续加油!")))
		})

	engine.OnRegex(`填糖果堆\s*(\d+)`).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			luckyMoney, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
			selfMoney := wallet.GetWalletOf(uid)
			if selfMoney < luckyMoney {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你的糖果不够!")))
				return
			}
			err := wallet.InsertWalletOf(uid, -luckyMoney)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			err = wallet.InsertWalletOf(zero.BotConfig.SuperUsers[2], luckyMoney)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			pool := wallet.GetWalletOf(zero.BotConfig.SuperUsers[2])
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("发糖果咯!\n糖果堆有", pool, "个糖果!\n发送 摸糖果堆 来摸糖果!")))
		})

	engine.OnFullMatch("摸糖果堆").SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			pool := wallet.GetWalletOf(zero.BotConfig.SuperUsers[2])
			if pool <= 0 {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("糖果堆没糖了!")))
				return
			}
			num := 0
			if pool <= 10 {
				num = pool
			} else {
				r := rand.Intn(pool / 10)
				num = rand.Intn(pool/2) - r
			}

			if num < 0 {
				num = -num
			}
			if num > pool {
				num = pool
			}
			err := wallet.InsertWalletOf(zero.BotConfig.SuperUsers[2], -num)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			err = wallet.InsertWalletOf(uid, num)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			if pool-num > 0 {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你摸到了", num, "个糖果!\n糖果堆剩余", pool-num, "个糖果!")))
			} else {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你摸到了", num, "个糖果,糖果堆被你摸空了!")))
			}

		})

	engine.OnRegex(`^下糖果雨\s*(\d+)份\s*(\d+)个$`, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			uid := ctx.Event.UserID
			redPackNum, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
			moneyRedPack, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
			selfMoney := wallet.GetWalletOf(uid)
			if selfMoney < moneyRedPack {
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("你的糖果不够!")))
				return
			}
			err := wallet.InsertWalletOf(uid, -moneyRedPack)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("出现了错误!")))
				return
			}
			moneyRemain := moneyRedPack
			numRemain := redPackNum
			//receiveMap := make(map[int]int, 0)
			ctx.SendChain(message.Text("糖果雨开始, 60秒未抢完会自动退回!\n发送 抢糖果雨 来抢糖果雨!"))
			repeatMap := make(map[int64]int, 0)
			next := zero.NewFutureEvent("message", 999, false, zero.RegexRule(`^抢糖果雨$`), zero.OnlyGroup, zero.CheckGroup(ctx.Event.GroupID))
			receive, cancel := next.Repeat()
			defer cancel()

			after := time.NewTimer(60 * time.Second)
		EXIT:
			for {
				select {
				case <-after.C:
					ctx.SendChain(message.Text("糖果雨结束了,剩余", moneyRemain, "未领取,已退回。"))
					err := wallet.InsertWalletOf(uid, moneyRemain)
					if err != nil {
						ctx.SendChain(message.Text("出现了错误,请联系管理员!"))
						return
					}
					break EXIT
				case c := <-receive:
					uid := c.Event.UserID
					if _, ok := repeatMap[uid]; !ok {
						if moneyRemain != 0 {
							x := DoubleAverage(numRemain, moneyRemain)
							repeatMap[uid] = x
							moneyRemain -= x
							numRemain--
							err := wallet.InsertWalletOf(uid, x)
							if err != nil {
								ctx.SendChain(message.Text("出现了错误,请联系管理员!"))
							} else {
								if moneyRemain == 0 || numRemain == 0 {
									ctx.Send(message.ReplyWithMessage(c.Event.MessageID, message.Text("你抢到了", x, "个糖果!\n糖果雨抢完了~")))
								} else {
									ctx.Send(message.ReplyWithMessage(c.Event.MessageID, message.Text("你抢到了", x, "个糖果!\n剩余", numRemain, "份!\n剩余", moneyRemain, "个糖果~")))
								}
							}
						}
					} else {
						ctx.Send(message.ReplyWithMessage(c.Event.MessageID, message.Text("你抢过一次了!")))
					}
				default:
					if moneyRemain == 0 || numRemain == 0 {
						break EXIT
					}
				}
			}

		})

	engine.OnRegex(`^(\d+)d(\d+)$`, zero.AdminPermission, zero.OnlyToMe, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			num, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[1])
			faceRange, _ := strconv.Atoi(ctx.State["regex_matched"].([]string)[2])
			var str = ""
			if num <= 1 {
				ctx.SendChain(message.Text("个数不合法!"))
				return
			}
			if faceRange <= 1 {
				ctx.SendChain(message.Text("面数不合法!"))
				return
			}
			for i := 0; i < num; i++ {
				face := rand.Intn(faceRange) + 1
				str += fmt.Sprintf("%d. %d\n", i+1, face)
			}
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text(str)))
		})
}

// DoubleAverage 二倍均值算法
func DoubleAverage(count, amount int) int {
	// 能抢到的最小金额
	var min = 1
	if count == 1 {
		return amount
	}
	//计算出最大可用金额
	max := amount - min*count
	//计算出最大可用平均值
	avg := max / count
	//二倍均值基础上再加上最小金额 防止出现金额为0
	avg2 := 2*avg + min
	//随机红包金额序列元素，把二倍均值作为随机的最大数
	rand.Seed(time.Now().UnixNano())
	x := rand.Intn(avg2) + min
	return x
}

func isSuperUserPermission(uid int64) bool {
	for _, su := range zero.BotConfig.SuperUsers {
		if su == uid {
			return true
		}
	}
	return false
}
