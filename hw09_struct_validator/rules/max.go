package rules

type MaxRule struct{ ValidationLimit }

func (m MaxRule) Validate(value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m MaxRule) GetLimit() ValidationLimit {
	//TODO implement me
	panic("implement me")
}

func (m MaxRule) GetError() error {
	//TODO implement me
	panic("implement me")
}
