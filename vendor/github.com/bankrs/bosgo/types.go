package bosgo

import (
	"time"
)

type DeveloperCredentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type DeveloperProfile struct {
	Company             string `json:"company"`
	HasProductionAccess bool   `json:"has_production_access"`
}

type ApplicationPage struct {
	Applications []ApplicationMetadata `json:"applications,omitempty"`
}

type ApplicationMetadata struct {
	ApplicationID string    `json:"application_id,omitempty"`
	Label         string    `json:"label,omitempty"`
	CreatedAt     time.Time `json:"created_at,omitempty"`
}

type StatsPeriod struct {
	From   string `json:"from_date"`
	To     string `json:"to_date"`
	Domain string `json:"domain"`
}

type UsersStats struct {
	StatsPeriod
	UsersTotal StatsValue        `json:"users_total"` // with weekly relative change
	UsersToday StatsValue        `json:"users_today"` // with daily relative change
	Stats      []DailyUsersStats `json:"stats"`
}

type StatsValue struct {
	Value int64 `json:"value"`
}

type DailyUsersStats struct {
	Date        string `json:"date"`
	UsersTotal  int64  `json:"users_total"`
	NewUsers    int64  `json:"new_users"`
	ActiveUsers int64  `json:"active_users"`
}

type TransfersStats struct {
	StatsPeriod
	TotalOut StatsMoneyAmount      `json:"total_out"`
	TodayOut StatsMoneyAmount      `json:"today_out"`
	Stats    []DailyTransfersStats `json:"stats"`
}

type DailyTransfersStats struct {
	Date string           `json:"date"`
	Out  StatsMoneyAmount `json:"out"`
}

type MerchantsStats struct {
	StatsPeriod
	Stats []DailyMerchantsStats `json:"stats"`
}

type DailyMerchantsStats struct {
	Date      string      `json:"date"`
	Merchants []NameValue `json:"merchants"`
}

type ProvidersStats struct {
	StatsPeriod
	Stats []DailyProvidersStats `json:"stats"`
}

type DailyProvidersStats struct {
	Date      string      `json:"date"`
	Providers []NameValue `json:"providers"`
}

type StatsMoneyAmount struct {
	Value    float64 `json:"value"`
	Currency string  `json:"currency"`
}

type NameValue struct {
	Name  string `json:"name"`
	Value int64  `json:"value"`
}

type RequestsStats struct {
	StatsPeriod
	RequestsTotal StatsValue           `json:"requests_total"`
	RequestsToday StatsValue           `json:"requests_today"`
	Stats         []DailyRequestsStats `json:"stats"`
}

type DailyRequestsStats struct {
	Date          string `json:"date"`
	RequestsTotal int64  `json:"requests_total"`
}

type CategoryList []Category

type Category struct {
	ID    int64             `json:"id"`
	Names map[string]string `json:"names"`
	Group string            `json:"group"`
}

type ProviderSearchResults []ProviderSearchResult

type ProviderSearchResult struct {
	Score    float64  `json:"score"`
	Provider Provider `json:"provider"`
}

type Provider struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Country     string          `json:"country"`
	URL         string          `json:"url"`
	Address     string          `json:"address"`
	PostalCode  string          `json:"postal_code"`
	Challenges  []ChallengeSpec `json:"challenges"`
}

type ChallengeSpec struct {
	ID          string            `json:"id"`
	Description string            `json:"description"`
	Type        ChallengeType     `json:"type"`
	Secure      bool              `json:"secure"`
	UnStoreable bool              `json:"unstoreable"`
	Options     map[string]string `json:"options,omitempty"`
}

type ChallengeType string

const (
	ChallengeTypeAlpha        ChallengeType = "alpha"
	ChallengeTypeNumeric      ChallengeType = "numeric"
	ChallengeTypeAlphaNumeric ChallengeType = "alphanumeric"
)

type ChallengeAnswerList []ChallengeAnswer

type ChallengeAnswer struct {
	ID    string `json:"id"`
	Value string `json:"value"`
	Store bool   `json:"store"`
}

