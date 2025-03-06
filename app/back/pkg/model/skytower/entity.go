package skytower

type DefaultResponse struct {
	Error      int    `json:"error"`
	Message    string `json:"message"`
	StatusCode int    `json:"status_code"`
}

type MessageForTelegram struct {
	Text           string  `json:"text"`
	IdUser         []int64 `json:"id_user"`
	SendToEveryone bool    `json:"sendToEveryone"`
}

type UserInventory struct {
	CountFree  int64 `json:"count_free"`
	ListObject []struct {
		//BodyParts []struct {
		//	Chkey string `json:"chkey"`
		//	Id    int    `json:"id"`
		//	Value string `json:"value"`
		//} `json:"body_parts"`
		Characteristic []ObjectCharacteristic `json:"characteristic"`
		//Cost        int     `json:"cost"`
		//CostSale    int     `json:"cost_sale"`
		Counts int64 `json:"counts"`
		//Description string  `json:"description"`
		IdObject int64 `json:"id_object"`
		//Image    string `json:"image"`
		//	IsSendClan  bool    `json:"isSendClan"`
		Level     *int64  `json:"level"`
		Name      string  `json:"name"`
		NameRaces *string `json:"name_races"`
		//	ParentId    int     `json:"parent_id"`
		Type             string `json:"type"`
		IsNotTransmitted bool   `json:"not_transmitted" db:"not_transmitted"`
	} `json:"list_object"`
}

type ObjectCharacteristic struct {
	IdCharacteristic int64  `json:"id_characteristic"`
	Name             string `json:"name"`
	Value            string `json:"value"`
}

type Object struct {
	Cost           int64  `json:"cost"`
	CostSale       int64  `json:"cost_sale"`
	Id             int64  `json:"id"`
	IdRaces        int64  `json:"id_races"`
	Image          string `json:"image"`
	IsDeleted      bool   `json:"isdeleted"`
	Level          int64  `json:"level"`
	Name           string `json:"name"`
	NameRaces      string `json:"name_races"`
	NotTransmitted bool   `json:"not_transmitted"`
	ParentId       int64  `json:"parent_id"`
	Type           string `json:"type"`
}

type InventoryAndItem struct {
	FreeSlot int64           `json:"free_slot"`
	Items    []InventoryItem `json:"items"`
}

type InventoryItem struct {
	Id    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	Race  string `json:"race" db:"race"`
	Count int64  `json:"count" db:"count"`
}

type ObjectsInventory struct {
	Object []Object `json:"objects"`
}

type OnlyId struct {
	Id int64 `json:"id" db:"id"`
}

type Player struct {
	DateCr     string  `json:"datecr"`
	DateCrVip  *string `json:"datecr_vip" `
	DateEndVip *string `json:"dateend_vip"`
	FirstName  string  `json:"first_name"`
	Gender     *bool   `json:"gender"`
	Id         int64   `json:"id"`
	IdVip      *int64  `json:"id_vip"`
	LastName   string  `json:"last_name"`
	Level      *int    `json:"level"`
	Nickname   *string `json:"nickname"`
	Username   string  `json:"username"`
	RaceStr    string  `json:"name_races"`
}

type Monster struct {
	Id        int64  `json:"id"`
	IsDeleted bool   `json:"isdeleted"`
	Name      string `json:"name"`
	Types     string `json:"types"`
	UrlPhoto  string `json:"url_photo"`
}
