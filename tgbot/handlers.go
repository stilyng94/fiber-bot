package tgbot

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// [/start] handler
func onStart(ctx context.Context, b *bot.Bot, update *models.Update, webappURL string) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   "You are welcome to my <b>Thrift store</b>",
		ReplyMarkup: models.InlineKeyboardMarkup{
			InlineKeyboard: [][]models.InlineKeyboardButton{{
				{Text: "Start shopping", WebApp: &models.WebAppInfo{URL: webappURL}},
			}},
		},
		ParseMode: "HTML"})

	if err != nil {
		fmt.Printf("failed to send start message: %v", err)
	}
}

func onCheckout(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.AnswerPreCheckoutQuery(ctx, &bot.AnswerPreCheckoutQueryParams{OK: true, PreCheckoutQueryID: update.PreCheckoutQuery.ID})
	if err != nil {
		fmt.Printf("failed to checkout: %v", err)
	}
}

func onPayment(ctx context.Context, b *bot.Bot, update *models.Update) {
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		Text:      "Thank you for your purchase!",
		ParseMode: models.ParseModeHTML,
	})
	if err != nil {
		fmt.Printf("failed to acknowledge payment: %v", err)
	}
}

// [/help] handler
func onHelp(ctx context.Context, b *bot.Bot, update *models.Update) {
	commands, err := b.GetMyCommands(ctx, nil)
	if err != nil {
		fmt.Printf("failed to getting commands: %v", err)
	}
	var cmds strings.Builder
	for _, cmd := range commands {
		cmds.WriteString(fmt.Sprintf("/%s - %s\n", cmd.Command, cmd.Description))
	}

	_, err = b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:    update.Message.Chat.ID,
		ParseMode: models.ParseModeHTML,
		Text:      fmt.Sprintf("Here are all my commands:\n%s", cmds.String()),
	})
	if err != nil {
		fmt.Printf("failed to get commands: %v", err)
	}
}

// // Create a matcher which only matches text which is not a command.
// func noCommands(msg *gotgbot.Message) bool {
// 	return message.Text(msg) && !message.Command(msg)
// }

// // cancel cancels the conversation.
// func cancel(b *gotgbot.Bot, ctx *ext.Context) error {
// 	_, err := ctx.EffectiveMessage.Reply(b, "Oh, goodbye!", &gotgbot.SendMessageOpts{
// 		ParseMode: "html",
// 	})
// 	if err != nil {
// 		return fmt.Printf("failed to send cancel message: %v", err)
// 	}
// 	defer resetState(ctx.EffectiveChat.Id)
// 	return handlers.EndConversation()
// }

// // restart restarts the conversation.
// func restart(b *gotgbot.Bot, ctx *ext.Context) error {
// 	resetState(ctx.EffectiveChat.Id)
// 	_, err := ctx.EffectiveMessage.Reply(b, "Let's start uploading a product. Please enter the <b>product name</b>:", &gotgbot.SendMessageOpts{ParseMode: "html"})
// 	if err != nil {
// 		return fmt.Printf("failed initAddProduct: %v", err)
// 	}
// 	return handlers.NextConversationState(string(StateName))
// }

// func initAddProduct(b *gotgbot.Bot, ctx *ext.Context, envConfig config.EnvConfig) error {
// 	resetState(ctx.EffectiveChat.Id)
// 	if !slices.Contains(envConfig.AllowedTelegramAdmins, fmt.Sprint(ctx.EffectiveSender.User.Id)) {
// 		_, err := ctx.EffectiveMessage.Reply(b, "Wrong command\n. Please enter for /help available commands", nil)
// 		if err != nil {
// 			return fmt.Printf("failed initAddProduct: %v", err)
// 		}
// 		return handlers.EndConversation()
// 	}
// 	_, err := ctx.EffectiveMessage.Reply(b, "Let's start uploading a product. Please enter the <b>product name</b>:", &gotgbot.SendMessageOpts{ParseMode: "html"})
// 	if err != nil {
// 		return fmt.Printf("failed onUploadProduct: %v", err)
// 	}
// 	return handlers.NextConversationState(string(StateName))
// }

// // [productName] sets the product's name.
// func productName(b *gotgbot.Bot, ctx *ext.Context) error {
// 	inputName := strings.TrimSpace(ctx.EffectiveMessage.Text)
// 	if len(inputName) == 0 {
// 		// Retry
// 		return handlers.NextConversationState(string(StateName))
// 	}
// 	_, err := ctx.EffectiveMessage.Reply(b, "Please enter the <b>product description</b>:", &gotgbot.SendMessageOpts{
// 		ParseMode: "html",
// 	})
// 	if err != nil {
// 		return fmt.Printf("failed to send name message: %v", err)
// 	}
// 	product := getProductData(ctx.EffectiveChat.Id)
// 	product.Name = inputName
// 	productData[ctx.EffectiveChat.Id] = product
// 	return handlers.NextConversationState(string(StateDescription))
// }

// // [productDescription] sets the product's description.
// func productDescription(b *gotgbot.Bot, ctx *ext.Context) error {
// 	inputDescription := strings.TrimSpace(ctx.EffectiveMessage.Text)
// 	if len(inputDescription) == 0 {
// 		// If the number is not valid, try again!
// 		ctx.EffectiveMessage.Reply(b, "Description cannot be empty. Could you repeat?", &gotgbot.SendMessageOpts{
// 			ParseMode: "html",
// 		})
// 		// We try the age handler again
// 		return handlers.NextConversationState(string(StateDescription))
// 	}
// 	_, err := ctx.EffectiveMessage.Reply(b, "Please select at <b>most 1 image</b> of the product:", &gotgbot.SendMessageOpts{
// 		ParseMode: "html",
// 	})
// 	if err != nil {
// 		return fmt.Printf("failed to send description message: %v", err)
// 	}
// 	product := getProductData(ctx.EffectiveChat.Id)
// 	product.Description = inputDescription
// 	productData[ctx.EffectiveChat.Id] = product
// 	return handlers.NextConversationState(string(StateImage))
// }