type UserListPage struct {
	Users      []string `json:"users,omitempty"`
	NextCursor string   `json:"next,omitempty"`
}

type UserCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserToken struct {
	ID    string `json:"id"`    // globally unique identifier for a user
	Token string `json:"token"` // session token
}

type AccessPage struct {
	Accesses []Access `json:"accesses"`
}

type Access struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Enabled    bool      `json:"enabled"`
	IsPinSaved bool      `json:"is_pin_saved"`
	ProviderID string    `json:"provider_id"`
	Accounts   []Account `json:"accounts,omitempty"`
}

type Account struct {
	ID               int64               `json:"id"`
	ProviderID       string              `json:"provider_id"`
	BankAccessID     int64               `json:"bank_access_id"`
	Name             string              `json:"name"`
	Type             AccountType         `json:"type"`
	Number           string              `json:"number"`
	Balance          string              `json:"balance"`
	BalanceDate      time.Time           `json:"balance_date"`
	AvailableBalance string              `json:"available_balance"`
	CreditLine       string              `json:"credit_line"`
	Enabled          bool                `json:"enabled"`
	Currency         string              `json:"currency"`
	IBAN             string              `json:"iban"`
	Supported        bool                `json:"supported"`
	Alias            string              `json:"alias"`
	Capabilities     AccountCapabilities `json:"capabilities" `
	Bin              string              `json:"bin"`
}

type AccountCapabilities struct {
	AccountStatement  []string `json:"account_statement"`
	Transfer          []string `json:"transfer"`
	RecurringTransfer []string `json:"recurring_transfer"`
}

type AccountType string

const (
	AccountTypeBank       AccountType = "bank"
	AccountTypeCreditCard AccountType = "credit_card"
	AccountTypeLoan       AccountType = "loan"
)

type Job struct {
	URI string `json:"uri"`
}

type JobStatus struct {
	Finished  bool       `json:"finished"`
	Stage     JobStage   `json:"stage"`
	Challenge *Challenge `json:"challenge,omitempty"`
	URI       string     `json:"uri,omitempty"`
	Errors    []Problem  `json:"errors,omitempty"`
	Access    *JobAccess `json:"access,omitempty"`
}

type JobStage string

const (
	JobStageUnauthenticated JobStage = "unauthenticated"
	JobStageAuthenticated   JobStage = "authenticated"
	JobStageChallenge       JobStage = "challenge"
	JobStageImported        JobStage = "imported"
	JobStageCancelled       JobStage = "cancelled"
	JobStageProblem         JobStage = "problem"
)

type Challenge struct {
	CanContinue    bool             `json:"can_continue"`
	MaxSteps       uint             `json:"max_steps"`
	CurStep        uint             `json:"cur_step"`
	NextChallenges []ChallengeField `json:"next_challenges"`
	LastProblems   []Problem        `json:"last_problems"`
	Hint           string           `json:"hint"`
}

type ChallengeField struct {
	ID            string   `json:"id"`
	Description   string   `json:"description"`
	ChallengeType string   `json:"type"`
	Previous      string   `json:"previous"`
	Stored        bool     `json:"stored"`
	Reset         bool     `json:"reset"`
	Secure        bool     `json:"secure"`
	Optional      bool     `json:"optional"`
	UnStoreable   bool     `json:"unstoreable"`
	Methods       []string `json:"methods"`
}

type Problem struct {
	Domain string                 `json:"domain"`
	Code   string                 `json:"code"`
	Info   map[string]interface{} `json:"info"`
}

type JobAccess struct {
	ID         int64        `json:"id,omitempty"`
	ProviderID string       `json:"provider_id,omitempty"`
	Name       string       `json:"name,omitempty"`
	Accounts   []JobAccount `json:"accounts,omitempty"`
}

type JobAccount struct {
	ID        int64  `json:"id,omitempty"`
	Name      string `json:"name"`
	Supported bool   `json:"supported"`
	Number    string `json:"number"`
	IBAN      string `json:"iban"`
}

