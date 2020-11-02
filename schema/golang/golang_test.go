package golang

import (
	"reflect"
	"testing"

	"github.com/webrpc/webrpc/schema"
)

func Test_interfaceAllMethodNames(t *testing.T) {
	type args struct {
		goType string
	}
	tests := []struct {
		name     string
		args     args
		expected []string
	}{
		{
			name: "List of methods for interface",
			args: args{
				goType: "{BorrowBook(ctx context.Context, BookID int64) error; GetBookAuthor(ctx context.Context, BookID int64) (Author, map[string]string, regexp.Regexp, error); GetBooks(ctx context.Context) ([]Book, string, error)}",
			},
			expected: []string{"BorrowBook", "GetBookAuthor", "GetBooks"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := interfaceAllMethodNames(tt.args.goType)
			if !reflect.DeepEqual(tt.expected, got) {
				t.Errorf("Split() = %v, expected %v", got, tt.expected)
			}
		})
	}
}

func Test_buildArgumentsList(t *testing.T) {
	type args struct {
		s         *schema.WebRPCSchema
		goType    string
		method    string
		checkType string
	}
	tests := []struct {
		name     string
		args     args
		expected []*schema.MethodArgument
		wantErr  bool
	}{
		{
			name: "build Argument input list",
			args: args{
				goType:    "{BorrowBook(ctx context.Context, BookID int64) error; GetBookAuthor(ctx context.Context, BookID int64) (Author, map[string]string, regexp.Regexp, error); GetBooks(ctx context.Context) ([]Book, string, error)}",
				method:    "BorrowBook",
				checkType: "isInputArgs",
			},
			wantErr: false,
		},
		{
			name: "build Argument output list",
			args: args{
				goType:    "{BorrowBook(ctx context.Context, BookID int64) error; GetBookAuthor(ctx context.Context, BookID int64) (Author, map[string]string, regexp.Regexp, error); GetBooks(ctx context.Context) ([]Book, string, error)}",
				method:    "GetBooks",
				checkType: "isOutputArgs",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildArgumentsList(tt.args.s, tt.args.goType, tt.args.method, tt.args.checkType)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildArgumentsList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) == 0 {
				t.Errorf("buildArgumentsList() = %v. We expect some value here.", got)
			}
		})
	}
}

func Test_fieldsOfStruct(t *testing.T) {
	type args struct {
		goType string
	}
	tests := []struct {
		name     string
		args     args
		expected []string
	}{
		{
			name: "Author struct field",
			args: args{
				goType: "{ID BookID; Name string; Authors []Author}",
			},
			expected: []string{"ID BookID", "Name string", "Authors []Author"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := fieldsOfStruct(tt.args.goType)
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("fieldsOfStruct() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestParser_goparse(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		p       *Parser
		args    args
		want    *schema.WebRPCSchema
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.p.goparse(tt.args.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parser.goparse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parser.goparse() = %v, want %v", got, tt.want)
			}
		})
	}
}
