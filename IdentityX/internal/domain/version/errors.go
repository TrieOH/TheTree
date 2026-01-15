package version

type ErrRegisterOnVersionDraft struct{}

func (ErrRegisterOnVersionDraft) Error() string {
	return "can't register to a draft schema version"
}

type ErrRegisterOnVersionArchive struct{}

func (ErrRegisterOnVersionArchive) Error() string {
	return "can't register to an archived schema version"
}
