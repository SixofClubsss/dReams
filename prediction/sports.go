package prediction

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/SixofClubsss/dReams/menu"
	"github.com/SixofClubsss/dReams/rpc"
	"github.com/SixofClubsss/dReams/table"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type sportsItems struct {
	Contract       string
	Contract_list  []string
	Favorites_list []string
	Sports_list    *widget.List
	connected_box  *widget.Check
}

var SportsControl sportsItems

func SportsConnectedBox() fyne.Widget {
	SportsControl.connected_box = widget.NewCheck("", func(b bool) {})
	SportsControl.connected_box.Disable()

	return SportsControl.connected_box
}

func SportsContractEntry() fyne.Widget {
	options := []string{""}
	table.Actions.S_contract = widget.NewSelectEntry(options)
	table.Actions.S_contract.PlaceHolder = "Contract Address: "
	table.Actions.S_contract.OnCursorChanged = func() {
		if rpc.Signal.Daemon {
			yes, _ := rpc.CheckBetContract(SportsControl.Contract)
			if yes {
				SportsControl.connected_box.SetChecked(true)
			} else {
				SportsControl.connected_box.SetChecked(false)
			}
		}
	}

	this := binding.BindString(&SportsControl.Contract)
	table.Actions.S_contract.Bind(this)

	return table.Actions.S_contract
}

func SportsBox() fyne.CanvasObject {
	table.Actions.Game_select = widget.NewSelect(table.Actions.Game_options, func(s string) {
		split := strings.Split(s, "   ")
		a, b := menu.GetSportsTeams(SportsControl.Contract, split[0])
		if table.Actions.Game_select.SelectedIndex() >= 0 {
			table.Actions.Multi.Show()
			table.Actions.ButtonA.Show()
			table.Actions.ButtonB.Show()
			table.Actions.ButtonA.Text = a
			table.Actions.ButtonA.Refresh()
			table.Actions.ButtonB.Text = b
			table.Actions.ButtonB.Refresh()
		} else {
			table.Actions.Multi.Hide()
			table.Actions.ButtonA.Hide()
			table.Actions.ButtonB.Hide()
		}
	})

	table.Actions.Game_select.PlaceHolder = "Select Game #"
	table.Actions.Game_select.Hide()

	var Multi_options = []string{"1x", "3x", "5x"}
	table.Actions.Multi = widget.NewRadioGroup(Multi_options, func(s string) {

	})
	table.Actions.Multi.Horizontal = true
	table.Actions.Multi.Hide()

	table.Actions.ButtonA = widget.NewButton("TEAM A", func() {
		if len(SportsControl.Contract) == 64 {
			confirmPopUp(3, table.Actions.ButtonA.Text, table.Actions.ButtonB.Text)
		}
	})
	table.Actions.ButtonA.Hide()

	table.Actions.ButtonB = widget.NewButton("TEAM B", func() {
		if len(SportsControl.Contract) == 64 {
			confirmPopUp(4, table.Actions.ButtonA.Text, table.Actions.ButtonB.Text)
		}
	})
	table.Actions.ButtonB.Hide()

	sports_muli := container.NewCenter(table.Actions.Multi)
	sports_actions := container.NewVBox(
		sports_muli,
		table.Actions.Game_select,
		table.Actions.ButtonA,
		table.Actions.ButtonB)

	table.Actions.Sports_box = sports_actions
	table.Actions.Sports_box.Hide()

	return table.Actions.Sports_box

}

func SportsListings() fyne.CanvasObject { /// sports contract list
	SportsControl.Sports_list = widget.NewList(
		func() int {
			return len(SportsControl.Contract_list)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(SportsControl.Contract_list[i])
		})

	var item string

	SportsControl.Sports_list.OnSelected = func(id widget.ListItemID) {
		if id != 0 && rpc.Wallet.Connect {
			table.Actions.Game_select.ClearSelected()
			table.Actions.Game_select.Options = []string{}
			table.Actions.Game_select.Refresh()
			split := strings.Split(SportsControl.Contract_list[id], "   ")
			trimmed := strings.Trim(split[2], " ")
			table.Actions.Sports_box.Show()
			if len(trimmed) == 64 {
				item = SportsControl.Contract_list[id]
				table.Actions.S_contract.SetText(trimmed)
			}
		} else {
			table.Actions.Sports_box.Hide()
		}
	}

	save := widget.NewButton("Favorite", func() {
		SportsControl.Favorites_list = append(SportsControl.Favorites_list, item)
		sort.Strings(SportsControl.Favorites_list)
	})

	cont := container.NewBorder(
		nil,
		container.NewBorder(nil, nil, nil, save, layout.NewSpacer()),
		nil,
		nil,
		SportsControl.Sports_list)

	return cont
}

func SportsFavorites() fyne.CanvasObject {
	favorites := widget.NewList(
		func() int {
			return len(SportsControl.Favorites_list)
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(SportsControl.Favorites_list[i])
		})

	var item string

	favorites.OnSelected = func(id widget.ListItemID) {
		split := strings.Split(SportsControl.Favorites_list[id], "   ")
		if len(split) >= 3 {
			trimmed := strings.Trim(split[2], " ")
			if len(trimmed) == 64 {
				item = SportsControl.Favorites_list[id]
				table.Actions.S_contract.SetText(trimmed)
			}
		}
	}

	remove := widget.NewButton("Remove", func() {
		if len(SportsControl.Favorites_list) > 0 {
			favorites.UnselectAll()
			new := SportsControl.Favorites_list
			favorites.UnselectAll()
			for i := range new {
				if new[i] == item {
					copy(new[i:], new[i+1:])
					new[len(new)-1] = ""
					new = new[:len(new)-1]
					SportsControl.Favorites_list = new
					break
				}
			}
		}
		favorites.Refresh()
		sort.Strings(SportsControl.Favorites_list)
	})

	cont := container.NewBorder(
		nil,
		container.NewBorder(nil, nil, nil, remove, layout.NewSpacer()),
		nil,
		nil,
		favorites)

	return cont
}

func sports(league string) (api string) {
	switch league {
	case "NHL":
		api = "http://site.api.espn.com/apis/site/v2/sports/hockey/nhl/scoreboard"
	case "FIFA":
		api = "http://site.api.espn.com/apis/site/v2/sports/soccer/fifa.world/scoreboard"
	case "NFL":
		api = "http://site.api.espn.com/apis/site/v2/sports/football/nfl/scoreboard"
	case "NBA":
		api = "http://site.api.espn.com/apis/site/v2/sports/basketball/nba/scoreboard"
	default:
		api = ""
	}

	return api
}

