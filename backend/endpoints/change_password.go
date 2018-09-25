package endpoints

import (
	"context"
	"errors"
	"io/ioutil"
	"math"
	"regexp"
	"strings"

	"github.com/go-chi/jwtauth"
	"github.com/go-kit/kit/endpoint"
	yaml "gopkg.in/yaml.v2"

	jwtAuth "github.com/dwarvesf/smithy/backend/auth"
	backendConfig "github.com/dwarvesf/smithy/backend/config"

	"github.com/dwarvesf/smithy/backend/service"
)

// ChangePasswordRequest store change password structer
type ChangePasswordRequest struct {
	OldPassword   string `json:"old_password"`
	NewPassword   string `json:"new_password"`
	ReNewPassword string `json:"re_new_password"`
}

// ChangePasswordResponse store change password respone
type ChangePasswordResponse struct {
	Complexity string `json:"complexity"`
}

const (
	veryWeakStr  = "Very Weak"
	weakStr      = "Weak"
	goodStr      = "Good"
	strongStr    = "Strong"
	verStrongStr = "Very Strong"
)

func makeChangePasswordEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, ok := request.(ChangePasswordRequest)

		var (
			oldPassword   = req.OldPassword
			newPassword   = req.NewPassword
			reNewPassword = req.ReNewPassword
		)

		_, claims, _ := jwtauth.FromContext(ctx)
		cfg := s.SyncConfig()
		userMap := cfg.ConvertUserListToMap()
		userInfo, ok := userMap[claims["username"].(string)]

		if !ok {
			return nil, errors.New("username is invalid")
		}

		if userInfo.Password != oldPassword {
			return nil, jwtAuth.ErrOldPasswordInvalid
		}

		if newPassword != reNewPassword {
			return nil, jwtAuth.ErrRePasswordIsNotMatch
		}

		complexity := checkPassword(reNewPassword)
		if complexity == veryWeakStr || complexity == weakStr {
			return nil, jwtAuth.ErrPassWordIsVeryWeak
		}

		tmpCfg := cloneConfig(cfg, userInfo, newPassword)

		buff, err := yaml.Marshal(tmpCfg)
		if err != nil {
			return nil, err
		}

		err = ioutil.WriteFile("example_dashboard_config.yaml", buff, 0644)
		if err != nil {
			return nil, err
		}

		return ChangePasswordResponse{complexity}, nil
	}
}

func cloneConfig(cfg *backendConfig.Config, userInfo backendConfig.User, newPassword string) *backendConfig.Config {
	tmpCfg := &backendConfig.Config{}
	tmpCfg.SerectKey = cfg.SerectKey
	tmpCfg.AgentURL = cfg.AgentURL
	tmpCfg.PersistenceFileName = cfg.PersistenceFileName
	tmpCfg.PersistenceSupport = cfg.PersistenceSupport

	for i, user := range cfg.Authentication.UserList {
		if user.Username == userInfo.Username {
			cfg.Authentication.UserList[i].Password = newPassword
		}
	}
	tmpCfg.Authentication = cfg.Authentication

	return tmpCfg
}

