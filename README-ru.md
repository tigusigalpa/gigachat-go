# 🚀 GigaChat Go SDK

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigusigalpa/gigachat-go?style=flat-square)](https://goreportcard.com/report/github.com/tigusigalpa/gigachat-go)

Полнофункциональный Go SDK для работы с Sber GigaChat API. Простой, мощный и идиоматичный Go-клиент для работы с AI моделями GigaChat, включая поддержку streaming и генерации изображений.

**🌐 Язык:** [English](README.md) | Русский

## ✨ Возможности

- 🔌 **Простая интеграция** с GigaChat API
- 🔐 **Автоматическое управление OAuth и токенами** с потокобезопасным кешированием
- 🎯 **Поддержка всех моделей** (GigaChat-2, GigaChat-2-Pro, GigaChat-2-Max)
- 📝 **Поддержка диалогов** с управлением контекстом
- ⚡ **Streaming поддержка** для ответов в реальном времени
- 🎨 **Генерация изображений** с помощью функции text2image
- 🖼️ **Автоматическое скачивание изображений** и кодирование в base64
- 🎭 **Стилизация изображений** через системные промпты
- 🔧 **Паттерн функциональных опций** для чистого API
- 🧪 **Полные примеры** для всех функций
- 📚 **Полная документация** с примерами кода
- 🛡️ **Типобезопасность** с правильной обработкой ошибок

## 📦 Установка

```bash
go get github.com/tigusigalpa/gigachat-go
```

## ⚙️ Настройка

### 1. Получение авторизационных данных

Для работы с GigaChat API необходимо получить авторизационные данные:

1. Зарегистрируйтесь в [личном кабинете Sber AI](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)
2. Создайте проект и получите **Client ID** и **Client Secret**
3. Сгенерируйте **Authorization Key** (Base64 от "Client ID:Client Secret")

> 💡 **Подробная инструкция**: [Создание проекта и получение ключей](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)

### 2. Настройка окружения

Настройте переменные окружения:

```bash
# Способ 1: Используя готовый Authorization Key
export GIGACHAT_AUTH_KEY=your_base64_encoded_auth_key

# Способ 2: Используя Client ID и Client Secret (автоматически сгенерирует auth_key)
export GIGACHAT_CLIENT_ID=your_client_id
export GIGACHAT_CLIENT_SECRET=your_client_secret

# Дополнительные настройки
export GIGACHAT_SCOPE=GIGACHAT_API_PERS
```

## 💡 Использование

### Базовое использование

```go
package main

import (
    "encoding/base64"
    "fmt"
    "log"
    "os"

    gigachat "github.com/tigusigalpa/gigachat-go"
)

func main() {
    // Создание ключа авторизации из учетных данных
    authKey := base64.StdEncoding.EncodeToString(
        []byte(os.Getenv("GIGACHAT_CLIENT_ID") + ":" + os.Getenv("GIGACHAT_CLIENT_SECRET")),
    )

    // Или использование готового ключа
    // authKey := os.Getenv("GIGACHAT_AUTH_KEY")

    // Создание менеджера токенов
    tokenManager := gigachat.NewTokenManager(authKey)

    // Создание клиента
    client := gigachat.NewClient(tokenManager)

    // Получение доступных моделей
    models, err := client.Models()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("Доступные модели:")
    for _, model := range models.Data {
        fmt.Printf("- %s\n", model.ID)
    }

    // Отправка сообщения
    messages := []gigachat.Message{
        {Role: "user", Content: "Привет! Расскажи анекдот"},
    }

    response, err := client.Chat(messages)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(response.Choices[0].Message.Content)
}
```

### Работа с диалогами

```go
// Создание диалога с системным промптом
conversation := gigachat.Conversation(
    "Ты полезный помощник программиста на Go",
    "Как создать REST API в Go?",
)

response, err := client.Chat(conversation)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Ассистент:", gigachat.ExtractContent(response))

// Продолжение диалога
conversation = append(conversation, gigachat.Message{
    Role:    "assistant",
    Content: response.Choices[0].Message.Content,
})

conversation = gigachat.ContinueChat(conversation, "А как добавить middleware для аутентификации?")

response, err = client.Chat(conversation)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Ассистент:", gigachat.ExtractContent(response))
```

### Streaming запросы

```go
messages := []gigachat.Message{
    {Role: "user", Content: "Напиши длинную историю о космосе"},
}

fmt.Println("Потоковый ответ:")

err := client.ChatStream(messages, func(event *gigachat.ChatResponse, done bool, err error) {
    if err != nil {
        log.Printf("Ошибка: %v", err)
        return
    }

    if done {
        fmt.Println("\n✅ Готово!")
        return
    }

    if len(event.Choices) > 0 && event.Choices[0].Delta.Content != "" {
        fmt.Print(event.Choices[0].Delta.Content)
    }
})

if err != nil {
    log.Fatal(err)
}
```

### Продвинутое использование с опциями

```go
// Настройка менеджера токенов
tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithScope("GIGACHAT_API_PERS"),
    gigachat.WithInsecureSkipVerify(false), // Установите true для отключения проверки SSL
)

// Настройка клиента
client := gigachat.NewClient(
    tokenManager,
    gigachat.WithDefaultModel(gigachat.GigaChat2Pro),
    gigachat.WithBaseURI("https://gigachat.devices.sberbank.ru"),
)

// Отправка сообщения с опциями
messages := []gigachat.Message{
    {Role: "system", Content: "Ты эксперт по квантовой физике"},
    {Role: "user", Content: "Объясни квантовую суперпозицию простыми словами"},
}

response, err := client.Chat(
    messages,
    gigachat.WithModel(gigachat.GigaChat2Pro),
    gigachat.WithTemperature(0.7),
    gigachat.WithMaxTokens(500),
    gigachat.WithTopP(0.9),
)

if err != nil {
    log.Fatal(err)
}

fmt.Println(response.Choices[0].Message.Content)
fmt.Printf("Использовано токенов: %d\n", response.Usage.TotalTokens)
```

## 🤖 Доступные модели

GigaChat поддерживает несколько моделей для различных задач. Актуальный список моделей доступен в [официальной документации](https://developers.sber.ru/docs/ru/gigachat/models).

### Модели для генерации текста

| Модель | Описание | Использование |
|--------|----------|---------------|
| **GigaChat-2** | Базовая модель второго поколения | Общие задачи, диалоги |
| **GigaChat-2-Pro** | Продвинутая модель с улучшенными возможностями | Сложные задачи, креативное письмо |
| **GigaChat-2-Max** | Максимальная модель для самых сложных задач | Профессиональные задачи, анализ |

### Модели для эмбеддингов

| Модель | Описание | Использование |
|--------|----------|---------------|
| **Embeddings** | Базовая модель для векторного представления | Поиск по смыслу, кластеризация |
| **EmbeddingsGigaR** | Улучшенная модель для создания эмбеддингов | Точный поиск, семантический анализ |

### Использование констант моделей

```go
// Использование констант для генерации
response, err := client.Chat(
    messages,
    gigachat.WithModel(gigachat.GigaChat2Pro),
)

// Получение списка доступных моделей
generationModels := gigachat.GetGenerationModels()
embeddingModels := gigachat.GetEmbeddingModels()

// Проверка валидности модели
if gigachat.IsValidGenerationModel(gigachat.GigaChat2) {
    // Модель валидна для генерации
}
```

## 🔧 Параметры генерации

Доступные параметры для настройки генерации:

```go
response, err := client.Chat(
    messages,
    gigachat.WithModel(gigachat.GigaChat2Pro),      // Модель для использования
    gigachat.WithTemperature(0.7),                   // Креативность (0.0 - 2.0)
    gigachat.WithTopP(0.9),                          // Nucleus sampling (0.0 - 1.0)
    gigachat.WithMaxTokens(1000),                    // Максимальное количество токенов
    gigachat.WithRepetitionPenalty(1.1),             // Штраф за повторения (0.0 - 2.0)
    gigachat.WithUpdateInterval(0),                  // Интервал обновления для streaming
)
```

## 🎨 Генерация изображений

GigaChat поддерживает генерацию изображений с помощью встроенной функции text2image. Для создания изображений используйте глагол "нарисуй" в промпте и параметр `function_call: auto`.

### Базовое использование

```go
// Простая генерация изображения
result, err := client.CreateImage("Нарисуй красивый закат над морем")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Изображение создано! File ID: %s\n", result.FileID)

// Декодирование и сохранение изображения
imageData, err := base64.StdEncoding.DecodeString(result.Content)
if err != nil {
    log.Fatal(err)
}

err = os.WriteFile("sunset.jpg", imageData, 0644)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Изображение сохранено в sunset.jpg")
```

### Генерация с системным промптом (стилизация)

```go
// Генерация в стиле конкретного художника
result, err := client.CreateImage(
    "Нарисуй розового кота",
    gigachat.WithSystemMessage("Ты — Василий Кандинский"),
)

// Генерация в определенном стиле
result, err := client.CreateImage(
    "Нарисуй космический корабль",
    gigachat.WithSystemMessage("Ты — художник-концептуалист научной фантастики"),
    gigachat.WithImageTemperature(0.8),
)
```

### Раздельная генерация и скачивание

```go
// Генерация изображения (возвращает ответ с ID изображения)
response, err := client.GenerateImage("Нарисуй футуристический город")
if err != nil {
    log.Fatal(err)
}

// Извлечение ID изображения из ответа
content := response.Choices[0].Message.Content
// Ответ содержит: <img src="file-id" fuse="true"/>

// Отдельное скачивание изображения
imageData, err := client.DownloadImage(fileID)
if err != nil {
    log.Fatal(err)
}

// Сохранение файла
decodedData, _ := base64.StdEncoding.DecodeString(imageData)
os.WriteFile("city.jpg", decodedData, 0644)
```

### Доступные методы для работы с изображениями

| Метод | Описание | Возвращает |
|-------|----------|------------|
| `GenerateImage(prompt, options...)` | Генерирует изображение и возвращает ответ API | `*ChatResponse` |
| `DownloadImage(fileID)` | Скачивает изображение по ID | `string` (base64) |
| `CreateImage(prompt, options...)` | Генерирует и скачивает изображение в одном вызове | `*ImageResult` |

### Опции генерации изображений

```go
result, err := client.CreateImage(
    "Нарисуй дракона",
    gigachat.WithSystemMessage("Ты — художник фэнтези"),
    gigachat.WithImageModel(gigachat.GigaChat2Pro),
    gigachat.WithImageTemperature(0.8),
)
```

### Обработка ошибок при генерации изображений

```go
result, err := client.CreateImage("Нарисуй дракона")
if err != nil {
    switch e := err.(type) {
    case *gigachat.ValidationError:
        fmt.Printf("Ошибка валидации: %v\n", e)
    case *gigachat.GigaChatError:
        fmt.Printf("Ошибка GigaChat API (код %d): %v\n", e.Code, e)
    case *gigachat.AuthenticationError:
        fmt.Printf("Ошибка аутентификации: %v\n", e)
    default:
        fmt.Printf("Неизвестная ошибка: %v\n", err)
    }
    return
}

fmt.Printf("Изображение сохранено с ID: %s\n", result.FileID)
```

### Примеры промптов для генерации

```go
// Хорошие промпты (содержат "нарисуй")
client.CreateImage("Нарисуй закат в горах")
client.CreateImage("Нарисуй портрет кота в стиле ренессанса")
client.CreateImage("Нарисуй абстрактную композицию")

// Стилизация через system_message
client.CreateImage(
    "Нарисуй цветы",
    gigachat.WithSystemMessage("Ты — Клод Моне, рисуешь в стиле импрессионизма"),
)

client.CreateImage(
    "Нарисуй робота",
    gigachat.WithSystemMessage("Ты — концепт-художник для научно-фантастических фильмов"),
)
```

> **Важно**: Для генерации изображений промпт должен содержать глагол "нарисуй" или аналогичные команды рисования. API автоматически определяет необходимость вызова функции text2image при наличии параметра `function_call: auto`.

## ⚠️ Обработка ошибок

SDK предоставляет специализированные типы ошибок для различных сценариев:

```go
response, err := client.Chat(messages)
if err != nil {
    switch e := err.(type) {
    case *gigachat.AuthenticationError:
        // Ошибки авторизации (неверные ключи, истекший токен)
        fmt.Printf("Ошибка авторизации: %v\n", e)
    case *gigachat.ValidationError:
        // Ошибки валидации (неверный формат сообщений)
        fmt.Printf("Ошибка валидации: %v\n", e)
    case *gigachat.GigaChatError:
        // Общие ошибки GigaChat API
        fmt.Printf("Ошибка GigaChat (код %d): %v\n", e.Code, e)
    default:
        fmt.Printf("Неизвестная ошибка: %v\n", err)
    }
    return
}
```

> 📖 **Подробнее об ошибках**: [Официальная документация GigaChat API](https://developers.sber.ru/docs/ru/gigachat/api/errors-description)

## 📚 Примеры

Репозиторий включает полные примеры:

- **[Базовое использование](examples/basic/main.go)** - Простой чат и список моделей
- **[Streaming](examples/streaming/main.go)** - Потоковые ответы в реальном времени
- **[Генерация изображений](examples/image_generation/main.go)** - Генерация и сохранение изображений
- **[Диалоги](examples/conversation/main.go)** - Многоходовые диалоги
- **[Продвинутое](examples/advanced/main.go)** - Продвинутая настройка и опции

Запуск примеров:

```bash
cd examples/basic
go run main.go

cd ../streaming
go run main.go

cd ../image_generation
go run main.go
```

## 🔧 Опции конфигурации

### Опции менеджера токенов

```go
tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithScope("GIGACHAT_API_PERS"),           // Scope API
    gigachat.WithOAuthURI("https://..."),              // Пользовательский OAuth URI
    gigachat.WithInsecureSkipVerify(true),             // Пропустить проверку SSL
    gigachat.WithHTTPClient(customHTTPClient),         // Пользовательский HTTP клиент
)
```

### Опции клиента

```go
client := gigachat.NewClient(
    tokenManager,
    gigachat.WithBaseURI("https://..."),               // Пользовательский base URI
    gigachat.WithDefaultModel(gigachat.GigaChat2Pro),  // Модель по умолчанию
    gigachat.WithClientInsecureSkipVerify(true),       // Пропустить проверку SSL
    gigachat.WithHTTPClient(customHTTPClient),         // Пользовательский HTTP клиент
)
```

## 🧪 Тестирование

```bash
# Запуск тестов
go test -v ./...

# Запуск тестов с покрытием
go test -v -cover ./...

# Генерация отчета о покрытии
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📖 Документация

- [Официальная документация GigaChat API](https://developers.sber.ru/docs/ru/gigachat/api/overview)
- [Модели GigaChat](https://developers.sber.ru/docs/ru/gigachat/models)
- [Справочник API](https://developers.sber.ru/docs/ru/gigachat/api/reference)
- [Коды ошибок](https://developers.sber.ru/docs/ru/gigachat/api/errors-description)
- [Руководство по быстрому старту](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)

## 🤝 Вклад в проект

Вклад приветствуется! Пожалуйста, не стесняйтесь отправлять Pull Request.

1. Сделайте Fork репозитория
2. Создайте ветку для новой функции (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add some amazing feature'`)
4. Отправьте в ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

## 📝 Лицензия

Этот проект лицензирован под лицензией MIT - см. файл [LICENSE](LICENSE) для деталей.

## 👤 Автор

**Igor Sazonov**

- GitHub: [@tigusigalpa](https://github.com/tigusigalpa)
- Email: sovletig@gmail.com

## 🙏 Благодарности

- [Sber GigaChat](https://developers.sber.ru/docs/ru/gigachat/api/overview) за потрясающий AI API
- Сообществу Go за отличные инструменты и библиотеки

## 📊 Статус проекта

Проект активно поддерживается. Если вы столкнулись с проблемами или у вас есть предложения, пожалуйста, [откройте issue](https://github.com/tigusigalpa/gigachat-go/issues).

---

Сделано с ❤️ [Igor Sazonov](https://github.com/tigusigalpa)
