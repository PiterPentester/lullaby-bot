package bot

import (
	"fmt"
	"log"
	"time"

	"github.com/ckayt/lullaby/internal/config"
	"github.com/ckayt/lullaby/internal/system"
	tele "gopkg.in/telebot.v3"
)

type Bot struct {
	teleBot *tele.Bot
	sys     system.Manager
	config  *config.Config
}

func New(cfg *config.Config, sys system.Manager) (*Bot, error) {
	pref := tele.Settings{
		Token:  cfg.Token,
		Poller: &tele.LongPoller{Timeout: 10 * time.Second},
	}

	b, err := tele.NewBot(pref)
	if err != nil {
		return nil, err
	}

	return &Bot{
		teleBot: b,
		sys:     sys,
		config:  cfg,
	}, nil
}

func (b *Bot) Start() {
	// Middleware: Concise Logger
	b.teleBot.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			update := c.Update()
			if update.Message != nil {
				log.Printf("User [%d] sent message: %s", c.Sender().ID, update.Message.Text)
			} else if c.Callback() != nil {
				// Telebot v3 stores the inline button identifier in Unique (often prefixed with \f)
				log.Printf("User [%d] pressed button: %q", c.Sender().ID, c.Callback().Unique)
			}
			return next(c)
		}
	})

	// Middleware: Authorization
	b.teleBot.Use(func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if !b.config.IsAuthorized(c.Sender().ID) {
				log.Printf("Unauthorized access attempt from user ID: %d", c.Sender().ID)
				return c.Send("⛔ Access Denied. Your ID is not in the authorized list.")
			}
			return next(c)
		}
	})

	// Keyboards
	menu := &tele.ReplyMarkup{}
	btnReboot := menu.Data("🔄 Reboot", "reboot_confirm")
	btnPowerOff := menu.Data("🔌 Power Off", "poweroff_confirm")
	menu.Inline(
		menu.Row(btnReboot, btnPowerOff),
	)

	confirmMenu := &tele.ReplyMarkup{}
	btnYesReboot := confirmMenu.Data("✅ Yes, Reboot", "reboot_now")
	btnYesPowerOff := confirmMenu.Data("✅ Yes, Power Off", "poweroff_now")
	btnCancel := confirmMenu.Data("❌ Cancel", "cancel")

	// Handlers
	b.teleBot.Handle("/start", func(c tele.Context) error {
		return c.Send("Orange Pi 5 Manager\nSelect an action:", menu)
	})

	// Reboot Flow
	b.teleBot.Handle(&btnReboot, func(c tele.Context) error {
		confirmMenu.Inline(confirmMenu.Row(btnYesReboot, btnCancel))
		return c.Edit("Are you sure you want to **REBOOT**?", confirmMenu, tele.ModeMarkdown)
	})

	b.teleBot.Handle(&btnYesReboot, func(c tele.Context) error {
		_ = c.Send("🚀 Rebooting now...")
		err := b.sys.Reboot()
		if err != nil {
			return c.Send(fmt.Sprintf("❌ Error during reboot: %v", err))
		}
		return nil
	})

	// Power Off Flow
	b.teleBot.Handle(&btnPowerOff, func(c tele.Context) error {
		confirmMenu.Inline(confirmMenu.Row(btnYesPowerOff, btnCancel))
		return c.Edit("Are you sure you want to **POWER OFF**?", confirmMenu, tele.ModeMarkdown)
	})

	b.teleBot.Handle(&btnYesPowerOff, func(c tele.Context) error {
		_ = c.Send("💤 Powering off now...")
		err := b.sys.PowerOff()
		if err != nil {
			return c.Send(fmt.Sprintf("❌ Error during power off: %v", err))
		}
		return nil
	})

	// Cancel
	b.teleBot.Handle(&btnCancel, func(c tele.Context) error {
		return c.Edit("Action cancelled.\nSelect an action:", menu)
	})

	log.Println("Bot started successfully")
	b.teleBot.Start()
}

func (b *Bot) Stop() {
	b.teleBot.Stop()
}
