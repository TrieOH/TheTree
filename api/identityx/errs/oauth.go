package errs

// ProviderNotConfiguredError is a known provider with no credentials in the
// scope: the platform's env, or the project's provider row. It stays by
// hand because ForProject, which isn't a member, picks the message.
type ProviderNotConfiguredError struct {
	Provider   string
	ForProject bool
}

func (e *ProviderNotConfiguredError) Error() string {
	if e.ForProject {
		return "provider not configured for this project: " + e.Provider
	}
	return "provider not configured: " + e.Provider
}