func checkPassword(pwd string) string {
	var (
		nScore, nLength                                                                                                                                   = 0, 0
		nMultMidChar                                                                                                                                      = 2
		nConsecAlphaUC, nConsecAlphaLC, nConsecNumber, nConsecSymbol, nConsecCharType                                                                     = 0, 0, 0, 0, 0
		nAlphaUC, nAlphaLC, nNumber, nSymbol, nMidChar, nRequirements, nReqChar, nRepInc, nSeqAlpha, nSeqNumber, nSeqSymbol, nSeqChar, nRepChar, nUnqChar = 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0
		nTmpAlphaUC, nTmpAlphaLC, nTmpNumber, nTmpSymbol                                                                                                  = -1, -1, -1, -1
		sAlphas                                                                                                                                           = "abcdefghijklmnopqrstuvwxyz"
		sNumerics                                                                                                                                         = "01234567890"
		sSymbols                                                                                                                                          = ")!@#$%^&*()"
		nMultLength, nMultNumber                                                                                                                          = 4, 4
		nMultSymbol                                                                                                                                       = 6
		sComplexity                                                                                                                                       = "Too Short"
		nMinPwdLen                                                                                                                                        = 8
	)

	if pwd != "" {
		nLength = len(pwd)
		nScore = int(nLength * nMultLength)
		arrPwd := strings.Split(pwd, "")
		var arrPwdLen = len(arrPwd)
		var validUpcase = regexp.MustCompile(`/[A-Z]/g`)
		var validlowCase = regexp.MustCompile(`/[A-Z]/g`)
		var validNumber = regexp.MustCompile(`/[A-Z]/g`)
		var validUpOthers = regexp.MustCompile(`/[^a-zA-Z0-9_]/g`)

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
					nRepInc = nRepInc + int(math.Abs(float64(arrPwdLen/(j-i))))
				}
			}
			if bCharExists {
				nRepChar++
				nUnqChar = arrPwdLen - nRepChar
				if nUnqChar != 0 {
					nRepInc = int(math.Ceil(float64(nRepInc / nUnqChar)))
				} else {
					nRepInc = int(math.Ceil(float64(nRepInc)))
				}
			}
		}

		/* Check for sequential alpha string patterns (forward and reverse) */
		for s := 0; s < 23; s++ {
			var sFwd = sAlphas[s:(s + 3)]
			var sRev = Reverse(sFwd)
			if strings.ContainsAny(strings.ToLower(pwd), sFwd) || strings.ContainsAny(strings.ToLower(pwd), sRev) {
				nSeqAlpha++
				nSeqChar++
			}
		}

		/* Check for sequential numeric string patterns (forward and reverse) */
		for s := 0; s < 8; s++ {
			var sFwd = sNumerics[s:(s + 3)]
			var sRev = Reverse(sFwd)
			if strings.ContainsAny(strings.ToLower(pwd), sFwd) || strings.ContainsAny(strings.ToLower(pwd), sRev) {
				nSeqNumber++
				nSeqChar++
			}
		}

		/* Check for sequential symbol string patterns (forward and reverse)  */
		for s := 0; s < 8; s++ {
			var sFwd = sSymbols[s:(s + 3)]
			var sRev = Reverse(sFwd)
			if strings.ContainsAny(strings.ToLower(pwd), sFwd) || strings.ContainsAny(strings.ToLower(pwd), sRev) {
				nSeqSymbol++
				nSeqChar++
			}
		}

		if nAlphaUC > 0 && nAlphaUC < nLength {
			nScore = int(nScore + ((nLength - nAlphaUC) * 2))
		}
		if nAlphaLC > 0 && nAlphaLC < nLength {
			nScore = int(nScore + ((nLength - nAlphaLC) * 2))
		}
		if nNumber > 0 && nNumber < nLength {
			nScore = int(nScore + (nNumber * nMultNumber))
		}
		if nSymbol > 0 {
			nScore = int(nScore + (nSymbol * nMultSymbol))
		}
		if nMidChar > 0 {
			nScore = int(nScore + (nMidChar * nMultMidChar))
		}

		var arrChars = [6]int{nLength, nAlphaUC, nAlphaLC, nNumber, nSymbol}
		var arrCharsIds = [6]string{"nLength", "nAlphaUC", "nAlphaLC", "nNumber", "nSymbol"}
		var arrCharsLen = len(arrChars)
		for c := 0; c < arrCharsLen; c++ {
			var minVal = 0
			if arrCharsIds[c] == "nLength" {
				minVal = int(nMinPwdLen - 1)
			}
			if arrChars[c] == int(minVal+1) {
				nReqChar++
			} else if arrChars[c] > int(minVal+1) {
				nReqChar++
			}
		}
		nRequirements = nReqChar
		var nMinReqChars = 0
		if len(pwd) >= nMinPwdLen {
			nMinReqChars = 3
		} else {
			nMinReqChars = 4
		}
		if nRequirements > nMinReqChars { // One or more required characters exist
			nScore = int(nScore + (nRequirements * 2))
		}

		/* Determine complexity based on overall score */
		if nScore > 100 {
			nScore = 100
		} else if nScore < 0 {
			nScore = 0
		}
		if nScore >= 0 && nScore < 20 {
			sComplexity = "Very Weak"
		} else if nScore >= 20 && nScore < 40 {
			sComplexity = "Weak"
		} else if nScore >= 40 && nScore < 60 {
			sComplexity = "Good"
		} else if nScore >= 60 && nScore < 80 {
			sComplexity = "Strong"
		} else if nScore >= 80 && nScore <= 100 {
			sComplexity = "Very Strong"
		}
	}

	return sComplexity
}

func Reverse(s string) string {
	var reverse string
	for i := len(s) - 1; i >= 0; i-- {
		reverse += string(s[i])
	}
	return reverse
}
