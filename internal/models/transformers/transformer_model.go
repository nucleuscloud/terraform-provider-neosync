package transformer_model

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/types"
	mgmtv1alpha1 "github.com/nucleuscloud/neosync/backend/gen/go/protos/mgmt/v1alpha1"
)

type Transformer struct {
	Config *TransformerConfig `tfsdk:"config"`
}

type TransformerConfig struct {
	GenerateEmail              *TransformerEmpty          `tfsdk:"generate_email"`
	TransformEmail             *TransformEmail            `tfsdk:"transform_email"`
	GenerateBool               *TransformerEmpty          `tfsdk:"generate_bool"`
	GenerateCardNumber         *GenerateCardNumber        `tfsdk:"generate_card_number"`
	GenerateCity               *TransformerEmpty          `tfsdk:"generate_city"`
	GenerateE164PhoneNumber    *GenerateE164PhoneNumber   `tfsdk:"generate_e164_phone_number"`
	GenerateFirstName          *TransformerEmpty          `tfsdk:"generate_firstname"`
	GenerateFloat64            *GenerateFloat64           `tfsdk:"generate_float64"`
	GenerateFullAddress        *TransformerEmpty          `tfsdk:"generate_full_address"`
	GenerateFullName           *TransformerEmpty          `tfsdk:"generate_fullname"`
	GenerateGender             *GenerateGender            `tfsdk:"generate_gender"`
	GenerateInt64PhoneNumber   *TransformerEmpty          `tfsdk:"generate_int64_phone_number"`
	GenerateInt64              *GenerateInt64             `tfsdk:"generate_int64"`
	GenerateLastName           *TransformerEmpty          `tfsdk:"generate_lastname"`
	GenerateSha256Hash         *TransformerEmpty          `tfsdk:"generate_sha256"`
	GenerateSsn                *TransformerEmpty          `tfsdk:"generate_ssn"`
	GenerateState              *TransformerEmpty          `tfsdk:"generate_state"`
	GenerateStreetAddress      *TransformerEmpty          `tfsdk:"generate_street_address"`
	GenerateStringPhoneNumber  *GenerateStringPhoneNumber `tfsdk:"generate_string_phone_number"`
	GenerateString             *GenerateString            `tfsdk:"generate_string"`
	GenerateUnixtimestamp      *TransformerEmpty          `tfsdk:"generate_unix_timestamp"`
	GenerateUsername           *TransformerEmpty          `tfsdk:"generate_username"`
	GenerateUtctimestamp       *TransformerEmpty          `tfsdk:"generate_utc_timestamp"`
	GenerateUuid               *GenerateUuid              `tfsdk:"generate_uuid"`
	GenerateZipcode            *TransformerEmpty          `tfsdk:"generate_zipcode"`
	TransformE164PhoneNumber   *TransformE164PhoneNumber  `tfsdk:"transform_e164_phone_number"`
	TransformFirstName         *TransformFirstName        `tfsdk:"transform_firstname"`
	TransformFloat64           *TransformFloat64          `tfsdk:"transform_float64"`
	TransformFullName          *TransformFullName         `tfsdk:"transform_fullname"`
	TransformInt64PhoneNumber  *TransformInt64PhoneNumber `tfsdk:"transform_int64_phone_number"`
	TransformInt64             *TransformInt64            `tfsdk:"transform_int64"`
	TransformLastName          *TransformLastName         `tfsdk:"transform_lastname"`
	TransformPhoneNumber       *TransformPhoneNumber      `tfsdk:"transform_phone_number"`
	TransformString            *TransformString           `tfsdk:"transform_string"`
	Passthrough                *TransformerEmpty          `tfsdk:"passthrough"`
	Null                       *TransformerEmpty          `tfsdk:"null"`
	UserDefinedTransformer     *UserDefinedTransformer    `tfsdk:"user_defined_transformer"`
	GenerateDefault            *TransformerEmpty          `tfsdk:"generate_default"`
	TransformJavascript        *TransformJavascript       `tfsdk:"transform_javascript"`
	GenerateJavascript         *GenerateJavascript        `tfsdk:"generate_javascript"`
	GenerateCategorical        *GenerateCategorical       `tfsdk:"generate_categorical"`
	TransformCharacterScramble *TransformerEmpty          `tfsdk:"transform_character_scramble"`
}
type TransformerEmpty struct{}
type TransformEmail struct {
	PreserveDomain types.Bool `tfsdk:"preserve_domain"`
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type GenerateCardNumber struct {
	ValidLuhn types.Bool `tfsdk:"valid_luhn"`
}
type GenerateE164PhoneNumber struct {
	Min types.Int64 `tfsdk:"min"`
	Max types.Int64 `tfsdk:"max"`
}
type GenerateFloat64 struct {
	RandomizeSign types.Bool    `tfsdk:"randomize_sign"`
	Min           types.Float64 `tfsdk:"min"`
	Max           types.Float64 `tfsdk:"max"`
	Precision     types.Int64   `tfsdk:"precision"`
}
type GenerateGender struct {
	Abbreviate types.Bool `tfsdk:"abbreviate"`
}
type GenerateInt64 struct {
	RandomizeSign types.Bool  `tfsdk:"randomize_sign"`
	Min           types.Int64 `tfsdk:"min"`
	Max           types.Int64 `tfsdk:"max"`
}
type GenerateStringPhoneNumber struct {
	Min types.Int64 `tfsdk:"min"`
	Max types.Int64 `tfsdk:"max"`
}
type GenerateString struct {
	Min types.Int64 `tfsdk:"min"`
	Max types.Int64 `tfsdk:"max"`
}
type GenerateUuid struct {
	IncludeHyphens types.Bool `tfsdk:"include_hyphens"`
}
type TransformE164PhoneNumber struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformFirstName struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformFloat64 struct {
	RandomizationRangeMin types.Float64 `tfsdk:"randomization_range_min"`
	RandomizationRangeMax types.Float64 `tfsdk:"randomization_range_max"`
}
type TransformFullName struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformInt64PhoneNumber struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformInt64 struct {
	RandomizationRangeMin types.Int64 `tfsdk:"randomization_range_min"`
	RandomizationRangeMax types.Int64 `tfsdk:"randomization_range_max"`
}
type TransformLastName struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformPhoneNumber struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformString struct {
	PreserveLength types.Bool `tfsdk:"preserve_length"`
}
type TransformJavascript struct {
	Code types.String `tfsdk:"code"`
}
type GenerateJavascript struct {
	Code types.String `tfsdk:"code"`
}
type UserDefinedTransformer struct {
	Id types.String `tfsdk:"id"`
}
type GenerateCategorical struct {
	Categories types.String `tfsdk:"categories"`
}

func (t *Transformer) FromDto(dto *mgmtv1alpha1.TransformerConfig) error {
	if t == nil {
		return errors.New("transformer is nil")
	}

	if dto == nil {
		return errors.New("dto is nil")
	}

	t.Config = &TransformerConfig{}

	err := t.Config.FromDto(dto)
	if err != nil {
		return err
	}

	return nil
}

func (tc *TransformerConfig) FromDto(dto *mgmtv1alpha1.TransformerConfig) error {
	if tc == nil {
		return errors.New("transformer config is nil")
	}

	if dto == nil {
		return errors.New("dto is nil")
	}

	switch config := dto.GetConfig().(type) {
	case *mgmtv1alpha1.TransformerConfig_GenerateEmailConfig:
		tc.GenerateEmail = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_TransformEmailConfig:
		tc.TransformEmail = &TransformEmail{
			PreserveDomain: types.BoolPointerValue(config.TransformEmailConfig.PreserveDomain),
			PreserveLength: types.BoolPointerValue(config.TransformEmailConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateBoolConfig:
		tc.GenerateBool = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateCardNumberConfig:
		tc.GenerateCardNumber = &GenerateCardNumber{
			ValidLuhn: types.BoolPointerValue(config.GenerateCardNumberConfig.ValidLuhn),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateCityConfig:
		tc.GenerateCity = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateE164PhoneNumberConfig:
		tc.GenerateE164PhoneNumber = &GenerateE164PhoneNumber{
			Min: types.Int64PointerValue(config.GenerateE164PhoneNumberConfig.Min),
			Max: types.Int64PointerValue(config.GenerateE164PhoneNumberConfig.Max),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateFirstNameConfig:
		tc.GenerateFirstName = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateFloat64Config:
		tc.GenerateFloat64 = &GenerateFloat64{
			RandomizeSign: types.BoolPointerValue(config.GenerateFloat64Config.RandomizeSign),
			Min:           types.Float64PointerValue(config.GenerateFloat64Config.Min),
			Max:           types.Float64PointerValue(config.GenerateFloat64Config.Max),
			Precision:     types.Int64PointerValue(config.GenerateFloat64Config.Precision),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateFullAddressConfig:
		tc.GenerateFullAddress = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateFullNameConfig:
		tc.GenerateFullName = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateGenderConfig:
		tc.GenerateGender = &GenerateGender{
			Abbreviate: types.BoolPointerValue(config.GenerateGenderConfig.Abbreviate),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateInt64PhoneNumberConfig:
		tc.GenerateInt64PhoneNumber = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateInt64Config:
		tc.GenerateInt64 = &GenerateInt64{
			RandomizeSign: types.BoolPointerValue(config.GenerateInt64Config.RandomizeSign),
			Min:           types.Int64PointerValue(config.GenerateInt64Config.Min),
			Max:           types.Int64PointerValue(config.GenerateInt64Config.Max),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateLastNameConfig:
		tc.GenerateLastName = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateSha256HashConfig:
		tc.GenerateSha256Hash = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateSsnConfig:
		tc.GenerateSsn = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateStateConfig:
		tc.GenerateState = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateStreetAddressConfig:
		tc.GenerateStreetAddress = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateStringPhoneNumberConfig:
		tc.GenerateStringPhoneNumber = &GenerateStringPhoneNumber{
			Min: types.Int64PointerValue(config.GenerateStringPhoneNumberConfig.Min),
			Max: types.Int64PointerValue(config.GenerateStringPhoneNumberConfig.Max),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateStringConfig:
		tc.GenerateString = &GenerateString{
			Min: types.Int64PointerValue(config.GenerateStringConfig.Min),
			Max: types.Int64PointerValue(config.GenerateStringConfig.Max),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateUnixtimestampConfig:
		tc.GenerateUnixtimestamp = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateUsernameConfig:
		tc.GenerateUsername = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateUtctimestampConfig:
		tc.GenerateUtctimestamp = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateUuidConfig:
		tc.GenerateUuid = &GenerateUuid{
			IncludeHyphens: types.BoolPointerValue(config.GenerateUuidConfig.IncludeHyphens),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateZipcodeConfig:
		tc.GenerateZipcode = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_TransformE164PhoneNumberConfig:
		tc.TransformE164PhoneNumber = &TransformE164PhoneNumber{
			PreserveLength: types.BoolPointerValue(config.TransformE164PhoneNumberConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformFirstNameConfig:
		tc.TransformFirstName = &TransformFirstName{
			PreserveLength: types.BoolPointerValue(config.TransformFirstNameConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformFloat64Config:
		tc.TransformFloat64 = &TransformFloat64{
			RandomizationRangeMin: types.Float64PointerValue(config.TransformFloat64Config.RandomizationRangeMin),
			RandomizationRangeMax: types.Float64PointerValue(config.TransformFloat64Config.RandomizationRangeMax),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformFullNameConfig:
		tc.TransformFullName = &TransformFullName{
			PreserveLength: types.BoolPointerValue(config.TransformFullNameConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformInt64PhoneNumberConfig:
		tc.TransformInt64PhoneNumber = &TransformInt64PhoneNumber{
			PreserveLength: types.BoolPointerValue(config.TransformInt64PhoneNumberConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformInt64Config:
		tc.TransformInt64 = &TransformInt64{
			RandomizationRangeMin: types.Int64PointerValue(config.TransformInt64Config.RandomizationRangeMin),
			RandomizationRangeMax: types.Int64PointerValue(config.TransformInt64Config.RandomizationRangeMax),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformLastNameConfig:
		tc.TransformLastName = &TransformLastName{
			PreserveLength: types.BoolPointerValue(config.TransformLastNameConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformPhoneNumberConfig:
		tc.TransformPhoneNumber = &TransformPhoneNumber{
			PreserveLength: types.BoolPointerValue(config.TransformPhoneNumberConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformStringConfig:
		tc.TransformString = &TransformString{
			PreserveLength: types.BoolPointerValue(config.TransformStringConfig.PreserveLength),
		}
	case *mgmtv1alpha1.TransformerConfig_PassthroughConfig:
		tc.Passthrough = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_Nullconfig:
		tc.Null = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_UserDefinedTransformerConfig:
		tc.UserDefinedTransformer = &UserDefinedTransformer{
			Id: types.StringValue(config.UserDefinedTransformerConfig.Id),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateDefaultConfig:
		tc.GenerateDefault = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_TransformJavascriptConfig:
		tc.TransformJavascript = &TransformJavascript{
			Code: types.StringValue(config.TransformJavascriptConfig.Code),
		}
	case *mgmtv1alpha1.TransformerConfig_GenerateCategoricalConfig:
		tc.GenerateCategorical = &GenerateCategorical{
			Categories: types.StringPointerValue(config.GenerateCategoricalConfig.Categories),
		}
	case *mgmtv1alpha1.TransformerConfig_TransformCharacterScrambleConfig:
		tc.TransformCharacterScramble = &TransformerEmpty{}
	case *mgmtv1alpha1.TransformerConfig_GenerateJavascriptConfig:
		tc.GenerateJavascript = &GenerateJavascript{
			Code: types.StringValue(config.GenerateJavascriptConfig.Code),
		}
	default:
		return fmt.Errorf("this job mapping transformer is not currently supported by this provider: %w", errors.ErrUnsupported)
	}
	return nil
}

func (t *Transformer) ToDto() (*mgmtv1alpha1.TransformerConfig, error) {
	if t == nil {
		return nil, errors.New("transformer is nil")
	}

	if t.Config == nil {
		return nil, errors.New("transformer config is nil")
	}

	return t.Config.ToDto()
}

func (tc *TransformerConfig) ToDto() (*mgmtv1alpha1.TransformerConfig, error) {
	if tc == nil {
		return nil, errors.New("transformer config is nil")
	}

	dto := &mgmtv1alpha1.TransformerConfig{}

	if tc.GenerateEmail != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateEmailConfig{
			GenerateEmailConfig: &mgmtv1alpha1.GenerateEmail{},
		}
	} else if tc.TransformEmail != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformEmailConfig{
			TransformEmailConfig: &mgmtv1alpha1.TransformEmail{
				PreserveDomain: tc.TransformEmail.PreserveDomain.ValueBoolPointer(),
				PreserveLength: tc.TransformEmail.PreserveLength.ValueBoolPointer(),
			},
		}
	} else if tc.GenerateBool != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateBoolConfig{
			GenerateBoolConfig: &mgmtv1alpha1.GenerateBool{},
		}
	} else if tc.GenerateCardNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateCardNumberConfig{
			GenerateCardNumberConfig: &mgmtv1alpha1.GenerateCardNumber{
				ValidLuhn: tc.GenerateCardNumber.ValidLuhn.ValueBoolPointer(),
			},
		}
	} else if tc.GenerateCity != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateCityConfig{
			GenerateCityConfig: &mgmtv1alpha1.GenerateCity{},
		}
	} else if tc.GenerateE164PhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateE164PhoneNumberConfig{
			GenerateE164PhoneNumberConfig: &mgmtv1alpha1.GenerateE164PhoneNumber{
				Min: tc.GenerateE164PhoneNumber.Min.ValueInt64Pointer(),
				Max: tc.GenerateE164PhoneNumber.Max.ValueInt64Pointer(),
			},
		}
	} else if tc.GenerateFirstName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFirstNameConfig{
			GenerateFirstNameConfig: &mgmtv1alpha1.GenerateFirstName{},
		}
	} else if tc.GenerateFloat64 != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFloat64Config{
			GenerateFloat64Config: &mgmtv1alpha1.GenerateFloat64{
				RandomizeSign: tc.GenerateFloat64.RandomizeSign.ValueBoolPointer(),
				Min:           tc.GenerateFloat64.Min.ValueFloat64Pointer(),
				Max:           tc.GenerateFloat64.Max.ValueFloat64Pointer(),
			},
		}
	} else if tc.GenerateFullAddress != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFullAddressConfig{
			GenerateFullAddressConfig: &mgmtv1alpha1.GenerateFullAddress{},
		}
	} else if tc.GenerateFullName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFullNameConfig{
			GenerateFullNameConfig: &mgmtv1alpha1.GenerateFullName{},
		}
	} else if tc.GenerateGender != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateGenderConfig{
			GenerateGenderConfig: &mgmtv1alpha1.GenerateGender{
				Abbreviate: tc.GenerateGender.Abbreviate.ValueBoolPointer(),
			},
		}
	} else if tc.GenerateInt64PhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateInt64PhoneNumberConfig{
			GenerateInt64PhoneNumberConfig: &mgmtv1alpha1.GenerateInt64PhoneNumber{},
		}
	} else if tc.GenerateInt64 != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateInt64Config{
			GenerateInt64Config: &mgmtv1alpha1.GenerateInt64{
				RandomizeSign: tc.GenerateInt64.RandomizeSign.ValueBoolPointer(),
				Min:           tc.GenerateInt64.Min.ValueInt64Pointer(),
				Max:           tc.GenerateInt64.Max.ValueInt64Pointer(),
			},
		}
	} else if tc.GenerateLastName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateLastNameConfig{
			GenerateLastNameConfig: &mgmtv1alpha1.GenerateLastName{},
		}
	} else if tc.GenerateSha256Hash != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateSha256HashConfig{
			GenerateSha256HashConfig: &mgmtv1alpha1.GenerateSha256Hash{},
		}
	} else if tc.GenerateSsn != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateSsnConfig{
			GenerateSsnConfig: &mgmtv1alpha1.GenerateSSN{},
		}
	} else if tc.GenerateState != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateStateConfig{
			GenerateStateConfig: &mgmtv1alpha1.GenerateState{},
		}
	} else if tc.GenerateStreetAddress != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateStreetAddressConfig{
			GenerateStreetAddressConfig: &mgmtv1alpha1.GenerateStreetAddress{},
		}
	} else if tc.GenerateStringPhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateStringPhoneNumberConfig{
			GenerateStringPhoneNumberConfig: &mgmtv1alpha1.GenerateStringPhoneNumber{
				Min: tc.GenerateStringPhoneNumber.Min.ValueInt64Pointer(),
				Max: tc.GenerateStringPhoneNumber.Max.ValueInt64Pointer(),
			},
		}
	} else if tc.GenerateString != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateStringConfig{
			GenerateStringConfig: &mgmtv1alpha1.GenerateString{
				Min: tc.GenerateString.Min.ValueInt64Pointer(),
				Max: tc.GenerateString.Max.ValueInt64Pointer(),
			},
		}
	} else if tc.GenerateUnixtimestamp != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateUnixtimestampConfig{
			GenerateUnixtimestampConfig: &mgmtv1alpha1.GenerateUnixTimestamp{},
		}
	} else if tc.GenerateUsername != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateUsernameConfig{
			GenerateUsernameConfig: &mgmtv1alpha1.GenerateUsername{},
		}
	} else if tc.GenerateUtctimestamp != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateUtctimestampConfig{
			GenerateUtctimestampConfig: &mgmtv1alpha1.GenerateUtcTimestamp{},
		}
	} else if tc.GenerateUuid != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateUuidConfig{
			GenerateUuidConfig: &mgmtv1alpha1.GenerateUuid{
				IncludeHyphens: tc.GenerateUuid.IncludeHyphens.ValueBoolPointer(),
			},
		}
	} else if tc.GenerateZipcode != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateZipcodeConfig{
			GenerateZipcodeConfig: &mgmtv1alpha1.GenerateZipcode{},
		}
	} else if tc.TransformE164PhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateE164PhoneNumberConfig{
			GenerateE164PhoneNumberConfig: &mgmtv1alpha1.GenerateE164PhoneNumber{
				Min: tc.GenerateE164PhoneNumber.Min.ValueInt64Pointer(),
				Max: tc.GenerateE164PhoneNumber.Max.ValueInt64Pointer(),
			},
		}
	} else if tc.TransformFirstName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateFirstNameConfig{
			GenerateFirstNameConfig: &mgmtv1alpha1.GenerateFirstName{},
		}
	} else if tc.TransformFloat64 != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformFloat64Config{
			TransformFloat64Config: &mgmtv1alpha1.TransformFloat64{
				RandomizationRangeMin: tc.TransformFloat64.RandomizationRangeMin.ValueFloat64Pointer(),
				RandomizationRangeMax: tc.TransformFloat64.RandomizationRangeMax.ValueFloat64Pointer(),
			},
		}
	} else if tc.TransformFullName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformFullNameConfig{
			TransformFullNameConfig: &mgmtv1alpha1.TransformFullName{
				PreserveLength: tc.TransformFullName.PreserveLength.ValueBoolPointer(),
			},
		}
	} else if tc.TransformInt64PhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformInt64PhoneNumberConfig{
			TransformInt64PhoneNumberConfig: &mgmtv1alpha1.TransformInt64PhoneNumber{
				PreserveLength: tc.TransformInt64PhoneNumber.PreserveLength.ValueBoolPointer(),
			},
		}
	} else if tc.TransformInt64 != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateInt64Config{
			GenerateInt64Config: &mgmtv1alpha1.GenerateInt64{
				RandomizeSign: tc.GenerateInt64.RandomizeSign.ValueBoolPointer(),
				Min:           tc.GenerateInt64.Min.ValueInt64Pointer(),
				Max:           tc.GenerateInt64.Max.ValueInt64Pointer(),
			},
		}
	} else if tc.TransformLastName != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformLastNameConfig{
			TransformLastNameConfig: &mgmtv1alpha1.TransformLastName{
				PreserveLength: tc.TransformLastName.PreserveLength.ValueBoolPointer(),
			},
		}
	} else if tc.TransformPhoneNumber != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformPhoneNumberConfig{
			TransformPhoneNumberConfig: &mgmtv1alpha1.TransformPhoneNumber{
				PreserveLength: tc.TransformPhoneNumber.PreserveLength.ValueBoolPointer(),
			},
		}
	} else if tc.TransformString != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformStringConfig{
			TransformStringConfig: &mgmtv1alpha1.TransformString{
				PreserveLength: tc.TransformString.PreserveLength.ValueBoolPointer(),
			},
		}
	} else if tc.Passthrough != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_PassthroughConfig{
			PassthroughConfig: &mgmtv1alpha1.Passthrough{},
		}
	} else if tc.Null != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_Nullconfig{
			Nullconfig: &mgmtv1alpha1.Null{},
		}
	} else if tc.UserDefinedTransformer != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_UserDefinedTransformerConfig{
			UserDefinedTransformerConfig: &mgmtv1alpha1.UserDefinedTransformerConfig{
				Id: tc.UserDefinedTransformer.Id.ValueString(),
			},
		}
	} else if tc.GenerateDefault != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateDefaultConfig{
			GenerateDefaultConfig: &mgmtv1alpha1.GenerateDefault{},
		}
	} else if tc.TransformJavascript != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformJavascriptConfig{
			TransformJavascriptConfig: &mgmtv1alpha1.TransformJavascript{
				Code: tc.TransformJavascript.Code.ValueString(),
			},
		}
	} else if tc.GenerateCategorical != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateCategoricalConfig{
			GenerateCategoricalConfig: &mgmtv1alpha1.GenerateCategorical{
				Categories: tc.GenerateCategorical.Categories.ValueStringPointer(),
			},
		}
	} else if tc.TransformCharacterScramble != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_TransformCharacterScrambleConfig{
			TransformCharacterScrambleConfig: &mgmtv1alpha1.TransformCharacterScramble{},
		}
	} else if tc.GenerateJavascript != nil {
		dto.Config = &mgmtv1alpha1.TransformerConfig_GenerateJavascriptConfig{
			GenerateJavascriptConfig: &mgmtv1alpha1.GenerateJavascript{
				Code: tc.GenerateJavascript.Code.ValueString(),
			},
		}
	} else {
		return nil, fmt.Errorf("transformer config is not currently supported by this provider: %w", errors.ErrUnsupported)
	}

	return dto, nil
}
