package main

import (
	"consumer/config"
	"consumer/repository"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

)

func main() {
	cfg := config.MustLoadConfig()

	db, err := repository.ConnectDb(repository.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		Username: cfg.DB.Username,
		Password: cfg.DB.Password,
		DBname:   cfg.DB.DBname,
		SSL:      cfg.DB.SSL,
	})
	if err != nil {
		fmt.Println("Ошибка подключения к базе данных:", err)
		return
	}
	fmt.Println(cfg.DB.Host, cfg.DB.Port, cfg.DB.Username, cfg.DB.Password, cfg.DB.DBname, cfg.DB.SSL)

	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		fmt.Println("Ошибка создания экземпляра драйвера базы данных:", err)
		return
	}

	// Создание экземпляра объекта Migrate
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", cfg.App.MigrationPath), "postgres", driver)
	if err != nil {
		fmt.Println("Ошибка создания экземпляра объекта Migrate:", err)
		return
	}

	// Применение всех миграций
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		fmt.Println("Ошибка применения миграций:", err)
		return
	}

	fmt.Println("Все миграции успешно применены")
}
