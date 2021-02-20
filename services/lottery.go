package services

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/garbein/lottery-golang/apps"
	"github.com/garbein/lottery-golang/models"
	"github.com/garbein/lottery-golang/responses"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const random_max = 1000

// 开始抽奖
func StartLottery(userId int) (responses.ResponseBody, error) {
	// 一天只能抽奖一次
	redisConn := apps.App.Redis.Get()
	defer redisConn.Close()
	h, _ := redis.Bool(redisConn.Do("sadd", "draw:record:"+time.Now().Format("20060102"), userId))
	fmt.Println("h:", h)
	if !h {
		return responses.ErrorResponse("一天只有一次抽奖机会"), nil
	}
	prizeList, err := getPrizeList()
	if err != nil {
		return responses.ErrorResponse("未中奖"), nil
	}
	// 随机抽奖
	lotteryPrize := randomLottery(prizeList)
	if lotteryPrize == nil {
		return responses.ErrorResponse("未中奖"), nil
	}

	// 执行抽奖规则
	rule := execRule(userId, lotteryPrize)
	if !rule {
		return responses.ErrorResponse("未中奖"), nil
	}

	// 检测库存
	if !checkStock(lotteryPrize) {
		apps.App.Logger.Info("奖品库存不够", zap.String("prizeName", lotteryPrize.Name))
		return responses.ErrorResponse("未中奖"), nil
	}

	var prize responses.Prize
	prize.Name = lotteryPrize.Name
	prize.Number = 1
	afterLottery(userId, lotteryPrize)
	return responses.SuccessResponse(prize), nil
}

// 执行抽奖规则
func execRule(userId int, prize *models.FormatPrize) bool {
	if prize.Rule == "" {
		return false
	}
	rule := make(map[string]int)
	json.Unmarshal([]byte(prize.Rule), &rule)
	redisConn := apps.App.Redis.Get()
	defer redisConn.Close()
	// 奖品每天最多中奖限制
	if dayLimit, ok := rule["day_limit"]; ok {
		dayNum, err := redis.Int(redisConn.Do("incr", fmt.Sprintf("prize:day:limit:%d:%s", prize.Id, time.Now().Format("20060102"))))
		if err != nil || dayNum > dayLimit {
			apps.App.Logger.Info("超过每天最多中奖规则限制", zap.String("prizeName", prize.Name))
			return false
		}
	}
	//有序集合中元素是用户ID,分数是中奖数量
	if userLimit, ok := rule["user_limit"]; ok {
		userNum, err := redis.Int(redisConn.Do("zscore", fmt.Sprintf("prize:user:num:%d", prize.Id), userId))
		if err != nil || userNum > userLimit {
			apps.App.Logger.Info("超过最多中奖规则限制", zap.String("prizeName", prize.Name), zap.Int("userId", userId))
			return false
		}
	}
	return true
}

// 库存检查
func checkStock(prize *models.FormatPrize) bool {
	rule := make(map[string]int)
	json.Unmarshal([]byte(prize.Rule), &rule)
	unlimit, ok := rule["unlimit"]
	if ok && unlimit > 0 {
		return true
	}
	redisConn := apps.App.Redis.Get()
	defer redisConn.Close()
	stock, err := redis.Int(redisConn.Do("hincrby", "prize:stock", prize.Id, 1))
	apps.App.Logger.Info("stock", zap.Int("stock", stock))
	if err != nil {
		return false
	}
	if stock >= prize.TotalStock {
		return false
	}
	return false
}

// 抽中奖后
func afterLottery(userId int, prize *models.FormatPrize) {
	apps.App.DB.Transaction(func(tx *gorm.DB) error {
		if _, err := prize.UpdateUsedStock(); err != nil {
			return err
		}
		if _, err := addUserPrize(userId, prize.Id); err != nil {
			return err
		}
		redisConn := apps.App.Redis.Get()
		defer redisConn.Close()
		if _, err := redisConn.Do("zincrby", fmt.Sprintf("prize:user:num:%d", prize.Id), 1, userId); err != nil {
			return err
		}
		return nil
	})
}

// 记录用户中奖奖品
func addUserPrize(userId int, prizeId int) (int, error) {
	userPrize := models.UserPrize{UserId: userId, PrizeId: prizeId, Status: 1}
	return userPrize.Create()
}

// 随机抽奖
func randomLottery(prizeList []*models.FormatPrize) *models.FormatPrize {
	rs := rand.NewSource(time.Now().UnixNano())
	r := rand.New(rs)
	num := r.Intn(random_max * 100)
	for _, p := range prizeList {
		if num > p.Low && num <= p.High {
			return p
		}
	}
	return nil
}

// 取参与活动的奖品
func getPrizeList() ([]*models.FormatPrize, error) {
	var formatPrizeList []*models.FormatPrize
	redisConn := apps.App.Redis.Get()
	defer redisConn.Close()
	cache, err := redis.Bytes(redisConn.Do("GET", "prize:list"))
	if err == nil && cache != nil {
		if err := json.Unmarshal(cache, &formatPrizeList); err == nil {
			return formatPrizeList, nil
		}
	}
	var prize models.Prize
	prizeList, err := prize.GetPrizeList()
	if err != nil {
		return nil, err
	}
	high := 0
	for _, prize := range prizeList {
		var fp models.FormatPrize
		fp.Prize = *prize
		num := random_max * prize.LotteryPercent
		fp.Low = high + 1
		fp.High = high + num
		formatPrizeList = append(formatPrizeList, &fp)
	}
	j, err := json.Marshal(formatPrizeList)
	if err == nil {
		redisConn := apps.App.Redis.Get()
		defer redisConn.Close()
		redisConn.Do("SET", "prize:list", j)
	}
	return formatPrizeList, nil
}
