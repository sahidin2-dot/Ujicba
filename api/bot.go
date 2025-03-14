package api

import (
	"encoding/json"
	"fmt"
	"github.com/Abishnoi69/Force-Sub-Bot/FallenSub/modules"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/PaulSonOfLars/gotgbot/v2"
)

var (
	allowedTokens    = strings.Split(os.Getenv("BOT_TOKEN"), " ")
	lenAllowedTokens = len(allowedTokens)
)

const (
	statusCodeSuccess = 200
)

// Bot handles all incoming traffic from webhooks.
func Handler(w http.ResponseWriter, r *http.Request) {
    // ... (kode fungsi tetap sama)
}

	// Validate URL path
	if len(split) < 2 {
		http.Error(w, "URL path too short", http.StatusBadRequest)
		return
	}

	botToken := split[len(split)-2]

	// Create the bot object
	bot, err := gotgbot.NewBot(botToken, &gotgbot.BotOpts{DisableTokenCheck: false})
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating bot: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if the bot token is allowed
	if lenAllowedTokens > 0 && allowedTokens[0] != "" && !findInStringSlice(allowedTokens, botToken) {
		_, _ = bot.DeleteWebhook(&gotgbot.DeleteWebhookOpts{DropPendingUpdates: true}) // It doesn't matter if it errors
		w.WriteHeader(statusCodeSuccess)
		return
	}

	// Read and parse the update
	var update gotgbot.Update
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading request body: %v", err), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &update)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error unmarshalling body: %v", err), http.StatusBadRequest)
		return
	}

	// Set the bot username
	bot.Username = split[len(split)-1]

	// Dispatch the update to the appropriate handler
	err = modules.Dispatcher.ProcessUpdate(bot, &update, map[string]any{})
	if err != nil {
		fmt.Printf("Error while processing update: %v\n", err)
	}

	w.WriteHeader(statusCodeSuccess)
}

// Helper function to check if a value is present in a string slice
func findInStringSlice(slice []string, val string) bool {
    for _, item := range slice {
        if item == val {
            return true
        }
    }
    return false
}

