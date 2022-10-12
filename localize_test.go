package localize

import (
	"encoding/json"
	"testing"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"
)

func TestDefaultValues(t *testing.T) {
	t.Run("Проверка работы установки значений по умолчанию", func(t *testing.T) {
		newDefaultLanguage := language.English
		SetDefaultLanguage(newDefaultLanguage)
		assert.Equal(
			t,
			defaultLanguage,
			newDefaultLanguage,
			"Не произошло изменение языка по умолчнию",
		)

		newDefaultFileFormat := "yaml"
		SetDefaultFileFormat(newDefaultFileFormat)
		assert.Equal(
			t,
			defaultFileFormat,
			newDefaultFileFormat,
			"Не произошло изменение формата файлов по умолчанию",
		)

		newDefaultFilesPath := "./test"
		SetDefaultFilesPath(newDefaultFilesPath)
		assert.Equal(
			t,
			defaultFilesPath,
			newDefaultFilesPath,
			"Не произошло изменение пути к файлам с переводом по умолчанию",
		)
	})
}

func TestTranslate(t *testing.T) {
	t.Run("Выполнение перевода", func(t *testing.T) {
		SetDefaultFilesPath("./files/")
		SetDefaultLanguage(language.Russian)
		SetDefaultFileFormat("json")
		SetDefaultFileUnmarshal(json.Unmarshal)

		err := NewTranslater(language.English)
		assert.NoError(t, err)

		word, err := Localizer.Localize(&i18n.LocalizeConfig{MessageID: "login"})
		assert.NoError(t, err, "Ожидается что слово login обработается, т.к. присуствует в en файле")

		word, err = Localizer.WordTranslate("password")
		assert.Error(t, err, "Ожидается ошибка что password не найден в en раскладке")
		assert.NotEmpty(t, word, "Несмотря на то, что password не найден в en раскладке, он должен быть найден в ru раскладке")

		word = Localizer.MustWordTranslate("login")
		assert.NotEmpty(t, word, "Ожидается что слово login обработается, т.к. присуствует в en файле")

	})

}
