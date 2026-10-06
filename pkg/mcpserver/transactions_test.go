package mcpserver

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/ashwinath/moneybags/pkg/config"
	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/utils"

	"github.com/ashwinath/simple/framework"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// allTransactionTypes lists every type the database defines, so that the alias
// table can be checked for coverage against it.
var allTransactionTypes = []db.TransactionType{
	db.TypeReimburse,
	db.TypeSharedReimburse,
	db.TypeShared,
	db.TypeSpecialShared,
	db.TypeSpecialSharedReimburse,
	db.TypeOwn,
	db.TypeSpecialOwn,
	db.TypeCreditCard,
	db.TypeInsurance,
	db.TypeTithe,
	db.TypeTax,
	db.TypeSharedCCReimburse,
}

func TestResolveTransactionType(t *testing.T) {
	var tests = []struct {
		name     string
		input    string
		expected db.TransactionType
		wantOk   bool
	}{
		{
			name:     "canonical value",
			input:    "OWN",
			expected: db.TypeOwn,
			wantOk:   true,
		},
		{
			name:     "lowercase",
			input:    "own",
			expected: db.TypeOwn,
			wantOk:   true,
		},
		{
			name:     "mixed case",
			input:    "OwN",
			expected: db.TypeOwn,
			wantOk:   true,
		},
		{
			name:     "surrounding whitespace",
			input:    "  own  ",
			expected: db.TypeOwn,
			wantOk:   true,
		},
		{
			name:     "collapsed internal whitespace",
			input:    "credit    card",
			expected: db.TypeCreditCard,
			wantOk:   true,
		},
		{
			name:     "telegram shorthand for credit card",
			input:    "cc",
			expected: db.TypeCreditCard,
			wantOk:   true,
		},
		{
			name:     "underscores are folded into the canonical value",
			input:    "SHARED_CC_REIMBURSE",
			expected: db.TypeSharedCCReimburse,
			wantOk:   true,
		},
		{
			name:     "telegram spelling of shared cc reim",
			input:    "shared cc reim",
			expected: db.TypeSharedCCReimburse,
			wantOk:   true,
		},
		{
			name:     "dashes are folded",
			input:    "special-shared",
			expected: db.TypeSpecialShared,
			wantOk:   true,
		},
		{
			name:   "unknown type",
			input:  "nonsense",
			wantOk: false,
		},
		{
			name:   "empty type",
			input:  "",
			wantOk: false,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			transactionType, ok := resolveTransactionType(tt.input)
			assert.Equal(t, tt.wantOk, ok)
			if tt.wantOk {
				assert.Equal(t, tt.expected, transactionType)
			}
		})
	}
}

// TestEveryTransactionTypeIsResolvable guards against a transaction type being
// added to the database without a matching alias.
func TestEveryTransactionTypeIsResolvable(t *testing.T) {
	for _, transactionType := range allTransactionTypes {
		transactionType := transactionType
		t.Run(string(transactionType), func(t *testing.T) {
			resolved, ok := resolveTransactionType(string(transactionType))
			assert.True(t, ok)
			assert.Equal(t, transactionType, resolved)
		})
	}
}

// TestAliasTableHasNoStrayEntries keeps an alias from mapping onto a type that
// is not an alias at all.
func TestAliasTableHasNoStrayEntries(t *testing.T) {
	validTypes := make(map[db.TransactionType]struct{}, len(allTransactionTypes))
	for _, transactionType := range allTransactionTypes {
		validTypes[transactionType] = struct{}{}
	}

	for alias, transactionType := range transactionTypeAliases {
		alias := alias
		transactionType := transactionType
		t.Run(alias, func(t *testing.T) {
			_, ok := validTypes[transactionType]
			assert.True(t, ok)
			assert.Equal(t, alias, normaliseTransactionType(alias))
		})
	}
}

func TestValidTransactionTypesIsSorted(t *testing.T) {
	transactionTypes := validTransactionTypes()
	assert.Len(t, transactionTypes, len(transactionTypeAliases))
	assert.IsIncreasing(t, transactionTypes)
}

