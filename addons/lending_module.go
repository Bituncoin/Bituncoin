package addons

import (
	"errors"
	"sync"
	"time"
)

// LendingModule implements lending and borrowing features as an add-on
type LendingModule struct {
	name     string
	version  string
	status   ModuleStatus
	config   map[string]interface{}
	loans    map[string]*Loan
	offers   map[string]*LendingOffer
	mutex    sync.RWMutex
}

// Loan represents a lending/borrowing position
type Loan struct {
	ID              string    `json:"id"`
	Borrower        string    `json:"borrower"`
	Lender          string    `json:"lender"`
	Amount          float64   `json:"amount"`
	Collateral      float64   `json:"collateral"`
	InterestRate    float64   `json:"interestRate"`
	Duration        int64     `json:"duration"`
	StartTime       time.Time `json:"startTime"`
	EndTime         time.Time `json:"endTime"`
	Status          string    `json:"status"` // active, repaid, defaulted
	CollateralRatio float64   `json:"collateralRatio"`
}

// LendingOffer represents a lending offer
type LendingOffer struct {
	ID           string  `json:"id"`
	Lender       string  `json:"lender"`
	Amount       float64 `json:"amount"`
	InterestRate float64 `json:"interestRate"`
	MinDuration  int64   `json:"minDuration"`
	MaxDuration  int64   `json:"maxDuration"`
	MinCollateral float64 `json:"minCollateral"`
	Status       string  `json:"status"` // available, filled, cancelled
}

// NewLendingModule creates a new lending add-on module
func NewLendingModule() *LendingModule {
	return &LendingModule{
		name:    "DeFi Lending",
		version: "1.0.0",
		status:  StatusDisabled,
		loans:   make(map[string]*Loan),
		offers:  make(map[string]*LendingOffer),
	}
}

// GetName returns the module name
func (lm *LendingModule) GetName() string {
	return lm.name
}

// GetVersion returns the module version
func (lm *LendingModule) GetVersion() string {
	return lm.version
}

// GetCategory returns the module category
func (lm *LendingModule) GetCategory() ModuleCategory {
	return CategoryLending
}

// GetDescription returns module description
func (lm *LendingModule) GetDescription() string {
	return "Decentralized lending and borrowing with collateral-based loans"
}

// Initialize initializes the module
func (lm *LendingModule) Initialize(config map[string]interface{}) error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	
	lm.config = config
	return nil
}

// Start starts the module
func (lm *LendingModule) Start() error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	
	lm.status = StatusEnabled
	return nil
}

// Stop stops the module
func (lm *LendingModule) Stop() error {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	
	lm.status = StatusDisabled
	return nil
}

// GetStatus returns the current status
func (lm *LendingModule) GetStatus() ModuleStatus {
	lm.mutex.RLock()
	defer lm.mutex.RUnlock()
	
	return lm.status
}

// Execute executes a module-specific action
func (lm *LendingModule) Execute(action string, params map[string]interface{}) (interface{}, error) {
	lm.mutex.Lock()
	defer lm.mutex.Unlock()
	
	switch action {
	case "create_offer":
		return lm.createOffer(params)
		
	case "list_offers":
		return lm.listOffers(), nil
		
	case "create_loan":
		return lm.createLoan(params)
		
	case "list_loans":
		return lm.listLoans(), nil
		
	case "repay_loan":
		loanID, ok := params["loan_id"].(string)
		if !ok {
			return nil, errors.New("loan_id required")
		}
		return lm.repayLoan(loanID)
		
	default:
		return nil, errors.New("unknown action")
	}
}

func (lm *LendingModule) createOffer(params map[string]interface{}) (*LendingOffer, error) {
	lender, ok := params["lender"].(string)
	if !ok {
		return nil, errors.New("lender required")
	}
	
	amount, ok := params["amount"].(float64)
	if !ok {
		return nil, errors.New("amount required")
	}
	
	interestRate, ok := params["interest_rate"].(float64)
	if !ok {
		interestRate = 5.0 // Default 5%
	}
	
	offer := &LendingOffer{
		ID:            generateLoanID(),
		Lender:        lender,
		Amount:        amount,
		InterestRate:  interestRate,
		MinDuration:   7 * 24 * 60 * 60,  // 7 days
		MaxDuration:   365 * 24 * 60 * 60, // 1 year
		MinCollateral: 150.0, // 150% collateral ratio
		Status:        "available",
	}
	
	lm.offers[offer.ID] = offer
	return offer, nil
}

func (lm *LendingModule) listOffers() []LendingOffer {
	offers := make([]LendingOffer, 0, len(lm.offers))
	for _, offer := range lm.offers {
		if offer.Status == "available" {
			offers = append(offers, *offer)
		}
	}
	return offers
}

func (lm *LendingModule) createLoan(params map[string]interface{}) (*Loan, error) {
	borrower, ok := params["borrower"].(string)
	if !ok {
		return nil, errors.New("borrower required")
	}
	
	offerID, ok := params["offer_id"].(string)
	if !ok {
		return nil, errors.New("offer_id required")
	}
	
	offer, exists := lm.offers[offerID]
	if !exists {
		return nil, errors.New("offer not found")
	}
	
	if offer.Status != "available" {
		return nil, errors.New("offer not available")
	}
	
	collateral, ok := params["collateral"].(float64)
	if !ok {
		return nil, errors.New("collateral required")
	}
	
	// Calculate collateral ratio
	collateralRatio := (collateral / offer.Amount) * 100
	if collateralRatio < offer.MinCollateral {
		return nil, errors.New("insufficient collateral")
	}
	
	duration := int64(30 * 24 * 60 * 60) // Default 30 days
	if d, ok := params["duration"].(float64); ok {
		duration = int64(d)
	}
	
	now := time.Now()
	loan := &Loan{
		ID:              generateLoanID(),
		Borrower:        borrower,
		Lender:          offer.Lender,
		Amount:          offer.Amount,
		Collateral:      collateral,
		InterestRate:    offer.InterestRate,
		Duration:        duration,
		StartTime:       now,
		EndTime:         now.Add(time.Duration(duration) * time.Second),
		Status:          "active",
		CollateralRatio: collateralRatio,
	}
	
	lm.loans[loan.ID] = loan
	offer.Status = "filled"
	
	return loan, nil
}

func (lm *LendingModule) listLoans() []Loan {
	loans := make([]Loan, 0, len(lm.loans))
	for _, loan := range lm.loans {
		loans = append(loans, *loan)
	}
	return loans
}

func (lm *LendingModule) repayLoan(loanID string) (*Loan, error) {
	loan, exists := lm.loans[loanID]
	if !exists {
		return nil, errors.New("loan not found")
	}
	
	if loan.Status != "active" {
		return nil, errors.New("loan not active")
	}
	
	loan.Status = "repaid"
	return loan, nil
}

// generateLoanID generates a unique loan ID
func generateLoanID() string {
	return time.Now().Format("20060102150405")
}
