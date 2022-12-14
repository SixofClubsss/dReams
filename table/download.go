package table

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	coingecko "github.com/superoo7/go-gecko/v3"

	"github.com/SixofClubsss/dReams/rpc"

	"fyne.io/fyne/v2/canvas"
)

type sharedCards struct {
	//window   fyne.Window
	//progress *widget.ProgressBar
	//Back canvas.Image
	// Card1  canvas.Image
	// Card2  canvas.Image
	// Card3  canvas.Image
	// Card4  canvas.Image
	// Card5  canvas.Image
	// Card6  canvas.Image
	// Card7  canvas.Image
	// Card8  canvas.Image
	// Card9  canvas.Image
	// Card10 canvas.Image
	// Card11 canvas.Image
	// Card12 canvas.Image
	// Card13 canvas.Image
	// Card14 canvas.Image
	// Card15 canvas.Image
	// Card16 canvas.Image
	// Card17 canvas.Image
	// Card18 canvas.Image
	// Card19 canvas.Image
	// Card20 canvas.Image
	// Card21 canvas.Image
	// Card22 canvas.Image
	// Card23 canvas.Image
	// Card24 canvas.Image
	// Card25 canvas.Image
	// Card26 canvas.Image
	// Card27 canvas.Image
	// Card28 canvas.Image
	// Card29 canvas.Image
	// Card30 canvas.Image
	// Card31 canvas.Image
	// Card32 canvas.Image
	// Card33 canvas.Image
	// Card34 canvas.Image
	// Card35 canvas.Image
	// Card36 canvas.Image
	// Card37 canvas.Image
	// Card38 canvas.Image
	// Card39 canvas.Image
	// Card40 canvas.Image
	// Card41 canvas.Image
	// Card42 canvas.Image
	// Card43 canvas.Image
	// Card44 canvas.Image
	// Card45 canvas.Image
	// Card46 canvas.Image
	// Card47 canvas.Image
	// Card48 canvas.Image
	// Card49 canvas.Image
	// Card50 canvas.Image
	// Card51 canvas.Image
	// Card52 canvas.Image
	//Empty canvas.Image

	P1_avatar canvas.Image
	P2_avatar canvas.Image
	P3_avatar canvas.Image
	P4_avatar canvas.Image
	P5_avatar canvas.Image
	P6_avatar canvas.Image

	GotP1 bool
	GotP2 bool
	GotP3 bool
	GotP4 bool
	GotP5 bool
	GotP6 bool
}

func ClearShared() {
	rpc.Display.Res = ""
	rpc.Round.First_try = true
	rpc.Round.P1_url = ""
	rpc.Round.P2_url = ""
	rpc.Round.P3_url = ""
	rpc.Round.P4_url = ""
	rpc.Round.P5_url = ""
	rpc.Round.P6_url = ""
	rpc.Round.P1_name = ""
	rpc.Round.P2_name = ""
	rpc.Round.P3_name = ""
	rpc.Round.P4_name = ""
	rpc.Round.P5_name = ""
	rpc.Round.P6_name = ""
	rpc.Round.Bettor = ""
	rpc.Round.Raisor = ""
	rpc.Round.Last = 0
	rpc.Signal.Reveal = false
	rpc.Signal.Out1 = false
	Shared.GotP1 = false
	Shared.GotP2 = false
	Shared.GotP3 = false
	Shared.GotP4 = false
	Shared.GotP5 = false
	Shared.GotP6 = false
	Shared.P1_avatar = *canvas.NewImageFromImage(nil)
	Shared.P2_avatar = *canvas.NewImageFromImage(nil)
	Shared.P3_avatar = *canvas.NewImageFromImage(nil)
	Shared.P4_avatar = *canvas.NewImageFromImage(nil)
	Shared.P5_avatar = *canvas.NewImageFromImage(nil)
	Shared.P6_avatar = *canvas.NewImageFromImage(nil)

}

