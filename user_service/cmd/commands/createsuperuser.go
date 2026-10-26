package commands

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/FatAndreasbot/go_project/user_service/cmd/project_init"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/FatAndreasbot/go_project/user_service/infra/config"
	"github.com/google/uuid"
)

func createsuperuser(options *project_init.LaunchOptions) error {
	config, err := config.LoadConfig(options.EnvPath)
	if err != nil {
		fmt.Println(os.Getwd())
		log.Default().Println(err)
		return err
	}

	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("email:")
	scanner.Scan()
	email := scanner.Text()
	fmt.Println("password:")
	scanner.Scan()
	password := scanner.Text()

	userRepo, groupRepo, _, _, err := project_init.SetupOutgoingAdapters(config)
	if err != nil {
		log.Default().Println(err)
		return err
	}

	_, _, passwordEncrypter := project_init.SetupBaseHandlers(config, userRepo, groupRepo)

	adminGroup, err := groupRepo.GetGroupByName(ctx, "admin")
	if err != nil {
		log.Default().Println(err)
		return err
	}

	pwordHash, err := passwordEncrypter.EncryptPassword(password)

	userRepo.CreateUser(ctx, &models.UserData{
		Email:          email,
		Password:       password,
		IdempotencyKey: uuid.New(),
		GroupID:        adminGroup.ID,
	}, pwordHash)

	fmt.Println("user was created")
	return nil
}
