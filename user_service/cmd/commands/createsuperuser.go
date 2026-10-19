package commands

import (
	"bufio"
	"context"
	"fmt"
	"os"

	"github.com/FatAndreasbot/go_project/user_service/cmd/project_init"
	"github.com/FatAndreasbot/go_project/user_service/domain/models"
	"github.com/google/uuid"
)

func createsuperuser() error {
	ctx := context.Background()
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("email:")
	scanner.Scan()
	email := scanner.Text()
	fmt.Println("password:")
	scanner.Scan()
	password := scanner.Text()

	userRepo, groupRepo, _, _, err := project_init.SetupOutgoingAdapters()
	if err != nil {
		return err
	}

	adminGroup, err := groupRepo.GetGroupByName(ctx, "admin")
	if err != nil {
		return err
	}

	userRepo.CreateUser(ctx, &models.UserData{
		Email:          email,
		Password:       password,
		IdempotencyKey: uuid.New(),
		GroupID:        adminGroup.ID,
	})

	fmt.Println("user was created")
	return nil
}
