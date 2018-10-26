package endpoints

import (
	"context"
	"errors"
	"math"
	"regexp"
	"strings"

	"github.com/go-chi/jwtauth"
	"github.com/go-kit/kit/endpoint"

	jwtAuth "github.com/dwarvesf/smithy/backend/auth"
	"github.com/dwarvesf/smithy/backend/domain"
	"github.com/dwarvesf/smithy/backend/service"
)

// ChangePasswordRequest store change password structer
type ChangePasswordRequest struct {
	OldPassword             string `json:"old_password"`
	NewPassword             string `json:"new_password"`
	NewPasswordConfirmation string `json:"new_password_confirmation"`
}

// ChangePasswordResponse store change password respone
type ChangePasswordResponse struct {
	Complexity string `json:"complexity"`
}

const (
	TooShort   = "Too Short"
	VeryWeak   = "Very Weak"
	Weak       = "Weak"
	Good       = "Good"
	Strong     = "Strong"
	VeryStrong = "Very Strong"
)

func makeChangePasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ChangePasswordRequest)
		if !ok {
			return nil, errors.New("failed to make type assertion")
		}

		_, claims, _ := jwtauth.FromContext(ctx)

		var (
			userName                = claims["username"].(string)
			email                   = claims["email"].(string)
			isEmailAccount          = claims["is_email_account"].(bool)
			oldPassword             = req.OldPassword
			newPassword             = req.NewPassword
			newPasswordConfirmation = req.NewPasswordConfirmation
		)

		user := &domain.User{
			Username:       userName,
			Email:          email,
			IsEmailAccount: isEmailAccount,
		}

		user, err := s.UserService.Find(user)
		if err != nil {
			return nil, errors.New("username is invalid")
		}

		if user.Password != oldPassword {
			return nil, jwtAuth.ErrOldPasswordInvalid
		}

		if newPassword != newPasswordConfirmation {
			return nil, jwtAuth.ErrRePasswordIsNotMatch
		}

		complexity := checkPassword(newPasswordConfirmation)
		if complexity == VeryWeak || complexity == TooShort {
			return nil, jwtAuth.ErrPassWordIsVeryWeak
		}

		user.Password = newPassword
		_, err = s.UserService.Update(user)
		if err != nil {
			return nil, err
		}

		return ChangePasswordResponse{complexity}, nil
	}
}