// // [productPrice] sets the product's productPrice.
// func productPrice(b *gotgbot.Bot, ctx *ext.Context) error {
// 	inputPrice := strings.TrimSpace(ctx.EffectiveMessage.Text)
// 	price, err := strconv.Atoi(inputPrice)
// 	if err != nil {
// 		ctx.EffectiveMessage.Reply(b, "This doesn't seem to be a number. Could you repeat?", &gotgbot.SendMessageOpts{
// 			ParseMode: "html",
// 		})
// 		// We try the age handler again
// 		return handlers.NextConversationState(string(StatePrice))
// 	}

// 	product := getProductData(ctx.EffectiveChat.Id)
// 	product.Price = int64(price)
// 	productData[ctx.EffectiveChat.Id] = product

// 	_, err = ctx.EffectiveMessage.Reply(b, "Type [<b>yes</b>] to confirm your product data Or [<b>no</b>] to cancel:", &gotgbot.SendMessageOpts{
// 		ParseMode: "HTML",
// 	})
// 	if err != nil {
// 		return fmt.Printf("failed to send confirmation message: %v", err)
// 	}
// 	return handlers.NextConversationState(string(StateUpload))
// }

// // [productImage] sets the product's image.
// func productImage(b *gotgbot.Bot, ctx *ext.Context) error {
// 	img := ctx.EffectiveMessage.Photo[0]

// 	product := getProductData(ctx.EffectiveChat.Id)
// 	product.FileID = img.FileId
// 	productData[ctx.EffectiveChat.Id] = product

// 	_, err := ctx.EffectiveMessage.Reply(b, "Please enter the <b>product price</b>:", &gotgbot.SendMessageOpts{
// 		ParseMode: "html",
// 	})
// 	if err != nil {
// 		return fmt.Printf("failed to send price message: %v", err)
// 	}
// 	return handlers.NextConversationState(string(StatePrice))
// }

// // upload product.
// func productUpload(b *gotgbot.Bot, ctx *ext.Context, productRepo repos.ProductRepo, envConfig config.EnvConfig) error {
// 	inputUpload := strings.ToLower(strings.TrimSpace(ctx.EffectiveMessage.Text))

// 	if !slices.Contains([]string{"yes", "no"}, inputUpload) {
// 		_, err := ctx.EffectiveMessage.Reply(b, "Type [<b>yes</b>] to confirm your product data Or [<b>no</b>] to cancel:", &gotgbot.SendMessageOpts{
// 			ParseMode: "HTML",
// 		})
// 		if err != nil {
// 			return fmt.Printf("failed to send confirmation message: %v", err)
// 		}
// 		// We try the age handler again
// 		return handlers.NextConversationState(string(StatePrice))
// 	}

// 	if inputUpload != "yes" {
// 		ctx.EffectiveMessage.Reply(b, "Upload cancelled", &gotgbot.SendMessageOpts{
// 			ParseMode: "html",
// 		})
// 		resetState(ctx.EffectiveChat.Id)
// 		return handlers.EndConversation()
// 	}

// 	product := getProductData(ctx.EffectiveChat.Id)
// 	f, err := b.GetFile(product.FileID, nil)
// 	if err != nil {
// 		return fmt.Printf("failed to upload image: %v", err)
// 	}

// 	cld := config.InitCloudinary(envConfig)
// 	uploadCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()
// 	resp, err := cld.Upload.Upload(uploadCtx, f.URL(b, nil), uploader.UploadParams{
// 		UniqueFilename: api.Bool(false),
// 		Overwrite:      api.Bool(true)})
// 	if err != nil {
// 		return fmt.Printf("failed to upload image: %v", err)
// 	}

// 	var cmds strings.Builder
// 	cmds.WriteString(fmt.Sprintf("%s - <b>%s</b>\n", StateName, product.Name))
// 	cmds.WriteString(fmt.Sprintf("%s - <b>%s</b>\n", StateDescription, product.Description))
// 	cmds.WriteString(fmt.Sprintf("%s - <b>%v</b>\n", StatePrice, product.Price))

// 	_, err = ctx.EffectiveMessage.Reply(b, fmt.Sprintf("Here is your product data:\n%s", cmds.String()), &gotgbot.SendMessageOpts{
// 		ParseMode: "HTML",
// 	})
// 	if err != nil {
// 		return fmt.Printf("failed to commit product: %v", err)
// 	}
// 	queryCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
// 	defer cancel()
// 	defer resetState(ctx.EffectiveChat.Id)

// 	productRepo.CreateProduct(queryCtx, repos.HandleAddProductPayload{Title: product.Name, Price: product.Price,
// 		ImageUrl: resp.SecureURL, Description: product.Description})

// 	_, err = ctx.EffectiveMessage.Reply(b, "<b>Upload successful</b>", &gotgbot.SendMessageOpts{
// 		ParseMode: "HTML",
// 	})
// 	if err != nil {
// 		return fmt.Printf("failed to commit product: %v", err)
// 	}
// 	return handlers.EndConversation()
// }

// var productData = make(map[int64]ProductInfo)

// type ProductInfo struct {
// 	Name        string
// 	Description string
// 	Price       int64
// 	FileID      string
// }

// func getProductData(chatID int64) ProductInfo {
// 	return productData[chatID]

// }

// func resetState(chatID int64) {
// 	delete(productData, chatID)
// }

// //TODO: delete chat
