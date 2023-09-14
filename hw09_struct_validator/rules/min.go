package rules

type MinRule struct{ ValidationLimit }

func (m MinRule) Validate(value interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (m MinRule) GetLimit() ValidationLimit {
	//TODO implement me
	panic("implement me")
}

func (m MinRule) GetError() error {
	//TODO implement me
	panic("implement me")
}
