package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

const (
	title       = "いきものAZ エクスポーター"
	description = `これはいきもの特化型SNS「いきものAZ」のユーザデータの保存を行うツールです。
以下のように使ってください:
  1. いきものAZのマイページにいき、アドレス欄からURLをコピペする
  2. 「保存先フォルダを選ぶ」をクリック (いくつかファイルとフォルダができます)
  3. データを保存するフォルダを選択して"Open"を押す (英語ですみません)
  4. 「エクスポート開始」をクリック
  5. 待つ
  6. 保存先のindex.htmlをダブルクリック`

	captionChooseButton = "保存先フォルダを選ぶ"
	captionChooser      = "フォルダを選んでください"
	captionExportButton = "エクスポート開始"
)

type state struct {
	targetPath binding.String
	mypageURL  binding.String
	status     binding.String
	dotCount   int
}

type ExportFn = func(string, string) error

func exportButtonClicked(s *state, fn ExportFn, b1, b2 *widget.Button) {
	b1.Disable()
	b2.Disable()

	path, err := s.targetPath.Get()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}
	url, err := s.mypageURL.Get()
	if err != nil {
		fmt.Printf("%v\n", err)
		return
	}

	c := make(chan bool)
	go calculateStatus(s, c)
	err = fn(path, url)
	c <- true
	if err != nil {
		s.status.Set(fmt.Sprintf("%v", err))
	}

	b1.Enable()
	b2.Enable()
}

func prepareContent(w fyne.Window, s *state, fn ExportFn) *fyne.Container {
	// 保存先
	targetPathLabel := widget.NewLabelWithData(s.targetPath)

	// 保存先
	mypageURLEntry := widget.NewEntryWithData(s.mypageURL)
	mypageURLEntry.SetPlaceHolder("ここにマイページのアドレスを貼ってください")

	// 保存先選択ダイアログ
	chooserFn := func(uri fyne.ListableURI, err error) {
		if uri != nil {
			s.targetPath.Set(uri.Path())
		}
	}
	chooser := dialog.NewFolderOpen(chooserFn, w)
	// 保存先選択ボタン
	chooseButtonFn := func() {
		chooser.Show()
	}
	chooseButton := widget.NewButton(captionChooseButton, chooseButtonFn)

	// エクスポート開始ボタン
	exportButton := widget.NewButton("", func() {})
	exportButtonFn := func() {
		// fn := func(userID, targetDir string) error {
		// 	time.Sleep(10 * time.Second)
		// 	return nil
		// }
		exportButtonClicked(s, fn, chooseButton, exportButton)
	}
	exportButton = widget.NewButton(captionExportButton, exportButtonFn)

	// 実行ステータス
	statusLabel := widget.NewLabel("")
	statusLabel.Bind(s.status)

	mypageURLContainer := fyne.NewContainerWithLayout(
		layout.NewFormLayout(),
		widget.NewLabel("マイページURL:"),
		mypageURLEntry,
	)
	targetPathContainer := fyne.NewContainerWithLayout(
		layout.NewFormLayout(),
		widget.NewLabel("結果保存先:"),
		targetPathLabel,
	)
	statusContainer := fyne.NewContainerWithLayout(
		layout.NewFormLayout(),
		widget.NewLabel("状態:"),
		statusLabel,
	)
	buttonsContainer := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		layout.NewSpacer(),
		chooseButton,
		exportButton,
	)

	container := fyne.NewContainerWithLayout(
		layout.NewVBoxLayout(),
		widget.NewLabel(description),
		mypageURLContainer,
		targetPathContainer,
		statusContainer,
		buttonsContainer,
	)

	return container
}

func Start(fn ExportFn) {
	state := &state{
		targetPath: binding.NewString(),
		mypageURL:  binding.NewString(),
		status:     binding.NewString(),
		dotCount:   0,
	}

	app := app.New()
	app.Settings().SetTheme(&myTheme{})
	w := app.NewWindow(title)

	w.SetFixedSize(true)
	// w.Resize(fyne.NewSize(width, height))
	w.CenterOnScreen()
	w.SetContent(prepareContent(w, state, fn))

	w.ShowAndRun()
}
