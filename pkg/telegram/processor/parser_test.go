package telegramprocessor

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/ashwinath/moneybags/pkg/db"
	"github.com/stretchr/testify/assert"
)

func parseDateForced(t *testing.T, dateString string) *time.Time {
	loc, err := time.LoadLocation("Asia/Singapore")
	assert.Nil(t, err)

	parsed, err := time.ParseInLocation(time.DateOnly, dateString, loc)
	assert.Nil(t, err)

	return &parsed
}

func TestParser(t *testing.T) {
	var tests = []struct {
		name          string
		testString    string
		expected      Chunk
		expectedError error
	}{
		{
			name:       "Add reim",
			testString: "Add reim food 24.5",
			expected: Chunk{
				Instruction:    Add,
				Type:           db.TypeReimburse,
				Classification: "food",
				Amount:         24.5,
			},
			expectedError: nil,
		},
		{
			name:       "Add shared reim",
			testString: "Add shared reim christmas's dinner 60.5",
			expected: Chunk{
				Instruction:    Add,
				Type:           db.TypeSharedReimburse,
				Classification: "christmas's dinner",
				Amount:         60.5,
			},
			expectedError: nil,
		},
		{
			name:       "Add special shared reim",
			testString: "Add special shared reim furniture 2400.5",
			expected: Chunk{
				Instruction:    Add,
				Type:           db.TypeSpecialSharedReimburse,
				Classification: "furniture",
				Amount:         2400.5,
			},
			expectedError: nil,
		},
		{
			name:       "Add shared",
			testString: "Add shared lunch 10.5",
			expected: Chunk{
				Instruction:    Add,
				Type:           db.TypeShared,
				Classification: "lunch",
				Amount:         10.5,
			},
			expectedError: nil,
		},
		{
			name:       "Add special shared",
			testString: "Add special shared washing machine 810.5",
			expected: Chunk{
				Instruction:    Add,
				Type:           db.TypeSpecialShared,
				Classification: "washing machine",
				Amount:         810.5,
			},
			expectedError: nil,
		},
		{
			name:       "Add special own",
			testString: "Add special own holiday to japan 810.5",
			expected: Chunk{
				Instruction:    Add,
				Type:           db.TypeSpecialOwn,
				Classification: "holiday to japan",
				Amount:         810.5,
			},
			expectedError: nil,
		},
		{
			name:       "Add own",
			testString: "Add own computer and monitors 2010",
			expected: Chunk{
				Instruction:    Add,
				Type:           db.TypeOwn,
				Classification: "computer and monitors",
				Amount:         2010.0,
			},
			expectedError: nil,
		},
		{
			name:       "Add shared with date",
			testString: "Add shared dinner 12.5 2020-10-02",
			expected: Chunk{
				Instruction:    Add,
				Type:           db.TypeShared,
				Classification: "dinner",
				Amount:         12.5,
				Date:           parseDateForced(t, "2020-10-02"),
			},
			expectedError: nil,
		},
		{
			name:       "Delete transaction",
			testString: "Del 100",
			expected: Chunk{
				Instruction: Delete,
				ID:          100,
			},
			expectedError: nil,
		},
		{
			name:       "Help",
			testString: "help",
			expected: Chunk{
				Instruction: Help,
			},
			expectedError: nil,
		},
		{
			name:       "Generate",
			testString: "Gen March 2020",
			expected: Chunk{
				Instruction: Generate,
				StartDate:   parseDateForced(t, "2020-03-01"),
			},
			expectedError: nil,
		},
		{
			name:       "Add credit card",
			testString: "Add cc 1003.52",
			expected: Chunk{
				Instruction: Add,
				Type:        db.TypeCreditCard,
				Amount:      1003.52,
			},
			expectedError: nil,
		},
		{
			name:       "Add insurance",
			testString: "Add insurance 200.32",
			expected: Chunk{
				Instruction: Add,
				Type:        db.TypeInsurance,
				Amount:      200.32,
			},
			expectedError: nil,
		},
		{
			name:       "Add tithe",
			testString: "Add tithe 500",
			expected: Chunk{
				Instruction: Add,
				Type:        db.TypeTithe,
				Amount:      500,
			},
			expectedError: nil,
		},
		{
			name:       "Add tax",
			testString: "Add tax 45 2023-03-23",
			expected: Chunk{
				Instruction: Add,
				Type:        db.TypeTax,
				Amount:      45,
				Date:        parseDateForced(t, "2023-03-23"),
			},
			expectedError: nil,
		},
		{
			name:       "Add shared cc reim",
			testString: "Add shared cc reim petrol 20.45 2023-03-25",
			expected: Chunk{
				Instruction:    Add,
				Type:           db.TypeSharedCCReimburse,
				Classification: "petrol",
				Amount:         20.45,
				Date:           parseDateForced(t, "2023-03-25"),
			},
			expectedError: nil,
		},
		{
			name:          "wrong instruction",
			testString:    "update own computer and monitors 2010",
			expected:      Chunk{},
			expectedError: errors.New("Error: could not parse instruction token: update, Message: update own computer and monitors 2010"),
		},
		{
			name:          "wrong Add when should have been Del",
			testString:    "Add 100",
			expected:      Chunk{},
			expectedError: errors.New("Error: unable to parse type: 100, Message: Add 100"),
		},
		{
			name:          "wrong special but shared or own does not come after",
			testString:    "Add special something toy 100.4",
			expected:      Chunk{},
			expectedError: errors.New("Error: unable to parse type: special something, Message: Add special something toy 100.4"),
		},
		{
			name:          "wrong type",
			testString:    "Add property sale 1000000.4",
			expected:      Chunk{},
			expectedError: errors.New("Error: unable to parse type: property, Message: Add property sale 1000000.4"),
		},
		{
			name:          "classification is a number",
			testString:    "Add reim 100.4",
			expected:      Chunk{},
			expectedError: errors.New("Error: empty classification token, Message: Add reim 100.4"),
		},
		{
			name:          "empty classification",
			testString:    "Add reim",
			expected:      Chunk{},
			expectedError: errors.New("Error: empty classification token, Message: Add reim"),
		},
		{
			name:          "empty amount",
			testString:    "Add own bicycle",
			expected:      Chunk{},
			expectedError: errors.New("Error: empty amount token, Message: Add own bicycle"),
		},
		{
			name:          "empty string",
			testString:    "",
			expected:      Chunk{},
			expectedError: errors.New("Error: empty instruction token, Message: "),
		},
		{
			name:          "empty transaction type",
			testString:    "Add",
			expected:      Chunk{},
			expectedError: errors.New("Error: empty type token, Message: Add"),
		},
		{
			name:          "empty id for delete",
			testString:    "Del",
			expected:      Chunk{},
			expectedError: errors.New("Error: empty ID token, Message: Del"),
		},
		{
			name:          "delete cannot parse id",
			testString:    "Del you",
			expected:      Chunk{},
			expectedError: errors.New("Error: strconv.Atoi: parsing \"you\": invalid syntax, Message: Del you"),
		},
		{
			name:          "generate year empty",
			testString:    "Gen March",
			expected:      Chunk{},
			expectedError: errors.New("Error: empty year token, Message: Gen March"),
		},
		{
			name:          "generate year out of range (above 2200)",
			testString:    "Gen March 2201",
			expected:      Chunk{},
			expectedError: errors.New("Error: year out of range (1970 - 2200 allowed): 2201, Message: Gen March 2201"),
		},
		{
			name:          "generate year out of range (before 1970)",
			testString:    "Gen March 1969",
			expected:      Chunk{},
			expectedError: errors.New("Error: year out of range (1970 - 2200 allowed): 1969, Message: Gen March 1969"),
		},
		{
			name:          "generate month empty",
			testString:    "Gen",
			expected:      Chunk{},
			expectedError: errors.New("Error: empty month token, Message: Gen"),
		},
		{
			name:          "generate month no such month",
			testString:    "Gen Januahree 2020",
			expected:      Chunk{},
			expectedError: errors.New("Error: unable to parse month: Januahree, Message: Gen Januahree 2020"),
		},
		{
			name:          "Add shared special (order flipped)",
			testString:    "Add shared special washing machine 810.5",
			expected:      Chunk{},
			expectedError: fmt.Errorf("Error: %s, Message: Add shared special washing machine 810.5", errorSuggestAMoreSuitableSpecialSharedType),
		},
		{
			name:          "suggest a better shared cc reim",
			testString:    "Add shared cc 100",
			expected:      Chunk{},
			expectedError: errors.New("Error: did you mean shared cc reim?, Message: Add shared cc 100"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			actual, err := Parse(tt.testString)
			if tt.expectedError == nil {
				assert.Equal(t, tt.expected, *actual)
			} else {
				assert.Equal(t, tt.expectedError, err)
			}
		})
	}
}