func GetCurrentWeek(league string) {
	for i := 0; i < 8; i++ {
		now := time.Now().AddDate(0, 0, i)
		date := time.Unix(now.Unix(), 0).String()
		date = date[0:10]
		comp := date[0:4] + date[5:7] + date[8:10]
		switch league {
		case "FIFA":
			GetSoccer(comp)
		case "NBA":
			GetBasketball(comp)
		case "NFL":
			GetFootball(comp)
		case "NHL":
			GetHockey(comp)
		default:

		}
	}
}

func callSoccer(date, league string) (s *soccer) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", sports(league)+"?dates="+date, nil)
	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
	}

	json.Unmarshal(b, &s)

	return s
}

func callBasketball(date, league string) (bb *basketball) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", sports(league)+"?dates="+date, nil)
	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
	}

	json.Unmarshal(b, &bb)

	return bb
}

func callFootball(date, league string) (f *football) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", sports(league)+"?dates="+date, nil)
	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
	}

	json.Unmarshal(b, &f)

	return f
}

func callHockey(date, league string) (h *hockey) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", sports(league)+"?dates="+date, nil)
	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
	}

	json.Unmarshal(b, &h)

	return h
}

func GetGameEnd(date, game, league string) {
	var found hockey

	client := &http.Client{}
	req, err := http.NewRequest("GET", sports(league)+"?dates="+date, nil)

	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
	}

	json.Unmarshal(b, &found)

	for i := range found.Events {
		trimmed := strings.Trim(found.Events[i].Competitions[0].StartDate, "Z")
		utc_time, err := time.Parse("2006-01-02T15:04", trimmed)
		if err != nil {
			log.Println(err)
		}

		a := found.Events[i].Competitions[0].Competitors[0].Team.Abbreviation
		b := found.Events[i].Competitions[0].Competitors[1].Team.Abbreviation
		g := a + "-" + b

		if g == game {
			PS_Control.S_end.SetText(strconv.Itoa(int(utc_time.Unix())))
		}
	}
}

func callScores(date, league string) (s *scores) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", sports(league)+"?dates="+date, nil)
	if err != nil {
		log.Println(err.Error())
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err.Error())
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err.Error())
	}

	json.Unmarshal(b, &s)

	return s
}

func GetScores(label *widget.Label, league string) {
	var single bool
	for i := -1; i < 1; i++ {
		day := time.Now().AddDate(0, 0, i)
		date := time.Unix(day.Unix(), 0).String()
		date = date[0:10]
		comp := date[0:4] + date[5:7] + date[8:10]
		found := callScores(comp, league)
		if !single {
			label.SetText(found.Leagues[0].Abbreviation + "\n" + found.Day.Date + "\n")
		}

		for i := range found.Events {
			trimmed := strings.Trim(found.Events[i].Competitions[0].StartDate, "Z")
			utc_time, err := time.Parse("2006-01-02T15:04", trimmed)
			if err != nil {
				log.Println(err)
			}

			tz, _ := time.LoadLocation("Local")
			local := utc_time.In(tz).String()
			state := found.Events[i].Competitions[0].Status.Type.State
			team_a := found.Events[i].Competitions[0].Competitors[0].Team.Abbreviation
			team_b := found.Events[i].Competitions[0].Competitors[1].Team.Abbreviation
			score_a := found.Events[i].Competitions[0].Competitors[0].Score
			score_b := found.Events[i].Competitions[0].Competitors[1].Score
			period := found.Events[i].Status.Period
			clock := found.Events[i].Competitions[0].Status.DisplayClock
			complete := found.Events[i].Status.Type.Completed

			var format string
			switch league {
			case "FIFA":
				format = " Half "
			case "NBA":
				format = " Quarter "
			case "NFL":
				format = " Quarter "
			case "NHL":
				format = " Period "
			default:
			}

			var abv string
			switch period {
			case 0:
				abv = ""
			case 1:
				abv = "st "
			case 2:
				abv = "nd "
			case 3:
				abv = "rd "
			case 4:
				abv = "th "
			default:
				abv = "th "
			}
			if state == "pre" {
				label.SetText(label.Text + team_a + " - " + team_b + "\nStart time: " + local + "\nState: " + state + "\nComplete: " + strconv.FormatBool(complete) + "\n\n")
			} else {
				label.SetText(label.Text + team_a + " - " + team_b + "\nStart time: " + local + "\nState: " + state +
					"\n" + strconv.Itoa(period) + abv + format + " " + clock + "\n" + team_a + ": " + score_a + "\n" + team_b + ": " + score_b + "\nComplete: " + strconv.FormatBool(complete) + "\n\n")
			}

			single = true
		}
	}
	label.Refresh()
}

func GetHockey(date string) {
	found := callHockey(date, "NHL")
	for i := range found.Events {
		pregame := found.Events[i].Competitions[0].Status.Type.State
		trimmed := strings.Trim(found.Events[i].Competitions[0].StartDate, "Z")
		utc_time, err := time.Parse("2006-01-02T15:04", trimmed)
		if err != nil {
			log.Println(err)
		}

		tz, _ := time.LoadLocation("Local")

		teamA := found.Events[i].Competitions[0].Competitors[0].Team.Abbreviation
		teamB := found.Events[i].Competitions[0].Competitors[1].Team.Abbreviation

		if !found.Events[i].Status.Type.Completed && pregame == "pre" {
			current := PS_Control.S_game.Options
			new := append(current, utc_time.In(tz).String()[0:16]+"   "+teamA+"-"+teamB)
			PS_Control.S_game.Options = new
		}
	}
}

func GetSoccer(date string) {
	found := callSoccer(date, "FIFA")
	for i := range found.Events {
		pregame := found.Events[i].Competitions[0].Status.Type.State

		trimmed := strings.Trim(found.Events[i].Competitions[0].StartDate, "Z")
		utc_time, err := time.Parse("2006-01-02T15:04", trimmed)
		if err != nil {
			log.Println(err)
		}

		tz, _ := time.LoadLocation("Local")

		teamA := found.Events[i].Competitions[0].Competitors[0].Team.Abbreviation
		teamB := found.Events[i].Competitions[0].Competitors[1].Team.Abbreviation

		if !found.Events[i].Status.Type.Completed && pregame == "pre" {
			current := PS_Control.S_game.Options
			new := append(current, utc_time.In(tz).String()[0:16]+"   "+teamA+"-"+teamB)
			PS_Control.S_game.Options = new
		}
	}
}

