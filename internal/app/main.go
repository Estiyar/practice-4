package app

import (
	"context"
	"fmt"
	"time"

	"practice3go/internal/repository"
	"practice3go/internal/repository/_postgres"
	"practice3go/pkg/modules"
)

func Run() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := initPostgreConfig()
	pg := _postgres.NewPGXDialect(ctx, cfg)

	repos := repository.NewRepositories(pg)

	newID, err := repos.CreateUser(modules.User{
		Name:  "Alice",
		Email: "alice@mail.com",
		Age:   20,
	})
	if err != nil {
		fmt.Println("Create error:", err)
		return
	}

	fmt.Println("Created id:", newID)

	one, err := repos.GetUserByID(newID)
	fmt.Println("GetByID:", one, "err:", err)

	err = repos.UpdateUser(newID, modules.User{
		Name:  "Alice Updated",
		Email: "alice2@mail.com",
		Age:   21,
	})
	fmt.Println("Update err:", err)

	deleted, err := repos.DeleteUserByID(newID)
	fmt.Println("Delete affected:", deleted, "err:", err)
}

func initPostgreConfig() *modules.PostgreConfig {
	return &modules.PostgreConfig{
		Host:        "localhost",
		Port:        "5432",
		Username:    "postgres",
		Password:    "Esti2005",
		DBName:      "mydb",
		SSLMode:     "disable",
		ExecTimeout: 5 * time.Second,
	}
}
