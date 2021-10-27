package fastfiber

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/nerocho/fastfiber/utils/eventmanager"
)

func GraceRun(app *fiber.App) {
	go func() {
		addr := fmt.Sprintf(":%d", Conf.GetInt("System.Port"))
		if err := app.Listen(addr); err != nil {
			log.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //os.Interrupt, os.Kill, syscall.SIGQUIT,
	receive := <-quit
	Logger.Info().Str("signal", receive.String()).Msg("ProcessKilled")
	TaskWithTimeout(app.Shutdown, time.Second*10)

	fmt.Println("Running cleanup tasks...")
	eventmanager.CreateEventManageFactory().FuzzyCall(EventDestroyPrefix)
	fmt.Println("Application successful shutdown.")
}
