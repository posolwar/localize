package localize

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

type Translater struct {
	// Хранит в себе наборы языков и методы для трансляции.
	// Потокобезопасен для горутин (https://github.com/nicksnyder/go-i18n/issues/188).
	*i18n.Localizer
}

var (
	Localizer *Translater

	// Язык приложения по умолчанию
	defaultLanguage = language.Russian

	// Формат файлов, для файлов с переводом по умолчанию
	defaultFileFormat = "json"

	// Метод унмаршала файлов с переводом по умолчанию
	defaultUnmarshalFunc i18n.UnmarshalFunc = json.Unmarshal

	// Место хранение файлов с переводами по умолчанию
	defaultFilesPath = "resources/"
)

// Установить другой язык по умолчанию для приложения
func SetDefaultLanguage(selectLanguage language.Tag) {
	defaultLanguage = selectLanguage
}

// Установить формат файлов, используемых для хранения языка по умолчанию
func SetDefaultFileFormat(format string) {
	defaultFileFormat = format
}

// Установить метод унмаршала файлов по умолчанию
func SetDefaultFileUnmarshal(unmarshalFunc i18n.UnmarshalFunc) {
	defaultUnmarshalFunc = unmarshalFunc
}

// Установить путь к файлам с переводами по умолчанию
func SetDefaultFilesPath(path string) {
	defaultFilesPath = path
}

// Создаем Translater (обертка над Localizer), имеющий методы перевода слов
func NewTranslater(selectLanguage language.Tag) (err error) {
	// Создает новый bundle (хранилище набора сообщений и правила плюрализации) с дефолтным языком
	bundle := i18n.NewBundle(defaultLanguage)

	// Заполняем bundle словами из дефолтного и указанного языка(en.json, ru.json...)
	err = getFileWords(bundle, selectLanguage)
	if err != nil {
		return err
	}

	// Получем Localizer, который ищет сообщения в "bundle", в соответствии с языковыми предпочтениями в "selectLanguage".
	// Если слово не было найдено в указанном selectLanguage, то ищет его в defaultLanguage.
	localizer := i18n.NewLocalizer(bundle, selectLanguage.String(), defaultLanguage.String())

	Localizer = &Translater{localizer}

	return nil
}

// Загрузить файлы перевода в объект переводчика
func getFileWords(bundle *i18n.Bundle, selectLanguage language.Tag) error {
	// Выбираем формат и функцию унмаршала файла
	bundle.RegisterUnmarshalFunc(defaultFileFormat, defaultUnmarshalFunc)

	// Загружаем файл с выбранным языком
	_, err := bundle.LoadMessageFile(defaultFilesPath + selectLanguage.String() + "." + defaultFileFormat)
	if err != nil {
		return err
	}

	// Проверка на то, был ли уже добавлен язык по умолчанию
	if selectLanguage == defaultLanguage {
		return nil
	}

	// Загружаем файл с языком по умолчанию
	_, err = bundle.LoadMessageFile(defaultFilesPath + defaultLanguage.String() + "." + defaultFileFormat)

	return err
}

// Перевод слова без возврата ошибки. Панике при ошибке.
func (tr *Translater) MustWordTranslate(word string) string {
	localizeConf := i18n.LocalizeConfig{MessageID: word}

	return tr.Localizer.MustLocalize(&localizeConf)
}

// Перевод слова с возвратом ошибки.
func (tr *Translater) WordTranslate(word string) (string, error) {
	localizeConf := i18n.LocalizeConfig{MessageID: word}

	return tr.Localizer.Localize(&localizeConf)
}
