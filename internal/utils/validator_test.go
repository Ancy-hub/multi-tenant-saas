package utils

import "testing"

func TestIsValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{"Valid standard email", "test@example.com", true},
		{"Valid with plus", "test+tag@example.com", true},
		{"Valid with numbers", "user123@domain.co", true},
		{"Invalid missing @", "testexample.com", false},
		{"Invalid missing domain", "test@", false},
		{"Invalid missing TLD", "test@example", false},
		{"Invalid spaces", "test @example.com", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidEmail(tt.email)
			if result != tt.expected {
				t.Errorf("IsValidEmail(%q) = %v; want %v", tt.email, result, tt.expected)
			}
		})
	}
}

func TestIsStrongPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		expected bool
	}{
		{"Strong password", "StrongPass1!", true},
		{"Strong with symbols", "P@ssw0rd2024#", true},
		{"Too short", "Sh0rt!", false},
		{"Missing uppercase", "weakpass1!", false},
		{"Missing lowercase", "WEAKPASS1!", false},
		{"Missing number", "WeakPassword!", false},
		{"Missing special char", "WeakPass1234", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsStrongPassword(tt.password)
			if result != tt.expected {
				t.Errorf("IsStrongPassword(%q) = %v; want %v", tt.password, result, tt.expected)
			}
		})
	}
}
