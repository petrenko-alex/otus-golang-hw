package rules

type InRule struct{ ValidationLimit }

func (i InRule) Validate(value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (i InRule) GetLimit() ValidationLimit {
	//TODO implement me
	panic("implement me")
}

func (i InRule) GetError() error {
	//TODO implement me
	panic("implement me")
}
