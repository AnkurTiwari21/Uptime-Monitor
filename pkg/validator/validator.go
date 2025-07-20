package validator

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"unicode"

	"github.com/ankur12345678/uptime-monitor/pkg/logger"
	"github.com/nyaruka/phonenumbers"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	"golang.org/x/exp/slices"
)

func isValidPassword(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 7 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func isValidDOB(dob string) bool {
	layout := "01-02-2006" //  MM-DD-YYY

	// Parse the provided DOB string
	dobTime, err := time.Parse(layout, dob)
	if err != nil {
		return false // Invalid format
	}

	// Get the current date
	currentTime := time.Now()

	// Calculate the age by subtracting the DOB from the current date
	age := currentTime.Year() - dobTime.Year()

	// Adjust age if the birthday hasn't occurred yet this year
	if currentTime.YearDay() < dobTime.YearDay() {
		age--
	}

	// Check if the age is within a valid range (e.g., 18 to 100 years)
	return age >= 18 && age <= 100
}

func InitValidator() (*validator.Validate, ut.Translator, error) {
	translator := en.New()
	uni := ut.New(translator, translator)

	// this is usually known or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, found := uni.GetTranslator("en")
	if !found {
		logger.Fatal("translator not found")
	}

	v := validator.New()

	if err := enTranslations.RegisterDefaultTranslations(v, trans); err != nil {
		logger.Fatal(err)
	}

	err := v.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} is a required field", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphanumeric validation
	err = v.RegisterValidation("alphanumeric", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
		return re.MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphanumeric translation
	err = v.RegisterTranslation("alphanumeric", trans, func(ut ut.Translator) error {
		return ut.Add("alphanumeric", "{0} should be alphanumeric", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumeric", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphanumericifpresent validation
	err = v.RegisterValidation("alphanumericifpresent", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z0-9]*$`)
		if fl.Field().String() != "" {
			return re.MatchString(fl.Field().String())
		}
		return true
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphanumericifpresent translation
	err = v.RegisterTranslation("alphanumericifpresent", trans, func(ut ut.Translator) error {
		return ut.Add("alphanumericifpresent", "{0} should be alphanumeric", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumericifpresent", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphanumericwithhyphenunderscore validation
	err = v.RegisterValidation("alphanumericwithhyphenunderscoreifpresent", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z\d-_]+$`)
		if fl.Field().String() != "" {
			return re.MatchString(fl.Field().String())
		}
		return true
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphanumericwithhyphenunderscore translation
	err = v.RegisterTranslation("alphanumericwithhyphenunderscoreifpresent", trans, func(ut ut.Translator) error {
		return ut.Add("alphanumericwithhyphenunderscoreifpresent", "{0} should be alphanumericwithhyphenunderscore", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumericwithhyphenunderscoreifpresent", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphanumericwithhyphenunderscore validation
	err = v.RegisterValidation("alphanumericwithhyphenunderscore", func(fl validator.FieldLevel) bool {
		re := regexp.MustCompile(`^[a-zA-Z\d-_]+$`)
		return re.MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphanumericwithhyphenunderscore translation
	err = v.RegisterTranslation("alphanumericwithhyphenunderscore", trans, func(ut ut.Translator) error {
		return ut.Add("alphanumericwithhyphenunderscore", "{0} should be alphanumericwithhyphenunderscore", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumericwithhyphenunderscore", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// website validation
	err = v.RegisterValidation("website", func(fl validator.FieldLevel) bool {
		// Either with https site or non https site
		// https://uibakery.io/regex-library/url
		// return isUrl(fl.Field().String())
		return regexp.
			MustCompile(`^(https?://)?(www\.)?[a-zA-Z0-9.-]+\.[a-z]{2,}(:[0-9]{1,5})?(/[\w ./?%&=]*)?$`).
			MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// website translation
	err = v.RegisterTranslation("website", trans, func(ut ut.Translator) error {
		return ut.Add("website", "{0} is not in a valid website format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("website", fe.Field())
		return t
	})

	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphabets validation
	err = v.RegisterValidation("alphabet", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[a-zA-Z]+$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// alphabets translation
	err = v.RegisterTranslation("alphabet ", trans, func(ut ut.Translator) error {
		return ut.Add("alphabet", "{0} must have characters only from a-z,A-z", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphabet", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// notzerofloat validation
	err = v.RegisterValidation("notzerofloat", func(fl validator.FieldLevel) bool {
		val := fl.Field().Float()
		return val > 0
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// notzerofloat translation
	err = v.RegisterTranslation("notzerofloat ", trans, func(ut ut.Translator) error {
		return ut.Add("notzerofloat", "{0} must be greater than zero", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notzerofloat", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// notzerouint validation
	err = v.RegisterValidation("notzerouint", func(fl validator.FieldLevel) bool {
		val := fl.Field().Uint()
		return val > 0
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// notzerouint translation
	err = v.RegisterTranslation("notzerouint ", trans, func(ut ut.Translator) error {
		return ut.Add("notzerouint", "{0} must be greater than zero", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notzerouint", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// notzeroint validation
	err = v.RegisterValidation("notzeroint", func(fl validator.FieldLevel) bool {
		val := fl.Field().Int()
		return val > 0
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// notzeroint translation
	err = v.RegisterTranslation("notzeroint", trans, func(ut ut.Translator) error {
		return ut.Add("notzeroint", "{0} must be greater than zero", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("notzeroint", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// zipcode validation
	err = v.RegisterValidation("zipcode", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^\d{5}$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// zipcode translator
	err = v.RegisterTranslation("zipcode", trans, func(ut ut.Translator) error {
		return ut.Add("zipcode ", "{0} is invalid ", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("zipcode", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// websiteifpresent validation
	err = v.RegisterValidation("websiteifpresent", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			return regexp.
				MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`).
				MatchString(fl.Field().String())
		}
		return true
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// websiteifpresent translation
	err = v.RegisterTranslation("websiteifpresent", trans, func(ut ut.Translator) error {
		return ut.Add("websiteifpresent", "{0} is not in a valid website format", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("websiteifpresent", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// email if exists
	err = v.RegisterValidation("emailifpresent", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			return regexp.MustCompile(`[a-z0-9]+@[a-z]+\.[a-z]{2,3}`).
				MatchString(fl.Field().String())
		}
		return true
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// email translation
	err = v.RegisterTranslation("emailifpresent", trans, func(ut ut.Translator) error {
		return ut.Add("emailifpresent", "{0} is not a valid email", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("emailifpresent", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// phonenumber validation
	err = v.RegisterValidation("phonenumber", func(fl validator.FieldLevel) bool {

		num, err := phonenumbers.Parse(fl.Field().String(), "US")
		if err != nil {
			logger.Error("error in parsing phone number ", err)
			return false
		}
		return phonenumbers.IsValidNumber(num)
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// phonenumber translation
	err = v.RegisterTranslation("phonenumber", trans, func(ut ut.Translator) error {
		return ut.Add("phonenumber", "{0} is not a valid phonenumber", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phonenumber", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// phonenumber validation
	err = v.RegisterValidation("phonenumberifpresent", func(fl validator.FieldLevel) bool {

		num, err := phonenumbers.Parse(fl.Field().String(), "US")
		if err != nil {
			logger.Error("error in parsing phone number ", err)
			return false
		}
		return phonenumbers.IsValidNumber(num)

	})

	// userid validation
	err = v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		return regexp.
			MustCompile(`^(1\s?)?(\d{3}|\(\d{3}\))[\s\-]?\d{3}[\s\-]?\d{4}$`).
			MatchString(fl.Field().String()) || regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z0-9.-]{2,}$`).
			MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// phonenumber translation
	err = v.RegisterTranslation("phonenumberifpresent", trans, func(ut ut.Translator) error {
		return ut.Add("phonenumberifpresent", "{0} is not a valid phonenumber", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("phonenumberifpresent", fe.Field())
		return t
	})

	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}
	// userid translation
	err = v.RegisterTranslation("username", trans, func(ut ut.Translator) error {
		return ut.Add("username", "{0} is not a valid email/phonenumber", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("username", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// phonecountrycode validation
	err = v.RegisterValidation("phonecountrycode", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`\+1$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// phonecountrycode translation
	err = v.RegisterTranslation("phonecountrycode", trans, func(ut ut.Translator) error {
		return ut.Add("countrycode", "{0} is not supported", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("countrycode", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	// currency validation
	err = v.RegisterValidation("currency", func(fl validator.FieldLevel) bool {
		//TODO: will be fetched in future from db or some other source
		var currenciesAccepted = []string{
			"usd",
		}
		return slices.Contains(currenciesAccepted, fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// currency translation
	err = v.RegisterTranslation("currency", trans, func(ut ut.Translator) error {
		return ut.Add("currency", "{0} is not accepted", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("currency", fe.Value().(string))
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translation", err)
	}

	// numeric validation
	err = v.RegisterValidation("numeric", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[0-9]*$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// numeric translation
	err = v.RegisterTranslation("numeric", trans, func(ut ut.Translator) error {
		return ut.Add("numeric", "{0} is not valid numeric value", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("numeric", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translation", err)
	}

	// alphanumericwithspace
	err = v.RegisterValidation("alphanumericwithspace", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[A-Za-z0-9 ]+$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// alphanumericwithspace
	err = v.RegisterTranslation("alphanumericwithspace", trans, func(ut ut.Translator) error {
		return ut.Add("alphanumericwithspace", "{0} must only contain alphanumeric characters with spaces", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumericwithspace", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translation", err)
	}

	// alphawithspace
	err = v.RegisterValidation("alphawithspace", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^[A-Za-z ]+$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	// alphanumericwithspace
	err = v.RegisterTranslation("alphawithspace", trans, func(ut ut.Translator) error {
		return ut.Add("alphawithspace", "{0} must only contain alphabet characters with spaces", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphawithspace", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translation", err)
	}

	// alphanumericwithspaceifexists
	err = v.RegisterValidation("alphanumericwithspaceifexists", func(fl validator.FieldLevel) bool {
		if fl.Field().String() != "" {
			return regexp.MustCompile(`^[A-Za-z0-9 ]+$`).MatchString(fl.Field().String())
		} else {
			return true
		}
	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	err = v.RegisterTranslation("alphanumericwithspaceifexists", trans, func(ut ut.Translator) error {
		return ut.Add("alphanumericwithspaceifexists", "{0} must only contain alphanumeric characters with spaces", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("alphanumericwithspaceifexists", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translation", err)
	}

	// password validator
	err = v.RegisterValidation("password", func(fl validator.FieldLevel) bool {
		return isValidPassword(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("unable to register password validation ", err)
	}

	err = v.RegisterTranslation("password", trans, func(ut ut.Translator) error {
		return ut.Add("password", `{0} must have atleast 1 uppercase letter, 1 lowercase letter,1 number , 1 special character and 8 to 30 letter
		`, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("password", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register required translator ", err)
	}

	err = v.RegisterValidation("accountnumber", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`\W*\d{8,17}\b`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("unable to register account number validation ", err)
	}

	err = v.RegisterTranslation("accountnumber", trans, func(ut ut.Translator) error {
		return ut.Add("accountnumber", `Account number is invalid`, true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("accountnumber", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register required translator ", err)
	}

	err = v.RegisterValidation("tobetrue", func(fl validator.FieldLevel) bool {
		return fl.Field().Bool()
	})
	if err != nil {
		logger.Fatal("unable to register user validation ", err)
	}

	err = v.RegisterTranslation("tobetrue", trans, func(ut ut.Translator) error {
		return ut.Add("tobetrue", "{0} must be accepted", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("tobetrue", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register translation ", err)
	}

	err = v.RegisterValidation("dob", func(fl validator.FieldLevel) bool {
		return isValidDOB(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("unable to register validation for dob ", err)
	}

	err = v.RegisterTranslation("dob", trans, func(ut ut.Translator) error {
		return ut.Add("dob", "Invalid Birthday", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("dob", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register translation ", err)
	}

	err = v.RegisterValidation("dobifexists", func(fl validator.FieldLevel) bool {
		if fl.Field().String() == "" {
			return true
		}
		return isValidDOB(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("unable to register validation for dob ", err)
	}

	err = v.RegisterTranslation("dobifexists", trans, func(ut ut.Translator) error {
		return ut.Add("dobifexists", "Invalid Birthday", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("dobifexists", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register translation ", err)
	}

	err = v.RegisterValidation("date", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		layout := "01/02/2006"

		_, err := time.Parse(layout, val)

		return err == nil

	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	//date validation
	err = v.RegisterTranslation("date", trans, func(ut ut.Translator) error {
		return ut.Add("date", "{0} must be in format mm/dd/yyyy", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("date", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	err = v.RegisterValidation("future_date", func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		layout := "01/02/2006"

		parsedTime, _ := time.Parse(layout, val)

		currentTime := time.Now()

		isFuture := parsedTime.After(currentTime)

		return !isFuture

	})
	if err != nil {
		logger.Fatal("Unable to register required validator", err)
	}

	err = v.RegisterTranslation("future_date", trans, func(ut ut.Translator) error {
		return ut.Add("future_date", "{0} is of future date", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("future_date", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("Unable to register required translator", err)
	}

	err = v.RegisterValidation("onetimecode", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`/^[0-9]{6,6}$/g`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("unable to register validation ", err)
	}

	err = v.RegisterTranslation("onetimecode", trans, func(ut ut.Translator) error {
		return ut.Add("onetimecode", "Invalid code", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("onetimecode", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("error in validating one time code ", err)
	}

	// employeridentificationnumber
	err = v.RegisterValidation("employeridentificationnumber", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^\d{9}$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("unable to register validation ", err)
	}

	err = v.RegisterTranslation("employeridentificationnumber", trans, func(ut ut.Translator) error {
		return ut.Add("employeridentificationnumber", "{0} is not a valid employer identification number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("employeridentificationnumber", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register translation ", err)
	}

	// routingnumber
	err = v.RegisterValidation("routingnumber", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^\d{9}$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("unable to register validation ", err)
	}

	err = v.RegisterTranslation("routingnumber", trans, func(ut ut.Translator) error {
		return ut.Add("routingnumber", "{0} is not a valid routing number", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("routingnumber", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register translation ", err)
	}

	err = v.RegisterValidation("ssn", func(fl validator.FieldLevel) bool {
		return regexp.MustCompile(`^\d{3}[- ]?\d{2}[- ]?\d{4}$`).MatchString(fl.Field().String())
	})
	if err != nil {
		logger.Fatal("unable to register validation ", err)
	}

	err = v.RegisterTranslation("ssn", trans, func(ut ut.Translator) error {
		return ut.Add("ssn", "Invalid ssn", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("ssn", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register validation ", err)
	}

	err = v.RegisterValidation("lessthan150", func(fl validator.FieldLevel) bool {
		return len(strings.TrimSpace(fl.Field().String())) < 150
	})
	if err != nil {
		logger.Fatal("unable to register validation ", err)
	}

	err = v.RegisterTranslation("lessthan150", trans, func(ut ut.Translator) error {
		return ut.Add("lessthan150", "{0} must be less than 150 characters", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("lessthan150", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register validation ", err)
	}

	err = v.RegisterValidation("startWithAlphanumericCharacter", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()
		if str == "" {
			return true
		}
		return unicode.IsLetter(rune(str[0])) || unicode.IsDigit(rune(str[0]))
	})
	if err != nil {
		logger.Fatal("unable to register validation ", err)
	}

	err = v.RegisterTranslation("startWithAlphanumericCharacter", trans, func(ut ut.Translator) error {
		return ut.Add("startWithAlphanumericCharacter", "{0} should start with alpha numeric", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("startWithAlphanumericCharacter", fe.Field())
		return t
	})
	if err != nil {
		logger.Fatal("unable to register validation ", err)
	}

	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})

	return v, trans, nil
}
