package rules

type RegexpRule struct{ ValidationLimit }

func (r RegexpRule) Validate(value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (r RegexpRule) GetLimit() ValidationLimit {
	//TODO implement me
	panic("implement me")
}

func (r RegexpRule) GetError() error {
	//TODO implement me
	panic("implement me")
}