func TestSearchTransactionsInputSearchOptions(t *testing.T) {
	var tests = []struct {
		name          string
		input         searchTransactionsInput
		expected      *db.SearchTransactionsOptions
		expectErr     bool
		errorMentions string
	}{
		{
			name:     "description only",
			input:    searchTransactionsInput{Description: "coffee"},
			expected: &db.SearchTransactionsOptions{Description: "coffee"},
		},
		{
			name:     "type only",
			input:    searchTransactionsInput{Type: "cc"},
			expected: &db.SearchTransactionsOptions{Types: []db.TransactionType{db.TypeCreditCard}},
		},
		{
			name: "both",
			input: searchTransactionsInput{
				Description: "coffee",
				Type:        "shared cc reim",
			},
			expected: &db.SearchTransactionsOptions{
				Description: "coffee",
				Types:       []db.TransactionType{db.TypeSharedCCReimburse},
			},
		},
		{
			name:     "neither",
			input:    searchTransactionsInput{},
			expected: &db.SearchTransactionsOptions{},
		},
		{
			name:          "unknown type is rejected",
			input:         searchTransactionsInput{Type: "nonsense"},
			expectErr:     true,
			errorMentions: "nonsense",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			options, err := tt.input.searchOptions()
			if tt.expectErr {
				assert.NotNil(t, err)
				assert.Contains(t, err.Error(), tt.errorMentions)
				// The error has to name the alternatives for the caller to recover.
				assert.Contains(t, err.Error(), string(db.TypeOwn))
				return
			}

			assert.Nil(t, err)
			assert.Equal(t, tt.expected, options)
		})
	}
}

func TestNewServerFailsWithoutTransactionDB(t *testing.T) {
	fw := createFW(t, nil)

	server, err := NewServer(fw)
	assert.Nil(t, server)
	assert.NotNil(t, err)
}

func TestSearchTransactionsToolDescriptionListsEveryType(t *testing.T) {
	description := searchTransactionsToolDescription
	types := validTransactionTypes()
	description = strings.Replace(description, "%s", strings.Join(types, ", "), 1)

	for _, transactionType := range types {
		assert.Contains(t, description, transactionType)
	}
}

func TestSearchTransactionsHandler(t *testing.T) {
	var tests = []struct {
		name            string
		inserts         []db.Transaction
		input           searchTransactionsInput
		expectedDescrip []string
		expectErr       bool
	}{
		{
			name: "matches ignoring case",
			inserts: []db.Transaction{
				{
					Date:           testDate(t, "2023-04-03"),
					Type:           db.TypeOwn,
					Classification: "Starbucks Coffee",
					Amount:         5.40,
				},
				{
					Date:           testDate(t, "2023-04-04"),
					Type:           db.TypeOwn,
					Classification: "lunch",
					Amount:         9.00,
				},
			},
			input:           searchTransactionsInput{Description: "cOfFeE"},
			expectedDescrip: []string{"Starbucks Coffee"},
		},
		{
			name: "filters by type alone",
			inserts: []db.Transaction{
				{
					Date:           testDate(t, "2023-04-03"),
					Type:           db.TypeOwn,
					Classification: "meal",
					Amount:         5.40,
				},
				{
					Date:           testDate(t, "2023-04-04"),
					Type:           db.TypeTithe,
					Classification: "",
					Amount:         200.00,
				},
			},
			input:           searchTransactionsInput{Type: "tithe"},
			expectedDescrip: []string{""},
		},
		{
			name: "unknown type surfaces an error",
			inserts: []db.Transaction{
				{
					Date:           testDate(t, "2023-04-03"),
					Type:           db.TypeOwn,
					Classification: "meal",
					Amount:         5.40,
				},
			},
			input:     searchTransactionsInput{Type: "nonsense"},
			expectErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			err := db.RunTest(func(database *db.DB) {
				fw := createFW(t, database)
				server, err := NewServer(fw)
				assert.Nil(t, err)

				txDB, err := db.NewTransactionDB(database)
				assert.Nil(t, err)
				assert.Nil(t, txDB.BulkAdd(tt.inserts))

				_, output, err := server.searchTransactions(context.Background(), nil, tt.input)
				if tt.expectErr {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)

				descriptions := make([]string, 0, len(output.Transactions))
				for _, transaction := range output.Transactions {
					descriptions = append(descriptions, transaction.Description)
				}
				assert.Equal(t, tt.expectedDescrip, descriptions)
				assert.Equal(t, len(descriptions), output.Count)
				assert.False(t, output.Truncated)
			})
			assert.Nil(t, err)
		})
	}
}