var Shared sharedCards

func GetUrls(face, back string) {
	if rpc.Round.ID != 1 {
		Settings.FaceUrl = face
		Settings.BackUrl = back
	}
}

func DownloadFile(Url, fileName string) (canvas.Image, error) {
	response, err := http.Get(Url)
	if err != nil {
		return *canvas.NewImageFromImage(nil), err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return *canvas.NewImageFromImage(nil), errors.New("received non 200 response code")
	}

	img := *canvas.NewImageFromReader(response.Body, fileName)

	return img, nil
}

func ShowAvatar(tab bool) {
	if tab {
		if rpc.Round.P1_url != "" {
			if !Shared.GotP1 {
				img1, _ := DownloadFile(rpc.Round.P1_url, "P1")
				Shared.P1_avatar = img1
				Shared.GotP1 = true
			}
		} else {
			Shared.GotP1 = false
		}

		if rpc.Round.P2_url != "" {
			if !Shared.GotP2 {
				img2, _ := DownloadFile(rpc.Round.P2_url, "P2")
				Shared.P2_avatar = img2
				Shared.GotP2 = true
			}
		} else {
			Shared.GotP2 = false
		}

		if rpc.Round.P3_url != "" {
			if !Shared.GotP3 {
				img3, _ := DownloadFile(rpc.Round.P3_url, "P3")
				Shared.P3_avatar = img3
				Shared.GotP3 = true
			}
		} else {
			Shared.GotP3 = false
		}

		if rpc.Round.P4_url != "" {
			if !Shared.GotP4 {
				img4, _ := DownloadFile(rpc.Round.P4_url, "P4")
				Shared.P4_avatar = img4
				Shared.GotP4 = true
			}
		} else {
			Shared.GotP4 = false
		}

		if rpc.Round.P5_url != "" {
			if !Shared.GotP5 {
				img5, _ := DownloadFile(rpc.Round.P5_url, "P5")
				Shared.P5_avatar = img5
				Shared.GotP5 = true
			}
		} else {
			Shared.GotP5 = false
		}

		if rpc.Round.P6_url != "" {
			if !Shared.GotP6 {
				img6, _ := DownloadFile(rpc.Round.P6_url, "P6")
				Shared.P6_avatar = img6
				Shared.GotP6 = true
			}
		} else {
			Shared.GotP6 = false
		}
	}
}

