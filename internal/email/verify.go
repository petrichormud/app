package email

import (
	"context"
	"fmt"
	"os"

	"github.com/google/uuid"
	redis "github.com/redis/go-redis/v9"
	resend "github.com/resendlabs/resend-go"

	"petrichormud.com/app/internal/shared"
)

const ThirtyMinutesInNanoseconds = 30 * 60 * 1000 * 1000 * 1000

func SendVerificationEmail(i *shared.Interfaces, id int64, email string) error {
	token := uuid.NewString()
	key := VerificationKey(token)
	if err := Cache(i.Redis, key, id); err != nil {
		return err
	}

	if os.Getenv("DISABLE_RESEND") == "true" {
		return nil
	}

	base := os.Getenv("BASE_URL")
	url := fmt.Sprintf("%s/verify?t=%s", base, token)
	params := &resend.SendEmailRequest{
		To:      []string{email},
		From:    "verify@petrichormud.com",
		Html:    fmt.Sprintf("Welcome to PetrichorMUD! Please <a href=%q>click here</a> to verify your email address.", url),
		Subject: fmt.Sprintf("[PetrichorMUD] Verify %s", email),
		ReplyTo: "support@petrichormud.com",
	}
	_, err := i.Resend.Emails.Send(params)
	if err != nil {
		return err
	}

	return nil
}

func VerificationKey(id string) string {
	return fmt.Sprintf("%s:%s", shared.VerifyEmailTokenKey, id)
}

func Cache(r *redis.Client, key string, id int64) error {
	return r.Set(context.Background(), key, id, ThirtyMinutesInNanoseconds).Err()
}
