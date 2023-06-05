package main

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"
)

func main() {
	var res string
	options := append(
		chromedp.DefaultExecAllocatorOptions[:],
		//不检查默认浏览器
		chromedp.NoDefaultBrowserCheck,
		//禁用chrome的handless(禁用无窗口模式，即开启窗口模式)
		chromedp.Flag("headless", false),
		//开启图像界面
		chromedp.Flag("blink-settings", "imageEnabled=true"),
		//忽略错误
		chromedp.Flag("ignore-certificate-errors", true),
		//禁用网络安全标志
		chromedp.Flag("disable-web-security", true),
		//开启插件支持
		chromedp.Flag("disable-extensions", true),
		chromedp.Flag("disable-default-apps", true),
		//设置网站不是首次运行
		chromedp.NoFirstRun,
		//设置窗口大小
		chromedp.WindowSize(1900, 1024),
	)
	allocator, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	print(cancel)
	ctx, cancel := chromedp.NewContext(
		allocator,
		chromedp.WithLogf(log.Printf),
	)
	//设置超时时间
	ctx, cancel = context.WithTimeout(ctx, 10*time.Second)
	//运行chromedp，操作浏览器
	//chromedp监听网页上弹出alert对话框
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		//time.Sleep(time.Second)
		if ev, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			fmt.Println("closing alert:", ev.Message)
			go func() {
				fmt.Println("进入了协程")
				//自动关闭alert对话框
				if err := chromedp.Run(ctx,
					//注释掉下一行可以更清楚地看到效果
					chromedp.Sleep(time.Second),
					page.HandleJavaScriptDialog(true),
				); err != nil {
					panic(err)
				}
			}()
		}
	})
	chromedp.Run(
		ctx,
		//跳转到目标页面
		chromedp.Navigate("http://jwgl.hist.edu.cn/cas/login.action"),
		//输入账号密码
		//等页面加载出该元素后再执行操作
		chromedp.WaitVisible(`document.getElementById("username1")`, chromedp.ByJSPath),
		chromedp.SetValue(`document.getElementById("username1")`, "xxx", chromedp.ByJSPath),
		chromedp.SetValue(`document.getElementById("username")`, "xxx", chromedp.ByJSPath),
		//填入密码
		chromedp.WaitVisible(`document.getElementById("password1")`, chromedp.ByJSPath),
		chromedp.SetValue(`document.getElementById("password1")`, "xxx", chromedp.ByJSPath),
		chromedp.SetValue(`document.getElementById("password")`, "xxx", chromedp.ByJSPath),
		//点击登录按钮
		chromedp.Click(`document.getElementById("login")`, chromedp.ByJSPath),
		//点击教学评价
		chromedp.Click(`document.querySelectorAll("#normal_use_menu li")[0]`, chromedp.ByJSPath),
		//等待显示
		chromedp.WaitVisible(`document.getElementById("frmDesk").contentDocument.getElementById("frame_1").contentDocument.getElementById("selPJLC")`, chromedp.ByJSPath),
		chromedp.Click(`document.getElementById("frmDesk").contentDocument.getElementById("frame_1").contentDocument.querySelector("#selPJLC ")`, chromedp.ByJSPath),
		chromedp.SetValue(`document.getElementById("frmDesk").contentDocument.getElementById("frame_1").contentDocument.querySelector("#selPJLC ")`, "{\"pjfsbz\":\"0\",\"lcjc\":\"第二阶段评价\",\"sfkpsj\":\"1\",\"sfzbpj\":\"1\",\"xn\":\"2022\",\"xq_m\":\"1\",\"lcqc\":\"2022-2023学年第二学期\",\"jsrq\":\"2023-06-06 23:59\",\"sfwjpj\":\"1\",\"lcdm\":\"2022102\",\"qsrq\":\"2023-06-01 00:00\"}", chromedp.ByJSPath),
		//点击评价按钮
		chromedp.WaitVisible(`document.getElementById("frmDesk").contentDocument.getElementById("frame_1").contentDocument.getElementById("frmReport").contentDocument.querySelector("#tr0_wjdc > a")`, chromedp.ByJSPath),
		chromedp.Click(`document.getElementById("frmDesk").contentDocument.getElementById("frame_1").contentDocument.getElementById("frmReport").contentDocument.querySelector("#tr0_wjdc > a")`, chromedp.ByJSPath),
		//第一部分
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			fmt.Println("进入第一部分")
			for i := 0; i < 10; i++ {
				target := fmt.Sprintf("document.querySelector(\"#dialog-frame\").contentDocument.querySelector(\"#wdt_0_%v_1\")", i)
				err = chromedp.WaitVisible(target, chromedp.ByJSPath).Do(ctx)
				if err != nil {
					fmt.Println(fmt.Sprintf("第一部分，第%v题等待其按钮出现失败， 具体标签为：%v", i, target))
					return err
				}
				err = chromedp.Click(target, chromedp.ByJSPath).Do(ctx)
				if err != nil {
					fmt.Println(fmt.Sprintf("第二部分，第%v题等待其按钮出现失败，具体标签为：%v", i, target))
					return err
				}

			}
			return err
		}),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			fmt.Println("进入第二部分")
			for i := 0; i < 11; i++ {

				target := fmt.Sprintf("document.querySelector(\"#dialog-frame\").contentDocument.querySelector(\"#radio%v_0\")", i)
				if i == 3 {
					target = fmt.Sprintf("document.querySelector(\"#dialog-frame\").contentDocument.querySelector(\"#radio%v_2\")", i)
				}
				err = chromedp.WaitVisible(target, chromedp.ByJSPath).Do(ctx)
				if err != nil {
					fmt.Println(fmt.Sprintf("第一部分，第%v题等待其按钮出现失败， 具体标签为：%v", i, target))
					return err
				}
				err = chromedp.Click(target, chromedp.ByJSPath).Do(ctx)
				if err != nil {
					fmt.Println(fmt.Sprintf("第二部分，第%v题等待其按钮出现失败，具体标签为：%v", i, target))
					return err
				}

			}
			return err
		}),
		//第三部分是一个填空题
		chromedp.WaitVisible("document.querySelector(\"#dialog-frame\").contentDocument.querySelector(\"#area11\")", chromedp.ByJSPath),
		chromedp.SetValue("document.querySelector(\"#dialog-frame\").contentDocument.querySelector(\"#area11\")", "好", chromedp.ByJSPath),
		chromedp.Sleep(time.Second),
		//点击保存按钮
		chromedp.ActionFunc(func(ctx context.Context) error {
			fmt.Println("1234564")
			return nil
		}),
		//保险起见，点击取消按钮
		//chromedp.Click("document.querySelector(\"#dialog-frame\").contentDocument.querySelector(\"#butSave\").nextElementSibling", chromedp.ByJSPath),
		//下面是点击提交按钮
		chromedp.Click("document.querySelector(\"#dialog-frame\").contentDocument.querySelector(\"#butSave\")", chromedp.ByJSPath),
	)


	fmt.Println("res" + res)
}

