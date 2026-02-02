package testing

import (
	"GoAuth/internal/domain/permissions"
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/MintzyG/fail"
)

func testConditions(t *testing.T) {
	t.Run("ValidateConditions", func(t *testing.T) {
		tests := []struct {
			name      string
			condition permissions.Condition
			wantErr   bool
			errMsg    string
			wantPath  string
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
			{
				name: "deeply nested path",
				condition: permissions.Condition{
					Path:  "resource.metadata.audit.created_by.id",
					Op:    permissions.OpEq,
					Value: "user-123",
				},
				wantErr: false,
			},

			// ==================== SAD PATHS ====================
			{
				name: "empty AND operator",
				condition: permissions.Condition{
					And: &[]permissions.Condition{},
				},
				wantErr:  true,
				errMsg:   "and: AND conditions cannot be empty",
				wantPath: "and",
			},
			{
				name: "empty OR operator",
				condition: permissions.Condition{
					Or: &[]permissions.Condition{},
				},
				wantErr:  true,
				errMsg:   "or: OR conditions cannot be empty",
				wantPath: "or",
			},
			{
				name: "missing operator",
				condition: permissions.Condition{
					Path:  "resource.status",
					Value: "active",
				},
				wantErr:  true,
				errMsg:   "op: operator is required",
				wantPath: "'.'",
			},
			{
				name: "invalid operator name",
				condition: permissions.Condition{
					Path:  "resource.status",
					Op:    "invalid_op",
					Value: "active",
				},
				wantErr:  true,
				errMsg:   "op: unsupported operator 'invalid_op'",
				wantPath: "'.'",
			},
			{
				name: "temporal condition missing field",
				condition: permissions.Condition{
					Op:     permissions.OpGraceBefore,
					Margin: "15m",
				},
				wantErr:  true,
				errMsg:   "field: required for 'grace_before' operator",
				wantPath: "'.'",
			},
			{
				name: "temporal condition missing margin",
				condition: permissions.Condition{
					Field: "resource.start_time",
					Op:    permissions.OpGraceBefore,
				},
				wantErr:  true,
				errMsg:   "margin: required for 'grace_before' operator",
				wantPath: "'.'",
			},
			{
				name: "temporal condition invalid duration",
				condition: permissions.Condition{
					Field:  "resource.start_time",
					Op:     permissions.OpGraceBefore,
					Margin: "invalid_duration",
				},
				wantErr:  true,
				errMsg:   "margin: invalid duration",
				wantPath: "'.'",
			},
			{
				name: "grace_duration missing field_start",
				condition: permissions.Condition{
					FieldEnd: "resource.end_time",
					Op:       permissions.OpGraceDuration,
					Margin:   "1h",
				},
				wantErr:  true,
				errMsg:   "field_start: required for 'grace_duration' operator",
				wantPath: "'.'",
			},
			{
				name: "grace_duration missing field_end",
				condition: permissions.Condition{
					FieldStart: "resource.start_time",
					Op:         permissions.OpGraceDuration,
					Margin:     "1h",
				},
				wantErr:  true,
				errMsg:   "field_end: required for 'grace_duration' operator",
				wantPath: "'.'",
			},
			{
				name: "predicate missing path",
				condition: permissions.Condition{
					Op:    permissions.OpEq,
					Value: "active",
				},
				wantErr:  true,
				errMsg:   "path: required for predicate operator",
				wantPath: "'.'",
			},
			{
				name: "predicate missing value and ref",
				condition: permissions.Condition{
					Path: "resource.status",
					Op:   permissions.OpEq,
				},
				wantErr:  true,
				errMsg:   "either value or ref must be provided",
				wantPath: "'.'",
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
				wantErr:  true,
				errMsg:   "unexpected fields for grace operator",
				wantPath: "'.'",
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
				wantErr:  true,
				errMsg:   "temporal fields not allowed for predicate operator",
				wantPath: "'.'",
			},
			{
				name: "exists operator with value should fail",
				condition: permissions.Condition{
					Path:  "resource.field",
					Op:    permissions.OpExists,
					Value: "should_not_be_here",
				},
				wantErr:  true,
				errMsg:   "value and ref should not be provided for 'exists' operator",
				wantPath: "'.'",
			},
			{
				name: "exists operator with ref should fail",
				condition: permissions.Condition{
					Path: "resource.field",
					Op:   permissions.OpExists,
					Ref:  "should_not_be_here",
				},
				wantErr:  true,
				errMsg:   "value and ref should not be provided for 'exists' operator",
				wantPath: "'.'",
			},
			{
				name: "nested validation error",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{Path: "resource.status", Op: permissions.OpEq, Value: "active"},
						{Path: "resource.count", Op: permissions.OpGt}, // Missing value/ref
					},
				},
				wantErr:  true,
				errMsg:   "either value or ref must be provided",
				wantPath: "and[1]",
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
				if tt.wantErr && err != nil {
					if tt.wantPath != "" && !contains(err.Error(), tt.wantPath) {
						t.Errorf("ValidateCondition() error = %v, want path %v", err, tt.wantPath)
					}

					if tt.errMsg != "" {
						if contains(err.Error(), tt.errMsg) {
							return
						}

						var failErr *fail.Error
						if errors.As(err, &failErr) && failErr.Meta != nil {
							traces := toStringSlice(failErr.Meta["traces"])
							for _, trace := range traces {
								if contains(trace, tt.errMsg) {
									return
								}
							}
						}

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
			wantPath   string
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
				wantErr:  true,
				errMsg:   "op: operator is required",
				wantPath: "'.'",
			},
			{
				name: "invalid temporal condition",
				rawJSON: `{
				"field": "resource.start_time",
				"op": "grace_before",
				"margin": "invalid_duration"
			}`,
				wantErr:  true,
				errMsg:   "invalid duration",
				wantPath: "'.'",
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
				if tt.wantErr && err != nil {
					if tt.wantPath != "" && !contains(err.Error(), tt.wantPath) {
						t.Errorf("DecodeAndValidateCondition() error = %v, want path %v", err, tt.wantPath)
					}

					if tt.errMsg != "" {
						if contains(err.Error(), tt.errMsg) {
							return
						}

						var failErr *fail.Error
						if errors.As(err, &failErr) && failErr.Meta != nil {
							traces := toStringSlice(failErr.Meta["traces"])
							for _, trace := range traces {
								if contains(trace, tt.errMsg) {
									return
								}
							}
						}

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

	t.Run("Complex Resource Tests", func(t *testing.T) {
		now := time.Now().UTC()

		tests := []struct {
			name        string
			condition   permissions.Condition
			context     permissions.ConditionContext
			shouldPass  bool
			description string
		}{
			// ==================== PASSING TESTS (Complex Valid Cases) ====================
			{
				name: "complex_event_access_with_grace_and_ownership",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							// User owns the event OR is an admin
							Or: &[]permissions.Condition{
								{
									Path: "subject.id",
									Op:   permissions.OpEq,
									Ref:  "resource.event.owner.id",
								},
								{
									Path:  "subject.role",
									Op:    permissions.OpEq,
									Value: "admin",
								},
							},
						},
						{
							// Within 30m grace period before start
							Field:  "resource.event.start_time",
							Op:     permissions.OpGraceBefore,
							Margin: "30m",
						},
						{
							// Event capacity check
							Path: "resource.event.attendees",
							Op:   permissions.OpLt,
							Ref:  "resource.event.max_capacity",
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{
						"id":   "user-123",
						"role": "admin",
					},
					Resource: map[string]interface{}{
						"event": map[string]interface{}{
							"start_time": now.Add(15 * time.Minute),
							"owner": map[string]interface{}{
								"id": "user-456",
							},
							"attendees":    45,
							"max_capacity": 50,
						},
					},
					Environment: map[string]interface{}{
						"now": now,
					},
				},
				shouldPass:  true,
				description: "Admin user within grace period, capacity available - should pass",
			},
			{
				name: "activity_enrollment_with_prerequisites_and_department",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							// Check if user has prerequisite course
							Path: "resource.activity.prerequisites",
							Op:   permissions.OpContains,
							Ref:  "subject.completed_courses[0]",
						},
						{
							// Check department match
							Or: &[]permissions.Condition{
								{
									Path:  "subject.department",
									Op:    permissions.OpEq,
									Value: "engineering",
								},
								{
									Path:  "subject.department",
									Op:    permissions.OpEq,
									Value: "training",
								},
							},
						},
						{
							// Activity not full
							Path:  "resource.activity.enrolled",
							Op:    permissions.OpLt,
							Value: 30,
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{
						"id":                "user-789",
						"department":        "engineering",
						"completed_courses": []string{"basic-programming", "git-fundamentals"},
					},
					Resource: map[string]interface{}{
						"activity": map[string]interface{}{
							"prerequisites": []string{"basic-programming"},
							"enrolled":      25,
						},
					},
				},
				shouldPass:  true,
				description: "User has prerequisites and correct department - should pass",
			},
			{
				name: "resource_grace_duration_with_location_constraint",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							// Access during event duration with grace
							FieldStart: "resource.event.start_time",
							FieldEnd:   "resource.event.end_time",
							Op:         permissions.OpGraceDuration,
							Margin:     "15m",
						},
						{
							// Specific location requirement
							Path:  "resource.event.location.building",
							Op:    permissions.OpEq,
							Value: "HQ",
						},
						{
							Path:  "resource.event.location.floor",
							Op:    permissions.OpIn,
							Value: []interface{}{"1", "2", "3"},
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{"id": "user-111"},
					Resource: map[string]interface{}{
						"event": map[string]interface{}{
							"start_time": now.Add(-10 * time.Minute),
							"end_time":   now.Add(1 * time.Hour),
							"location": map[string]interface{}{
								"building": "HQ",
								"floor":    "2",
							},
						},
					},
					Environment: map[string]interface{}{
						"now": now,
					},
				},
				shouldPass:  true,
				description: "Within grace duration and correct location - should pass",
			},
			{
				name: "nested_logical_with_array_matching",
				condition: permissions.Condition{
					Or: &[]permissions.Condition{
						{
							And: &[]permissions.Condition{
								{
									Path:  "resource.project.status",
									Op:    permissions.OpEq,
									Value: "active",
								},
								{
									Path:  "subject.teams",
									Op:    permissions.OpContains,
									Value: "team-alpha",
								},
							},
						},
						{
							Not: &permissions.Condition{
								Path:  "resource.project.archived",
								Op:    permissions.OpEq,
								Value: true,
							},
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{
						"id":    "user-222",
						"teams": []string{"team-beta", "team-alpha"},
					},
					Resource: map[string]interface{}{
						"project": map[string]interface{}{
							"status":   "active",
							"archived": false,
						},
					},
				},
				shouldPass:  true,
				description: "User in team-alpha on active project - should pass",
			},
			{
				name: "complex_string_pattern_matching",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							Path:  "resource.document.name",
							Op:    permissions.OpStartsWith,
							Value: "CONFIDENTIAL_",
						},
						{
							Path:  "resource.document.name",
							Op:    permissions.OpEndsWith,
							Value: "_INTERNAL.pdf",
						},
						{
							Or: &[]permissions.Condition{
								{
									Path:  "subject.clearance_level",
									Op:    permissions.OpGte,
									Value: 3,
								},
								{
									Path:  "subject.department",
									Op:    permissions.OpEq,
									Value: "legal",
								},
							},
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{
						"id":              "user-333",
						"clearance_level": 4,
						"department":      "engineering",
					},
					Resource: map[string]interface{}{
						"document": map[string]interface{}{
							"name": "CONFIDENTIAL_Q3_REPORT_INTERNAL.pdf",
						},
					},
				},
				shouldPass:  true,
				description: "High clearance level with matching filename pattern - should pass",
			},

			// ==================== FAILING TESTS (Complex Invalid Cases) ====================
			{
				name: "event_access_outside_grace_period",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							Field:  "resource.event.start_time",
							Op:     permissions.OpGraceBefore,
							Margin: "30m",
						},
						{
							Path:  "resource.event.status",
							Op:    permissions.OpEq,
							Value: "scheduled",
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{"id": "user-444"},
					Resource: map[string]interface{}{
						"event": map[string]interface{}{
							"start_time": now.Add(45 * time.Minute), // Outside 30m grace
							"status":     "scheduled",
						},
					},
					Environment: map[string]interface{}{
						"now": now,
					},
				},
				shouldPass:  false,
				description: "Outside grace period - should fail",
			},
			{
				name: "activity_missing_prerequisites",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							Path:  "resource.activity.prerequisites",
							Op:    permissions.OpContainsAll,
							Value: []interface{}{"advanced-programming", "system-design"},
						},
						{
							Path:  "subject.completed_courses",
							Op:    permissions.OpContains,
							Value: "advanced-programming",
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{
						"id":                "user-555",
						"completed_courses": []string{"basic-programming", "git-fundamentals"},
					},
					Resource: map[string]interface{}{
						"activity": map[string]interface{}{
							"prerequisites": []string{"advanced-programming", "system-design"},
						},
					},
				},
				shouldPass:  false,
				description: "Missing system-design prerequisite - should fail",
			},
			{
				name: "location_mismatch_with_capacity_full",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							FieldStart: "resource.event.start_time",
							FieldEnd:   "resource.event.end_time",
							Op:         permissions.OpGraceDuration,
							Margin:     "10m",
						},
						{
							Path:  "resource.event.location.room",
							Op:    permissions.OpEq,
							Value: "A-101",
						},
						{
							Path:  "resource.event.attendees",
							Op:    permissions.OpLt,
							Value: 50,
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{"id": "user-666"},
					Resource: map[string]interface{}{
						"event": map[string]interface{}{
							"start_time": now.Add(-5 * time.Minute),
							"end_time":   now.Add(1 * time.Hour),
							"location": map[string]interface{}{
								"room": "B-205", // Wrong room
							},
							"attendees": 60, // Over capacity
						},
					},
					Environment: map[string]interface{}{
						"now": now,
					},
				},
				shouldPass:  false,
				description: "Wrong room and over capacity - should fail",
			},
			{
				name: "mixed_nested_conditions_failing_all",
				condition: permissions.Condition{
					Or: &[]permissions.Condition{
						{
							And: &[]permissions.Condition{
								{
									Path:  "resource.project.visibility",
									Op:    permissions.OpEq,
									Value: "public",
								},
								{
									Not: &permissions.Condition{
										Path:  "resource.project.restricted",
										Op:    permissions.OpEq,
										Value: true,
									},
								},
							},
						},
						{
							And: &[]permissions.Condition{
								{
									Path:  "subject.department",
									Op:    permissions.OpEq,
									Value: "engineering",
								},
								{
									Path:  "subject.security_clearance",
									Op:    permissions.OpGte,
									Value: 2,
								},
							},
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{
						"id":                 "user-777",
						"department":         "marketing", // Wrong dept
						"security_clearance": 1,           // Too low
					},
					Resource: map[string]interface{}{
						"project": map[string]interface{}{
							"visibility": "private", // Not public
							"restricted": true,      // Restricted
						},
					},
				},
				shouldPass:  false,
				description: "Neither OR branch conditions met - should fail",
			},
			{
				name: "string_pattern_mismatch_with_insufficient_clearance",
				condition: permissions.Condition{
					And: &[]permissions.Condition{
						{
							Path:  "resource.file.path",
							Op:    permissions.OpMatches,
							Value: `^/secure/[a-z]+/classified_.*\.pdf$`,
						},
						{
							Or: &[]permissions.Condition{
								{
									Path:  "subject.role",
									Op:    permissions.OpEq,
									Value: "security-admin",
								},
								{
									Path:  "subject.clearance_level",
									Op:    permissions.OpGte,
									Value: 5,
								},
							},
						},
					},
				},
				context: permissions.ConditionContext{
					Subject: map[string]interface{}{
						"id":              "user-888",
						"clearance_level": 3,               // Too low
						"role":            "standard-user", // Wrong role
					},
					Resource: map[string]interface{}{
						"file": map[string]interface{}{
							"path": "/secure/finance/classified_budget.pdf",
						},
					},
				},
				shouldPass:  false,
				description: "Pattern matches but insufficient clearance - should fail",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				ok, motive, err := tt.condition.Evaluate(context.Background(), &tt.context)
				if err != nil {
					t.Fatalf("Evaluate() unexpected error: %v", err)
				}

				if ok != tt.shouldPass {
					t.Errorf("Evaluate() = %v, want %v. Description: %s, Motive: %+v",
						ok, tt.shouldPass, tt.description, motive)
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

func toStringSlice(v any) []string {
	if v == nil {
		return []string{}
	}

	if s, ok := v.([]string); ok {
		return s
	}

	if arr, ok := v.([]any); ok {
		result := make([]string, 0, len(arr))
		for _, item := range arr {
			if s, ok := item.(string); ok {
				result = append(result, s)
			}
		}
		return result
	}

	return []string{}
}
