package main

import (
	"crypto/md5"
	"embed"
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"howett.net/plist"

	"github.com/unknwon/i18n"
)

var version = 221

var hosts = []string{"https://idea.jeter.eu.org", "http://129.154.205.7:7191", "http://jetbra.serv00.net:7191", "http://ba.serv00.net:7191"}
var host = hosts[0]
var githubPath = "https://ghp.ci/https://github.com/kingparks/jetbra-activate/releases/download/latest/"
var err error

var green = "\033[32m%s\033[0m\n"
var yellow = "\033[33m%s\033[0m\n"
var hGreen = "\033[1;32m%s\033[0m"
var dGreen = "\033[4;32m%s\033[0m\n"
var red = "\033[31m%s\033[0m\n"
var defaultColor = "%s"
var lang, _ = getLocale()
var deviceID = getMacMD5()
var client = Client{Hosts: hosts}

//go:embed all:script
var scriptFS embed.FS

//go:embed all:locales
var localeFS embed.FS

type Tr struct {
	i18n.Locale
}

var tr *Tr

func main() {
	language := flag.String("l", lang, "set language, eg: zh, en, nl, ru, hu, tr")
	flag.Parse()

	localeFileEn, _ := localeFS.ReadFile("locales/en.ini")
	_ = i18n.SetMessage("en", localeFileEn)
	localeFileNl, _ := localeFS.ReadFile("locales/nl.ini")
	_ = i18n.SetMessage("nl", localeFileNl)
	localeFileRu, _ := localeFS.ReadFile("locales/ru.ini")
	_ = i18n.SetMessage("ru", localeFileRu)
	localeFileHu, _ := localeFS.ReadFile("locales/hu.ini")
	_ = i18n.SetMessage("hu", localeFileHu)
	localeFileTr, _ := localeFS.ReadFile("locales/tr.ini")
	_ = i18n.SetMessage("tr", localeFileTr)
	localeFileEs, _ := localeFS.ReadFile("locales/es.ini")
	_ = i18n.SetMessage("es", localeFileEs)
	lang = *language
	switch lang {
	case "zh":
		tr = &Tr{Locale: i18n.Locale{Lang: "zh"}}
	case "nl":
		tr = &Tr{Locale: i18n.Locale{Lang: "nl"}}
	case "ru":
		tr = &Tr{Locale: i18n.Locale{Lang: "ru"}}
	case "hu":
		tr = &Tr{Locale: i18n.Locale{Lang: "hu"}}
	case "tr":
		tr = &Tr{Locale: i18n.Locale{Lang: "tr"}}
	case "es":
		tr = &Tr{Locale: i18n.Locale{Lang: "es"}}
	default:
		tr = &Tr{Locale: i18n.Locale{Lang: "en"}}
	}

	fmt.Printf(green, tr.Tr("IntelliJ 授权")+` v`+strings.Join(strings.Split(fmt.Sprint(version), ""), "."))
	client.SetProxy(lang)
	sCount, sPayCount, _, _, exp := client.GetMyInfo(deviceID)
	fmt.Printf(green, tr.Tr("设备码")+":"+deviceID)
	expTime, _ := time.ParseInLocation("2006-01-02 15:04:05", exp, time.Local)
	fmt.Printf(green, tr.Tr("付费到期时间")+":"+exp)
	fmt.Printf("\033[32m%s\033[0m\u001B[1;32m %s \u001B[0m\033[32m%s\033[0m\u001B[1;32m %s \u001B[0m\u001B[32m%s\u001B[0m\n",
		tr.Tr("推广命令：(已推广"), sCount, tr.Tr("人,推广已付费"), sPayCount, tr.Tr("人；每推广10人或推广付费2人可获得一年授权)"))
	fmt.Printf(hGreen, "bash <(curl "+githubPath+"install.sh) "+deviceID+"\n")

	printAD()
	checkUpdate(version)
	fmt.Println()
	fmt.Printf(defaultColor, tr.Tr("选择要授权的产品："))
	jbProduct := []string{"IntelliJ IDEA", "CLion", "PhpStorm", "Goland", "PyCharm", "WebStorm", "Rider", "DataGrip", "DataSpell"}
	jbProductChoice := []string{"idea", "clion", "phpstorm", "goland", "pycharm", "webstorm", "rider", "datagrip", "dataspell"}
	for i, v := range jbProduct {
		fmt.Printf(hGreen, fmt.Sprintf("%d. %s\t", i+1, v))
	}
	fmt.Println()
	fmt.Print(tr.Tr("请输入产品编号（直接回车默认为1）："))
	productIndex := 1
	_, _ = fmt.Scanln(&productIndex)
	if productIndex < 1 || productIndex > len(jbProduct) {
		fmt.Println(tr.Tr("输入有误"))
		return
	}
	fmt.Println(tr.Tr("选择的产品为：") + jbProduct[productIndex-1])
	fmt.Println()

	periodIndex := 2
	// 到期了
	_ = []time.Duration{100 * 24 * time.Hour, 367 * 24 * time.Hour}
	if expTime.Before(time.Now()) {
		fmt.Printf(defaultColor, tr.Tr("选择有效期："))
		jbPeriod := []string{tr.Tr("两天(免费)"), tr.Tr("一年(购买)")}
		for i, v := range jbPeriod {
			fmt.Printf(hGreen, fmt.Sprintf("%d. %s\t", i+1, v))
		}
		fmt.Println()
		fmt.Printf("%s", tr.Tr("请输入有效期编号（直接回车默认为1）："))
		_, _ = fmt.Scanln(&periodIndex)
		if periodIndex < 1 || periodIndex > len(jbPeriod) {
			fmt.Println(tr.Tr("输入有误"))
			return
		}
		fmt.Println(tr.Tr("选择的有效期为：") + jbPeriod[periodIndex-1])
		fmt.Println()
	}

	lic := ""
	for i := 0; i < 50; i++ {
		if i == 20 {
			Clean()
			Active(jbProductChoice[productIndex-1])
		}
		time.Sleep(20 * time.Millisecond)
		h := strings.Repeat("=", i) + strings.Repeat(" ", 49-i)
		fmt.Printf("\r%.0f%%[%s]", float64(i)/49*100, h)
	}
	fmt.Println()
	fmt.Println()

	switch periodIndex {
	case 1:
		isOk, result := client.GetLic(jbProductChoice[productIndex-1], periodIndex-1)
		if !isOk {
			fmt.Printf(red, result)
			return
		}
		lic = result
	case 2:
		// 没到期
		if expTime.After(time.Now()) {
			isOk, result := client.GetLic(jbProductChoice[productIndex-1], periodIndex-1)
			if !isOk {
				fmt.Printf(red, result)
				return
			}
			lic = result
			fmt.Println()
			goto Process
		}
		// 到期了
		payUrl, orderID := client.GetPayUrl()
		isCopyText := ""
		errClip := clipboard.WriteAll(payUrl)
		if errClip == nil {
			isCopyText = tr.Tr("（已复制到剪贴板）")
		}
		fmt.Println(tr.Tr("付费已到期,捐赠以获取一年期授权") + isCopyText)
		fmt.Printf(dGreen, payUrl)
		fmt.Println(tr.Tr("捐赠完成后请回车"))
		//检测控制台回车
	checkPay:
		_, _ = fmt.Scanln()
		isPay := client.PayCheck(orderID, deviceID)
		if !isPay {
			fmt.Println(tr.Tr("未捐赠,请捐赠完成后回车"))
			goto checkPay
		}
		isOk, result := client.GetLic(jbProductChoice[productIndex-1], periodIndex-1)
		if !isOk {
			fmt.Printf(red, result)
			return
		}
		fmt.Println()
	}

Process:
	isCopyText := ""
	err = clipboard.WriteAll(lic)
	if err == nil {
		isCopyText = tr.Tr("（已复制到剪贴板）")
	}
	fmt.Printf(yellow, tr.Tr("首次执行请重启IDE，然后填入下面授权码；非首次执行直接填入下面授权码即可")+isCopyText)
	switch runtime.GOOS {
	case "windows":
		_ = exec.Command("taskkill", "/IM", jbProductChoice[productIndex-1]+".exe", "/F").Run()
		_ = exec.Command("taskkill", "/IM", jbProductChoice[productIndex-1]+"64.exe", "/F").Run()
	case "darwin":
		_ = exec.Command("killall", "-9", jbProductChoice[productIndex-1]).Run()
	case "linux":
		_ = exec.Command("killall", "-9", "java").Run()
	}
	fmt.Println()
	fmt.Printf(hGreen, lic)
	fmt.Println()
	for i := 0; i < 4; i++ {
		_, _ = fmt.Scanln()
	}
}

