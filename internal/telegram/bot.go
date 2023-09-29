package telegram

import (
	"fmt"
	"github.com/Vivirinter/telegram-bot-url/internal/checker"
	"github.com/Vivirinter/telegram-bot-url/pkg/models"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/time/rate"
	"log"
	"net/url"
	"strings"
)

type Bot struct {
	Token   string
	bot     *tgbotapi.BotAPI
	checker *checker.Checker
}

const (
	startMessage        = "Hello! I'm a bot for checking HTTPS in URLs. Use the /check <URL> command to check a URL."
	helpMessage         = "Commands I can help with:\n/start - start working with the bot\n/check <URL> - check URL for HTTPS usage\n/help - get help about available commands."
	tooManyRequests     = "Too many requests. Please try again later."
	enterURL            = "Please specify a URL for checking. For example, /check https://example.com"
	errorCheckHTTPS     = "Let's say an error occurred while checking HTTPS. Please check your URL and try again."
	unrecognizedCommand = "Sorry, I didn't recognize this command. Use /help for a list of available commands."
	invalidURL          = "The text you specified is not a URL. Please try again."
)

var limiter = rate.NewLimiter(10, 10)

func isURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != ""
}
func NewBot(token string) (*Bot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	return &Bot{
		Token:   token,
		bot:     bot,
		checker: checker.NewChecker(),
	}, nil
}

func (b *Bot) handleCommands(update *tgbotapi.Update) {
	message := update.Message
	var err error
	switch message.Command() {
	case "start":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, startMessage)
		_, err = b.bot.Send(msg)
	case "help":
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, helpMessage)
		_, err = b.bot.Send(msg)
	case "check":
		if !limiter.Allow() {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, tooManyRequests)
			_, err = b.bot.Send(msg)
			return
		}

		urlStr := message.CommandArguments()
		if urlStr == "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, enterURL)
			_, err = b.bot.Send(msg)
		} else if !isURL(urlStr) {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, invalidURL)
			_, err = b.bot.Send(msg)
		} else {
			link := models.Link{URL: urlStr}
			if err := b.checker.CheckHTTPS(&link); err != nil {
				log.Printf("Error checking HTTPS for '%v': %v\n", link.URL, err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, errorCheckHTTPS)
				_, err = b.bot.Send(msg)
				return
			}
			if link.IsHTTPS {
				_ = checker.CheckCert(&link)
				link.CertificateInfo.IsCertValid = !link.CertificateInfo.IsSelfSigned
			}

			selectedHeaders := []string{"Content-Type", "Date", "Server"}
			for _, h := range selectedHeaders {
				if v, ok := link.ResponseInfo.Headers[h]; ok {
					link.ResponseInfo.SelectedHeaders = append(link.ResponseInfo.SelectedHeaders, fmt.Sprintf("\t%s: %s", h, strings.Join(v, ", ")))
				}
			}

			msg := tgbotapi.NewMessage(update.Message.Chat.ID, link.String())
			_, err = b.bot.Send(msg)
		}

	default:
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, unrecognizedCommand)
		_, err = b.bot.Send(msg)

		if err != nil {
			log.Printf("Error sending message: %v", err)
		}
	}

}

func (b *Bot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		b.handleCommands(&update)
	}

	return nil
}
