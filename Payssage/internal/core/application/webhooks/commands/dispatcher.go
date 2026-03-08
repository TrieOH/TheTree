package commands

import "context"

func (uc *CommandService) Dispatch(ctx context.Context, provider, intentID, event string) error {
	return uc.HandleProviderWebhook(ctx, provider, intentID, event)
}
