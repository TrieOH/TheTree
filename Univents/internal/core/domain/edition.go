package domain

import (
	"encoding/json"
	"fmt"
	"time"
	"univents/internal/shared/errx"
	"univents/internal/shared/validation"

	"github.com/google/uuid"
	"github.com/hibiken/asynq"
)

type EditionType string

const (
	EditionTypeYear    EditionType = "year"
	EditionTypeSeason  EditionType = "season"
	EditionTypeNumber  EditionType = "number"
	EditionTypeOrdinal EditionType = "ordinal"
	EditionTypeCustom  EditionType = "custom"
)

type EditionStatus string

const (
	EditionStatusDraft     EditionStatus = "draft"
	EditionStatusAnnounced EditionStatus = "announced"
	EditionStatusOpen      EditionStatus = "open"
	EditionStatusOngoing   EditionStatus = "ongoing"
	EditionStatusFinished  EditionStatus = "finished"
	EditionStatusCompleted EditionStatus = "completed"
	EditionStatusCancelled EditionStatus = "cancelled"
	EditionStatusPostponed EditionStatus = "postponed"
)

type EditionMonetaryType string

const (
	EditionMonetaryTypeFree  EditionMonetaryType = "free"
	EditionMonetaryTypePaid  EditionMonetaryType = "paid"
	EditionMonetaryTypeMixed EditionMonetaryType = "mixed"
)

type Edition struct {
	ID                            uuid.UUID           `json:"id"`
	EventID                       uuid.UUID           `json:"event_id"`
	GoauthScopeID                 uuid.UUID           `json:"goauth_scope_id"`
	Type                          EditionType         `json:"type"`
	EditionName                   string              `json:"edition_name"`
	Tagline                       *string             `json:"tagline"`
	Description                   *string             `json:"description"`
	Status                        EditionStatus       `json:"status"`
	MonetaryType                  EditionMonetaryType `json:"monetary_type"`
	RegistrationOpensAt           *time.Time          `json:"registration_opens_at"`
	RegistrationClosesAt          *time.Time          `json:"registration_closes_at"`
	StartsAt                      time.Time           `json:"starts_at"`
	EndsAt                        time.Time           `json:"ends_at"`
	Timezone                      string              `json:"timezone"`
	LocationName                  string              `json:"location_name"`
	LocationAddress               string              `json:"location_address"`
	LogoUrl                       *string             `json:"logo_url"`
	BannerUrl                     *string             `json:"banner_url"`
	ContactEmail                  *string             `json:"contact_email"`
	ContactPhone                  *string             `json:"contact_phone"`
	OrganizerName                 *string             `json:"organizer_name"`
	TriePaymentsCredentialID      *uuid.UUID          `json:"trie_payments_credential_id"`
	TriePaymentsProvider          *string             `json:"trie_payments_provider"`
	TriePaymentsProviderPublicKey *string             `json:"trie_payments_provider_public_key"`
	CreatedBy                     uuid.UUID           `json:"created_by"`
	CreatedAt                     time.Time           `json:"created_at"`
	UpdatedAt                     time.Time           `json:"updated_at"`
	DeletedAt                     *time.Time          `json:"deleted_at"`
}

type CreateEditionSpec struct {
	EventID              uuid.UUID           `json:"event_id"`
	Type                 EditionType         `json:"type"`
	EditionName          string              `json:"edition_name"`
	Tagline              *string             `json:"tagline"`
	Description          *string             `json:"description"`
	Status               EditionStatus       `json:"status"`
	MonetaryType         EditionMonetaryType `json:"monetary_type"`
	RegistrationOpensAt  *time.Time          `json:"registration_opens_at"`
	RegistrationClosesAt *time.Time          `json:"registration_closes_at"`
	StartsAt             time.Time           `json:"starts_at"`
	EndsAt               time.Time           `json:"ends_at"`
	Timezone             string              `json:"timezone"`
	LocationName         string              `json:"location_name"`
	LocationAddress      string              `json:"location_address"`
	LogoUrl              *string             `json:"logo_url"`
	BannerUrl            *string             `json:"banner_url"`
	ContactEmail         *string             `json:"contact_email"`
	ContactPhone         *string             `json:"contact_phone"`
	OrganizerName        *string             `json:"organizer_name"`
	CreatedBy            uuid.UUID           `json:"created_by"`
}

