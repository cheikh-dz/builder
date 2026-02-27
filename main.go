package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2"
)

var (
	shell32          = syscall.NewLazyDLL("shell32.dll")
	emptyRecycleBin  = shell32.NewProc("SHEmptyRecycleBinW")
)

const (
	SHERB_NOCONFIRMATION = 0x00000001
	SHERB_NOPROGRESSUI   = 0x00000002
	SHERB_NOSOUND        = 0x00000004
)

// EmptyRecycleBin يقوم بتفريغ سلة المهملات على Windows
func EmptyRecycleBin() error {
	// استدعاء دالة Windows API
	ret, _, err := emptyRecycleBin.Call(
		uintptr(0),                   // hwnd (لا توجد نافذة رئيسية)
		uintptr(0),                   // pszRootPath (NULL = جميع الأقراص)
		SHERB_NOCONFIRMATION|SHERB_NOPROGRESSUI|SHERB_NOSOUND,
	)
	
	if ret == 0 {
		return fmt.Errorf("فشل في تفريغ سلة المهملات: %v", err)
	}
	return nil
}

func main() {
	// إنشاء التطبيق والنافذة
	myApp := app.New()
	myWindow := myApp.NewWindow("تنظيف سلة المهملات")
	myWindow.Resize(fyne.NewSize(400, 200))

	// إنشاء التسمية
	label := widget.NewLabel("اضغط على الزر لتنظيف سلة المهملات")
	label.Alignment = fyne.TextAlignCenter

	// إنشاء الزر
	button := widget.NewButton("🗑️ تنظيف سلة المهملات", func() {
		// عرض نافذة تأكيد
		dialog.ShowConfirm(
			"تأكيد",
			"هل أنت متأكد من رغبتك في تفريغ سلة المهملات؟",
			func(confirmed bool) {
				if confirmed {
					err := EmptyRecycleBin()
					if err != nil {
						dialog.ShowError(err, myWindow)
					} else {
						dialog.ShowInformation("نجاح", "تم تفريغ سلة المهملات بنجاح!", myWindow)
					}
				}
			},
			myWindow,
		)
	})
	button.Importance = widget.HighImportance

	// ترتيب العناصر
	content := container.NewVBox(
		label,
		widget.NewSeparator(),
		button,
	)

	myWindow.SetContent(content)
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}