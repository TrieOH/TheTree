package schema

type ErrRegisterOnSchemaDraft struct{}

func (ErrRegisterOnSchemaDraft) Error() string {
	return "can't register to a draft schema"
}

type ErrRegisterOnSchemaArchive struct{}

func (ErrRegisterOnSchemaArchive) Error() string {
	return "can't register to an archived schema"
}

type ErrSchemaNoPublishedVersion struct{}

func (ErrSchemaNoPublishedVersion) Error() string {
	return "schema has no published version"
}

type ErrSchemaVersionMismatch struct{}

func (ErrSchemaVersionMismatch) Error() string {
	return "schema version and supplied version mismatch"
}
