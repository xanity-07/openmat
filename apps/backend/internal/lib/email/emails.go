package email

func (c *Client) SendWelcomeEmail(to string, firstName string) error {
	data := map[string]string{
		"UserFirstName": firstName,
	}

	return c.SendEmail(
		to,
		"Welcome to OpenMat!",
		TemplateWelcome,
		data,
	)
}
