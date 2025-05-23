package terx

func FilterCommand(command string) func(c *Ctx) bool {
	return func(c *Ctx) bool {
		return c.Command == command
	}
}

func FilterNotCommand() func(c *Ctx) bool {
	return func(c *Ctx) bool {
		return c.Command == ""
	}
}

func FilterIsMessage() func(c *Ctx) bool {
	return func(c *Ctx) bool {
		return c.Update.Message != nil
	}
}

func FilterIsCallback() func(c *Ctx) bool {
	return func(c *Ctx) bool {
		return c.Update.CallbackQuery != nil
	}
}
