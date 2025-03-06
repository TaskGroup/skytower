package skytower

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/TaskGroup/skytower/app/back/pkg/model/skytower"
	"net/http"
	"net/url"
	"strconv"
)

var (
	//  направляет сообщение пользователю
	UrlSendMessage = "/message/send"
	//  предметы в инвентаре
	UrlUserInventory = "/user/objects"
	//  проверяет возможности создания клана
	UrlClanCreateCheck = "/clan/create/check"
	//  создает клан игрока
	UrlClanCreate = "/clan/create"
	// Обновляет бафы в башне по игроку башни
	UrlPlayerBuffUpdate = "/clan/buff"
	// объекты башни
	UrlObjects = "/objects"
)

// Отправка сообщения игроку в Небесную Башню
func (a *RequestToSkyTower) SendMessageInTelegram(p skytower.MessageForTelegram) error {
	type Res struct {
		skytower.DefaultResponse
	}
	apiRes := Res{}
	bodyJson, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("cannot do Marshal SendMessageInTelegram %s: %w", "UrlSendMessage", err)
	}

	return a.sendAndHandleRequest(http.MethodPost, UrlSendMessage, bodyJson, nil, &apiRes)
}

// Объекты в Небесной Башне
func (a *RequestToSkyTower) Objects(ctx context.Context) ([]skytower.Object, error) {
	type Res struct {
		skytower.DefaultResponse
		Data []skytower.Object `json:"data"`
	}
	apiRes := Res{}
	err := a.sendAndHandleRequest(http.MethodGet, UrlObjects, nil, nil, &apiRes)
	if err != nil {
		return nil, err
	}
	return apiRes.Data, nil
}

// инвентарь пользователя в башне
func (a *RequestToSkyTower) UserInventory(idObjSync int64) (*skytower.UserInventory, error) {
	type Res struct {
		skytower.DefaultResponse
		Data skytower.UserInventory `json:"data"`
	}
	apiRes := Res{}
	query := url.Values{}
	query.Add("user_id", strconv.FormatInt(idObjSync, 10))
	err := a.sendAndHandleRequest(http.MethodGet, UrlUserInventory, nil, query, &apiRes)
	if err != nil {
		return nil, err
	}
	return &apiRes.Data, err
}

// получение свободных слотов и кол-во указанного предмета в инвентаре по его названию
func (a *RequestToSkyTower) UserInventoryAndOneItem(ctx context.Context, userObjSync int64, itemName string) (skytower.InventoryAndItem, error) {
	res := skytower.InventoryAndItem{}
	items, err := a.UserInventory(userObjSync)
	if err != nil {
		return res, fmt.Errorf("cannot do UserInventoryAndOneItem %s: %w", itemName, err)
	} else if items == nil {
		return res, fmt.Errorf("UserInventoryAndOneItem %s not found", itemName)
	}
	res.FreeSlot = items.CountFree
	for _, item := range items.ListObject {
		if item.Name == itemName {
			inventoryItem := skytower.InventoryItem{
				Id:    item.IdObject,
				Name:  item.Name,
				Count: item.Counts,
			}
			if item.NameRaces != nil {
				inventoryItem.Race = *item.NameRaces
			}
			res.Items = append(res.Items, inventoryItem)
		}
	}
	return res, nil
}

// Количество монет в инвентаре игрока
func (a *RequestToSkyTower) MoneyByPlayerId(ctx context.Context, idObjSync int64) (int64, error) {
	inventory, err := a.UserInventory(idObjSync)
	if err != nil {
		return 0, err
	} else if inventory == nil {
		return 0, fmt.Errorf("UserInventory returned nil")
	}
	for _, v := range inventory.ListObject {
		if v.Type == "money" {
			return v.Counts, nil
		}
	}
	return 0, nil
}