func GetWinner(sport, game, league string) (string, string) {
	for i := -2; i < 1; i++ {
		day := time.Now().AddDate(0, 0, i)
		date := time.Unix(day.Unix(), 0).String()
		date = date[0:10]
		comp := date[0:4] + date[5:7] + date[8:10]

		found := callScores(comp, league)

		for i := range found.Events {
			a := found.Events[i].Competitions[0].Competitors[0].Team.Abbreviation
			b := found.Events[i].Competitions[0].Competitors[1].Team.Abbreviation
			g := a + "-" + b

			if g == game {
				if found.Events[i].Status.Type.Completed {
					teamA := found.Events[i].Competitions[0].Competitors[0].Team.Abbreviation
					a_win := found.Events[i].Competitions[0].Competitors[0].Winner

					teamB := found.Events[i].Competitions[0].Competitors[1].Team.Abbreviation
					b_win := found.Events[i].Competitions[0].Competitors[1].Winner

					if a_win && !b_win {
						return "team_a", teamA
					} else if b_win && !a_win {
						return "team_b", teamB
					} else {
						return "", ""
					}
				}
			}
		}
	}
	return "", ""
}

func GetFootball(date string) {
	found := callFootball(date, "NFL")
	for i := range found.Events {
		pregame := found.Events[i].Competitions[0].Status.Type.State
		trimmed := strings.Trim(found.Events[i].Competitions[0].StartDate, "Z")
		utc_time, err := time.Parse("2006-01-02T15:04", trimmed)
		if err != nil {
			log.Println(err)
		}

		tz, _ := time.LoadLocation("Local")

		teamA := found.Events[i].Competitions[0].Competitors[0].Team.Abbreviation
		teamB := found.Events[i].Competitions[0].Competitors[1].Team.Abbreviation

		if !found.Events[i].Status.Type.Completed && pregame == "pre" {
			current := PS_Control.S_game.Options
			new := append(current, utc_time.In(tz).String()[0:16]+"   "+teamA+"-"+teamB)
			PS_Control.S_game.Options = new
		}
	}
}

func GetBasketball(date string) {
	found := callBasketball(date, "NBA")
	for i := range found.Events {
		pregame := found.Events[i].Competitions[0].Status.Type.State
		trimmed := strings.Trim(found.Events[i].Competitions[0].StartDate, "Z")
		utc_time, err := time.Parse("2006-01-02T15:04", trimmed)
		if err != nil {
			log.Println(err)
		}

		tz, _ := time.LoadLocation("Local")

		teamA := found.Events[i].Competitions[0].Competitors[0].Team.Abbreviation
		teamB := found.Events[i].Competitions[0].Competitors[1].Team.Abbreviation

		if !found.Events[i].Status.Type.Completed && pregame == "pre" {
			current := PS_Control.S_game.Options
			new := append(current, utc_time.In(tz).String()[0:16]+"   "+teamA+"-"+teamB)
			PS_Control.S_game.Options = new
		}
	}
}

type scores struct {
	Leagues []struct {
		ID           string `json:"id"`
		UID          string `json:"uid"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
		MidsizeName  string `json:"midsizeName"`
		Slug         string `json:"slug"`
		Season       struct {
			Year      int    `json:"year"`
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Type      struct {
				ID           string `json:"id"`
				Type         int    `json:"type"`
				Name         string `json:"name"`
				Abbreviation string `json:"abbreviation"`
			} `json:"type"`
		} `json:"season"`
		Logos               []struct{} `json:"logos"`
		CalendarType        string     `json:"calendarType"`
		CalendarIsWhitelist bool       `json:"calendarIsWhitelist"`
		CalendarStartDate   string     `json:"calendarStartDate"`
		CalendarEndDate     string     `json:"calendarEndDate"`
		Calendar            []string   `json:"calendar"`
	} `json:"leagues"`
	Season struct {
		Type int `json:"type"`
		Year int `json:"year"`
	} `json:"season"`
	Day struct {
		Date string `json:"date"`
	} `json:"day"`
	Events []struct {
		ID        string `json:"id"`
		UID       string `json:"uid"`
		Date      string `json:"date"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
		Season    struct {
			Year int    `json:"year"`
			Type int    `json:"type"`
			Slug string `json:"slug"`
		} `json:"season"`
		Competitions []struct {
			ID         string `json:"id"`
			UID        string `json:"uid"`
			Date       string `json:"date"`
			StartDate  string `json:"startDate"`
			Attendance int    `json:"attendance"`
			TimeValid  bool   `json:"timeValid"`
			Recent     bool   `json:"recent"`
			Status     struct {
				Clock        float64 `json:"clock"`
				DisplayClock string  `json:"displayClock"`
				Period       int     `json:"period"`
				Type         struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					State       string `json:"state"`
					Completed   bool   `json:"completed"`
					Description string `json:"description"`
					Detail      string `json:"detail"`
					ShortDetail string `json:"shortDetail"`
				} `json:"type"`
			} `json:"status"`
			Venue         struct{}      `json:"venue"`
			Format        struct{}      `json:"format"`
			Notes         []interface{} `json:"notes"`
			GeoBroadcasts []interface{} `json:"geoBroadcasts"`
			Broadcasts    []interface{} `json:"broadcasts"`
			Competitors   []struct {
				ID       string     `json:"id"`
				UID      string     `json:"uid"`
				Type     string     `json:"type"`
				Order    int        `json:"order"`
				HomeAway string     `json:"homeAway"`
				Winner   bool       `json:"winner"`
				Form     string     `json:"form"`
				Score    string     `json:"score"`
				Records  []struct{} `json:"records"`
				Team     struct {
					ID               string     `json:"id"`
					UID              string     `json:"uid"`
					Abbreviation     string     `json:"abbreviation"`
					DisplayName      string     `json:"displayName"`
					ShortDisplayName string     `json:"shortDisplayName"`
					Name             string     `json:"name"`
					Location         string     `json:"location"`
					Color            string     `json:"color"`
					AlternateColor   string     `json:"alternateColor"`
					IsActive         bool       `json:"isActive"`
					Logo             string     `json:"logo"`
					Links            []struct{} `json:"links"`
					Venue            struct{}   `json:"venue"`
				} `json:"team,omitempty"`
				Statistics []struct{} `json:"statistics"`
			} `json:"competitors"`
			Details   []struct{} `json:"details"`
			Headlines []struct{} `json:"headlines"`
		} `json:"competitions"`
		Status struct {
			Clock        float64 `json:"clock"`
			DisplayClock string  `json:"displayClock"`
			Period       int     `json:"period"`
			Type         struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				State       string `json:"state"`
				Completed   bool   `json:"completed"`
				Description string `json:"description"`
				Detail      string `json:"detail"`
				ShortDetail string `json:"shortDetail"`
			} `json:"type"`
		} `json:"status"`
		Links []struct{} `json:"links"`
	} `json:"events"`
}

