package testing

import (
	"GoAuth/internal/domain/permissions"
	"encoding/json"
	"testing"
)

func TestConditions(t *testing.T) {
	t.Run("ValidateConditions", func(t *testing.T) {
		tests := []struct {
			name      string
			condition permissions.Condition
			wantErr   bool
			errMsg    string
		}{
			// ==================== HAPPY PATHS ====================
			{
				name: "simple eq predicate with value",
				condition: permissions.Condition{
					Path:  "resource.status",
					Op:    permissions.OpEq,
					Value: "active",
				},
				wantErr: false,
			},
			{
				name: "simple gt predicate with ref",
				condition: permissions.Condition{
					Path: "resource.count",
					Op:   permissions.OpGt,
					Ref:  "subject.quota",
				},
				wantErr: false,
			},
			{
				name: "exists operator",
				condition: permissions.Condition{
					Path: "resource.optional_field",
					Op:   permissions.OpExists,
				},
				wantErr: false,
			},
			{
				name: "grace_before temporal condition",
				condition: permissions.Condition{
					Field:  "resource.start_time",
					Op:     permissions.OpGraceBefore,
					Margin: "15m",
				},
				wantErr: false,
			},
			{
				name: "grace_after temporal condition",
				condition: permissions.Condition{
					Field:  "resource.end_time",
					Op:     permissions.OpGraceAfter,
					Margin: "1h",
				},
				wantErr: false,
			},
			{
				name: "grace_around temporal condition",
				condition: permissions.Condition{
					Field:  "resource.deadline",
					Op:     permissions.OpGraceAround,
					Margin: "30m",
				},
				wantErr: false,
			},
			{
				name: "grace_duration temporal condition",
				condition: permissions.Condition{
					FieldStart: "resource.start_time",
					FieldEnd:   "resource.end_time",
					Op:         permissions.OpGraceDuration,
					Margin:     "2h",
				},
				wantErr: false,
			},
			{
				name: "nested AND conditions",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{Path: "resource.status", Op: permissions.OpEq, Value: "active"},
						{Path: "resource.count", Op: permissions.OpGt, Value: 10},
					},
				},
				wantErr: false,
			},
			{
				name: "nested OR conditions",
				condition: permissions.Condition{
					Or: &[]permissions.Condition{
						{Path: "resource.type", Op: permissions.OpEq, Value: "admin"},
						{Path: "resource.type", Op: permissions.OpEq, Value: "moderator"},
					},
				},
				wantErr: false,
			},
			{
				name: "nested NOT condition",
				condition: permissions.Condition{
					Not: &permissions.Condition{
						Path:  "resource.deleted",
						Op:    permissions.OpEq,
						Value: true,
					},
				},
				wantErr: false,
			},
			{
				name: "deeply nested logical conditions",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							Or: &[]permissions.Condition{
								{Path: "resource.type", Op: permissions.OpEq, Value: "premium"},
								{Path: "resource.type", Op: permissions.OpEq, Value: "enterprise"},
							},
						},
						{
							Not: &permissions.Condition{
								Path:  "resource.banned",
								Op:    permissions.OpEq,
								Value: true,
							},
						},
					},
				},
				wantErr: false,
			},
			{
				name: "string operators",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{Path: "resource.name", Op: permissions.OpStartsWith, Value: "prefix_"},
						{Path: "resource.email", Op: permissions.OpEndsWith, Value: "@example.com"},
						{Path: "resource.code", Op: permissions.OpMatches, Value: "^[A-Z]{3}-\\d{4}$"},
					},
				},
				wantErr: false,
			},
			{
				name: "array operators",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{Path: "resource.tags", Op: permissions.OpIn, Value: "urgent"},
						{Path: "resource.roles", Op: permissions.OpContains, Value: "admin"},
					},
				},
				wantErr: false,
			},

			// ==================== SAD PATHS ====================
			{
				name: "empty AND operator",
				condition: permissions.Condition{
					And: &[]permissions.Condition{},
				},
				wantErr: true,
				errMsg:  "and: AND conditions cannot be empty",
			},
			{
				name: "empty OR operator",
				condition: permissions.Condition{
					Or: &[]permissions.Condition{},
				},
				wantErr: true,
				errMsg:  "or: OR conditions cannot be empty",
			},
			{
				name: "missing operator",
				condition: permissions.Condition{
					Path:  "resource.status",
					Value: "active",
				},
				wantErr: true,
				errMsg:  "op: operator is required",
			},
			{
				name: "invalid operator name",
				condition: permissions.Condition{
					Path:  "resource.status",
					Op:    "invalid_op",
					Value: "active",
				},
				wantErr: true,
				errMsg:  "op: unsupported operator 'invalid_op'",
			},
			{
				name: "temporal condition missing field",
				condition: permissions.Condition{
					Op:     permissions.OpGraceBefore,
					Margin: "15m",
				},
				wantErr: true,
				errMsg:  "field: required for 'grace_before' operator",
			},
			{
				name: "temporal condition missing margin",
				condition: permissions.Condition{
					Field: "resource.start_time",
					Op:    permissions.OpGraceBefore,
				},
				wantErr: true,
				errMsg:  "margin: required for 'grace_before' operator",
			},
			{
				name: "temporal condition invalid duration",
				condition: permissions.Condition{
					Field:  "resource.start_time",
					Op:     permissions.OpGraceBefore,
					Margin: "invalid_duration",
				},
				wantErr: true,
				errMsg:  "margin: invalid duration",
			},
			{
				name: "grace_duration missing field_start",
				condition: permissions.Condition{
					FieldEnd: "resource.end_time",
					Op:       permissions.OpGraceDuration,
					Margin:   "1h",
				},
				wantErr: true,
				errMsg:  "field_start: required for 'grace_duration' operator",
			},
			{
				name: "grace_duration missing field_end",
				condition: permissions.Condition{
					FieldStart: "resource.start_time",
					Op:         permissions.OpGraceDuration,
					Margin:     "1h",
				},
				wantErr: true,
				errMsg:  "field_end: required for 'grace_duration' operator",
			},
			{
				name: "predicate missing path",
				condition: permissions.Condition{
					Op:    permissions.OpEq,
					Value: "active",
				},
				wantErr: true,
				errMsg:  "path: required for predicate operator",
			},
			{
				name: "predicate missing value and ref",
				condition: permissions.Condition{
					Path: "resource.status",
					Op:   permissions.OpEq,
				},
				wantErr: true,
				errMsg:  "either value or ref must be provided",
			},
			{
				name: "temporal condition with conflicting predicate fields",
				condition: permissions.Condition{
					Field:  "resource.start_time",
					Path:   "resource.status",
					Op:     permissions.OpGraceBefore,
					Margin: "15m",
					Value:  "active",
				},
				wantErr: true,
				errMsg:  "unexpected fields for grace operator",
			},
			{
				name: "predicate with conflicting temporal fields",
				condition: permissions.Condition{
					Path:   "resource.status",
					Op:     permissions.OpEq,
					Value:  "active",
					Field:  "resource.start_time",
					Margin: "15m",
				},
				wantErr: true,
				errMsg:  "temporal fields not allowed for predicate operator",
			},
			{
				name: "exists operator with value should fail",
				condition: permissions.Condition{
					Path:  "resource.field",
					Op:    permissions.OpExists,
					Value: "should_not_be_here",
				},
				wantErr: true,
				errMsg:  "value and ref should not be provided for 'exists' operator",
			},
			{
				name: "exists operator with ref should fail",
				condition: permissions.Condition{
					Path: "resource.field",
					Op:   permissions.OpExists,
					Ref:  "should_not_be_here",
				},
				wantErr: true,
				errMsg:  "value and ref should not be provided for 'exists' operator",
			},
			{
				name: "nested validation error",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{Path: "resource.status", Op: permissions.OpEq, Value: "active"},
						{Path: "resource.count", Op: permissions.OpGt}, // Missing value/ref
					},
				},
				wantErr: true,
				errMsg:  "and[1]: either value or ref must be provided",
			},

			// ==================== EDGE CASES ====================
			{
				name: "zero values",
				condition: permissions.Condition{
					Path:  "resource.count",
					Op:    permissions.OpEq,
					Value: 0,
				},
				wantErr: false,
			},
			{
				name: "empty string value",
				condition: permissions.Condition{
					Path:  "resource.name",
					Op:    permissions.OpEq,
					Value: "",
				},
				wantErr: false,
			},
			{
				name: "nested single condition in AND",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{Path: "resource.status", Op: permissions.OpEq, Value: "active"},
					},
				},
				wantErr: false,
			},
			{
				name: "multiple levels of nesting",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							Or: &[]permissions.Condition{
								{
									And: &[]permissions.Condition{
										{Path: "a", Op: permissions.OpEq, Value: 1},
										{Path: "b", Op: permissions.OpEq, Value: 2},
									},
								},
								{Path: "c", Op: permissions.OpEq, Value: 3},
							},
						},
					},
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := permissions.ValidateCondition(tt.condition)
				if (err != nil) != tt.wantErr {
					t.Errorf("ValidateCondition() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErr && tt.errMsg != "" && err != nil {
					if !contains(err.Error(), tt.errMsg) {
						t.Errorf("ValidateCondition() error = %v, want error containing %v", err, tt.errMsg)
					}
				}
			})
		}
	})
	t.Run("EncodeDecodeConditions", func(t *testing.T) {
		tests := []struct {
			name       string
			condition  *permissions.Condition
			wantErr    bool
			validateFn func(t *testing.T, decoded *permissions.Condition)
		}{
			// ==================== HAPPY PATHS ====================
			{
				name: "encode/decode simple predicate",
				condition: &permissions.Condition{
					Path:  "resource.status",
					Op:    permissions.OpEq,
					Value: "active",
				},
				wantErr: false,
				validateFn: func(t *testing.T, decoded *permissions.Condition) {
					if decoded.Path != "resource.status" {
						t.Errorf("Path = %v, want %v", decoded.Path, "resource.status")
					}
					if decoded.Op != permissions.OpEq {
						t.Errorf("Op = %v, want %v", decoded.Op, permissions.OpEq)
					}
					if decoded.Value != "active" {
						t.Errorf("Value = %v, want %v", decoded.Value, "active")
					}
				},
			},
			{
				name: "encode/decode temporal condition",
				condition: &permissions.Condition{
					Field:  "resource.start_time",
					Op:     permissions.OpGraceBefore,
					Margin: "15m",
				},
				wantErr: false,
				validateFn: func(t *testing.T, decoded *permissions.Condition) {
					if decoded.Field != "resource.start_time" {
						t.Errorf("Field = %v, want %v", decoded.Field, "resource.start_time")
					}
					if decoded.Margin != "15m" {
						t.Errorf("Margin = %v, want %v", decoded.Margin, "15m")
					}
				},
			},
			{
				name: "encode/decode complex nested condition",
				condition: &permissions.Condition{
					And: &[]permissions.Condition{
						{
							Or: &[]permissions.Condition{
								{Path: "resource.type", Op: permissions.OpEq, Value: "premium"},
								{Path: "resource.type", Op: permissions.OpEq, Value: "enterprise"},
							},
						},
						{
							Not: &permissions.Condition{
								Path:  "resource.banned",
								Op:    permissions.OpEq,
								Value: true,
							},
						},
						{
							Field:  "resource.start_time",
							Op:     permissions.OpGraceBefore,
							Margin: "30m",
						},
					},
				},
				wantErr: false,
				validateFn: func(t *testing.T, decoded *permissions.Condition) {
					if decoded.And == nil || len(*decoded.And) != 3 {
						t.Errorf("AND conditions count = %v, want %v", len(*decoded.And), 3)
					}
				},
			},
			{
				name:      "encode/decode nil condition",
				condition: nil,
				wantErr:   false,
				validateFn: func(t *testing.T, decoded *permissions.Condition) {
					if decoded != nil {
						t.Errorf("Decoded = %v, want nil", decoded)
					}
				},
			},
			{
				name: "encode/decode with numeric values",
				condition: &permissions.Condition{
					And: &[]permissions.Condition{
						{Path: "resource.count", Op: permissions.OpGt, Value: 100},
						{Path: "resource.price", Op: permissions.OpLte, Value: 99.99},
					},
				},
				wantErr: false,
				validateFn: func(t *testing.T, decoded *permissions.Condition) {
					countCond := (*decoded.And)[0]
					if countCond.Value != float64(100) { // JSON unmarshals numbers as float64
						t.Errorf("Count value = %v, want %v", countCond.Value, 100)
					}
				},
			},
			{
				name: "encode/decode with boolean values",
				condition: &permissions.Condition{
					Path:  "resource.active",
					Op:    permissions.OpEq,
					Value: true,
				},
				wantErr: false,
				validateFn: func(t *testing.T, decoded *permissions.Condition) {
					if decoded.Value != true {
						t.Errorf("Value = %v, want true", decoded.Value)
					}
				},
			},

			// ==================== EDGE CASES ====================
			{
				name: "encode/decode empty strings",
				condition: &permissions.Condition{
					Path:  "resource.name",
					Op:    permissions.OpEq,
					Value: "",
				},
				wantErr: false,
				validateFn: func(t *testing.T, decoded *permissions.Condition) {
					if decoded.Value != "" {
						t.Errorf("Value = %v, want empty string", decoded.Value)
					}
				},
			},
			{
				name: "encode/decode special characters in values",
				condition: &permissions.Condition{
					Path:  "resource.pattern",
					Op:    permissions.OpMatches,
					Value: "^[a-zA-Z0-9_\\-\\\\s]+$",
				},
				wantErr: false,
				validateFn: func(t *testing.T, decoded *permissions.Condition) {
					if decoded.Value != "^[a-zA-Z0-9_\\-\\\\s]+$" {
						t.Errorf("Value = %v, want special char pattern", decoded.Value)
					}
				},
			},
			{
				name: "encode/decode condition with all optional fields",
				condition: &permissions.Condition{
					Path:       "test.path",
					Op:         permissions.OpEq,
					Value:      "test",
					Ref:        "", // Empty but present
					Field:      "",
					Margin:     "",
					FieldStart: "",
					FieldEnd:   "",
				},
				wantErr: false,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				// Test encoding
				encoded, err := permissions.EncodeCondition(tt.condition)
				if (err != nil) != tt.wantErr {
					t.Errorf("EncodeCondition() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// Test decoding
				decoded, err := permissions.DecodeCondition(encoded)
				if (err != nil) != tt.wantErr {
					t.Errorf("DecodeCondition() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				// Run custom validation
				if tt.validateFn != nil {
					tt.validateFn(t, decoded)
				}

				// If original was not nil, verify round-trip equality
				if tt.condition != nil && !tt.wantErr {
					// Re-encode decoded to compare
					reEncoded, err := permissions.EncodeCondition(decoded)
					if err != nil {
						t.Errorf("Re-encode failed: %v", err)
						return
					}

					// Compare JSON strings
					if string(*encoded) != string(*reEncoded) {
						t.Errorf("Round-trip encoding mismatch:\nOriginal: %s\nRe-encoded: %s", *encoded, *reEncoded)
					}
				}
			})
		}
	})

	t.Run("DecodeInvalidJSON", func(t *testing.T) {
		tests := []struct {
			name    string
			rawJSON string
			wantErr bool
			errMsg  string
		}{
			{
				name:    "invalid JSON syntax",
				rawJSON: `{"path": "test", "op": "eq", "value": }`,
				wantErr: true,
				errMsg:  "failed to decode condition",
			},
			{
				name:    "malformed JSON",
				rawJSON: `not json at all`,
				wantErr: true,
				errMsg:  "failed to decode condition",
			},
			{
				name:    "JSON with wrong types",
				rawJSON: `{"path": 123, "op": true, "value": "test"}`,
				wantErr: true,
				errMsg:  "cannot unmarshal number into Go struct field Condition.path",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				raw := json.RawMessage(tt.rawJSON)
				_, err := permissions.DecodeCondition(&raw)
				if (err != nil) != tt.wantErr {
					t.Errorf("DecodeCondition() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErr && tt.errMsg != "" && err != nil {
					if !contains(err.Error(), tt.errMsg) {
						t.Errorf("error = %v, want error containing %v", err, tt.errMsg)
					}
				}
			})
		}
	})

	t.Run("DecodeAndValidate", func(t *testing.T) {
		tests := []struct {
			name       string
			rawJSON    string
			wantErr    bool
			errMsg     string
			validateFn func(t *testing.T, c *permissions.Condition)
		}{
			{
				name: "valid simple condition",
				rawJSON: `{
				"path": "resource.status",
				"op": "eq",
				"value": "active"
			}`,
				wantErr: false,
				validateFn: func(t *testing.T, c *permissions.Condition) {
					if c.Path != "resource.status" {
						t.Errorf("Path = %v, want %v", c.Path, "resource.status")
					}
				},
			},
			{
				name: "valid nested condition",
				rawJSON: `{
				"and": [
					{"path": "resource.type", "op": "eq", "value": "premium"},
					{"path": "resource.age", "op": "gte", "value": 18}
				]
			}`,
				wantErr: false,
			},
			{
				name: "invalid condition missing operator",
				rawJSON: `{
				"path": "resource.status",
				"value": "active"
			}`,
				wantErr: true,
				errMsg:  "op: operator is required",
			},
			{
				name: "invalid temporal condition",
				rawJSON: `{
				"field": "resource.start_time",
				"op": "grace_before",
				"margin": "invalid_duration"
			}`,
				wantErr: true,
				errMsg:  "invalid duration",
			},
			{
				name:    "malformed JSON should fail decode",
				rawJSON: `{"invalid": json}`,
				wantErr: true,
				errMsg:  "failed to decode condition",
			},
			{
				name:    "null JSON should return nil",
				rawJSON: `null`,
				wantErr: false,
				validateFn: func(t *testing.T, c *permissions.Condition) {
					if c != nil {
						t.Errorf("Expected nil for null JSON, got %v", c)
					}
				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				var raw *json.RawMessage
				if tt.rawJSON != "" {
					tmp := json.RawMessage(tt.rawJSON)
					raw = &tmp
				}

				c, err := permissions.DecodeAndValidateCondition(raw)
				if (err != nil) != tt.wantErr {
					t.Errorf("DecodeAndValidateCondition() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				if tt.wantErr && tt.errMsg != "" && err != nil {
					if !contains(err.Error(), tt.errMsg) {
						t.Errorf("error = %v, want error containing %v", err, tt.errMsg)
					}
				}
				if tt.validateFn != nil {
					tt.validateFn(t, c)
				}
			})
		}
	})

	t.Run("TestCondition_RoundTripRealWorldExamples", func(t *testing.T) {
		examples := []struct {
			name        string
			description string
			condition   permissions.Condition
		}{
			{
				name:        "own_resource_only",
				description: "User can only access their own resources",
				condition: permissions.Condition{
					Path: "subject.id",
					Op:   permissions.OpEq,
					Ref:  "resource.owner_id",
				},
			},
			{
				name:        "check_in_grace_period",
				description: "Check-in allowed 15m before event start only",
				condition: permissions.Condition{
					Field:  "resource.start_time",
					Op:     permissions.OpGraceBefore,
					Margin: "15m",
				},
			},
			{
				name:        "event_access_with_duration_grace",
				description: "Access during event with 30m grace on both sides",
				condition: permissions.Condition{
					FieldStart: "resource.start_time",
					FieldEnd:   "resource.end_time",
					Op:         permissions.OpGraceDuration,
					Margin:     "30m",
				},
			},
			{
				name:        "premium_user_or_admin",
				description: "Premium users or admins get access",
				condition: permissions.Condition{
					Or: &[]permissions.Condition{
						{Path: "resource.user_type", Op: permissions.OpEq, Value: "premium"},
						{Path: "resource.role", Op: permissions.OpEq, Value: "admin"},
					},
				},
			},
			{
				name:        "complex_event_access",
				description: "Complex: premium users during grace period, or admins anytime",
				condition: permissions.Condition{
					Or: &[]permissions.Condition{
						{
							And: &[]permissions.Condition{
								{Path: "resource.user_type", Op: permissions.OpEq, Value: "premium"},
								{
									Field:  "resource.start_time",
									Op:     permissions.OpGraceAround,
									Margin: "30m",
								},
							},
						},
						{Path: "resource.role", Op: permissions.OpEq, Value: "admin"},
					},
				},
			},
		}

		for _, ex := range examples {
			t.Run(ex.name, func(t *testing.T) {
				// First validate
				if err := permissions.ValidateCondition(ex.condition); err != nil {
					t.Fatalf("Real world example failed validation: %v", err)
				}

				// Then encode/decode
				encoded, err := permissions.EncodeCondition(&ex.condition)
				if err != nil {
					t.Fatalf("Failed to encode real world example: %v", err)
				}

				decoded, err := permissions.DecodeCondition(encoded)
				if err != nil {
					t.Fatalf("Failed to decode real world example: %v", err)
				}

				// Validate decoded
				if err := permissions.ValidateCondition(*decoded); err != nil {
					t.Errorf("Decoded condition failed validation: %v", err)
				}

				// Verify structure is preserved
				if ex.condition.Op != "" && decoded.Op != ex.condition.Op {
					t.Errorf("Op mismatch after round-trip")
				}
			})
		}
	})
}

// Helper function
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && indexOfSubstring(s, substr) >= 0)
}

func indexOfSubstring(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
