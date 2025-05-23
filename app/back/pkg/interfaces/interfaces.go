package interfaces

import (
	"context"
	"github.com/TaskGroup/skytower/app/back/pkg/model/skytower"
)

type ISkyTowerProvider interface {
	SendMessageInTelegram(p skytower.MessageForTelegram) error
	UserInventory(idObjSync int64) (*skytower.UserInventory, error)
	MoneyByPlayerId(ctx context.Context, idObjSync int64) (int64, error)
	UserInventoryAndOneItem(ctx context.Context, userObjSync int64, itemName string) (skytower.InventoryAndItem, error)
	Objects(ctx context.Context) ([]skytower.Object, error)
	Players(ctx context.Context) ([]skytower.Player, error)
	Monsters(ctx context.Context) ([]skytower.Monster, error)
	UpdateBuffByPlayer(playerId int64) error
}
