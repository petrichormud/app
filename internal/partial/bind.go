package partial

import (
	"fmt"
	"html/template"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
	"petrichormud.com/app/internal/email"
)

var BindLoginErr = fiber.Map{
	"NoticeSectionID": "register-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"The username and password you entered couldn't be verified.",
		"Please try again.",
	},
}

var BindRegisterErrInternal = fiber.Map{
	"NoticeSectionID": "register-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"Something's gone terribly wrong.",
	},
	"RefreshButton": true,
}

var BindRegisterErrInvalidUsername = fiber.Map{
	"NoticeSectionID": "register-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid username.",
		"Please follow the prompts and try again.",
	},
}

var BindRegisterErrInvalidPassword = fiber.Map{
	"NoticeSectionID": "register-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid password.",
		"Please follow the prompts and try again.",
	},
}

var BindRegisterErrInvalidConfirmPassword = fiber.Map{
	"NoticeSectionID": "register-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"That password and password confirmation don't match.",
		"Please re-enter your password confirmation.",
	},
}

var BindRegisterErrConflict = fiber.Map{
	"NoticeSectionID": "register-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"Sorry! That username is already taken.",
		"Please try a different username.",
	},
}

var BindRecoverUsernameErrInvalid = fiber.Map{
	"NoticeSectionID": "recover-username-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid email address.",
		"Please try again.",
	},
}

var BindRecoverUsernameErrInternal = fiber.Map{
	"NoticeSectionID": "recover-username-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"Something's gone horribly wrong.",
	},
	"RefreshButton": true,
}

var BindRecoverPasswordErrInvalidEmail = fiber.Map{
	"NoticeSectionID": "recover-password-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid email address.",
		"Please try again.",
	},
}

var BindRecoverPasswordErrInvalidUsername = fiber.Map{
	"NoticeSectionID": "recover-password-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"What you entered isn't a valid username.",
		"Please try again.",
	},
}

var BindRecoverPasswordErrInternal = fiber.Map{
	"NoticeSectionID": "recover-password-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"Something's gone horribly wrong.",
	},
	"RefreshButton": true,
}

var BindResetPasswordErr = fiber.Map{
	"NoticeSectionID": "login-err",
	"SectionClass":    "pt-4",
	"NoticeText": []string{
		"The username and password you entered couldn't be verified.",
		"Please try again.",
	},
}

var BindProfileAddEmailErrUnauthorized = fiber.Map{
	"NoticeSectionID": "add-email-error",
	"SectionClass":    "pt-2 w-[60%]",
	"NoticeText": []string{
		"Your session has expired.",
	},
	"RefreshButton": true,
}

var BindProfileAddEmailErrInternal = fiber.Map{
	"NoticeSectionID": "add-email-error",
	"SectionClass":    "pt-2 w-[60%]",
	"NoticeText": []string{
		"Something's gone horribly wrong.",
	},
	"RefreshButton": true,
	"NoticeIcon":    true,
}

var BindProfileAddEmailErrInvalid = fiber.Map{
	"NoticeSectionID": "add-email-error",
	"SectionClass":    "pt-2 w-[60%]",
	"NoticeText": []string{
		"What you entered isn't a valid email address.",
		"Please try again.",
	},
	"NoticeIcon": true,
}

func BindProfileAddEmailErrTooMany() fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "You've already added the maximum %d emails.", email.MaxCount)
	return fiber.Map{
		"NoticeSectionID": "add-email-error",
		"SectionClass":    "pt-2 w-[60%]",
		"NoticeText": []string{
			sb.String(),
		},
		"NoticeIcon": true,
	}
}

func BindProfileAddEmailErrConflict(email string) fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "<span class=\"font-semibold\">%s</span> is already in use.", email)

	return fiber.Map{
		"NoticeSectionID": "add-email-error",
		"SectionClass":    "pt-2 w-[60%]",
		"NoticeText": []template.HTML{
			template.HTML(sb.String()),
			template.HTML("Please try a different address."),
		},
		"NoticeIcon": true,
	}
}

var BindProfileEditEmailErrUnauthorized = fiber.Map{
	"NoticeSectionID": "profile-email-error",
	"SectionClass":    "pt-2 w-[60%]",
	"NoticeText": []string{
		"Your session has expired.",
	},
	"RefreshButton": true,
	"NoticeIcon":    true,
}

var BindProfileEditEmailErrInternal = fiber.Map{
	"NoticeSectionID": "profile-email-error",
	"SectionClass":    "pt-2 w-[60%]",
	"NoticeText": []string{
		"Something's gone terribly wrong.",
	},
	"RefreshButton": true,
	"NoticeIcon":    true,
}

