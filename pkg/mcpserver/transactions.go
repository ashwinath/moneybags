package mcpserver

import (
	"context"
	"fmt"
	"sort"
	"strings"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/ashwinath/moneybags/pkg/utils"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

const searchTransactionsToolName = "search_transactions"

const searchTransactionsToolDescription = "Search recorded transactions by their free text description, " +
	"optionally restricted to a single type. The description is matched case-insensitively as a substring, " +
	"so it need not match the casing or the whole note: \"coffee\" also finds \"Starbucks Coffee\". " +
	"Omit description to search every description, or omit type to search every type. " +
	"Leaving both out returns the most recent transactions. Valid types are: %s."

const errorFormatUnknownTransactionType = "%q is not a valid transaction type, expected one of: %s"

// transactionTypeAliases maps every accepted spelling of a transaction type onto
// its canonical value. Keys are in normalised form, as produced by
// normaliseTransactionType, so that the database's SCREAMING_SNAKE values and
// the shorter Telegram spellings such as \"cc\" share one lookup table.
var transactionTypeAliases = map[string]db.TransactionType{
	"CC":                  db.TypeCreditCard,
	"CREDIT CARD":         db.TypeCreditCard,
	"INSURANCE":           db.TypeInsurance,
	"OWN":                 db.TypeOwn,
	"REIM":                db.TypeReimburse,
	"SHARED":              db.TypeShared,
	"SHARED CC REIM":      db.TypeSharedCCReimburse,
	"SHARED CC REIMBURSE": db.TypeSharedCCReimburse,
	"SHARED REIM":         db.TypeSharedReimburse,
	"SPECIAL OWN":         db.TypeSpecialOwn,
	"SPECIAL SHARED":      db.TypeSpecialShared,
	"SPECIAL SHARED REIM": db.TypeSpecialSharedReimburse,
	"TAX":                 db.TypeTax,
	"TITHE":               db.TypeTithe,
}

// registerTransactionTools adds the transaction tools to the server.
func (s *Server) registerTransactionTools() {
	mcp.AddTool(s.server, &mcp.Tool{
		Name: searchTransactionsToolName,
		Description: fmt.Sprintf(
			searchTransactionsToolDescription,
			strings.Join(validTransactionTypes(), ", "),
		),
	}, s.searchTransactions)
}

// searchTransactionsInput is the argument set for the search_transactions tool.
type searchTransactionsInput struct {
	Description string `json:"description,omitempty" jsonschema:"case-insensitive substring to look for inside a transaction's description, the free text note recorded with it"`
	Type        string `json:"type,omitempty" jsonschema:"optional transaction type to restrict the search to"`
}

// searchTransactionsOutput is the structured result of a transaction search.
type searchTransactionsOutput struct {
	Transactions []searchedTransaction `json:"transactions" jsonschema:"the matching transactions, most recent first"`
	Count        int                   `json:"count" jsonschema:"how many transactions this result holds"`
	Truncated    bool                  `json:"truncated" jsonschema:"true when more transactions matched than were returned, meaning this result is only the most recent slice of the matches"`
}

// searchedTransaction is a single transaction as handed back to MCP clients.
// The database model is not marshalled directly: it carries no JSON tags, so
// its date would render as a full timestamp and its audit columns would leak
// into the tool result.
type searchedTransaction struct {
	ID          uint    `json:"id" jsonschema:"the transaction's ID, which can be used to delete it"`
	Date        string  `json:"date" jsonschema:"the transaction date as yyyy-mm-dd"`
	Type        string  `json:"type" jsonschema:"the transaction type"`
	Description string  `json:"description" jsonschema:"the note recorded with the transaction, empty for types that do not take one"`
	Amount      float64 `json:"amount" jsonschema:"the amount as recorded, which is negative for reimbursements"`
}

// searchTransactions answers a search_transactions tool call.
func (s *Server) searchTransactions(
	_ context.Context,
	_ *mcp.CallToolRequest,
	input searchTransactionsInput,
) (*mcp.CallToolResult, *searchTransactionsOutput, error) {
	options, err := input.searchOptions()
	if err != nil {
		return nil, nil, err
	}

	var (
		transactions []db.Transaction
		truncated    bool
	)

	s.fw.TimeFunction("MCP search_transactions", func() {
		transactions, truncated, err = s.db.SearchTransactions(options)
	})
	if err != nil {
		return nil, nil, err
	}

	return nil, &searchTransactionsOutput{
		Transactions: searchResults(transactions),
		Count:        len(transactions),
		Truncated:    truncated,
	}, nil
}

// searchOptions turns the tool input into a database query, rejecting a type
// that does not name a real transaction type.
func (i searchTransactionsInput) searchOptions() (*db.SearchTransactionsOptions, error) {
	options := &db.SearchTransactionsOptions{
		Description: i.Description,
	}

	if i.Type == "" {
		return options, nil
	}

	transactionType, ok := resolveTransactionType(i.Type)
	if !ok {
		return nil, fmt.Errorf(
			errorFormatUnknownTransactionType,
			i.Type,
			strings.Join(validTransactionTypes(), ", "),
		)
	}

	options.Types = []db.TransactionType{transactionType}

	return options, nil
}

// searchResults converts transactions into the shape returned to clients.
func searchResults(transactions []db.Transaction) []searchedTransaction {
	results := make([]searchedTransaction, 0, len(transactions))
	for _, transaction := range transactions {
		results = append(results, searchedTransaction{
			ID:          transaction.ID,
			Date:        utils.SetDateToDateOnly(transaction.Date),
			Type:        string(transaction.Type),
			Description: transaction.Classification,
			Amount:      transaction.Amount,
		})
	}

	return results
}

// resolveTransactionType maps a caller supplied type onto a canonical
// transaction type, reporting whether it was recognised.
func resolveTransactionType(s string) (db.TransactionType, bool) {
	transactionType, ok := transactionTypeAliases[normaliseTransactionType(s)]

	return transactionType, ok
}

// normaliseTransactionType upper-cases a type and folds its separators, so
// that "shared_cc_reimburse", "shared cc reim" and "Shared CC Reimburse" all
// resolve to the same entry.
func normaliseTransactionType(s string) string {
	folded := strings.NewReplacer("_", " ", "-", " ").Replace(s)

	return strings.Join(strings.Fields(strings.ToUpper(folded)), " ")
}

// validTransactionTypes lists the accepted type spellings in a stable order, so
// that it can be embedded in the tool description and in error messages.
func validTransactionTypes() []string {
	transactionTypes := make([]string, 0, len(transactionTypeAliases))
	for alias := range transactionTypeAliases {
		transactionTypes = append(transactionTypes, alias)
	}
	sort.Strings(transactionTypes)

	return transactionTypes
}
