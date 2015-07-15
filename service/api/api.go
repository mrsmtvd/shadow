package api

func (s *ApiService) GetApiProcedures() []ApiProcedure {
	return []ApiProcedure{
		&PingProcedure{},
		&VersionProcedure{},
	}
}