func getMacMD5() string {
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("err:", err)
		return ""
	}
	var macAddress []string
	for _, inter := range interfaces {
		// 大于en6的排除
		if strings.HasPrefix(inter.Name, "en") {
			numStr := inter.Name[2:]
			num, _ := strconv.Atoi(numStr)
			if num > 6 {
				continue
			}
		}
		if strings.HasPrefix(inter.Name, "en") || strings.HasPrefix(inter.Name, "Ethernet") || strings.HasPrefix(inter.Name, "以太网") || strings.HasPrefix(inter.Name, "WLAN") {
			macAddress = append(macAddress, inter.HardwareAddr.String())
		}
	}
	sort.Strings(macAddress)
	return fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(macAddress, ","))))
}

func printAD() {
	ad := client.GetAD()
	if len(ad) == 0 {
		return
	}
	fmt.Printf(yellow, ad)
}

func checkUpdate(version int) {
	upUrl := client.CheckVersion(fmt.Sprint(version))
	if upUrl == "" {
		return
	}
	isCopyText := ""
	installCmd := `bash -c "$(curl -fsSL ` + githubPath + `install.sh)"`
	errClip := clipboard.WriteAll(installCmd)
	if errClip == nil {
		isCopyText = tr.Tr("（已复制到剪贴板）")
	}
	switch runtime.GOOS {
	case "windows":
		fmt.Printf(red, tr.Tr("有新版本，请关闭本窗口，将下面命令粘贴到GitBash窗口执行")+isCopyText+`：`)
	default:
		fmt.Printf(red, tr.Tr("有新版本，请关闭本窗口，将下面命令粘贴到新终端窗口执行")+isCopyText+`：`)
	}
	fmt.Printf(hGreen, installCmd)
	_, _ = fmt.Scanln()
	os.Exit(0)
	return
}