/*
func downloadMemoryDeck(url string) {
	var prog float64
	if url != "" {
		go func() {
			for i := 0; i < 53; i++ {
				float := float64(downloadSharedImages(url, i))
				prog = float / 53
				downloadPopUp(prog, i)
			}
			Settings.Shared = true
			Shared.window.Close()
		}()
	}
}

func downloadSharedImages(Url string, i int) int {
	fileName := "card" + strconv.Itoa(i) + ".png"
	log.Println("[dReams] Downloading ", Url+fileName)

	switch i {
	case 0:
		Shared.Back, _ = DownloadFile(Settings.BackUrl, fileName)
	case 1:
		Shared.Card1, _ = DownloadFile(Url+fileName, fileName)
	case 2:
		Shared.Card2, _ = DownloadFile(Url+fileName, fileName)
	case 3:
		Shared.Card3, _ = DownloadFile(Url+fileName, fileName)
	case 4:
		Shared.Card4, _ = DownloadFile(Url+fileName, fileName)
	case 5:
		Shared.Card5, _ = DownloadFile(Url+fileName, fileName)
	case 6:
		Shared.Card6, _ = DownloadFile(Url+fileName, fileName)
	case 7:
		Shared.Card7, _ = DownloadFile(Url+fileName, fileName)
	case 8:
		Shared.Card8, _ = DownloadFile(Url+fileName, fileName)
	case 9:
		Shared.Card9, _ = DownloadFile(Url+fileName, fileName)
	case 10:
		Shared.Card10, _ = DownloadFile(Url+fileName, fileName)
	case 11:
		Shared.Card11, _ = DownloadFile(Url+fileName, fileName)
	case 12:
		Shared.Card12, _ = DownloadFile(Url+fileName, fileName)
	case 13:
		Shared.Card13, _ = DownloadFile(Url+fileName, fileName)
	case 14:
		Shared.Card14, _ = DownloadFile(Url+fileName, fileName)
	case 15:
		Shared.Card15, _ = DownloadFile(Url+fileName, fileName)
	case 16:
		Shared.Card16, _ = DownloadFile(Url+fileName, fileName)
	case 17:
		Shared.Card17, _ = DownloadFile(Url+fileName, fileName)
	case 18:
		Shared.Card18, _ = DownloadFile(Url+fileName, fileName)
	case 19:
		Shared.Card19, _ = DownloadFile(Url+fileName, fileName)
	case 20:
		Shared.Card20, _ = DownloadFile(Url+fileName, fileName)
	case 21:
		Shared.Card21, _ = DownloadFile(Url+fileName, fileName)
	case 22:
		Shared.Card22, _ = DownloadFile(Url+fileName, fileName)
	case 23:
		Shared.Card23, _ = DownloadFile(Url+fileName, fileName)
	case 24:
		Shared.Card24, _ = DownloadFile(Url+fileName, fileName)
	case 25:
		Shared.Card25, _ = DownloadFile(Url+fileName, fileName)
	case 26:
		Shared.Card26, _ = DownloadFile(Url+fileName, fileName)
	case 27:
		Shared.Card27, _ = DownloadFile(Url+fileName, fileName)
	case 28:
		Shared.Card28, _ = DownloadFile(Url+fileName, fileName)
	case 29:
		Shared.Card29, _ = DownloadFile(Url+fileName, fileName)
	case 30:
		Shared.Card30, _ = DownloadFile(Url+fileName, fileName)
	case 31:
		Shared.Card31, _ = DownloadFile(Url+fileName, fileName)
	case 32:
		Shared.Card32, _ = DownloadFile(Url+fileName, fileName)
	case 33:
		Shared.Card33, _ = DownloadFile(Url+fileName, fileName)
	case 34:
		Shared.Card34, _ = DownloadFile(Url+fileName, fileName)
	case 35:
		Shared.Card35, _ = DownloadFile(Url+fileName, fileName)
	case 36:
		Shared.Card36, _ = DownloadFile(Url+fileName, fileName)
	case 37:
		Shared.Card37, _ = DownloadFile(Url+fileName, fileName)
	case 38:
		Shared.Card38, _ = DownloadFile(Url+fileName, fileName)
	case 39:
		Shared.Card39, _ = DownloadFile(Url+fileName, fileName)
	case 40:
		Shared.Card40, _ = DownloadFile(Url+fileName, fileName)
	case 41:
		Shared.Card41, _ = DownloadFile(Url+fileName, fileName)
	case 42:
		Shared.Card42, _ = DownloadFile(Url+fileName, fileName)
	case 43:
		Shared.Card43, _ = DownloadFile(Url+fileName, fileName)
	case 44:
		Shared.Card44, _ = DownloadFile(Url+fileName, fileName)
	case 45:
		Shared.Card45, _ = DownloadFile(Url+fileName, fileName)
	case 46:
		Shared.Card46, _ = DownloadFile(Url+fileName, fileName)
	case 47:
		Shared.Card47, _ = DownloadFile(Url+fileName, fileName)
	case 48:
		Shared.Card48, _ = DownloadFile(Url+fileName, fileName)
	case 49:
		Shared.Card49, _ = DownloadFile(Url+fileName, fileName)
	case 50:
		Shared.Card50, _ = DownloadFile(Url+fileName, fileName)
	case 51:
		Shared.Card51, _ = DownloadFile(Url+fileName, fileName)
	case 52:
		Shared.Card52, _ = DownloadFile(Url+fileName, fileName)
	}

	return i
}

// func SharedMemoryImage(c int) *canvas.Image {
// 	var card canvas.Image
// 	switch c {
// 	case 0:
// 		card = Shared.Back
// 	case 1:
// 		card = Shared.Card1
// 	case 2:
// 		card = Shared.Card2
// 	case 3:
// 		card = Shared.Card3
// 	case 4:
// 		card = Shared.Card4
// 	case 5:
// 		card = Shared.Card5
// 	case 6:
// 		card = Shared.Card6
// 	case 7:
// 		card = Shared.Card7
// 	case 8:
// 		card = Shared.Card8
// 	case 9:
// 		card = Shared.Card9
// 	case 10:
// 		card = Shared.Card10
// 	case 11:
// 		card = Shared.Card11
// 	case 12:
// 		card = Shared.Card12
// 	case 13:
// 		card = Shared.Card13
// 	case 14:
// 		card = Shared.Card14
// 	case 15:
// 		card = Shared.Card15
// 	case 16:
// 		card = Shared.Card16
// 	case 17:
// 		card = Shared.Card17
// 	case 18:
// 		card = Shared.Card18
// 	case 19:
// 		card = Shared.Card19
// 	case 20:
// 		card = Shared.Card20
// 	case 21:
// 		card = Shared.Card21
// 	case 22:
// 		card = Shared.Card22
// 	case 23:
// 		card = Shared.Card23
// 	case 24:
// 		card = Shared.Card24
// 	case 25:
// 		card = Shared.Card25
// 	case 26:
// 		card = Shared.Card26
// 	case 27:
// 		card = Shared.Card27
// 	case 28:
// 		card = Shared.Card28
// 	case 29:
// 		card = Shared.Card29
// 	case 30:
// 		card = Shared.Card30
// 	case 31:
// 		card = Shared.Card31
// 	case 32:
// 		card = Shared.Card32
// 	case 33:
// 		card = Shared.Card33
// 	case 34:
// 		card = Shared.Card34
// 	case 35:
// 		card = Shared.Card35
// 	case 36:
// 		card = Shared.Card36
// 	case 37:
// 		card = Shared.Card37
// 	case 38:
// 		card = Shared.Card38
// 	case 39:
// 		card = Shared.Card39
// 	case 40:
// 		card = Shared.Card40
// 	case 41:
// 		card = Shared.Card41
// 	case 42:
// 		card = Shared.Card42
// 	case 43:
// 		card = Shared.Card43
// 	case 44:
// 		card = Shared.Card44
// 	case 45:
// 		card = Shared.Card45
// 	case 46:
// 		card = Shared.Card46
// 	case 47:
// 		card = Shared.Card47
// 	case 48:
// 		card = Shared.Card48
// 	case 49:
// 		card = Shared.Card49
// 	case 50:
// 		card = Shared.Card50
// 	case 51:
// 		card = Shared.Card51
// 	case 52:
// 		card = Shared.Card52
// 	default:
// 		card = *canvas.NewImageFromFile("")
// 	}

// 	return &card
// }

func downloadProgress(p float64) fyne.Widget {
	Shared.progress = widget.NewProgressBar()
	this := binding.BindFloat(&p)
	Shared.progress.Bind(this)

	return Shared.progress
}

func downloadPopUp(p float64, i int) { /// pop up for loading progress
	if i == 0 {
		Shared.window = fyne.CurrentApp().NewWindow("Loading Custom Deck")
		Shared.window.Resize(fyne.NewSize(300, 30))
		Shared.window.SetFixedSize(true)
		Shared.window.SetIcon(nil)
		content := container.NewMax(downloadProgress(p))
		Shared.window.SetContent(content)
		Shared.window.Show()
	} else {
		content := container.NewMax(downloadProgress(p))
		Shared.window.SetContent(content)
	}
}
*/
/*
// for on demand

func SharedImage(c string) *canvas.Image {
	var card canvas.Image
	switch c {
	case "card0.png":
		card, _ = DownloadFile(table.Settings.BackUrl, c)
	case "card1.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card2.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card3.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card4.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card5.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card6.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card7.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card8.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card9.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card10.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card11.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card12.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card13.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card14.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card15.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card16.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card17.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card18.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card19.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card20.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card21.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card22.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card23.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card24.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card25.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card26.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card27.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card28.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card29.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card30.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card31.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card32.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card33.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card34.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card35.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card36.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card37.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card38.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card39.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card40.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card41.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card42.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card43.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card44.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card45.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card46.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card47.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card48.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card49.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card50.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card51.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	case "card52.png":
		card, _ = DownloadFile(table.Settings.FaceUrl+c, c)
	default:
		card = *canvas.NewImageFromFile("")
	}

	fmt.Println(card)

	return &card
}
*/