// Получение всех монстров из Башни
func (a *RequestToSkyTower) Monsters(ctx context.Context) ([]skytower.Monster, error) {
	type Res struct {
		skytower.DefaultResponse
		Data []skytower.Monster `json:"data"`
	}
	apiRes := Res{}
	err := a.sendAndHandleRequest(http.MethodGet, "/enemy/all/get", nil, nil, &apiRes)
	if err != nil {
		return nil, err
	}
	return apiRes.Data, nil
}

// Забирает предмет из инвентаря пользователя
func (a *RequestToSkyTower) UserInventoryDel(idSyncPlayer int64, p skytower.ObjectsInventory) error {
	type Res struct {
		Message string `json:"message"`
		Error   int    `json:"error"`
	}
	apiRes := Res{}
	bodyJson, err := json.Marshal(p)
	if err != nil {
		return fmt.Errorf("cannot do Marshal UserInventoryDel %s: %w", "UrlUserInventory", err)
	}
	urlForSend := UrlUserInventory + "/" + strconv.FormatInt(idSyncPlayer, 10)

	return a.sendAndHandleRequest(http.MethodDelete, urlForSend, bodyJson, nil, &apiRes)
}

// Запрос в SkyTower на списание ресурсов при создании клана
func (a *RequestToSkyTower) CreateClan(idObjSync int64) error {
	if err := a.CreateClanCheck(idObjSync); err != nil {
		return err
	}
	type Res struct {
		Message string `json:"message"`
		Error   int    `json:"error"`
		Data    int    `json:"data"`
	}
	apiRes := Res{}
	query := url.Values{}
	query.Add("id", strconv.FormatInt(idObjSync, 10))
	return a.sendAndHandleRequest(http.MethodPost, UrlClanCreate, nil, query, &apiRes)
}

// обновление бафов в башне по игрокам
func (a *RequestToSkyTower) UpdateBuffsForPlayer(ObjectSyncListId []skytower.OnlyId) []error {
	var res []error
	for _, v := range ObjectSyncListId {
		if err := a.UpdateBuffByPlayer(v.Id); err != nil {
			res = append(res, err)
		}
	}
	return res
}

// Запрос в SkyTower на проверку возможности создания клана
func (a *RequestToSkyTower) CreateClanCheck(idObjSync int64) error {
	type Res struct {
		Message string `json:"message"`
		Error   int    `json:"error"`
		Data    bool   `json:"data"`
	}
	apiRes := Res{}
	query := url.Values{}
	query.Add("id", strconv.FormatInt(idObjSync, 10))
	if err := a.sendRequest(http.MethodGet, UrlClanCreateCheck, nil, query, &apiRes); err != nil {
		return err
	}
	if apiRes.Error == 1 {
		return fmt.Errorf("createClanCheck: error code = %d, error_msg = %s", apiRes.Error, apiRes.Message)
	} else if apiRes.Data == false {
		return fmt.Errorf("недостаточно ресурсов или опыта для создания гильдии. Требуется 50 уровень игрока и 1 500 000 монет (игрок %d)", idObjSync)
	}

	return nil
}

// Обновляет бафы по игрокам
func (a *RequestToSkyTower) UpdateBuffByPlayer(playerId int64) error {
	type Res struct {
		Message string `json:"message"`
		Error   int    `json:"error"`
	}
	apiRes := Res{}
	urlForSend := UrlPlayerBuffUpdate + "/" + strconv.FormatInt(playerId, 10)
	return a.sendAndHandleRequest(http.MethodPost, urlForSend, nil, nil, &apiRes)
}

func (a *RequestToSkyTower) Players(ctx context.Context) ([]skytower.Player, error) {
	type Res struct {
		Error   int               `json:"error"`
		Message *string           `json:"message"`
		Data    []skytower.Player `json:"data"`
	}
	apiRes := Res{}
	if err := a.sendRequest(http.MethodGet, "/user/get", nil, nil, &apiRes); err != nil {
		return nil, err
	}
	if apiRes.Error == 1 {
		return nil, fmt.Errorf("Players: error code = %d, error_msg = %s", apiRes.Error, apiRes.Message)
	}
	return apiRes.Data, nil
}
