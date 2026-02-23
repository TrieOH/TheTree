package domain

import (
	"encoding/json"
	"fmt"
	"time"
	"univents/internal/shared/errx"

	"github.com/MintzyG/fail/v3"
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
	ID                   uuid.UUID           `json:"id"`
	EventID              uuid.UUID           `json:"event_id"`
	GoauthScopeID        uuid.UUID           `json:"goauth_scope_id"`
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
	CreatedAt            time.Time           `json:"created_at"`
	UpdatedAt            time.Time           `json:"updated_at"`
	DeletedAt            *time.Time          `json:"deleted_at"`
}

type CreateEditionSpec struct {
	EventID              uuid.UUID           `json:"event_id"`
	GoAuthEventScopeID   uuid.UUID           `json:"go_auth_event_scope_id"`
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
		return nil, fail.New(errx.SYSUUIDV7GenerationError).WithArgs("NewEdition")
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

// FIXME make me unique errors
func (e *Edition) validate() error {
	if e.EventID == uuid.Nil {
		return fail.New(errx.EditionInvalidID).WithArgs(uuid.Nil.String())
	}

	if e.EditionName == "" {
		return fail.New(errx.EditionValidationFailed).Trace("edition_name is required")
	}

	if e.StartsAt.IsZero() || e.EndsAt.IsZero() {
		return fail.New(errx.EditionValidationFailed).Trace("start and end times are required")
	}

	if !e.StartsAt.Before(e.EndsAt) {
		return fail.New(errx.EditionValidationFailed).Trace("starts_must_be_before_ends")
	}

	if e.RegistrationOpensAt != nil {
		if e.RegistrationClosesAt != nil {
			if e.RegistrationClosesAt.Before(*e.RegistrationOpensAt) {
				return fail.New(errx.EditionValidationFailed).Trace("registration closing time cant be before registration opening time")
			}
		}

		if e.StartsAt.Before(*e.RegistrationOpensAt) {
			return fail.New(errx.EditionValidationFailed).Trace("registration cant open after event start")
		}
	}

	if e.RegistrationClosesAt != nil {
		if e.EndsAt.Before(*e.RegistrationClosesAt) {
			return fail.New(errx.EditionValidationFailed).Trace("registration cant close after event end")
		}
	}

	if e.Timezone == "" {
		return fail.New(errx.EditionValidationFailed).Trace("timezone is required")
	}

	return nil
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
const AsynqEditionClose = "edition:close"
const AsynqEditionStart = "edition:start"
const AsynqEditionEnd = "edition:end"

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
	), nil
}

func NewEndEditionTask(editionID uuid.UUID, endsAt time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(EditionPayload{EditionID: editionID})
	if err != nil {
		return nil, err
	}

	return asynq.NewTask(AsynqEditionEnd, payload,
		asynq.TaskID(fmt.Sprintf("%s:%s", editionID, AsynqEditionEnd)),
		asynq.ProcessAt(endsAt),
		asynq.Unique(time.Hour),
	), nil
}
