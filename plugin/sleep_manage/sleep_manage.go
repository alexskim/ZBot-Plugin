// Package sleepmanage 睡眠管理
package sleepmanage

import (
	"fmt"
	"github.com/FloatTech/AnimeAPI/wallet"
	"math/rand"
	"time"

	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/message"

	ctrl "github.com/FloatTech/zbpctrl"
	"github.com/FloatTech/zbputils/control"
)

func init() {
	engine := control.Register("sleepmanage", &ctrl.Options[*zero.Ctx]{
		DisableOnDefault:  false,
		Brief:             "睡眠小助手",
		Help:              "- 早安\n- 晚安",
		PrivateDataFolder: "sleep",
	})
	go func() {
		sdb = initialize(engine.DataFolder() + "manage.db")
	}()
	engine.OnFullMatch("早安", isMorning, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			position, getUpTime, first := sdb.getUp(ctx.Event.GroupID, ctx.Event.UserID)
			log.Debugln(position, getUpTime)
			hour, minute, second := timeDuration(getUpTime)
			if (hour == 0 && minute == 0 && second == 0) || hour >= 24 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("早安成功！你是今天第%d个起床的", position)))
			} else {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("早安成功！你的睡眠时长为%d时%d分%d秒,你是今天第%d个起床的", hour, minute, second, position)))
			}
			if first == 0 {
				time.Sleep(time.Second * 1)
				morning(position, ctx)
			}
		})
	engine.OnFullMatch("晚安", isEvening, zero.OnlyGroup).SetBlock(true).
		Handle(func(ctx *zero.Ctx) {
			position, sleepTime, first := sdb.sleep(ctx.Event.GroupID, ctx.Event.UserID)
			log.Debugln(position, sleepTime)
			hour, minute, second := timeDuration(sleepTime)
			if (hour == 0 && minute == 0 && second == 0) || hour >= 24 {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("晚安成功！你是今天第%d个睡觉的", position)))
			} else {
				ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("晚安成功！你的清醒时长为%d时%d分%d秒,你是今天第%d个睡觉的", hour, minute, second, position)))
			}
			if first == 0 {
				time.Sleep(time.Second * 1)
				evening(position, ctx)
			}
		})
}

func timeDuration(time time.Duration) (hour, minute, second int64) {
	hour = int64(time) / (1000 * 1000 * 1000 * 60 * 60)
	minute = (int64(time) - hour*(1000*1000*1000*60*60)) / (1000 * 1000 * 1000 * 60)
	second = (int64(time) - hour*(1000*1000*1000*60*60) - minute*(1000*1000*1000*60)) / (1000 * 1000 * 1000)
	return hour, minute, second
}

// 只统计6点到12点的早安
func isMorning(ctx *zero.Ctx) bool {
	now := time.Now().Hour()
	return now >= 6 && now <= 12
}

// 只统计21点到凌晨3点的晚安
func isEvening(ctx *zero.Ctx) bool {
	now := time.Now().Hour()
	return now >= 21 || now <= 3
}

func morning(position int, ctx *zero.Ctx) {
	if position == 1 {
		r := rand.Intn(70) + 10
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("你是今天第一个起床的,送你 %d 块糖果", r)))
		err := wallet.InsertWalletOf(ctx.Event.UserID, r)
		if err != nil {
			ctx.SendChain(message.Text("[sleep_mamager] 钱包坏掉力:\n", err))
			return
		}
	}
	if position == 2 {
		r := rand.Intn(50) + 10
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("你是今天第二个起床的,送你 %d 块糖果", r)))
		err := wallet.InsertWalletOf(ctx.Event.UserID, r)
		if err != nil {
			ctx.SendChain(message.Text("[sleep_mamager] 钱包坏掉力:\n", err))
			return
		}
	}
	if position == 3 {
		r := rand.Intn(30) + 10
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("你是今天第三个起床的,送你 %d 块糖果", r)))
		err := wallet.InsertWalletOf(ctx.Event.UserID, r)
		if err != nil {
			ctx.SendChain(message.Text("[sleep_mamager] 钱包坏掉力:\n", err))
			return
		}
	}
	if position > 3 {
		r := rand.Intn(10) + 10
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("送你 %d 块糖果", r)))
		err := wallet.InsertWalletOf(ctx.Event.UserID, r)
		if err != nil {
			ctx.SendChain(message.Text("[sleep_mamager] 钱包坏掉力:\n", err))
			return
		}
	}
}

func evening(position int, ctx *zero.Ctx) {
	if position == 1 {
		r := rand.Intn(70) + 10
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("你是今天第一个睡觉的,送你 %d 块糖果", r)))
		err := wallet.InsertWalletOf(ctx.Event.UserID, r)
		if err != nil {
			ctx.SendChain(message.Text("[sleep_mamager] 钱包坏掉力:\n", err))
			return
		}
	}
	if position == 2 {
		r := rand.Intn(50) + 10
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("你是今天第二个睡觉的,送你 %d 块糖果", r)))
		err := wallet.InsertWalletOf(ctx.Event.UserID, r)
		if err != nil {
			ctx.SendChain(message.Text("[sleep_mamager] 钱包坏掉力:\n", err))
			return
		}
	}
	if position == 3 {
		r := rand.Intn(30) + 10
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("你是今天第三个睡觉的,送你 %d 块糖果", r)))
		err := wallet.InsertWalletOf(ctx.Event.UserID, r)
		if err != nil {
			ctx.SendChain(message.Text("[sleep_mamager] 钱包坏掉力:\n", err))
			return
		}
	}
	if position > 3 {
		r := rand.Intn(10) + 10
		ctx.SendChain(message.Reply(ctx.Event.MessageID), message.Text(fmt.Sprintf("送你 %d 块糖果", r)))
		err := wallet.InsertWalletOf(ctx.Event.UserID, r)
		if err != nil {
			ctx.SendChain(message.Text("[sleep_mamager] 钱包坏掉力:\n", err))
			return
		}
	}
}