type ogreFeed struct {
	Success      bool   `json:"success"`
	Initialprice string `json:"initialprice"`
	Price        string `json:"price"`
	High         string `json:"high"`
	Low          string `json:"low"`
	Volume       string `json:"volume"`
	Bid          string `json:"bid"`
	Ask          string `json:"ask"`
}

type kuFeed struct {
	Code string `json:"code"`
	Data struct {
		Time        int64  `json:"time"`
		Sequence    string `json:"sequence"`
		Price       string `json:"price"`
		Size        string `json:"size"`
		BestBid     string `json:"bestBid"`
		BestBidSize string `json:"bestBidSize"`
		BestAsk     string `json:"bestAsk"`
		BestAskSize string `json:"bestAskSize"`
	} `json:"data"`
}

func CoinDecimal(ticker string) int {
	split := strings.Split(ticker, "-")

	if len(split) == 2 {
		switch split[1] {
		case "BTC":
			return 8
		default:
			return 2
		}
	}
	return 2
}

func GetPrice(coin string) (price float64, display string) {
	var t float64
	var k float64
	var g float64
	priceT := getOgre(coin)
	priceK := getKucoin(coin)
	priceG := getGeko(coin)

	if CoinDecimal(coin) == 8 {
		if tf, err := strconv.ParseFloat(priceT, 64); err == nil {
			t = tf * 100000000
		}

		if kf, err := strconv.ParseFloat(priceK, 64); err == nil {
			k = kf * 100000000
		}

		if gf, err := strconv.ParseFloat(priceG, 64); err == nil {
			g = gf * 100000000
		}
	} else {
		if tf, err := strconv.ParseFloat(priceT, 64); err == nil {
			t = tf * 100
		}

		if kf, err := strconv.ParseFloat(priceK, 64); err == nil {
			k = kf * 100
		}

		if gf, err := strconv.ParseFloat(priceG, 64); err == nil {
			g = gf * 100
		}
	}

	if t > 0 && k > 0 && g > 0 {
		price = (t + k + g) / 3
	} else if t > 0 && k > 0 {
		price = (t + k) / 2
	} else if k > 0 && g > 0 {
		price = (k + g) / 2
	} else if t > 0 && g > 0 {
		price = (t + g) / 2
	} else if t > 0 {
		price = t
	} else if k > 0 {
		price = k
	} else if g > 0 {
		price = g
	} else {
		price = 0
		log.Println("[dReams] Error getting price feed")
	}

	if CoinDecimal(coin) == 8 {
		display = fmt.Sprintf("%.8f", price/100000000)
	} else {
		display = fmt.Sprintf("%.2f", price/100)
	}

	return
}