func TestSearchTransactionsHandlerReportsTruncation(t *testing.T) {
	err := db.RunTest(func(database *db.DB) {
		fw := createFW(t, database)
		server, err := NewServer(fw)
		assert.Nil(t, err)

		inserts := make([]db.Transaction, 0, db.MaxTransactionSearchResults+1)
		for i := range db.MaxTransactionSearchResults + 1 {
			inserts = append(inserts, db.Transaction{
				Date:           testDate(t, "2023-04-03"),
				Type:           db.TypeOwn,
				Classification: "match",
				Amount:         float64(i),
			})
		}

		txDB, err := db.NewTransactionDB(database)
		assert.Nil(t, err)
		assert.Nil(t, txDB.BulkAdd(inserts))

		_, output, err := server.searchTransactions(
			context.Background(),
			nil,
			searchTransactionsInput{Description: "match"},
		)
		assert.Nil(t, err)
		assert.Equal(t, db.MaxTransactionSearchResults, output.Count)
		assert.True(t, output.Truncated)
	})
	assert.Nil(t, err)
}

func TestSearchResultsFormatsDates(t *testing.T) {
	// The database hands timestamps back in UTC, so the conversion to Singapore
	// time is what has to be asserted, not the Singapore wall clock.
	singapore, err := time.LoadLocation("Asia/Singapore")
	assert.Nil(t, err)
	midnight := time.Date(2023, 4, 3, 0, 0, 0, 0, singapore).UTC()

	transactions := []db.Transaction{
		{
			ID:             7,
			Date:           midnight,
			Type:           db.TypeOwn,
			Classification: "meal",
			Amount:         12.345,
		},
		{
			ID:   8,
			Date: time.Date(2023, 4, 30, 20, 0, 0, 0, time.UTC),
			Type: db.TypeOwn,
		},
	}

	assert.Equal(
		t,
		[]searchedTransaction{
			{
				ID:          7,
				Date:        "2023-04-03",
				Type:        string(db.TypeOwn),
				Description: "meal",
				Amount:      12.345,
			},
			{
				ID:          8,
				Date:        "2023-05-01",
				Type:        string(db.TypeOwn),
				Description: "",
			},
		},
		searchResults(transactions),
	)
	assert.Equal(t, []searchedTransaction{}, searchResults(nil))
}

func createFW(t *testing.T, database *db.DB) framework.FW {
	t.Helper()

	logger, _ := zap.NewProduction()
	sugar := logger.Sugar()

	configPath := utils.GetLocalRepoLocation() + "/pkg/config/testdata/config.yaml"
	c, err := config.New(configPath)
	assert.Nil(t, err)

	databases := map[string]any{}
	if database != nil {
		txDB, err := db.NewTransactionDB(database)
		assert.Nil(t, err)
		databases[db.TransactionDatabaseName] = txDB
	}

	return framework.New(c, sugar, databases, map[string]any{})
}

func testDate(t *testing.T, dateString string) time.Time {
	t.Helper()

	loc, err := time.LoadLocation("Asia/Singapore")
	assert.Nil(t, err)

	parsed, err := time.ParseInLocation(time.DateOnly, dateString, loc)
	assert.Nil(t, err)

	return parsed
}