var BindProfileEditEmailErrInvalid = fiber.Map{
	"NoticeSectionID": "profile-email-error",
	"SectionClass":    "pt-2 w-[60%]",
	"NoticeText": []string{
		"What you entered isn't a valid email address.",
		"Please try again.",
	},
	"NoticeIcon": true,
}

func BindProfileEditEmailErrConflict(email string) fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "<span class=\"font-semibold\">%s</span> is already in use.", email)

	return fiber.Map{
		"NoticeSectionID": "profile-email-error",
		"SectionClass":    "pt-2 w-[60%]",
		"NoticeText": []template.HTML{
			template.HTML(sb.String()),
			template.HTML("Please try a different address."),
		},
		"NoticeIcon": true,
	}
}

func BindProfileEditEmailErrConflictSame(email string) fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "This email is already set to <span class=\"font-semibold\">%s</span>.", email)

	return fiber.Map{
		"NoticeSectionID": "profile-email-error",
		"SectionClass":    "pt-2 w-[60%]",
		"NoticeText": []template.HTML{
			template.HTML(sb.String()),
			template.HTML("If you'd like to edit this email, choose a different address."),
		},
		"NoticeIcon": true,
	}
}

var BindProfileDeleteEmailErrUnauthorized = fiber.Map{
	"NoticeSectionID": "profile-email-error",
	"SectionClass":    "pt-2 w-[60%]",
	"NoticeText": []string{
		"Your session has expired.",
	},
	"RefreshButton": true,
	"NoticeIcon":    true,
}

var BindProfileDeleteEmailErrInternal = fiber.Map{
	"NoticeSectionID": "profile-email-error",
	"SectionClass":    "pt-2 w-[60%]",
	"NoticeText": []string{
		"Something's gone terribly wrong.",
	},
	"RefreshButton": true,
	"NoticeIcon":    true,
}

var BindProfileEmailResendVerificationErrNoID = fiber.Map{
	"NoticeSectionID": "profile-email-error",
	"SectionClass":    "htmx-trade-indicator-out",
	"NoticeText": []string{
		"Something's gone terribly wrong.",
	},
	"RefreshButton": true,
	"NoticeIcon":    true,
}

func BindProfileEmailResendVerificationErrInternal(id int64) fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "email-verified-status-%d", id)
	return fiber.Map{
		"NoticeSectionID": sb.String(),
		"SectionClass":    "htmx-trade-indicator-out",
		"NoticeText": []string{
			"Something's gone terribly wrong.",
		},
		"RefreshButton": true,
		"NoticeIcon":    true,
	}
}

func BindProfileEmailResendVerificationErrForbiddenAlreadyVerified(id int64) fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "email-verified-status-%d", id)
	return fiber.Map{
		"NoticeSectionID": sb.String(),
		"SectionClass":    "htmx-trade-indicator-out",
		"NoticeText": []string{
			"This email is already verified by another user.",
		},
		"RefreshButton": true,
		"NoticeIcon":    true,
	}
}

func BindProfileEmailResendVerificationInfoConflict(id int64) fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "email-verified-status-%d", id)
	return fiber.Map{
		"NoticeSectionID": sb.String(),
		"SectionClass":    "htmx-trade-indicator-out",
		"NoticeText": []string{
			"This email is already verified.",
		},
		"RefreshButton": true,
		"NoticeIcon":    true,
	}
}

func BindProfileEmailResendVerificationSuccess(id int64) fiber.Map {
	var sb strings.Builder
	fmt.Fprintf(&sb, "email-verified-status-%d", id)
	return fiber.Map{
		"NoticeSectionID": sb.String(),
		"SectionClass":    "htmx-trade-indicator-out",
		"NoticeText": []string{
			"Success!",
			"Please check your inbox for the verification link.",
		},
		"NoticeIcon": true,
	}
}

type BindNoticeSectionParams struct {
	SectionID     string
	SectionClass  string
	NoticeText    []string
	RefreshButton bool
	NoticeIcon    bool
	Success       bool
	Error         bool
	Warn          bool
}

func BindNoticeSection(p BindNoticeSectionParams) fiber.Map {
	return fiber.Map{
		"NoticeSectionID": p.SectionID,
		"SectionClass":    p.SectionClass,
		"NoticeText":      p.NoticeText,
		"RefreshButton":   p.RefreshButton,
		"NoticeIcon":      p.NoticeIcon,
		"Success":         p.Success,
		"Error":           p.Error,
		"Warn":            p.Warn,
	}
}