func getOgre(coin string) string {
	decimal := 2
	var url string
	var found ogreFeed
	switch coin {
	case "BTC-USDT":
		url = "https://tradeogre.com/api/v1/ticker/usdt-btc"
	case "DERO-USDT":
		url = "https://tradeogre.com/api/v1/ticker/usdt-dero"
	case "XMR-USDT":
		url = "https://tradeogre.com/api/v1/ticker/usdt-xmr"
	case "DERO-BTC":
		url = "https://tradeogre.com/api/v1/ticker/btc-dero"
		decimal = 8
	case "XMR-BTC":
		url = "https://tradeogre.com/api/v1/ticker/btc-xmr"
		decimal = 8
	default:
		return ""
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[getOgre]", err)
		return ""
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println("[getOgre]", err)
		return ""
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("[getOgre]", err)
		return ""
	}

	json.Unmarshal(b, &found)

	if s, err := strconv.ParseFloat(found.Price, 64); err == nil {
		if decimal == 8 {
			return fmt.Sprintf("%.8f", s)
		}
		return fmt.Sprintf("%.2f", s)
	}

	return found.Price
}

func getKucoin(coin string) string {
	decimal := 2
	var url string
	var found kuFeed
	switch coin {
	case "BTC-USDT":
		url = "https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=BTC-USDT"
	case "DERO-USDT":
		url = "https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=DERO-USDT"
	case "XMR-USDT":
		url = "https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=XMR-USDT"
	case "DERO-BTC":
		url = "https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=DERO-BTC"
		decimal = 8
	case "XMR-BTC":
		url = "https://api.kucoin.com/api/v1/market/orderbook/level1?symbol=XMR-BTC"
		decimal = 8
	default:
		return ""
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println("[getKucoin]", err)
		return ""
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println("[getKucoin]", err)
		return ""
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Println("[getKucoin]", err)
		return ""
	}

	json.Unmarshal(b, &found)

	if s, err := strconv.ParseFloat(found.Data.Price, 64); err == nil {
		if decimal == 8 {
			return fmt.Sprintf("%.8f", s)
		}
		return fmt.Sprintf("%.2f", s)
	}

	return found.Data.Price
}