type soccer struct {
	Leagues []struct {
		ID           string `json:"id"`
		UID          string `json:"uid"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
		MidsizeName  string `json:"midsizeName"`
		Slug         string `json:"slug"`
		Season       struct {
			Year      int    `json:"year"`
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Type      struct {
				ID           string `json:"id"`
				Type         int    `json:"type"`
				Name         string `json:"name"`
				Abbreviation string `json:"abbreviation"`
			} `json:"type"`
		} `json:"season"`
		Logos []struct {
			Href        string   `json:"href"`
			Width       int      `json:"width"`
			Height      int      `json:"height"`
			Alt         string   `json:"alt"`
			Rel         []string `json:"rel"`
			LastUpdated string   `json:"lastUpdated"`
		} `json:"logos"`
		CalendarType        string   `json:"calendarType"`
		CalendarIsWhitelist bool     `json:"calendarIsWhitelist"`
		CalendarStartDate   string   `json:"calendarStartDate"`
		CalendarEndDate     string   `json:"calendarEndDate"`
		Calendar            []string `json:"calendar"`
	} `json:"leagues"`
	Season struct {
		Type int `json:"type"`
		Year int `json:"year"`
	} `json:"season"`
	Day struct {
		Date string `json:"date"`
	} `json:"day"`
	Events []struct {
		ID        string `json:"id"`
		UID       string `json:"uid"`
		Date      string `json:"date"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
		Season    struct {
			Year int    `json:"year"`
			Type int    `json:"type"`
			Slug string `json:"slug"`
		} `json:"season"`
		Competitions []struct {
			ID         string `json:"id"`
			UID        string `json:"uid"`
			Date       string `json:"date"`
			StartDate  string `json:"startDate"`
			Attendance int    `json:"attendance"`
			TimeValid  bool   `json:"timeValid"`
			Recent     bool   `json:"recent"`
			Status     struct {
				Clock        float64 `json:"clock"`
				DisplayClock string  `json:"displayClock"`
				Period       int     `json:"period"`
				Type         struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					State       string `json:"state"`
					Completed   bool   `json:"completed"`
					Description string `json:"description"`
					Detail      string `json:"detail"`
					ShortDetail string `json:"shortDetail"`
				} `json:"type"`
			} `json:"status"`
			Venue struct {
				ID       string `json:"id"`
				FullName string `json:"fullName"`
				Address  struct {
					City    string `json:"city"`
					Country string `json:"country"`
				} `json:"address"`
			} `json:"venue"`
			Format struct {
				Regulation struct {
					Periods int `json:"periods"`
				} `json:"regulation"`
			} `json:"format"`
			Notes         []interface{} `json:"notes"`
			GeoBroadcasts []interface{} `json:"geoBroadcasts"`
			Broadcasts    []interface{} `json:"broadcasts"`
			Competitors   []struct {
				ID       string `json:"id"`
				UID      string `json:"uid"`
				Type     string `json:"type"`
				Order    int    `json:"order"`
				HomeAway string `json:"homeAway"`
				Winner   bool   `json:"winner"`
				Form     string `json:"form"`
				Score    string `json:"score"`
				Records  []struct {
					Name         string `json:"name"`
					Type         string `json:"type"`
					Summary      string `json:"summary"`
					Abbreviation string `json:"abbreviation"`
				} `json:"records"`
				Team struct {
					ID               string `json:"id"`
					UID              string `json:"uid"`
					Abbreviation     string `json:"abbreviation"`
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Name             string `json:"name"`
					Location         string `json:"location"`
					Color            string `json:"color"`
					AlternateColor   string `json:"alternateColor"`
					IsActive         bool   `json:"isActive"`
					Logo             string `json:"logo"`
					Links            []struct {
						Rel        []string `json:"rel"`
						Href       string   `json:"href"`
						Text       string   `json:"text"`
						IsExternal bool     `json:"isExternal"`
						IsPremium  bool     `json:"isPremium"`
					} `json:"links"`
					Venue struct {
						ID string `json:"id"`
					} `json:"venue"`
				} `json:"team,omitempty"`
				Statistics []struct {
					Name         string `json:"name"`
					Abbreviation string `json:"abbreviation"`
					DisplayValue string `json:"displayValue"`
				} `json:"statistics"`
			} `json:"competitors"`
			Details []struct {
				Type struct {
					ID   string `json:"id"`
					Text string `json:"text"`
				} `json:"type"`
				Clock struct {
					Value        float64 `json:"value"`
					DisplayValue string  `json:"displayValue"`
				} `json:"clock"`
				Team struct {
					ID string `json:"id"`
				} `json:"team"`
				ScoreValue       int  `json:"scoreValue"`
				ScoringPlay      bool `json:"scoringPlay"`
				RedCard          bool `json:"redCard"`
				YellowCard       bool `json:"yellowCard"`
				PenaltyKick      bool `json:"penaltyKick"`
				OwnGoal          bool `json:"ownGoal"`
				Shootout         bool `json:"shootout"`
				AthletesInvolved []struct {
					ID          string `json:"id"`
					DisplayName string `json:"displayName"`
					ShortName   string `json:"shortName"`
					FullName    string `json:"fullName"`
					Jersey      string `json:"jersey"`
					Team        struct {
						ID string `json:"id"`
					} `json:"team"`
					Links []struct {
						Rel  []string `json:"rel"`
						Href string   `json:"href"`
					} `json:"links"`
					Position string `json:"position"`
				} `json:"athletesInvolved,omitempty"`
			} `json:"details"`
			Headlines []struct {
				Description   string `json:"description"`
				Type          string `json:"type"`
				ShortLinkText string `json:"shortLinkText"`
			} `json:"headlines"`
		} `json:"competitions"`
		Status struct {
			Clock        float64 `json:"clock"`
			DisplayClock string  `json:"displayClock"`
			Period       int     `json:"period"`
			Type         struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				State       string `json:"state"`
				Completed   bool   `json:"completed"`
				Description string `json:"description"`
				Detail      string `json:"detail"`
				ShortDetail string `json:"shortDetail"`
			} `json:"type"`
		} `json:"status"`
		Links []struct {
			Language   string   `json:"language"`
			Rel        []string `json:"rel"`
			Href       string   `json:"href"`
			Text       string   `json:"text"`
			ShortText  string   `json:"shortText"`
			IsExternal bool     `json:"isExternal"`
			IsPremium  bool     `json:"isPremium"`
		} `json:"links"`
	} `json:"events"`
}

