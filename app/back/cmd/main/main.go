package main

import (
	"context"
	"fmt"
	"github.com/TaskGroup/skytower/app/back/config"
	"github.com/TaskGroup/skytower/app/back/pkg/model/skytower"
	doSky "github.com/TaskGroup/skytower/app/back/pkg/service/skytower"
	"log"
)

type SkyTower struct {
	Sky ISkyTower
}

type ISkyTower interface {
	Objects(ctx context.Context) ([]skytower.Object, error)
}

func main() {
	ctx := context.Background()
	cfg := config.MustLoad()
	sky, err := newSkyTower(cfg.ApiAuth.Host, cfg.ApiAuth.Login, cfg.ApiAuth.Password)
	if err != nil {
		log.Fatal(err.Error())
	}

	obList, err := sky.Sky.Objects(ctx)
	if err != nil {
		log.Fatal(err.Error())
		return
	}
	fmt.Println("len(obList): ", len(obList))
	fmt.Println("stopped example")
}

func newSkyTower(host, username, pwd string) (SkyTower, error) {
	st, err := doSky.New(host, username, pwd)
	if err != nil {
		err = fmt.Errorf("Error creating to_sky_tower: " + err.Error())
		return SkyTower{}, err
	}
	if st == nil {
		err = fmt.Errorf("Error creating to_sky_tower: nil to_sky_tower")
		return SkyTower{}, err
	}
	return SkyTower{
		Sky: st,
	}, err
}
