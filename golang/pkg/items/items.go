package items

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"sort"
	"time"
)

var (
	lastRequestTime time.Time
	legoItems       []LegoItem
	jsonRegexp      *regexp.Regexp = regexp.MustCompile(`var wlJson = (\{.+?\});\r?\n`)
)

func init() {
	lastRequestTime = time.Now()
	legoItems = []LegoItem{}
}

type App struct {
	Client *http.Client
}

type WantedParts map[string]LegoItem

type LegoSet struct {
	ID      int
	Name    string
	PartQty int
}

type Stock map[string]int

type LegoItem struct {
	ID        string
	ItemName  string
	ItemID    int
	ColorID   int
	ColorHex  string
	ColorName string
	WantedQty int
	InStock   int
	ImgURL    string
	LegoSets  []LegoSet
}

type WantedItems struct {
	WantedItems  []LegoItem
	TotalResults int
}

type WantedListInfo struct {
	ID   int
	Name string
}

type WantedLists struct {
	WantedLists []WantedListInfo
}

const (
	baseURL  string = "https://www.bricklink.com"
	pageSize int    = 100
)

func (wantedLists WantedLists) GetWantedItems(app App, ID int) (WantedItems, error) {
	wantedItems := WantedItems{}
	page := 1
	for {
		var itemsOnPage WantedItems

		response, err := app.Client.Get(baseURL + fmt.Sprintf("/v2/wanted/search.page?type=A&wantedMoreID=%d&sort=1&pageSize=%d&page=%d", ID, pageSize, page))
		if err != nil {
			return WantedItems{}, err
		}
		defer response.Body.Close()
		responseData, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return wantedItems, err
		}
		jsonData := jsonRegexp.FindStringSubmatch(string(responseData))[1]
		err = json.Unmarshal([]byte(jsonData), &itemsOnPage)
		if err != nil {
			return WantedItems{}, err
		}

		wantedItems.WantedItems = append(wantedItems.WantedItems, itemsOnPage.WantedItems...)
		wantedItems.TotalResults = itemsOnPage.TotalResults
		if len(wantedItems.WantedItems) == itemsOnPage.TotalResults {
			break
		}
		page = page + 1
	}

	return wantedItems, nil
}

func (wantedLists WantedLists) GetStock(app App) (Stock, error) {
	stock := Stock{}
	inStock, err := wantedLists.GetWantedItems(app, 7154778)
	if err != nil {
		return nil, err
	}

	for _, wantedItem := range inStock.WantedItems {
		stock[wantedItem.GetID()] = wantedItem.WantedQty
	}

	return stock, nil
}

func (legoItem LegoItem) GetID() string {
	return fmt.Sprintf("%d-%d", legoItem.ItemID, legoItem.ColorID)
}

func (app *App) login() error {
	client := app.Client
	loginURL := baseURL + "/ajax/renovate/login.ajax"

	data := url.Values{
		"pageId":   {"LOGIN"},
		"userid":   {"gergely.mark@gmail.com"},
		"password": {"As76Gaard"},
	}

	response, err := client.PostForm(loginURL, data)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	_, err = ioutil.ReadAll(response.Body)
	return err
}

func (app App) NeededItems() ([]LegoItem, error) {
	if len(legoItems) != 0 && time.Now().Before(lastRequestTime.Add(5*time.Minute)) {
		return legoItems, nil
	}

	response, err := app.Client.Get(baseURL + "/v2/wanted/list.page")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	r := regexp.MustCompile(`var wlJson = (\{.+?\});\r?\n`)
	match := r.FindStringSubmatch(string(responseData))
	jsonData := "{}"
	if len(match) > 0 {
		jsonData = match[1]
	}
	log.Println("Status:", response.Status)
	var wantedLists WantedLists
	err = json.Unmarshal([]byte(jsonData), &wantedLists)
	if err != nil {
		return nil, err
	}

	wantedParts := WantedParts{}

	stock, err := wantedLists.GetStock(app)
	if err != nil {
		return nil, err
	}

	for _, wantedList := range wantedLists.WantedLists {
		if wantedList.Name == "Default Wanted List" || wantedList.Name == "raktár készlet" {
			continue
		}

		wantedItems, err := wantedLists.GetWantedItems(app, wantedList.ID)
		if err != nil {
			return nil, err
		}

		for _, wantedItem := range wantedItems.WantedItems {
			legoItem, ok := wantedParts[wantedItem.GetID()]
			if !ok {
				legoItem = wantedItem
				legoItem.ID = legoItem.GetID()
				legoItem.InStock = stock[legoItem.GetID()]
			} else {
				legoItem.WantedQty = legoItem.WantedQty + wantedItem.WantedQty
			}
			legoItem.LegoSets = append(legoItem.LegoSets, LegoSet{
				ID:      wantedList.ID,
				Name:    wantedList.Name,
				PartQty: wantedItem.WantedQty,
			})
			wantedParts[wantedItem.GetID()] = legoItem
		}
	}

	legoItems = make([]LegoItem, 0, len(wantedParts))
	for _, item := range wantedParts {
		legoItems = append(legoItems, item)
	}
	sort.Slice(legoItems, func(i, j int) bool {
		return legoItems[i].WantedQty > legoItems[j].WantedQty
	})
	lastRequestTime = time.Now()
	return legoItems, nil
}

func New() (App, error) {
	jar, _ := cookiejar.New(nil)

	app := App{
		Client: &http.Client{Jar: jar},
	}

	err := app.login()
	if err != nil {
		return App{}, err
	}

	return app, nil
}