func checkPassword(pwd string) string {
	var (
		nScore, nRequirements                                                                                           int
		nLength                                                                                                         int
		nRepInc                                                                                                         float64
		nMultMidChar                                                                                                    = 2
		nMultConsecAlphaUC, nMultConsecAlphaLC, nMultConsecNumber                                                       = 2, 2, 2
		nMultSeqAlpha, nMultSeqNumber, nMultSeqSymbol                                                                   = 3, 3, 3
		nConsecAlphaUC, nConsecAlphaLC, nConsecNumber, nConsecSymbol, nConsecCharType                                   = 0, 0, 0, 0, 0
		nAlphaUC, nAlphaLC, nNumber, nSymbol, nMidChar, nReqChar, nSeqAlpha, nSeqNumber, nSeqSymbol, nSeqChar, nRepChar = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
		nTmpAlphaUC, nTmpAlphaLC, nTmpNumber, nTmpSymbol                                                                = -1, -1, -1, -1
		sAlphas                                                                                                         = "abcdefghijklmnopqrstuvwxyz"
		sNumerics                                                                                                       = "01234567890"
		sSymbols                                                                                                        = ")!@#$%^&*()"
		nMultLength, nMultNumber                                                                                        = 4, 4
		nMultSymbol                                                                                                     = 6
		sComplexity                                                                                                     = "Too Short"
		nMinPwdLen                                                                                                      = 8
	)

	if pwd != "" {
		nLength = len(pwd)
		nScore = int(nLength * nMultLength)
		arrPwd := strings.Split(pwd, "")
		var arrPwdLen = len(arrPwd)
		var validUpcase = regexp.MustCompile(`[A-Z]`)
		var validlowCase = regexp.MustCompile(`[a-z]`)
		var validNumber = regexp.MustCompile(`[0-9]`)
		var validUpOthers = regexp.MustCompile(`[^a-zA-Z0-9_]`)

		for i, charA := range arrPwd {
			if validUpcase.MatchString(charA) {
				if nTmpAlphaUC != -1 {
					if (nTmpAlphaUC + 1) == i {
						nConsecAlphaUC++
						nConsecCharType++
					}
				}
				nTmpAlphaUC = i
				nAlphaUC++
			} else if validlowCase.MatchString(charA) {
				if nTmpAlphaLC != -1 {
					if (nTmpAlphaLC + 1) == i {
						nConsecAlphaLC++
						nConsecCharType++
					}
				}
				nTmpAlphaLC = i
				nAlphaLC++
			} else if validNumber.MatchString(charA) {
				if i > 0 && i < (arrPwdLen-1) {
					nMidChar++
				}
				if nTmpNumber != -1 {
					if (nTmpNumber + 1) == i {
						nConsecNumber++
						nConsecCharType++
					}
				}
				nTmpNumber = i
				nNumber++
			} else if validUpOthers.MatchString(charA) {
				if i > 0 && i < (arrPwdLen-1) {
					nMidChar++
				}
				if nTmpSymbol != -1 {
					if (nTmpSymbol + 1) == i {
						nConsecSymbol++
						nConsecCharType++
					}
				}
				nTmpSymbol = i
				nSymbol++
			}

			/* Internal loop through password to check for repeat characters */
			var bCharExists = false
			for j, charB := range arrPwd {
				if charA == charB && i != j { /* repeat character exists */
					bCharExists = true
					/*
						Calculate icrement deduction based on proximity to identical characters
						Deduction is incremented each time a new match is discovered
						Deduction amount is based on total password length divided by the
						difference of distance between currently selected match
					*/
					nRepInc += math.Abs(float64(arrPwdLen / (j - i)))
				}
			}
			if bCharExists {
				nRepChar++
				nUnqChar := float64(arrPwdLen - nRepChar)
				if nUnqChar != 0 {
					nRepInc = math.Ceil(nRepInc / nUnqChar)
				} else {
					nRepInc = math.Ceil(nRepInc)
				}
			}
		}

		/* Check for sequential alpha string patterns (forward and reverse) */
		for s := 0; s < 23; s++ {
			var sFwd = sAlphas[s:(s + 3)]
			var sRev = Reverse(sFwd)
			if strings.Contains(strings.ToLower(pwd), sFwd) || strings.Contains(strings.ToLower(pwd), sRev) {
				nSeqAlpha++
				nSeqChar++
			}
		}
		/* Check for sequential numeric string patterns (forward and reverse) */
		for s := 0; s < 8; s++ {
			var sFwd = sNumerics[s:(s + 3)]
			var sRev = Reverse(sFwd)
			if strings.Contains(strings.ToLower(pwd), sFwd) || strings.Contains(strings.ToLower(pwd), sRev) {
				nSeqNumber++
				nSeqChar++
			}
		}

		/* Check for sequential symbol string patterns (forward and reverse)  */
		for s := 0; s < 8; s++ {
			var sFwd = sSymbols[s:(s + 3)]
			var sRev = Reverse(sFwd)
			if strings.Contains(strings.ToLower(pwd), sFwd) || strings.Contains(strings.ToLower(pwd), sRev) {
				nSeqSymbol++
				nSeqChar++
			}
		}

		if nAlphaUC > 0 && nAlphaUC < nLength {
			nScore = nScore + ((nLength - nAlphaUC) * 2)
		}
		if nAlphaLC > 0 && nAlphaLC < nLength {
			nScore = nScore + ((nLength - nAlphaLC) * 2)
		}
		if nNumber > 0 && nNumber < nLength {
			nScore = nScore + (nNumber * nMultNumber)
		}
		if nSymbol > 0 {
			nScore = nScore + (nSymbol * nMultSymbol)
		}
		if nMidChar > 0 {
			nScore = nScore + (nMidChar * nMultMidChar)
		}

		/* Point deductions for poor practices */
		if (nAlphaLC > 0 || nAlphaUC > 0) && nSymbol == 0 && nNumber == 0 { // Only Letters
			nScore = nScore - nLength
		}
		if nAlphaLC == 0 && nAlphaUC == 0 && nSymbol == 0 && nNumber > 0 { // Only Numbers
			nScore = nScore - nLength
		}
		if nRepChar > 0 { // Same character exists more than once
			nScore = nScore - int(nRepInc)
		}
		if nConsecAlphaUC > 0 { // Consecutive Uppercase Letters exist
			nScore = nScore - (nConsecAlphaUC * nMultConsecAlphaUC)
		}
		if nConsecAlphaLC > 0 { // Consecutive Lowercase Letters exist
			nScore = nScore - (nConsecAlphaLC * nMultConsecAlphaLC)
		}
		if nConsecNumber > 0 { // Consecutive Numbers exist
			nScore = nScore - (nConsecNumber * nMultConsecNumber)
		}
		if nSeqAlpha > 0 { // Sequential alpha strings exist (3 characters or more)
			nScore = nScore - (nSeqAlpha * nMultSeqAlpha)
		}
		if nSeqNumber > 0 { // Sequential numeric strings exist (3 characters or more)
			nScore = nScore - (nSeqNumber * nMultSeqNumber)
		}
		if nSeqSymbol > 0 { // Sequential symbol strings exist (3 characters or more)
			nScore = nScore - (nSeqSymbol * nMultSeqSymbol)
		}

		var arrChars = [6]int{nLength, nAlphaUC, nAlphaLC, nNumber, nSymbol}
		var arrCharsIds = [6]string{"nLength", "nAlphaUC", "nAlphaLC", "nNumber", "nSymbol"}
		var arrCharsLen = len(arrChars)
		for c := 0; c < arrCharsLen; c++ {
			var minVal = 0
			if arrCharsIds[c] == "nLength" {
				minVal = nMinPwdLen - 1
			}
			if arrChars[c] == minVal+1 || arrChars[c] > minVal+1 {
				nReqChar++
			}
		}

		nRequirements = nReqChar
		var nMinReqChars int
		if len(pwd) >= nMinPwdLen {
			nMinReqChars = 3
		} else {
			nMinReqChars = 4
		}
		if nRequirements > nMinReqChars { // One or more required characters exist
			nScore = nScore + (nRequirements * 2)
		}

		/* Determine complexity based on overall score */
		if nScore > 100 {
			nScore = 100
		} else if nScore < 0 {
			nScore = 0
		}
		if nScore >= 0 && nScore < 20 {
			sComplexity = VeryWeak
		} else if nScore >= 20 && nScore < 40 {
			sComplexity = Weak
		} else if nScore >= 40 && nScore < 60 {
			sComplexity = Good
		} else if nScore >= 60 && nScore < 80 {
			sComplexity = Strong
		} else if nScore >= 80 && nScore <= 100 {
			sComplexity = VeryStrong
		}
	}

	return sComplexity
}

//Reverse string
func Reverse(s string) string {
	var reverse string
	for i := len(s) - 1; i >= 0; i-- {
		reverse += string(s[i])
	}
	return reverse
}
