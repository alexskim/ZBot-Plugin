package aamy

import (
	"fmt"
	"github.com/FloatTech/AnimeAPI/wallet"
	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
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
		Help:             "",
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

	engine.OnRegex(`转账\[CQ:at,qq=(\d+)\]\s(\d+)`).SetBlock(true).
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
			err = wallet.InsertWalletOf(target, money)
			if err != nil {
				fmt.Println(err)
				ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("转入失败,请检查!")))
				return
			}
			ctx.Send(message.ReplyWithMessage(ctx.Event.MessageID, message.Text("糖果转账成功!")))
		})

	engine.OnRegex(`无中生有\s(\d+)`, zero.SuperUserPermission).SetBlock(true).
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

	engine.OnRegex(`捐银行(\d+)`).SetBlock(true).
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
			err := wallet.InsertWalletOf(zero.BotConfig.SuperUsers[0], -tranMoney)
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

	engine.OnRegex(`填糖果堆(\d+)`).SetBlock(true).
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
