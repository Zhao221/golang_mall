package timer

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"golang_mall/dao/mysql"
	"golang_mall/dao/redis"
	"strconv"
)

var (
	count, month, year int64
)
var ctx context.Context

// 在这个示例中，我们使用 robfig/cron/v3 库创建了一个支持秒级别的 Cron 定时器，
// 并添加了一个每天凌晨 12 点触发的任务（Cron 表达式为 "0 0 0 * * *"）。
// 然后启动定时器，主 Goroutine 使用 select {} 阻塞等待，定时器将持续运行。
// 如果需要停止定时器，可以在适当的时机调用 c.Stop() 方法。
// 注意：在实际应用中，根据需要，你可以选择合适的方式来停止定时器。如通过信号、关闭通道等。

func InitDayCheckinDayTimer() {
	c := cron.New(cron.WithSeconds()) // 创建一个支持秒级别的 Cron 定时器
	_, err := c.AddFunc("0 0 0 * * *", func() {
		fmt.Println("执行任务...")
		var Ids []uint
		err := mysql.NewUserDao(ctx).Table("user").Select("id").Where("daily_checkin = ?", true).Find(&Ids).Error
		if err != nil {
			return
		}
		err = mysql.NewUserDao(ctx).Table("user").Select("daily_checkin").
			Where("id IN ?", Ids).UpdateColumn("daily_checkin", false).
			Error
		if err != nil {
			return
		}
	})
	if err != nil {
		return
	}
	c.Start() // 启动定时器
	// 主 Goroutine 阻塞等待，这里仅作示例，实际应用中可以根据需要停止定时器
}

// 一个支持秒级别的 Cron 定时器，并添加了一个每个月的 25 日凌晨触发的任务。
// Cron 表达式为 "0 0 0 25 * *"，表示每个月的 25 日的 00:00:00 触发任务。
// 这样，定时任务将在每个月的 25 日凌晨触发。

func InitMonthCheckinTimer() {
	c := cron.New(cron.WithSeconds()) // 创建一个支持秒级别的 Cron 定时器
	// 添加定时任务
	_, err := c.AddFunc("0 0 0 1 * *", func() {
		fmt.Println("执行任务...")
		var Ids []uint
		err := mysql.NewUserDao(ctx).Table("user").Select("id").Find(&Ids).Error
		if err != nil {
			return
		}
		for i := 0; i < len(Ids); i++ {
			count, err = redis.RedisClient.BitCount(redis.RedisContext, "Checkin"+strconv.Itoa(int(Ids[i])), nil).Result()
			if err != nil {
				return
			}
			mysql.NewUserDao(ctx).Table("user").Select("monthly_checkin").Where("id = ?", Ids[i]).First(&month)
			mysql.NewUserDao(ctx).Table("user").Select("year_checkin").Where("id = ?", Ids[i]).First(&year)
			month += count
			mysql.NewUserDao(ctx).Table("user").Where("id = ?", Ids[i]).Update("monthly_checkin", month)
			year += count
			mysql.NewUserDao(ctx).Table("user").Where("id = ?", Ids[i]).Update("year_checkin", year)
		}

	})
	if err != nil {
		return
	}
	// 启动定时器
	c.Start()
	// 主 Goroutine 阻塞等待
}