func getGeko(coin string) string {
	client := &http.Client{Timeout: time.Second * 10}
	CG := coingecko.NewClient(client)

	pair := "usd"
	var url string
	switch coin {
	case "BTC-USDT":
		url = "bitcoin"
	case "DERO-USDT":
		url = "dero"
	case "XMR-USDT":
		url = "monero"
	case "DERO-BTC":
		url = "dero"
		pair = "btc"
	case "XMR-BTC":
		url = "monero"
		pair = "btc"
	default:
		return ""
	}

	price, err := CG.SimpleSinglePrice(url, pair)
	if err != nil {
		log.Println("[getGeko]", err)
		return ""
	}

	if pair == "btc" {
		return fmt.Sprintf("%.8f", price.MarketPrice)
	}

	return fmt.Sprintf("%.2f", price.MarketPrice)
}

func downloadFileLocal(filepath string, url string) (err error) {
	_, dir := os.Stat("cards")
	if os.IsNotExist(dir) {
		log.Println("[dReams] Creating Cards Dir")
		mkdir := os.Mkdir("cards", 0755)
		if mkdir != nil {
			log.Println("[dReams]", mkdir)
		} else {
			mksub := os.Mkdir("cards/backs", 0755)
			if mksub != nil {
				log.Println("[dReams]", mksub)
			}
		}

	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func GetZipDeck(face, url string) {
	downloadFileLocal("cards/"+face+".zip", url)
	files, err := Unzip("cards/"+face+".zip", "cards/"+face)

	if err != nil {
		log.Println("[GetZipDeck]", err)
	}

	log.Println("[dReams] Unzipped files:\n" + strings.Join(files, "\n"))
}

func Unzip(src string, destination string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}

	defer r.Close()

	for _, f := range r.File {
		fpath := filepath.Join(destination, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(destination)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s is an illegal filepath", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath,
			os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
			f.Mode())

		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}
