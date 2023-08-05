package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
)

var once sync.Once
var tickerWaitTime int32 = 5000

func main() {
	configFile := "config.json"
	config, _ := readConfig(configFile)

	var account string = config.Account
	var password string = config.Password
	var inputText string = config.Content
	var targetFloor int = config.Floor - 1
	var targetUrl string = config.TargetUrl
	var numRoutine int = config.NumRoutine

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(time.Duration(tickerWaitTime) * time.Millisecond)
	defer ticker.Stop()

	for i := 0; i < numRoutine; i++ {
		go func() {
			url := launcher.New().MustLaunch()
			browser := rod.New().ControlURL(url).MustConnect()
			page := browser.MustPage(`https://user.gamer.com.tw`)
			page.MustWaitLoad()
			
			Login(page, account, password)
			FightForTop(ctx, cancel, page, inputText, targetUrl, targetFloor, ticker)
		}()
	}

	<-ctx.Done()
	fmt.Println(`主程式結束`)
}

func Login(page *rod.Page, account string, password string) {

	page.MustElement(`#BH-top-data > div.TOP-my.TOP-nologin > ul > li:nth-child(1) > a`).MustClick()
	page.MustWaitLoad()

	err := rod.Try(func() {
		page.Timeout(3 * time.Second).MustElement("#form-login > div:nth-child(4) > div > div > div > iframe")
	})

	if errors.Is(err, context.DeadlineExceeded) {
		fmt.Println(`無機器人驗證`)
		page.MustElement(`#form-login > input:nth-child(1)`).MustInput(account)
		page.MustElement(`#form-login > div.password-box > input`).MustInput(password)
		var pageUrl string = page.MustInfo().URL
		for {
			page.MustElement(`#btn-login`).MustDoubleClick()
			time.Sleep(1 * time.Second)
			if pageUrl != page.MustInfo().URL {
				break
			}
		}
		fmt.Println(`登入成功`)
		page.MustWaitLoad()

	} else if err != nil {
		fmt.Println(`其他錯誤`)

	} else {
		fmt.Println(`有機器人驗證，請手動登入或換成手機網路再次嘗試`)
		var pageUrl string = page.MustInfo().URL
		for {
			time.Sleep(1 * time.Second)
			if pageUrl != page.MustInfo().URL {
				break
			}
		}
	
	}
	

	
}

func FightForTop(ctx context.Context, cancel context.CancelFunc, page *rod.Page, inputText string, targetUrl string, targetFloor int, ticker *time.Ticker) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			page.MustNavigate(targetUrl + `&last=1#down`)
			page.MustWaitLoad()

			elements := page.MustElements("a.floor")
			re := regexp.MustCompile(`\d+`)
			var curFloor string = re.FindString(elements[len(elements)-1].MustText())

			tickerWaitTimeCurrent := atomic.LoadInt32(&tickerWaitTime)

			switch curFloorInt, _ := strconv.Atoi(curFloor); {

			// 準備搶樓
			case targetFloor-curFloorInt > 100:
				atomic.AddInt32(&tickerWaitTime, 30*1000)
				ticker = time.NewTicker(time.Duration(tickerWaitTimeCurrent) * time.Millisecond)
				fmt.Println(`目標樓層:` + strconv.Itoa(targetFloor+1) + `, 目前樓層:` + curFloor + `, 時間:` + time.Now().Format("2006-01-02 15:04:05.999"))
			case targetFloor-curFloorInt > 50:
				atomic.AddInt32(&tickerWaitTime, 15*1000)
				ticker = time.NewTicker(time.Duration(tickerWaitTimeCurrent) * time.Millisecond)
				fmt.Println(`目標樓層:` + strconv.Itoa(targetFloor+1) + `, 目前樓層:` + curFloor + `, 時間:` + time.Now().Format("2006-01-02 15:04:05.999"))
			case targetFloor-curFloorInt > 25:
				atomic.AddInt32(&tickerWaitTime, 7.5*1000)
				ticker = time.NewTicker(time.Duration(tickerWaitTimeCurrent) * time.Millisecond)
				fmt.Println(`目標樓層:` + strconv.Itoa(targetFloor+1) + `, 目前樓層:` + curFloor + `, 時間:` + time.Now().Format("2006-01-02 15:04:05.999"))
			case targetFloor-curFloorInt > 12:
				atomic.AddInt32(&tickerWaitTime, 4*1000)
				ticker = time.NewTicker(time.Duration(tickerWaitTimeCurrent) * time.Millisecond)
				fmt.Println(`目標樓層:` + strconv.Itoa(targetFloor+1) + `, 目前樓層:` + curFloor + `, 時間:` + time.Now().Format("2006-01-02 15:04:05.999"))
			case targetFloor-curFloorInt > 5:
				atomic.AddInt32(&tickerWaitTime, 1*1000)
				ticker = time.NewTicker(time.Duration(tickerWaitTimeCurrent) * time.Millisecond)
				fmt.Println(`目標樓層:` + strconv.Itoa(targetFloor+1) + `, 目前樓層:` + curFloor + `, 時間:` + time.Now().Format("2006-01-02 15:04:05.999"))
			case targetFloor-curFloorInt > 0:
				atomic.AddInt32(&tickerWaitTime, 0.5*1000)
				ticker = time.NewTicker(time.Duration(tickerWaitTimeCurrent) * time.Millisecond)
				fmt.Println(`目標樓層:` + strconv.Itoa(targetFloor+1) + `, 目前樓層:` + curFloor + `, 時間:` + time.Now().Format("2006-01-02 15:04:05.999"))
			
			// 搶樓失敗
			case curFloorInt > targetFloor:
				fmt.Println(`目前樓層:` + curFloor + `, 搶樓失敗, 時間:` + time.Now().Format("2006-01-02 15:04:05.999"))
				cancel()
			
			// 搶樓成功
			case curFloorInt == targetFloor:
				once.Do(func() {
					page.MustElement(`#postTips`).MustDoubleClick()
					page.MustElement(`#editor`).MustFrame().MustElement("p").MustInput(inputText)
					page.MustElement(".btn--sm.btn--send.btn--normal").MustDoubleClick()
					page.MustElement(`.btn-insert.btn-primary`).MustClick()
					page.MustElement(`.btn-insert.btn-primary`).MustClick()
					fmt.Println(`目前樓層:` + curFloor + `, 搶樓成功, 時間:` + time.Now().Format("2006-01-02 15:04:05.999"))
					time.Sleep(5 * time.Second)
					page.MustWaitLoad()
					cancel()
				})
			}
		}
	}
}

type Config struct {
	Account    string `json:"account"`
	Password   string `json:"password"`
	Content    string `json:"content"`
	Floor      int    `json:"floor"`
	TargetUrl  string `json:"targetUrl"`
	NumRoutine int    `json:"numRoutine"`
}

func readConfig(filePath string) (Config, error) {
	var config Config

	configFile, err := os.Open(filePath)
	if err != nil {
		return config, err
	}
	defer configFile.Close()

	decoder := json.NewDecoder(configFile)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}