func NewEdition(creatorID uuid.UUID, spec CreateEditionSpec) (*Edition, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return nil, errx.Internal("edition").SetMessage("error generating uuid").SetCause(err)
	}

	e := &Edition{
		ID:                   id,
		EventID:              spec.EventID,
		Type:                 spec.Type,
		EditionName:          spec.EditionName,
		Tagline:              spec.Tagline,
		Description:          spec.Description,
		Status:               EditionStatusDraft,
		MonetaryType:         EditionMonetaryTypeFree,
		RegistrationOpensAt:  spec.RegistrationOpensAt,
		RegistrationClosesAt: spec.RegistrationClosesAt,
		StartsAt:             spec.StartsAt,
		EndsAt:               spec.EndsAt,
		Timezone:             spec.Timezone,
		LocationName:         spec.LocationName,
		LocationAddress:      spec.LocationAddress,
		LogoUrl:              spec.LogoUrl,
		BannerUrl:            spec.BannerUrl,
		ContactEmail:         spec.ContactEmail,
		ContactPhone:         spec.ContactPhone,
		OrganizerName:        spec.OrganizerName,
		CreatedBy:            creatorID,
	}

	if err := e.validate(); err != nil {
		return nil, err
	}

	return e, nil
}

func (e *Edition) validate() error {
	now := time.Now()
	return validation.Run(
		validation.RequireUUID("edition", "event_id", e.EventID),
		validation.RequireUUID("edition", "created_by", e.CreatedBy),
		validation.RequireString("edition", "name", e.EditionName),
		validation.RequireString("edition", "timezone", e.Timezone),
		validation.RequireTime("edition", "starts_at", e.StartsAt),
		validation.RequireTime("edition", "ends_at", e.EndsAt),
		validation.Assert("edition", e.StartsAt.After(now), "start at must not be before now, legacy editions are not supported for now"),
		validation.Assert("edition", e.EndsAt.After(now), "ends at must not be before now, legacy editions are not supported for now"),
		validation.Assert("edition", e.StartsAt.Before(e.EndsAt), "edition start must be before edition end"),
		validation.AssertIf("edition",
			func() bool { return e.RegistrationOpensAt != nil && e.RegistrationClosesAt != nil },
			func() bool { return e.RegistrationOpensAt.Before(*e.RegistrationClosesAt) },
			"registration closing time cant be before registration opening time",
		),
		validation.AssertIf("edition",
			func() bool { return e.RegistrationOpensAt != nil },
			func() bool { return e.RegistrationOpensAt.Before(e.StartsAt) },
			"registration cant open after event start",
		),
		validation.AssertIf("edition",
			func() bool { return e.RegistrationClosesAt != nil },
			func() bool { return e.RegistrationClosesAt.Before(e.EndsAt) },
			"registration cant close after event end",
		),
	)
}

func (e *Edition) AddScope(scopeID uuid.UUID) {
	e.GoauthScopeID = scopeID
}

type EditionAuditAction string

