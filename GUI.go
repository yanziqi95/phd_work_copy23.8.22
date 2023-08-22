package main

import (
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
	"math/big"
)

func GUI() {
	myApp := app.New()
	myWindow := myApp.NewWindow("Repu-wallet")

	//TAB 1
	entry1 := widget.NewEntry()
	textArea1 := widget.NewMultiLineEntry()

	submit := &widget.Form{
		Items: []*widget.FormItem{ // we can specify items in the constructor
			{Text: "Target", Widget: entry1}, {Text: "Review", Widget: textArea1}},
		OnSubmit: func() {
			//submit logic
			//log.Println("Form submitted:", entry.Text)
			//log.Println("multiline:", textArea.Text)
			//submitReivew("1", entry1.Text, textArea1.Text)
			submitReview(entry1.Text, textArea1.Text, 0)
			dialog.ShowInformation("Congrats", "Reviews is submitted", myWindow)
		},
	}
	label1 := widget.NewLabel("Wallet address:")
	walletAdd := widget.NewLabel("N/A")
	//label2 := widget.NewLabel("Private key:")
	//privateKey := widget.NewLabel("N/A")

	label3 := widget.NewLabel("Nickname :")
	nickname := widget.NewLabel("N/A")
	//创建私钥
	register := widget.NewLabel("Please enter file name:")
	//保存私钥并创建
	fileEntry := widget.NewEntry()
	registerAndSave := widget.NewButton("generate & save", func() {
		//wallets, _ := NewWallets()
		//address, private, _ := wallets.CreateWallet()
		wallet, priKey, pubKey := NewWallet()
		add := string(wallet.GetAddress())

		walletAdd.SetText(add)
		//publicKey.SetText(public)
		//privateKey.SetText(priKey.D.String())
		nickname.SetText(fileEntry.Text)
		fmt.Printf("pri:%x,pub:%x,wallet:%x", priKey.D, pubKey, add)
		//分别保存私钥，公钥，钱包地址
		SaveToFile(priKey.D.String(), fileEntry.Text+"Pri")
		SaveToFile(string(pubKey), fileEntry.Text+"Pub")
		SaveToFile(add, fileEntry.Text+"Wal")

		//wallets.SaveToFile(address, fileEntry.Text)
		dialog.ShowInformation("Success", add, myWindow)
	})
	//registerAndSave := widget.NewButton("generate & save", func() {
	//	add.SetText("789")
	//	newFileName := "./wallet" + "789"
	//	file, err := os.Create(newFileName)
	//	if err != nil {
	//		fmt.Printf(err.Error())
	//		return
	//	}
	//	defer file.Close()
	//	file.WriteString("789")
	//	dialog.ShowInformation("congrats", "export wallet success", myWindow)
	//})

	registerContainer := widget.NewVBox(register, fileEntry, registerAndSave)
	infoContainer := widget.NewVBox(label1, walletAdd, label3, nickname)

	//登录容器
	loginLable := widget.NewLabel("Load private key from following file:")
	loginEntry := widget.NewEntry()
	loginBtn := widget.NewButton("Load private key", func() {
		priData, err := LoadFileToHex(loginEntry.Text + "Pri")
		pubData, err := LoadFileToHex(loginEntry.Text + "Pub")
		walData, err := LoadFileToHex(loginEntry.Text + "Wal")
		if err != nil {
			dialog.ShowInformation("Failed", err.Error(), myWindow)
		}
		private, _ := new(big.Int).SetString(string(priData), 10)
		fmt.Printf("prikey:%x", private)
		fmt.Printf("pubkey:%x", pubData)
		fmt.Printf("walAdd:%x", walData)
		//privateKey.SetText(private.String())
		nickname.SetText(loginEntry.Text)
		walletAdd.SetText(string(walData))

	})
	loginContainer := widget.NewVBox(loginLable, loginEntry, loginBtn)
	//login := &widget.Form{
	//	Items: []*widget.FormItem{ // we can specify items in the constructor
	//		{Text: "Login", Widget: entry}, {Text: "Review", Widget: textArea}},
	//	OnSubmit: func() {
	//		//submit logic
	//		//log.Println("Form submitted:", entry.Text)
	//		//log.Println("multiline:", textArea.Text)
	//		submitReivew("1", entry.Text, textArea.Text)
	//		dialog.ShowInformation("Congrats", "You are using following wallet", myWindow)
	//	},
	//}

	//检查评论容器
	checkLable1 := widget.NewLabel("Enter the target")
	checkEntry1 := widget.NewEntry()
	checkBtn1 := widget.NewButton("Show reviews", func() {
		//从网站请求评论列表
		fmt.Printf("Reviews")
	})
	checkBtn2 := widget.NewButton("Check reviews", func() {
		//比对网站和区块链评论
		if checkReview(checkEntry1.Text) {
			dialog.ShowInformation("Congrats", "Reviews is integrated", myWindow)
		} else {
			dialog.ShowInformation("Sorry", "Someone tampered reviews", myWindow)
		}
		//fmt.Printf("Check reviews")
	})
	checkReviewContainer := widget.NewVBox(checkLable1, checkEntry1, checkBtn1, checkBtn2)

	//交易容器
	txLabel1 := widget.NewLabel("Target")
	txEntry1 := widget.NewEntry()
	txBtn1 := widget.NewButton("send", func() {
		fmt.Printf("send tokens")
	})
	txLabel2 := widget.NewLabel("Balance:")
	balance := widget.NewLabel("N/A")
	balEntry := widget.NewEntry()
	checkBalance := widget.NewButton("Check balance", func() {
		//fmt.Printf("check balance")
		addr := balEntry.Text
		fmt.Printf("opening %x and requesting balance", addr)
		bal := balanceReq(addr)
		balance.SetText(bal)
	})
	transactionContainer := widget.NewVBox(txLabel1, txEntry1, txBtn1, balance, txLabel2, balance, balEntry, checkBalance)

	//同步容器
	syncBtn1 := widget.NewButton("Sync. with full node", func() {
		fmt.Printf("snynchronizing with full node......")
	})
	syncContainer := widget.NewVBox(syncBtn1)
	tabs := container.NewAppTabs(
		container.NewTabItem("INFO", infoContainer),
		container.NewTabItem("Submit reviews", submit),
		container.NewTabItem("Check reviews", checkReviewContainer),
		container.NewTabItem("transaction", transactionContainer),
		container.NewTabItem("Sync.", syncContainer),
		container.NewTabItem("Register", registerContainer),
		container.NewTabItem("Login", loginContainer),
		//container.NewTabItem("Login", login),
		//container.NewTabItem("Sync",sync)
		container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("ReadMe: You can use this application to \n 1. Submit reviews in a democratic way \n 2. Earn tokens \n 3. Trade tokens")),
	)

	tabs.SetTabLocation(container.TabLocationLeading)
	myWindow.Resize(fyne.Size{800, 400})
	myWindow.SetContent(tabs)
	myWindow.ShowAndRun()
}