// 获取推广人
func getPromotion() (promotion string) {
	b, _ := os.ReadFile(os.Getenv("HOME") + "/.jetbrarc")
	promotion = strings.TrimSpace(string(b))
	if len(promotion) == 0 {
		if len(os.Args) > 1 {
			promotion = os.Args[1]
		}
	}
	return
}

func getLocale() (langRes, locRes string) {
	osHost := runtime.GOOS
	langRes = "en"
	locRes = "US"
	switch osHost {
	case "windows":
		// Exec powershell Get-Culture on Windows.
		cmd := exec.Command("powershell", "Get-Culture | select -exp Name")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "-")
			langRes = langLoc[0]
			langRes = strings.Split(langRes, "-")[0]
			locRes = langLoc[1]
			return
		}
	case "darwin":
		// Exec shell Get-Culture on MacOS.
		cmd := exec.Command("sh", "osascript -e 'user locale of (get system info)'")
		output, err := cmd.Output()
		if err == nil {
			langLocRaw := strings.TrimSpace(string(output))
			langLoc := strings.Split(langLocRaw, "_")
			langRes = langLoc[0]
			langRes = strings.Split(langRes, "-")[0]
			if len(langLoc) == 1 {
				return
			}
			locRes = langLoc[1]
			return
		}
		plistB, err := os.ReadFile(os.Getenv("HOME") + "/Library/Preferences/.GlobalPreferences.plist")
		if err != nil {
			panic(err)
		}
		var a map[string]interface{}
		_, err = plist.Unmarshal(plistB, &a)
		if err != nil {
			panic(err)
		}
		langLocRaw := a["AppleLocale"].(string)
		langLoc := strings.Split(langLocRaw, "_")
		langRes = langLoc[0]
		langRes = strings.Split(langRes, "-")[0]
		if len(langLoc) == 1 {
			return
		}
		locRes = langLoc[1]
		return
	case "linux":
		envlang, ok := os.LookupEnv("LANG")
		if ok {
			langLocRaw := strings.TrimSpace(envlang)
			langLocRaw = strings.Split(envlang, ".")[0]
			langLoc := strings.Split(langLocRaw, "_")
			langRes = langLoc[0]
			langRes = strings.Split(langRes, "-")[0]
			if len(langLoc) == 1 {
				return
			}
			locRes = langLoc[1]
			return
		}
	}
	if langRes == "" {
		langLocRaw := os.Getenv("LC_CTYPE")
		langLocRaw = strings.Split(langLocRaw, ".")[0]
		langLoc := strings.Split(langLocRaw, "_")
		langRes = langLoc[0]
		langRes = strings.Split(langRes, "-")[0]
		if len(langLoc) == 1 {
			return
		}
		locRes = langLoc[1]
		return
	}
	return
}
