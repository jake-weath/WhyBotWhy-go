package database_client

import (
	"github.com/jake-weath/whybotwhy_go/pkg/command/command_type"
	"github.com/jake-weath/whybotwhy_go/pkg/database_client/model"
	"gorm.io/gorm"
)

var baseCounters = []model.Counter{
	{Name: "deaths", Count: 0},
	{Name: "boops", Count: 0},
}

var baseCommandTypes = []model.CommandType{
	{Name: command_type.TextCommandType},
	{Name: command_type.IncrementCountCommandType},
	{Name: command_type.IncrementCountByUserCommandType},
	{Name: command_type.SetCountCommandType},
	{Name: command_type.AddTextCommandType},
	{Name: command_type.RemoveTextCommandType},
	{Name: command_type.UserEnteredTextCommandType},
	{Name: command_type.AddQuoteCommandType},
}

var baseCommands = []model.Command{
	{Name: "whyme",
		CommandType: model.CommandType{Name: command_type.TextCommandType},
		CommandTexts: []model.CommandText{
			{Text: "WHY {{.chatUserName}} WHY!?"},
		},
	},
	{Name: "death",
		CommandType: model.CommandType{Name: command_type.IncrementCountCommandType},
		CommandTexts: []model.CommandText{
			{Text: "{{.streamName}} has died embarrassingly {{.count}} times on stream!"},
		},
		Counter: model.Counter{Name: "deaths"},
	},
	{Name: "setdeaths",
		CommandType: model.CommandType{Name: command_type.SetCountCommandType},
		CommandTexts: []model.CommandText{
			{Text: "Deaths set to {{.count}}."},
		},
		Counter:         model.Counter{Name: "deaths"},
		IsModeratorOnly: true,
	},
	{Name: "boop",
		CommandType: model.CommandType{Name: command_type.IncrementCountByUserCommandType},
		CommandTexts: []model.CommandText{
			{Text: "{{.chatUserName}} booped the snoot! The snoot has been booped {{.count}} times."},
		},
		Counter: model.Counter{Name: "boops"},
	},
	{Name: "addcommand",
		CommandType: model.CommandType{Name: command_type.AddTextCommandType},
		CommandTexts: []model.CommandText{
			{Text: "Command added."},
		},
		IsModeratorOnly: true,
	},
	{Name: "removecommand",
		CommandType: model.CommandType{Name: command_type.RemoveTextCommandType},
		CommandTexts: []model.CommandText{
			{Text: "Command removed."},
		},
		IsModeratorOnly: true,
	},
	{Name: "rules",
		CommandType: model.CommandType{Name: command_type.TextCommandType},
		CommandTexts: []model.CommandText{
			{Text: "Please remember the channel rules:"}, //TODO: Come up with rules timer
			{Text: "1. Be kind"},
			{Text: "2. No politics or religion"},
			{Text: "3. No spam "},
			{Text: "4. Only backseat if I ask for it"},
		},
	},
	{Name: "commands",
		CommandType: model.CommandType{Name: command_type.TextCommandType},
		CommandTexts: []model.CommandText{
			{Text: `The current commands are: {{.commands}}`},
		},
	},
	{Name: "addquote",
		CommandType: model.CommandType{Name: command_type.AddQuoteCommandType},
		CommandTexts: []model.CommandText{
			{Text: `Quote added`},
		},
	},
	{Name: "quote", //TODO implement specific quotes
		CommandType: model.CommandType{Name: command_type.TextCommandType},
		CommandTexts: []model.CommandText{
			{Text: `{{.randomQuote}}`},
		},
	},
}

func CreateInitialDatabaseData(db *gorm.DB) error {
	db.AutoMigrate(&model.CounterByUser{})
	db.AutoMigrate(&model.Counter{})
	if err := createInitialDatabaseCountersIfNotExists(db, baseCounters); err != nil {
		return err
	}

	db.AutoMigrate(&model.CommandType{})
	if err := createInitialDatabaseCommandTypesIfNotExists(db, baseCommandTypes); err != nil {
		return err
	}

	db.AutoMigrate(&model.Command{})
	db.AutoMigrate(&model.CommandText{})
	if err := createInitialDatabaseCommandsIfNotExists(db, baseCommands); err != nil {
		return err
	}

	db.AutoMigrate(&model.Quote{})

	return nil
}

func createInitialDatabaseCountersIfNotExists(db *gorm.DB, baseCounters []model.Counter) error {
	for _, counter := range baseCounters {
		if err := db.FirstOrCreate(&counter, counter).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandTypesIfNotExists(db *gorm.DB, baseCommandTypes []model.CommandType) error {
	for _, commandType := range baseCommandTypes {
		if err := db.FirstOrCreate(&commandType, commandType).Error; err != nil {
			return err
		}
	}
	return nil
}

func createInitialDatabaseCommandsIfNotExists(db *gorm.DB, baseCommands []model.Command) error {
	for _, command := range baseCommands {
		if err := db.First(&command.CommandType, command.CommandType).Error; err != nil {
			return err
		}

		if !command.Counter.Equals(model.Counter{}) {
			if err := db.First(&command.Counter, command.Counter).Error; err != nil {
				return err
			}
		}

		if err := db.FirstOrCreate(&command, command).Error; err != nil {
			return err
		}
	}
	return nil
}
