package gui

import (
	"image/color"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type myTheme struct{}

var _ fyne.Theme = (*myTheme)(nil)

func (m myTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	return theme.LightTheme().Color(name, theme.VariantLight)
}

func (m myTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.LightTheme().Icon(name)
}

func (m myTheme) Font(style fyne.TextStyle) fyne.Resource {
	if style.Monospace {
		return theme.LightTheme().Font(style)
	}
	if style.Italic {
		return theme.LightTheme().Font(style)
	}
	if style.Bold {
		return resourceRoundedMplus1pBoldTtf
	}
	return resourceRoundedMplus1pRegularTtf
}

func (m myTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.LightTheme().Size(name)
}

const (
	title       = "いきものAZ エクスポーター"
	description = `これはいきもの特化型SNS「いきものAZ」のユーザデータの保存を行うツールです。
以下のように使ってください:
  1. 「保存先フォルダを選ぶ」をクリック (いくつかファイルとフォルダができます)
  2. データを保存するフォルダを選択して"Open"を押す (英語ですみません)
  3. 「エクスポート開始」をクリック
  4. 待つ
  5. 保存先のindex.htmlをダブルクリック`

	captionChooseButton = "保存先フォルダを選ぶ"
	captionChooser      = "フォルダを選んでください"
	captionExportButton = "エクスポート開始"
)

type state struct {
	targetPath binding.String
	status     binding.String
	dotCount   int
}

func export(s *state, fn func(), b1, b2 *widget.Button) {
	b1.Disable()
	b2.Disable()

	c := make(chan bool)
	go calculateStatus(s, c)
	fn()
	c <- true

	b1.Enable()
	b2.Enable()
}

func Start() {
	state := &state{
		targetPath: binding.NewString(),
		status:     binding.NewString(),
		dotCount:   0,
	}
	state.status.Set("")

	app := app.New()
	app.Settings().SetTheme(&myTheme{})
	w := app.NewWindow(title)

	w.SetFixedSize(true)
	//	w.Resize(fyne.NewSize(width, height))
	w.CenterOnScreen()

	targetPathLabel := widget.NewLabel("")
	targetPathLabel.Bind(state.targetPath)
	chooserFn := func(uri fyne.ListableURI, err error) {
		if uri != nil {
			state.targetPath.Set(uri.Path())
		}
	}
	chooser := dialog.NewFolderOpen(chooserFn, w)

	chooseButtonFn := func() {
		chooser.Show()
	}
	chooseButton := widget.NewButton(captionChooseButton, chooseButtonFn)

	exportButton := widget.NewButton("", func() {})
	exportButtonFn := func() {
		fn := func() {
			time.Sleep(10 * time.Second)
		}
		export(state, fn, chooseButton, exportButton)
	}
	exportButton = widget.NewButton(captionExportButton, exportButtonFn)

	statusLabel := widget.NewLabel("")
	statusLabel.Bind(state.status)

	targetPathContainer := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
		widget.NewLabel("結果保存先:"),
		targetPathLabel,
	)
	statusContainer := fyne.NewContainerWithLayout(
		layout.NewHBoxLayout(),
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
		statusContainer,
		targetPathContainer,
		buttonsContainer,
	)
	w.SetContent(container)

	w.ShowAndRun()
}
