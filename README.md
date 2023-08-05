# AniGamerFloorCompete
巴哈姆特搶樓機器人

## 環境
建議使用VPN，例如:學校的VPN，或是自己的手機網路，以防IP被鎖後，無法使用巴哈姆特

## 使用方法
1. 先設定好config.json
```
{
    "account":"xxxxxx", // 你的帳號
    "password":"xxxxx", // 你的密碼
    "content":"蓋",     // 想要回覆的內容
    "floor":50000,      // 想要搶的樓層
    "targetUrl":"https://forum.gamer.com.tw/C.php?bsn=60076&snA=4981464", // 想要搶的文章網址
    "numRoutine": 1     // 執行緒數量(瀏覽器數量)，建議 1~2
}
```
2. 直接執行或在CLI中執行`win10.exe`，可以加上參數`-rod=show`來可視化執行過程，若需要手動登入則需要可視化
```
win10.exe -rod=show
```

3. (Optional)若電腦無法執行，請自行編譯`main.go`，在[GO官網](https://go.dev/dl/)下載GO後，執行`go build -O main.exe main.go`，即可編譯出`main.exe`，再執行`main.exe`即可