type TransactionPage struct {
	Transactions []Transaction `json:"data"`
	Total        int           `json:"total"`
	Limit        int           `json:"limit"`
	Offset       int           `json:"offset"`
}

type Transaction struct {
	ID                    int64        `json:"id"`
	AccessID              int64        `json:"user_bank_access_id,omitempty"`
	UserAccountID         int64        `json:"user_bank_account_id,omitempty"`
	UserAccount           AccountRef   `json:"user_account,omitempty"`
	CategoryID            int64        `json:"category_id,omitempty"`
	RepeatedTransactionID int64        `json:"repeated_transaction_id,omitempty"`
	Counterparty          Counterparty `json:"counterparty,omitempty"`
	EntryDate             time.Time    `json:"entry_date,omitempty"`
	SettlementDate        time.Time    `json:"settlement_date,omitempty"`
	Amount                *MoneyAmount `json:"amount,omitempty"`
	Usage                 string       `json:"usage,omitempty"`
	TransactionType       string       `json:"transaction_type,omitempty"`
}

type AccountRef struct {
	ProviderID string `json:"provider_id"`
	IBAN       string `json:"iban,omitempty"`
	Label      string `json:"label,omitempty"`
	Number     string `json:"id,omitempty"`
	Type       string `json:"type,omitempty"`
}

type Merchant struct {
	Name string `json:"name"`
}

type Counterparty struct {
	Name     string     `json:"name"`
	Account  AccountRef `json:"account,omitempty"`
	Merchant *Merchant  `json:"merchant,omitempty"`
}

type RepeatedTransactionPage struct {
	Transactions []RepeatedTransaction `json:"data"`
	Total        int                   `json:"total"`
	Limit        int                   `json:"limit"`
	Offset       int                   `json:"offset"`
}

type RepeatedTransaction struct {
	ID            int64          `json:"id"`
	AccessID      int64          `json:"user_bank_access_id,omitempty"`
	UserAccountID int64          `json:"user_bank_account_id,omitempty"`
	UserAccount   AccountRef     `json:"user_account"`
	RemoteAccount AccountRef     `json:"remote_account"`
	RemoteID      string         `json:"remote_id"`
	Schedule      RecurrenceRule `json:"schedule"`
	Amount        *MoneyAmount   `json:"amount"`
	Usage         string         `json:"usage"`
}

type RecurrenceRule struct {
	Start     time.Time `json:"start"`
	Until     time.Time `json:"until"`
	Frequency Frequency `json:"frequency"`
	Interval  int       `json:"interval"`
	ByDay     int       `json:"by_day"`
}

type Frequency string

const (
	FrequencyDaily   Frequency = "daily"
	FrequencyWeekly  Frequency = "weekly"
	FrequencyMonthly Frequency = "monthly"
	FrequencyYearly  Frequency = "yearly"
)

type MoneyAmount struct {
	Currency string `json:"currency"`
	Value    string `json:"value"`
}

type TransferAddress struct {
	Name      string `json:"name"`
	IBAN      string `json:"iban"`
	AccessID  int64  `json:"bank_access_id,omitempty"`
	AccountID int64  `json:"bank_account_id,omitempty"`
}

type TransferType string

const (
	TransferTypeRecurring TransferType = "recurring"
	TransferTypeRegular   TransferType = "regular"
)

type ChallengeAnswerMap map[string]ChallengeAnswer

type Transfer struct {
	ID               string             `json:"id"`
	From             TransferAddress    `json:"from"`
	To               TransferAddress    `json:"to"`
	Amount           *MoneyAmount       `json:"amount"`
	Usage            string             `json:"usage"`
	Version          int                `json:"version"`
	Step             TransferStep       `json:"step"`
	State            TransferState      `json:"state"`
	Schedule         *RecurrenceRule    `json:"schedule,omitempty"`
	EntryDate        time.Time          `json:"booking_date,omitempty"`
	SettlementDate   time.Time          `json:"effective_date,omitempty"`
	Created          time.Time          `json:"created,omitempty"`
	Updated          time.Time          `json:"updated,omitempty"`
	ChallengeAnswers ChallengeAnswerMap `json:"challenge_answers,omitempty"`
	Errors           []Problem          `json:"errors"`
}