type hockey struct {
	Leagues []struct {
		ID           string `json:"id"`
		UID          string `json:"uid"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
		Slug         string `json:"slug"`
		Season       struct {
			Year      int    `json:"year"`
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Type      struct {
				ID           string `json:"id"`
				Type         int    `json:"type"`
				Name         string `json:"name"`
				Abbreviation string `json:"abbreviation"`
			} `json:"type"`
		} `json:"season"`
		Logos []struct {
			Href        string   `json:"href"`
			Width       int      `json:"width"`
			Height      int      `json:"height"`
			Alt         string   `json:"alt"`
			Rel         []string `json:"rel"`
			LastUpdated string   `json:"lastUpdated"`
		} `json:"logos"`
		CalendarType        string   `json:"calendarType"`
		CalendarIsWhitelist bool     `json:"calendarIsWhitelist"`
		CalendarStartDate   string   `json:"calendarStartDate"`
		CalendarEndDate     string   `json:"calendarEndDate"`
		Calendar            []string `json:"calendar"`
	} `json:"leagues"`
	Season struct {
		Type int `json:"type"`
		Year int `json:"year"`
	} `json:"season"`
	Day struct {
		Date string `json:"date"`
	} `json:"day"`
	Events []struct {
		ID        string `json:"id"`
		UID       string `json:"uid"`
		Date      string `json:"date"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
		Season    struct {
			Year int    `json:"year"`
			Type int    `json:"type"`
			Slug string `json:"slug"`
		} `json:"season"`
		Competitions []struct {
			ID         string `json:"id"`
			UID        string `json:"uid"`
			Date       string `json:"date"`
			Attendance int    `json:"attendance"`
			Type       struct {
				ID           string `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"type"`
			TimeValid   bool `json:"timeValid"`
			NeutralSite bool `json:"neutralSite"`
			Recent      bool `json:"recent"`
			Venue       struct {
				ID       string `json:"id"`
				FullName string `json:"fullName"`
				Address  struct {
					City    string `json:"city"`
					State   string `json:"state"`
					Country string `json:"country"`
				} `json:"address"`
				Capacity int  `json:"capacity"`
				Indoor   bool `json:"indoor"`
			} `json:"venue"`
			Competitors []struct {
				ID       string `json:"id"`
				UID      string `json:"uid"`
				Type     string `json:"type"`
				Order    int    `json:"order"`
				HomeAway string `json:"homeAway"`
				Winner   bool   `json:"winner"`
				Team     struct {
					ID               string `json:"id"`
					UID              string `json:"uid"`
					Location         string `json:"location"`
					Name             string `json:"name"`
					Abbreviation     string `json:"abbreviation"`
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Color            string `json:"color"`
					AlternateColor   string `json:"alternateColor"`
					IsActive         bool   `json:"isActive"`
					Venue            struct {
						ID string `json:"id"`
					} `json:"venue"`
					Links []struct {
						Rel        []string `json:"rel"`
						Href       string   `json:"href"`
						Text       string   `json:"text"`
						IsExternal bool     `json:"isExternal"`
						IsPremium  bool     `json:"isPremium"`
					} `json:"links"`
					Logo string `json:"logo"`
				} `json:"team"`
				Score      string `json:"score"`
				Linescores []struct {
					Value float64 `json:"value"`
				} `json:"linescores"`
				Statistics []struct {
					Name         string `json:"name"`
					Abbreviation string `json:"abbreviation"`
					DisplayValue string `json:"displayValue"`
				} `json:"statistics"`
				Leaders []struct {
					Name             string `json:"name"`
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Abbreviation     string `json:"abbreviation"`
					Leaders          []struct {
						DisplayValue string  `json:"displayValue"`
						Value        float64 `json:"value"`
						Athlete      struct {
							ID          string `json:"id"`
							FullName    string `json:"fullName"`
							DisplayName string `json:"displayName"`
							ShortName   string `json:"shortName"`
							Links       []struct {
								Rel  []string `json:"rel"`
								Href string   `json:"href"`
							} `json:"links"`
							Headshot string `json:"headshot"`
							Jersey   string `json:"jersey"`
							Position struct {
								Abbreviation string `json:"abbreviation"`
							} `json:"position"`
							Team struct {
								ID string `json:"id"`
							} `json:"team"`
							Active bool `json:"active"`
						} `json:"athlete"`
						Team struct {
							ID string `json:"id"`
						} `json:"team"`
					} `json:"leaders"`
				} `json:"leaders"`
				Probables []struct {
					Name             string `json:"name"`
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Abbreviation     string `json:"abbreviation"`
					PlayerID         int    `json:"playerId"`
					Athlete          struct {
						ID          string `json:"id"`
						FullName    string `json:"fullName"`
						DisplayName string `json:"displayName"`
						ShortName   string `json:"shortName"`
						Links       []struct {
							Rel  []string `json:"rel"`
							Href string   `json:"href"`
						} `json:"links"`
						Headshot string `json:"headshot"`
						Jersey   string `json:"jersey"`
						Position string `json:"position"`
						Team     struct {
							ID string `json:"id"`
						} `json:"team"`
					} `json:"athlete"`
					Status struct {
						ID           string `json:"id"`
						Name         string `json:"name"`
						Type         string `json:"type"`
						Abbreviation string `json:"abbreviation"`
					} `json:"status"`
					Statistics []interface{} `json:"statistics"`
				} `json:"probables"`
				Records []struct {
					Name         string `json:"name"`
					Abbreviation string `json:"abbreviation"`
					Type         string `json:"type"`
					Summary      string `json:"summary"`
				} `json:"records"`
			} `json:"competitors"`
			Notes  []interface{} `json:"notes"`
			Status struct {
				Clock        float64 `json:"clock"`
				DisplayClock string  `json:"displayClock"`
				Period       int     `json:"period"`
				Type         struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					State       string `json:"state"`
					Completed   bool   `json:"completed"`
					Description string `json:"description"`
					Detail      string `json:"detail"`
					ShortDetail string `json:"shortDetail"`
				} `json:"type"`
				FeaturedAthletes []struct {
					Name             string `json:"name"`
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Abbreviation     string `json:"abbreviation"`
					PlayerID         int    `json:"playerId"`
					Athlete          struct {
						ID          string `json:"id"`
						FullName    string `json:"fullName"`
						DisplayName string `json:"displayName"`
						ShortName   string `json:"shortName"`
						Links       []struct {
							Rel  []string `json:"rel"`
							Href string   `json:"href"`
						} `json:"links"`
						Headshot string `json:"headshot"`
						Jersey   string `json:"jersey"`
						Position string `json:"position"`
						Team     struct {
							ID string `json:"id"`
						} `json:"team"`
					} `json:"athlete"`
					Team struct {
						ID string `json:"id"`
					} `json:"team"`
					Statistics []struct {
						Name         string `json:"name"`
						Abbreviation string `json:"abbreviation"`
						DisplayValue string `json:"displayValue"`
					} `json:"statistics"`
				} `json:"featuredAthletes"`
			} `json:"status"`
			Broadcasts []struct {
				Market string   `json:"market"`
				Names  []string `json:"names"`
			} `json:"broadcasts"`
			Format struct {
				Regulation struct {
					Periods int `json:"periods"`
				} `json:"regulation"`
			} `json:"format"`
			StartDate     string `json:"startDate"`
			GeoBroadcasts []struct {
				Type struct {
					ID        string `json:"id"`
					ShortName string `json:"shortName"`
				} `json:"type"`
				Market struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"market"`
				Media struct {
					ShortName string `json:"shortName"`
				} `json:"media"`
				Lang   string `json:"lang"`
				Region string `json:"region"`
			} `json:"geoBroadcasts"`
			Headlines []struct {
				Description   string `json:"description"`
				Type          string `json:"type"`
				ShortLinkText string `json:"shortLinkText"`
				Video         []struct {
					ID        int    `json:"id"`
					Source    string `json:"source"`
					Headline  string `json:"headline"`
					Thumbnail string `json:"thumbnail"`
					Duration  int    `json:"duration"`
					Tracking  struct {
						SportName    string `json:"sportName"`
						LeagueName   string `json:"leagueName"`
						CoverageType string `json:"coverageType"`
						TrackingName string `json:"trackingName"`
						TrackingID   string `json:"trackingId"`
					} `json:"tracking"`
					DeviceRestrictions struct {
						Type    string   `json:"type"`
						Devices []string `json:"devices"`
					} `json:"deviceRestrictions"`
					GeoRestrictions struct {
						Type      string   `json:"type"`
						Countries []string `json:"countries"`
					} `json:"geoRestrictions"`
					Links struct {
						API struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							Artwork struct {
								Href string `json:"href"`
							} `json:"artwork"`
						} `json:"api"`
						Web struct {
							Href  string `json:"href"`
							Short struct {
								Href string `json:"href"`
							} `json:"short"`
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
						} `json:"web"`
						Source struct {
							Mezzanine struct {
								Href string `json:"href"`
							} `json:"mezzanine"`
							Flash struct {
								Href string `json:"href"`
							} `json:"flash"`
							Hds struct {
								Href string `json:"href"`
							} `json:"hds"`
							Hls struct {
								Href string `json:"href"`
								Hd   struct {
									Href string `json:"href"`
								} `json:"HD"`
							} `json:"HLS"`
							Hd struct {
								Href string `json:"href"`
							} `json:"HD"`
							Full struct {
								Href string `json:"href"`
							} `json:"full"`
							Href string `json:"href"`
						} `json:"source"`
						Mobile struct {
							Alert struct {
								Href string `json:"href"`
							} `json:"alert"`
							Source struct {
								Href string `json:"href"`
							} `json:"source"`
							Href      string `json:"href"`
							Streaming struct {
								Href string `json:"href"`
							} `json:"streaming"`
							ProgressiveDownload struct {
								Href string `json:"href"`
							} `json:"progressiveDownload"`
						} `json:"mobile"`
					} `json:"links"`
				} `json:"video"`
			} `json:"headlines"`
		} `json:"competitions"`
		Links []struct {
			Language   string   `json:"language"`
			Rel        []string `json:"rel"`
			Href       string   `json:"href"`
			Text       string   `json:"text"`
			ShortText  string   `json:"shortText"`
			IsExternal bool     `json:"isExternal"`
			IsPremium  bool     `json:"isPremium"`
		} `json:"links"`
		Status struct {
			Clock        float64 `json:"clock"`
			DisplayClock string  `json:"displayClock"`
			Period       int     `json:"period"`
			Type         struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				State       string `json:"state"`
				Completed   bool   `json:"completed"`
				Description string `json:"description"`
				Detail      string `json:"detail"`
				ShortDetail string `json:"shortDetail"`
			} `json:"type"`
		} `json:"status"`
	} `json:"events"`
}

type football struct {
	Leagues []struct {
		ID           string `json:"id"`
		UID          string `json:"uid"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
		Slug         string `json:"slug"`
		Season       struct {
			Year      int    `json:"year"`
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Type      struct {
				ID           string `json:"id"`
				Type         int    `json:"type"`
				Name         string `json:"name"`
				Abbreviation string `json:"abbreviation"`
			} `json:"type"`
		} `json:"season"`
		Logos []struct {
			Href        string   `json:"href"`
			Width       int      `json:"width"`
			Height      int      `json:"height"`
			Alt         string   `json:"alt"`
			Rel         []string `json:"rel"`
			LastUpdated string   `json:"lastUpdated"`
		} `json:"logos"`
		CalendarType        string `json:"calendarType"`
		CalendarIsWhitelist bool   `json:"calendarIsWhitelist"`
		CalendarStartDate   string `json:"calendarStartDate"`
		CalendarEndDate     string `json:"calendarEndDate"`
		Calendar            []struct {
			Label     string `json:"label"`
			Value     string `json:"value"`
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Entries   []struct {
				Label          string `json:"label"`
				AlternateLabel string `json:"alternateLabel"`
				Detail         string `json:"detail"`
				Value          string `json:"value"`
				StartDate      string `json:"startDate"`
				EndDate        string `json:"endDate"`
			} `json:"entries"`
		} `json:"calendar"`
	} `json:"leagues"`
	Season struct {
		Type int `json:"type"`
		Year int `json:"year"`
	} `json:"season"`
	Week struct {
		Number     int `json:"number"`
		TeamsOnBye []struct {
			ID               string `json:"id"`
			UID              string `json:"uid"`
			Location         string `json:"location"`
			Name             string `json:"name"`
			Abbreviation     string `json:"abbreviation"`
			DisplayName      string `json:"displayName"`
			ShortDisplayName string `json:"shortDisplayName"`
			IsActive         bool   `json:"isActive"`
			Links            []struct {
				Rel        []string `json:"rel"`
				Href       string   `json:"href"`
				Text       string   `json:"text"`
				IsExternal bool     `json:"isExternal"`
				IsPremium  bool     `json:"isPremium"`
			} `json:"links"`
			Logo string `json:"logo"`
		} `json:"teamsOnBye"`
	} `json:"week"`
	Events []struct {
		ID        string `json:"id"`
		UID       string `json:"uid"`
		Date      string `json:"date"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
		Season    struct {
			Year int    `json:"year"`
			Type int    `json:"type"`
			Slug string `json:"slug"`
		} `json:"season"`
		Week struct {
			Number int `json:"number"`
		} `json:"week"`
		Competitions []struct {
			ID         string `json:"id"`
			UID        string `json:"uid"`
			Date       string `json:"date"`
			Attendance int    `json:"attendance"`
			Type       struct {
				ID           string `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"type"`
			TimeValid             bool `json:"timeValid"`
			NeutralSite           bool `json:"neutralSite"`
			ConferenceCompetition bool `json:"conferenceCompetition"`
			Recent                bool `json:"recent"`
			Venue                 struct {
				ID       string `json:"id"`
				FullName string `json:"fullName"`
				Address  struct {
					City  string `json:"city"`
					State string `json:"state"`
				} `json:"address"`
				Capacity int  `json:"capacity"`
				Indoor   bool `json:"indoor"`
			} `json:"venue"`
			Competitors []struct {
				ID       string `json:"id"`
				UID      string `json:"uid"`
				Type     string `json:"type"`
				Order    int    `json:"order"`
				HomeAway string `json:"homeAway"`
				Winner   bool   `json:"winner"`
				Team     struct {
					ID               string `json:"id"`
					UID              string `json:"uid"`
					Location         string `json:"location"`
					Name             string `json:"name"`
					Abbreviation     string `json:"abbreviation"`
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Color            string `json:"color"`
					AlternateColor   string `json:"alternateColor"`
					IsActive         bool   `json:"isActive"`
					Venue            struct {
						ID string `json:"id"`
					} `json:"venue"`
					Links []struct {
						Rel        []string `json:"rel"`
						Href       string   `json:"href"`
						Text       string   `json:"text"`
						IsExternal bool     `json:"isExternal"`
						IsPremium  bool     `json:"isPremium"`
					} `json:"links"`
					Logo string `json:"logo"`
				} `json:"team"`
				Score      string `json:"score"`
				Linescores []struct {
					Value float64 `json:"value"`
				} `json:"linescores"`
				Statistics []interface{} `json:"statistics"`
				Records    []struct {
					Name         string `json:"name"`
					Abbreviation string `json:"abbreviation,omitempty"`
					Type         string `json:"type"`
					Summary      string `json:"summary"`
				} `json:"records"`
			} `json:"competitors"`
			Notes  []interface{} `json:"notes"`
			Status struct {
				Clock        float64 `json:"clock"`
				DisplayClock string  `json:"displayClock"`
				Period       int     `json:"period"`
				Type         struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					State       string `json:"state"`
					Completed   bool   `json:"completed"`
					Description string `json:"description"`
					Detail      string `json:"detail"`
					ShortDetail string `json:"shortDetail"`
				} `json:"type"`
			} `json:"status"`
			Broadcasts []struct {
				Market string   `json:"market"`
				Names  []string `json:"names"`
			} `json:"broadcasts"`
			Leaders []struct {
				Name             string `json:"name"`
				DisplayName      string `json:"displayName"`
				ShortDisplayName string `json:"shortDisplayName"`
				Abbreviation     string `json:"abbreviation"`
				Leaders          []struct {
					DisplayValue string  `json:"displayValue"`
					Value        float64 `json:"value"`
					Athlete      struct {
						ID          string `json:"id"`
						FullName    string `json:"fullName"`
						DisplayName string `json:"displayName"`
						ShortName   string `json:"shortName"`
						Links       []struct {
							Rel  []string `json:"rel"`
							Href string   `json:"href"`
						} `json:"links"`
						Headshot string `json:"headshot"`
						Jersey   string `json:"jersey"`
						Position struct {
							Abbreviation string `json:"abbreviation"`
						} `json:"position"`
						Team struct {
							ID string `json:"id"`
						} `json:"team"`
						Active bool `json:"active"`
					} `json:"athlete"`
					Team struct {
						ID string `json:"id"`
					} `json:"team"`
				} `json:"leaders"`
			} `json:"leaders"`
			Format struct {
				Regulation struct {
					Periods int `json:"periods"`
				} `json:"regulation"`
			} `json:"format"`
			StartDate     string `json:"startDate"`
			GeoBroadcasts []struct {
				Type struct {
					ID        string `json:"id"`
					ShortName string `json:"shortName"`
				} `json:"type"`
				Market struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"market"`
				Media struct {
					ShortName string `json:"shortName"`
				} `json:"media"`
				Lang   string `json:"lang"`
				Region string `json:"region"`
			} `json:"geoBroadcasts"`
			Headlines []struct {
				Description   string `json:"description"`
				Type          string `json:"type"`
				ShortLinkText string `json:"shortLinkText"`
			} `json:"headlines"`
		} `json:"competitions"`
		Links []struct {
			Language   string   `json:"language"`
			Rel        []string `json:"rel"`
			Href       string   `json:"href"`
			Text       string   `json:"text"`
			ShortText  string   `json:"shortText"`
			IsExternal bool     `json:"isExternal"`
			IsPremium  bool     `json:"isPremium"`
		} `json:"links"`
		Status struct {
			Clock        float64 `json:"clock"`
			DisplayClock string  `json:"displayClock"`
			Period       int     `json:"period"`
			Type         struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				State       string `json:"state"`
				Completed   bool   `json:"completed"`
				Description string `json:"description"`
				Detail      string `json:"detail"`
				ShortDetail string `json:"shortDetail"`
			} `json:"type"`
		} `json:"status"`
	} `json:"events"`
}

