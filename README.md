# Gloria Mailer
Send bulk emails using excel mail merge. (Go Language Version)

利用 Excel 合併的方式寄大量的電子郵件。

## 關於本程式

這個程式目的是為了幫朋友 Gloria 寄一些郵件。
網路上應該可以找到很多很好的解決方案，只是因為個人很習慣直接自己開發一個。

先前寫了一個 [Java](https://github.com/linuzilla/gloria-mailer) 的版本，
但要 Run Java 的 Jar 檔需要先安裝 Java 的 Runtime，索性用 Go 寫一個版本。
一樣有跨平台，而且可以直接編譯成 Windows 上可執行的檔案。

我放一個 Sample 跟 Windows 的執行檔在 [Sample](sample) 目錄下。
另外也有一些 [畫面和說明](sample/README.md)。

## 環境需求及運作原理

這次寫就直接以 Gmail 做為 SMTP 伺服器的概念來開發測試，基本上，其它的 Mail Relay 
應該也都可以跑。

需要準備的資料是一個 Excel 格式的郵遞及套資料的檔，欄位請一定要放在第一列。
另外，當然是郵件的樣版，.eml 格式的，製作方式是先把想要寄的郵件內容（含 Subject) 寄給
自己一份，並下載成 .eml 格式的檔。要填的欄位部份請使用 << >> 括起來，注意，請用半型。
如 Excel 上的欄位叫 email，在郵件樣版上寫 << email >>。

## 安裝與使用

安裝的部份基本上，下載後執行下列指令就可以了。但如果沒有 Go 語言的編譯器的話，要先去下載才能做編譯動作。
           
```shell
go mod tidy
go build
```

設定檔這次改採用 TOML 的格式，它是一個很簡單的格式 。
預設的設定檔名稱為 settings.conf，可以用命令列參數或環境變數更改。

```toml
[main]
# Email 樣版的路徑 (可用命令列參數 override)
template = 'email.eml'
# 寄件者的 Email
sender-email = 'send@from.someone.on.earth'
# 寄件者的名字 (文字)
sender-name = 'John Doe'
# 是不是真得要寄出 Email 還是只是試跑的開關 (可用命令列參數 override)
send-email = false

[debug]
# 檔開 debugging mode 時，寄到這個信箱
send-to = 'debugger@email.address.com'
# debugging mode 的開關 (可用命令列參數 override)
debugging = true

[excel]
# Excel 檔的路徑 (可用命令列參數 override)
file = 'people.xlsx'
# 在 Excel 的欄位中，email 的欄位
email-column = 'Email'
# 在 Excel 的欄位中，名字的欄位
name-column = 'name'

[smtp]
# 以下這幾個是 Mail relay 用的資料，以下是使用 Gmail
host = 'smtp.gmail.com'
port = 587
auth = true
# 如果郵件伺服器需要認證，填上帳號及密碼
# 使用 Gmail 的話，可以到 Google 申請 App password (選擇 App 的部份選 mail)
# 用好之後再刪掉 App password
user = 'user@gmail.com'
password = 'password'
```

## 執行
基本上就是直接執行 go-gloria-mailer，基本上，設定是放在設定檔的。
這個版本有提供命令列參數，(用 --help 可以詳列)，部份功能可以用命令列參數蓋過去，
不必一直修改設定檔。
```sh
go-gloria-mailer
```

## 已知問題

程式目前只分析樣版 email，可以寄出 multipart/alternative 中的 text/plain
跟 text/html 的部份，但目前不支援夾檔。

## 後記

覺得 Go 在處理 Mail 的部份沒有 Java Mail 那麼成熟，但 Go 比 Java 方便的地方是：
當編譯後 Go 是一個簡單的執行檔，Java 則會變成一個 JAR 檔，需要 Java Runtime
來執行這個 JAR 檔。