const (
	EditionAuditActionCreated                   EditionAuditAction = "created"
	EditionAuditActionEdited                    EditionAuditAction = "edited"
	EditionAuditActionAnnounced                 EditionAuditAction = "announced"
	EditionAuditActionOpened                    EditionAuditAction = "opened"
	EditionAuditActionStarted                   EditionAuditAction = "started"
	EditionAuditActionCompleted                 EditionAuditAction = "completed"
	EditionAuditActionCancelled                 EditionAuditAction = "cancelled"
	EditionAuditActionPostponed                 EditionAuditAction = "postponed"
	EditionAuditActionDeleted                   EditionAuditAction = "deleted"
	EditionAuditActionRestored                  EditionAuditAction = "restored"
	EditionAuditActionNameChanged               EditionAuditAction = "name_changed"
	EditionAuditActionTaglineChanged            EditionAuditAction = "tagline_changed"
	EditionAuditActionDescriptionChanged        EditionAuditAction = "description_changed"
	EditionAuditActionTypeChanged               EditionAuditAction = "type_changed"
	EditionAuditActionStatusChanged             EditionAuditAction = "status_changed"
	EditionAuditActionVisibleFromChanged        EditionAuditAction = "visible_from_changed"
	EditionAuditActionRegistrationOpensChanged  EditionAuditAction = "registration_opens_changed"
	EditionAuditActionRegistrationClosesChanged EditionAuditAction = "registration_closes_changed"
	EditionAuditActionMonetaryTypeChanged       EditionAuditAction = "monetary_type_changed"
	EditionAuditActionTicketsEnabled            EditionAuditAction = "tickets_enabled"
	EditionAuditActionTicketsDisabled           EditionAuditAction = "tickets_disabled"
	EditionAuditActionTicketTiersEnabled        EditionAuditAction = "ticket_tiers_enabled"
	EditionAuditActionTicketTiersDisabled       EditionAuditAction = "ticket_tiers_disabled"
	EditionAuditActionCapacityChanged           EditionAuditAction = "capacity_changed"
	EditionAuditActionCapacityEnabled           EditionAuditAction = "capacity_enabled"
	EditionAuditActionCapacityDisabled          EditionAuditAction = "capacity_disabled"
	EditionAuditActionProductsEnabled           EditionAuditAction = "products_enabled"
	EditionAuditActionProductsDisabled          EditionAuditAction = "products_disabled"
	EditionAuditActionBundlesEnabled            EditionAuditAction = "bundles_enabled"
	EditionAuditActionBundlesDisabled           EditionAuditAction = "bundles_disabled"
	EditionAuditActionTokensEnabled             EditionAuditAction = "tokens_enabled"
	EditionAuditActionTokensDisabled            EditionAuditAction = "tokens_disabled"
	EditionAuditActionMaxTokensChanged          EditionAuditAction = "max_tokens_changed"
	EditionAuditActionDatesChanged              EditionAuditAction = "dates_changed"
	EditionAuditActionTimezoneChanged           EditionAuditAction = "timezone_changed"
	EditionAuditActionLocationChanged           EditionAuditAction = "location_changed"
	EditionAuditActionLogoUpdated               EditionAuditAction = "logo_updated"
	EditionAuditActionBannerUpdated             EditionAuditAction = "banner_updated"
	EditionAuditActionGalleryEnabled            EditionAuditAction = "gallery_enabled"
	EditionAuditActionGalleryDisabled           EditionAuditAction = "gallery_disabled"
	EditionAuditActionGalleryUpdated            EditionAuditAction = "gallery_updated"
	EditionAuditActionScheduleEnabled           EditionAuditAction = "schedule_enabled"
	EditionAuditActionScheduleDisabled          EditionAuditAction = "schedule_disabled"
	EditionAuditActionActivitiesEnabled         EditionAuditAction = "activities_enabled"
	EditionAuditActionActivitiesDisabled        EditionAuditAction = "activities_disabled"
	EditionAuditActionActivityInterestEnabled   EditionAuditAction = "activity_interest_enabled"
	EditionAuditActionActivityInterestDisabled  EditionAuditAction = "activity_interest_disabled"
	EditionAuditActionPaidActivitiesEnabled     EditionAuditAction = "paid_activities_enabled"
	EditionAuditActionPaidActivitiesDisabled    EditionAuditAction = "paid_activities_disabled"
	EditionAuditActionInterestListEnabled       EditionAuditAction = "interest_list_enabled"
	EditionAuditActionInterestListDisabled      EditionAuditAction = "interest_list_disabled"
	EditionAuditActionInterestListOpensChanged  EditionAuditAction = "interest_list_opens_changed"
	EditionAuditActionCheckoutEnabled           EditionAuditAction = "checkout_enabled"
	EditionAuditActionCheckoutDisabled          EditionAuditAction = "checkout_disabled"
	EditionAuditActionContactUpdated            EditionAuditAction = "contact_updated"
	EditionAuditActionUserRegistered            EditionAuditAction = "user_registered"
	EditionAuditActionRegistrationCancelled     EditionAuditAction = "registration_cancelled"
	EditionAuditActionUserCheckedIn             EditionAuditAction = "user_checked_in"
	EditionAuditActionUserCheckedOut            EditionAuditAction = "user_checked_out"
	EditionAuditActionAttendanceMarked          EditionAuditAction = "attendance_marked"
	EditionAuditActionStatusManuallyChanged     EditionAuditAction = "status_manually_changed"
	EditionAuditActionOwnershipTransferred      EditionAuditAction = "ownership_transferred"
)

const AsynqEditionOpen = "edition:open"
const AsynqEditionStart = "edition:start"
const AsynqEditionFinish = "edition:finish"

type EditionPayload struct {
	EditionID uuid.UUID `json:"edition_id"`
}

func NewOpenEditionTask(editionID uuid.UUID, registrationOpensAt time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(EditionPayload{EditionID: editionID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(AsynqEditionOpen, payload,
		asynq.TaskID(fmt.Sprintf("%s:%s", editionID, AsynqEditionOpen)),
		asynq.ProcessAt(registrationOpensAt),
		asynq.Unique(time.Hour),
		asynq.Retention(7*24*time.Hour),
	), nil
}

func NewStartEditionTask(editionID uuid.UUID, startsAt time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(EditionPayload{EditionID: editionID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(AsynqEditionStart, payload,
		asynq.TaskID(fmt.Sprintf("%s:%s", editionID, AsynqEditionStart)),
		asynq.ProcessAt(startsAt),
		asynq.Unique(time.Hour),
		asynq.Retention(7*24*time.Hour),
	), nil
}

func NewFinishEditionTask(editionID uuid.UUID, endsAt time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(EditionPayload{EditionID: editionID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(AsynqEditionFinish, payload,
		asynq.TaskID(fmt.Sprintf("%s:%s", editionID, AsynqEditionFinish)),
		asynq.ProcessAt(endsAt),
		asynq.Unique(time.Hour),
		asynq.Retention(7*24*time.Hour),
	), nil
}