type basketball struct {
	Leagues []struct {
		ID           string `json:"id"`
		UID          string `json:"uid"`
		Name         string `json:"name"`
		Abbreviation string `json:"abbreviation"`
		Slug         string `json:"slug"`
		Season       struct {
			Year      int    `json:"year"`
			StartDate string `json:"startDate"`
			EndDate   string `json:"endDate"`
			Type      struct {
				ID           string `json:"id"`
				Type         int    `json:"type"`
				Name         string `json:"name"`
				Abbreviation string `json:"abbreviation"`
			} `json:"type"`
		} `json:"season"`
		Logos []struct {
			Href        string   `json:"href"`
			Width       int      `json:"width"`
			Height      int      `json:"height"`
			Alt         string   `json:"alt"`
			Rel         []string `json:"rel"`
			LastUpdated string   `json:"lastUpdated"`
		} `json:"logos"`
		CalendarType        string   `json:"calendarType"`
		CalendarIsWhitelist bool     `json:"calendarIsWhitelist"`
		CalendarStartDate   string   `json:"calendarStartDate"`
		CalendarEndDate     string   `json:"calendarEndDate"`
		Calendar            []string `json:"calendar"`
	} `json:"leagues"`
	Season struct {
		Type int `json:"type"`
		Year int `json:"year"`
	} `json:"season"`
	Day struct {
		Date string `json:"date"`
	} `json:"day"`
	Events []struct {
		ID        string `json:"id"`
		UID       string `json:"uid"`
		Date      string `json:"date"`
		Name      string `json:"name"`
		ShortName string `json:"shortName"`
		Season    struct {
			Year int    `json:"year"`
			Type int    `json:"type"`
			Slug string `json:"slug"`
		} `json:"season"`
		Competitions []struct {
			ID         string `json:"id"`
			UID        string `json:"uid"`
			Date       string `json:"date"`
			Attendance int    `json:"attendance"`
			Type       struct {
				ID           string `json:"id"`
				Abbreviation string `json:"abbreviation"`
			} `json:"type"`
			TimeValid             bool `json:"timeValid"`
			NeutralSite           bool `json:"neutralSite"`
			ConferenceCompetition bool `json:"conferenceCompetition"`
			Recent                bool `json:"recent"`
			Venue                 struct {
				ID       string `json:"id"`
				FullName string `json:"fullName"`
				Address  struct {
					City  string `json:"city"`
					State string `json:"state"`
				} `json:"address"`
				Capacity int  `json:"capacity"`
				Indoor   bool `json:"indoor"`
			} `json:"venue"`
			Competitors []struct {
				ID       string `json:"id"`
				UID      string `json:"uid"`
				Type     string `json:"type"`
				Order    int    `json:"order"`
				HomeAway string `json:"homeAway"`
				Team     struct {
					ID               string `json:"id"`
					UID              string `json:"uid"`
					Location         string `json:"location"`
					Name             string `json:"name"`
					Abbreviation     string `json:"abbreviation"`
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Color            string `json:"color"`
					AlternateColor   string `json:"alternateColor"`
					IsActive         bool   `json:"isActive"`
					Venue            struct {
						ID string `json:"id"`
					} `json:"venue"`
					Links []struct {
						Rel        []string `json:"rel"`
						Href       string   `json:"href"`
						Text       string   `json:"text"`
						IsExternal bool     `json:"isExternal"`
						IsPremium  bool     `json:"isPremium"`
					} `json:"links"`
					Logo string `json:"logo"`
				} `json:"team"`
				Score      string `json:"score"`
				Statistics []struct {
					Name             string `json:"name"`
					Abbreviation     string `json:"abbreviation"`
					DisplayValue     string `json:"displayValue"`
					RankDisplayValue string `json:"rankDisplayValue,omitempty"`
				} `json:"statistics"`
				Records []struct {
					Name         string `json:"name"`
					Abbreviation string `json:"abbreviation,omitempty"`
					Type         string `json:"type"`
					Summary      string `json:"summary"`
				} `json:"records"`
				Leaders []struct {
					Name             string `json:"name"`
					DisplayName      string `json:"displayName"`
					ShortDisplayName string `json:"shortDisplayName"`
					Abbreviation     string `json:"abbreviation"`
					Leaders          []struct {
						DisplayValue string  `json:"displayValue"`
						Value        float64 `json:"value"`
						Athlete      struct {
							ID          string `json:"id"`
							FullName    string `json:"fullName"`
							DisplayName string `json:"displayName"`
							ShortName   string `json:"shortName"`
							Links       []struct {
								Rel  []string `json:"rel"`
								Href string   `json:"href"`
							} `json:"links"`
							Headshot string `json:"headshot"`
							Jersey   string `json:"jersey"`
							Position struct {
								Abbreviation string `json:"abbreviation"`
							} `json:"position"`
							Team struct {
								ID string `json:"id"`
							} `json:"team"`
							Active bool `json:"active"`
						} `json:"athlete"`
						Team struct {
							ID string `json:"id"`
						} `json:"team"`
					} `json:"leaders"`
				} `json:"leaders"`
			} `json:"competitors"`
			Notes  []interface{} `json:"notes"`
			Status struct {
				Clock        float64 `json:"clock"`
				DisplayClock string  `json:"displayClock"`
				Period       int     `json:"period"`
				Type         struct {
					ID          string `json:"id"`
					Name        string `json:"name"`
					State       string `json:"state"`
					Completed   bool   `json:"completed"`
					Description string `json:"description"`
					Detail      string `json:"detail"`
					ShortDetail string `json:"shortDetail"`
				} `json:"type"`
			} `json:"status"`
			Broadcasts []struct {
				Market string   `json:"market"`
				Names  []string `json:"names"`
			} `json:"broadcasts"`
			Format struct {
				Regulation struct {
					Periods int `json:"periods"`
				} `json:"regulation"`
			} `json:"format"`
			Tickets []struct {
				Summary         string `json:"summary"`
				NumberAvailable int    `json:"numberAvailable"`
				Links           []struct {
					Href string `json:"href"`
				} `json:"links"`
			} `json:"tickets"`
			StartDate     string `json:"startDate"`
			GeoBroadcasts []struct {
				Type struct {
					ID        string `json:"id"`
					ShortName string `json:"shortName"`
				} `json:"type"`
				Market struct {
					ID   string `json:"id"`
					Type string `json:"type"`
				} `json:"market"`
				Media struct {
					ShortName string `json:"shortName"`
				} `json:"media"`
				Lang   string `json:"lang"`
				Region string `json:"region"`
			} `json:"geoBroadcasts"`
			Odds []struct {
				Provider struct {
					ID       string `json:"id"`
					Name     string `json:"name"`
					Priority int    `json:"priority"`
				} `json:"provider"`
				Details   string  `json:"details"`
				OverUnder float64 `json:"overUnder"`
			} `json:"odds"`
		} `json:"competitions"`
		Links []struct {
			Language   string   `json:"language"`
			Rel        []string `json:"rel"`
			Href       string   `json:"href"`
			Text       string   `json:"text"`
			ShortText  string   `json:"shortText"`
			IsExternal bool     `json:"isExternal"`
			IsPremium  bool     `json:"isPremium"`
		} `json:"links"`
		Status struct {
			Clock        float64 `json:"clock"`
			DisplayClock string  `json:"displayClock"`
			Period       int     `json:"period"`
			Type         struct {
				ID          string `json:"id"`
				Name        string `json:"name"`
				State       string `json:"state"`
				Completed   bool   `json:"completed"`
				Description string `json:"description"`
				Detail      string `json:"detail"`
				ShortDetail string `json:"shortDetail"`
			} `json:"type"`
		} `json:"status"`
	} `json:"events"`
}
