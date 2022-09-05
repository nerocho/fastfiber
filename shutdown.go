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

//Grace shutdown with timeout
func GraceRun(app *fiber.App, timeout time.Duration) {
	go func() {
		addr := fmt.Sprintf(":%d", Conf.GetInt("System.Port"))
		if err := app.Listen(addr); err != nil {
			log.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM) //os.Interrupt, os.Kill, syscall.SIGQUIT,
	receive := <-quit
	Logger.Info("signal=", receive.String(), " ProcessKilled")
	TaskWithTimeout(app.Shutdown, time.Second*timeout)

	fmt.Println("Running cleanup tasks...")
	eventmanager.CreateEventManageFactory().FuzzyCall(EventDestroyPrefix)
	fmt.Println("Application successful shutdown.")
}