func ExampleListenTarget_acceptAlert() {
	//内置http测试服务器，用于在网页上显示alert按钮
	ts := httptest.NewServer(writeHTML(`
<input id='alert' type='button' value='alert' onclick='alert("alert text");'/>please等5秒后，自动点击Alert,并自动关闭alert对话框。
    `))
	defer ts.Close()
	//fmt.Println(ts.URL)
	//增加选项，允许chrome窗口显示出来
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("hide-scrollbars", false),
		chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	//创建chrome窗口
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	//chromedp监听网页上弹出alert对话框
	chromedp.ListenTarget(ctx, func(ev interface{}) {
		if ev, ok := ev.(*page.EventJavascriptDialogOpening); ok {
			fmt.Println("closing alert:", ev.Message)
			go func() {
				//自动关闭alert对话框
				if err := chromedp.Run(ctx,
					//注释掉下一行可以更清楚地看到效果
					page.HandleJavaScriptDialog(true),
				); err != nil {
					panic(err)
				}
			}()
		}
	})

	if err := chromedp.Run(ctx,
		chromedp.Navigate(ts.URL),
		chromedp.Sleep(5*time.Second),
		//自动点击页面上的alert按钮，弹出alert对话框
		chromedp.Click("#alert", chromedp.ByID),
	); err != nil {
		panic(err)
	}

	// Output:
	// closing alert: alert text
}
func writeHTML(content string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//在这里设置utf-8,避免乱码
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		io.WriteString(w, strings.TrimSpace(content))
	})
}
