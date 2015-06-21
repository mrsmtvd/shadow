package slack

func (s *SlackService) GetSlackCommands() []SlackCommand {
	return []SlackCommand{
		&VersionCommand{},
		&UnknownCommand{},
		&HelpCommand{},
		&HelloCommand{},
		&PingCommand{},
	}
}