type RecurringTransfer struct {
	ID               string             `json:"id"`
	From             TransferAddress    `json:"from"`
	To               TransferAddress    `json:"to"`
	Amount           MoneyAmount        `json:"amount"`
	Usage            string             `json:"usage"`
	Version          int                `json:"version"`
	Step             TransferStep       `json:"step"`
	State            TransferState      `json:"state"`
	Schedule         *RecurrenceRule    `json:"schedule,omitempty"`
	ChallengeAnswers ChallengeAnswerMap `json:"challenge_answers,omitempty"`
	Errors           []Problem          `json:"errors,omitempty"`
}

type TransferState string

const (
	TransferStateOngoing   TransferState = "ongoing"
	TransferStateSucceeded TransferState = "succeeded"
	TransferStateFailed    TransferState = "failed"
	TransferStateCancelled TransferState = "cancelled"
)

type TransferIntent string

const (
	TransferIntentProvidePIN             TransferIntent = "provide_pin"
	TransferIntentSelectAuthMethod       TransferIntent = "select_auth_method"
	TransferIntentProvideChallengeAnswer TransferIntent = "provide_challenge_answer"
	TransferIntentConfirmSimilarTransfer TransferIntent = "confirm_similar_transfer"
)

type PaymentTransferCancelParams struct {
	ID      string `json:"id"`
	Version int    `json:"version"`
}

type TransferStep struct {
	Intent TransferIntent    `json:"intent,omitempty"`
	Data   *TransferStepData `json:"data,omitempty"`
}

type AuthMethod struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

type TransferStepData struct {
	AuthMethods      []AuthMethod `json:"auth_methods,omitempty"`      // TAN Options
	Challenge        string       `json:"challenge,omitempty"`         // TAN Challenge
	ChallengeMessage string       `json:"challenge_message,omitempty"` // TAN Challenge Message
	TANType          TANType      `json:"tan_type,omitempty"`          // Type of the TAN (optical, itan, unknown)
	Confirm          bool         `json:"confirm,omitempty"`           // Confirm (similar transfer)
	Transfers        []Transfer   `json:"transfers,omitempty"`         // Transfer list (similar transfers)
}

type TANType string

const (
	// TANTypeOptical indicates an optical TAN such as flickering barcodes
	TANTypeOptical TANType = "optical"
	// TANTypeITAN indicates an iTAN (aka indexed TAN) such as a list of TAN numbers with a sequence
	TANTypeITAN TANType = "itan"
	// TANTypeMobile indicates a mobileTAN such as an SMS with a passcode
	TANTypeMobile TANType = "mobile"
	// TANTypeChip indicates a chipTAN provided from a calculator device
	TANTypeChip TANType = "chip"
	// TANTypePush indicates a push push notification to a mobile app
	TANTypePush TANType = "push"
	// TANTypeOTP indicates a one-time password
	TANTypeOTP TANType = "otp"
	// TANTypePhoto indicates a colorised matrix barcode
	TANTypePhoto TANType = "photo"

	TANTypeUnknown TANType = "unknown"
)

type DeletedUser struct {
	DeletedUserID string `json:"deleted_user_id"`
}

type IBANDetails struct {
	Account IBANAccount `json:"acc_ref"`
	Banks   []IBANBank  `json:"fis"`
}

type IBANBank struct {
	ID             string `json:"id"`              // the bank identity assigned by the identity provider
	Label          string `json:"label"`           // the bank name
	Country        string `json:"country"`         // the country (e.g. DE)
	Provider       string `json:"provider"`        //  the identity provider (e.g. BIC)
	ServiceContext string `json:"service_context"` // the service context, (e.g. SEPA)
}

type IBANAccount struct {
	IBAN     string `json:"IBAN"`     // the validated IBAN
	Provider string `json:"provider"` // the authoritative provider, IBO for IBANs
}
