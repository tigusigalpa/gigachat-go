# GigaChat Go SDK

![GigaChat Golang SDK](https://i.postimg.cc/YCWrRjhK/528512211-5cc6724d-fd74-4908-ba5d-4e473bd6e885.jpg)

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
[![License](https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/tigusigalpa/gigachat-go?style=flat-square)](https://goreportcard.com/report/github.com/tigusigalpa/gigachat-go)

Golang SDK для Sber GigaChat API. Поддерживает streaming, генерацию изображений и работу с диалогами.

**🌐 Язык:** [English](README-en.md) | Русский

**📦 Пакет:** [pkg.go.dev/github.com/tigusigalpa/gigachat-go](https://pkg.go.dev/github.com/tigusigalpa/gigachat-go)

## О проекте

Библиотека для интеграции GigaChat в Go-приложения. Подходит для чат-ботов, генерации контента и изображений.

### Зачем этот SDK?

- Работает с GigaChat от Сбера
- Поддержка streaming и больших объёмов данных
- Автоматическое управление токенами с потокобезопасным кешем
- Обработка ошибок и типобезопасность
- Функциональные опции для настройки

### Где применять

- **Чат-боты** — Telegram, Discord, Slack. SDK хранит контекст диалога.
- **Генерация текстов** — статьи, описания товаров, рассылки.
- **Изображения** — генерация по описанию, стилизация под художников.
- **Аналитика** — суммаризация документов, обработка отзывов.
- **Поиск** — семантический поиск через эмбеддинги.
- **Обучение** — интерактивные системы, проверка знаний.

### Сравнение

| Особенность | GigaChat Go SDK | Альтернативы |
|-------------|-----------------|----------------|
| Русский язык | Нативно | Ограниченно |
| Streaming | Да | Частично |
| Генерация изображений | Встроена | Отдельный API |
| Управление токенами | Автоматически | Вручную |
| Типобезопасность | Да | Частично |

## Возможности

- Интеграция с GigaChat API
- Автоматическое управление OAuth-токенами (потокобезопасный кеш)
- Все модели: GigaChat-2, GigaChat-2-Pro, GigaChat-2-Max
- Диалоги с контекстом
- Streaming-ответы
- Генерация изображений (text2image)
- Скачивание изображений в base64
- Стилизация через системные промпты
- Функциональные опции
- Типобезопасность и обработка ошибок

## Установка

```bash
go get github.com/tigusigalpa/gigachat-go
```

## Настройка

### 1. Получение авторизационных данных

Для работы с API нужны ключи:

1. Зарегистрируйтесь
   в [личном кабинете Sber AI](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)
2. Создайте проект и получите **Client ID** и **Client Secret**
3. Сгенерируйте **Authorization Key** (Base64 от "Client ID:Client Secret")

> Подробнее: [Создание проекта и получение ключей](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)

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

## Использование

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

## Доступные модели

Актуальный список: [документация](https://developers.sber.ru/docs/ru/gigachat/models).

### Модели для генерации текста

| Модель             | Описание                                       | Использование                     |
|--------------------|------------------------------------------------|-----------------------------------|
| **GigaChat-2**     | Базовая модель второго поколения               | Общие задачи, диалоги             |
| **GigaChat-2-Pro** | Продвинутая модель с улучшенными возможностями | Сложные задачи, креативное письмо |
| **GigaChat-2-Max** | Максимальная модель для самых сложных задач    | Профессиональные задачи, анализ   |

### Модели для эмбеддингов

| Модель              | Описание                                    | Использование                      |
|---------------------|---------------------------------------------|------------------------------------|
| **Embeddings**      | Базовая модель для векторного представления | Поиск по смыслу, кластеризация     |
| **EmbeddingsGigaR** | Улучшенная модель для создания эмбеддингов  | Точный поиск, семантический анализ |

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

## Параметры генерации

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

## Генерация изображений

Для создания изображений используйте глагол "нарисуй" в промпте.

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

| Метод                               | Описание                                          | Возвращает        |
|-------------------------------------|---------------------------------------------------|-------------------|
| `GenerateImage(prompt, options...)` | Генерирует изображение и возвращает ответ API     | `*ChatResponse`   |
| `DownloadImage(fileID)`             | Скачивает изображение по ID                       | `string` (base64) |
| `CreateImage(prompt, options...)`   | Генерирует и скачивает изображение в одном вызове | `*ImageResult`    |

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

> Промпт должен содержать "нарисуй" или аналогичную команду. API вызывает text2image автоматически при `function_call: auto`.

## Обработка ошибок

SDK предоставляет типизированные ошибки:

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

> Подробнее: [Коды ошибок GigaChat API](https://developers.sber.ru/docs/ru/gigachat/api/errors-description)

## Примеры

В репозитории есть примеры:

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

## Опции конфигурации

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

## Тестирование

Библиотека включает комплексный набор unit-тестов для всех основных компонентов.

### Структура тестов

- **`errors_test.go`** - Тесты для пользовательских типов ошибок (`GigaChatError`, `AuthenticationError`, `ValidationError`)
- **`models_test.go`** - Тесты для моделей данных, констант и JSON сериализации
- **`token_manager_test.go`** - Тесты управления токенами, кэширования и конкурентного доступа
- **`client_test.go`** - Тесты для всех методов клиента API с использованием mock HTTP серверов

### Запуск тестов

```bash
# Запуск всех тестов
go test -v

# Запуск конкретного файла тестов
go test -v -run TestClient
go test -v -run TestTokenManager

# Запуск конкретного теста
go test -v -run TestClient_Chat

# Запуск с покрытием кода
go test -cover

# Детальный отчет о покрытии
go test -coverprofile=coverage.out
go tool cover -html=coverage.out

# Запуск с race detector
go test -race -v
```

### Покрытие тестами

Тесты покрывают:
- Все публичные функции и методы
- Обработку ошибок и граничные случаи
- Валидацию входных данных
- Конкурентный доступ к токенам
- JSON маршалинг/анмаршалинг
- Mock HTTP серверы для изоляции от реальных API

## Документация

- [Официальная документация GigaChat API](https://developers.sber.ru/docs/ru/gigachat/api/overview)
- [Модели GigaChat](https://developers.sber.ru/docs/ru/gigachat/models)
- [Справочник API](https://developers.sber.ru/docs/ru/gigachat/api/reference)
- [Коды ошибок](https://developers.sber.ru/docs/ru/gigachat/api/errors-description)
- [Руководство по быстрому старту](https://developers.sber.ru/docs/ru/gigachat/quickstart/ind-create-project)

## Архитектура SDK

### Компоненты

#### TokenManager
Управляет OAuth-аутентификацией:
- Получение access token через OAuth 2.0
- Потокобезопасный кеш токенов
- Автообновление при истечении
- Поддержка custom HTTP-клиентов

#### Client
Интерфейс для работы с API:
- Чат, streaming, генерация изображений
- Управление моделями и параметрами
- Типизированные ошибки
- Функциональные опции

#### Models
Структуры данных:
- Константы моделей
- Валидация параметров
- Структуры ответов API
- Хелперы для сообщений

### Поток данных

```
Приложение → Client → TokenManager → OAuth Server
                ↓
           GigaChat API
                ↓
           Response → Client → Приложение
```

## Рекомендации

### Управление контекстом диалогов

```go
// Ограничивайте размер истории для экономии токенов
const maxHistoryMessages = 10

func trimHistory(messages []gigachat.Message) []gigachat.Message {
    if len(messages) <= maxHistoryMessages {
        return messages
    }
    // Сохраняем системное сообщение (если есть) и последние N сообщений
    systemMsg := []gigachat.Message{}
    if messages[0].Role == "system" {
        systemMsg = messages[:1]
        messages = messages[1:]
    }
    return append(systemMsg, messages[len(messages)-maxHistoryMessages:]...)
}
```

### Обработка rate limits

```go
import "time"

func chatWithRetry(client *gigachat.Client, messages []gigachat.Message, maxRetries int) (*gigachat.ChatResponse, error) {
    var response *gigachat.ChatResponse
    var err error
    
    for i := 0; i < maxRetries; i++ {
        response, err = client.Chat(messages)
        if err == nil {
            return response, nil
        }
        
        if gigachatErr, ok := err.(*gigachat.GigaChatError); ok {
            if gigachatErr.Code == 429 {
                // Экспоненциальная задержка при rate limit
                backoff := time.Duration(1<<uint(i)) * time.Second
                time.Sleep(backoff)
                continue
            }
        }
        return nil, err
    }
    return nil, err
}
```

### Оптимизация использования токенов

```go
// Используйте параметр max_tokens для контроля расходов
response, err := client.Chat(
    messages,
    gigachat.WithMaxTokens(500),  // Ограничение длины ответа
    gigachat.WithTemperature(0.3), // Меньше креативности = меньше токенов
)

// Мониторинг использования токенов
if err == nil {
    log.Printf("Использовано токенов: %d (prompt: %d, completion: %d)",
        response.Usage.TotalTokens,
        response.Usage.PromptTokens,
        response.Usage.CompletionTokens,
    )
}
```

### Безопасное хранение credentials

```go
// ❌ Плохо: хардкод в коде
authKey := "hardcoded_key_here"

// ✅ Хорошо: использование переменных окружения
authKey := os.Getenv("GIGACHAT_AUTH_KEY")
if authKey == "" {
    log.Fatal("GIGACHAT_AUTH_KEY не установлен")
}

// ✅ Еще лучше: использование секретных менеджеров
// Например, HashiCorp Vault, AWS Secrets Manager, etc.
```

### Graceful shutdown для streaming

```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

go func() {
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
    <-sigChan
    cancel()
}()

err := client.ChatStream(messages, func(event *gigachat.ChatResponse, done bool, err error) {
    select {
    case <-ctx.Done():
        return
    default:
        // Обработка события
    }
})
```

## Устранение неполадок

### Проблема: Ошибка аутентификации

**Симптомы:**
```
Authentication error: Can't decode 'Authorization' header
```

**Решение:**
1. Проверьте корректность формата Authorization Key (должен быть Base64)
2. Убедитесь, что используете правильный Client ID и Client Secret
3. Проверьте, что ключи не содержат лишних пробелов или символов

```go
// Правильное создание auth key
authKey := base64.StdEncoding.EncodeToString(
    []byte(strings.TrimSpace(clientID) + ":" + strings.TrimSpace(clientSecret)),
)
```

### Проблема: SSL/TLS ошибки

**Симптомы:**
```
x509: certificate signed by unknown authority
```

**Решение:**
```go
// Для разработки (НЕ для production!)
tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithInsecureSkipVerify(true),
)

// Для production: установите корректные сертификаты
```

### Проблема: Превышен лимит токенов

**Симптомы:**
```
HTTP 402: Payment Required
```

**Решение:**
1. Проверьте баланс токенов в личном кабинете
2. Оптимизируйте промпты для уменьшения расхода
3. Используйте параметр `max_tokens` для ограничения ответов
4. Рассмотрите переход на более высокий тарифный план

### Проблема: Медленные ответы

**Решение:**
1. Используйте streaming для больших ответов
2. Выберите подходящую модель (GigaChat-2 быстрее, чем GigaChat-2-Max)
3. Сократите размер промпта и истории диалога
4. Проверьте сетевое соединение

```go
// Используйте streaming для длинных ответов
err := client.ChatStream(messages, func(event *gigachat.ChatResponse, done bool, err error) {
    // Получайте части ответа по мере генерации
})
```

## FAQ

### Какая модель лучше для моей задачи?

- **GigaChat-2**: Подходит для простых диалогов, FAQ-ботов, базовой генерации текста
- **GigaChat-2-Pro**: Оптимален для сложных задач, креативного письма, анализа текстов
- **GigaChat-2-Max**: Используйте для профессиональных задач, требующих максимальной точности

### Можно ли использовать SDK в production?

Да:
- Обработка ошибок
- Потокобезопасность
- Автоматическое управление токенами
- Тесты

### Как ограничить расход токенов?

```go
response, err := client.Chat(
    messages,
    gigachat.WithMaxTokens(500),      // Максимум токенов в ответе
    gigachat.WithTemperature(0.3),    // Меньше вариативности
)
```

### Поддерживается ли работа через прокси?

Да, используйте custom HTTP client:

```go
proxyURL, _ := url.Parse("http://proxy:8080")
httpClient := &http.Client{
    Transport: &http.Transport{
        Proxy: http.ProxyURL(proxyURL),
    },
}

tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithHTTPClient(httpClient),
)
```

### Как обрабатывать большие файлы?

Разбивайте файлы на чанки и обрабатывайте последовательно:

```go
const chunkSize = 4000 // символов

func processLargeText(client *gigachat.Client, text string) ([]string, error) {
    var results []string
    
    for i := 0; i < len(text); i += chunkSize {
        end := i + chunkSize
        if end > len(text) {
            end = len(text)
        }
        
        chunk := text[i:end]
        messages := []gigachat.Message{
            {Role: "user", Content: "Суммируй текст: " + chunk},
        }
        
        response, err := client.Chat(messages)
        if err != nil {
            return nil, err
        }
        
        results = append(results, gigachat.ExtractContent(response))
    }
    
    return results, nil
}
```

### Можно ли использовать несколько моделей одновременно?

Да, создайте несколько клиентов или используйте опции:

```go
// Вариант 1: Разные клиенты
client1 := gigachat.NewClient(tm, gigachat.WithDefaultModel(gigachat.GigaChat2))
client2 := gigachat.NewClient(tm, gigachat.WithDefaultModel(gigachat.GigaChat2Pro))

// Вариант 2: Один клиент с опциями
response1, _ := client.Chat(messages, gigachat.WithModel(gigachat.GigaChat2))
response2, _ := client.Chat(messages, gigachat.WithModel(gigachat.GigaChat2Pro))
```

## Производительность

### Параллельная обработка запросов

```go
func processParallel(client *gigachat.Client, prompts []string) []string {
    results := make([]string, len(prompts))
    var wg sync.WaitGroup
    
    for i, prompt := range prompts {
        wg.Add(1)
        go func(idx int, p string) {
            defer wg.Done()
            
            messages := []gigachat.Message{{Role: "user", Content: p}}
            response, err := client.Chat(messages)
            if err != nil {
                log.Printf("Ошибка для промпта %d: %v", idx, err)
                return
            }
            results[idx] = gigachat.ExtractContent(response)
        }(i, prompt)
    }
    
    wg.Wait()
    return results
}
```

### Пул соединений

```go
// Настройка HTTP клиента с пулом соединений
httpClient := &http.Client{
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
    Timeout: 30 * time.Second,
}

tokenManager := gigachat.NewTokenManager(
    authKey,
    gigachat.WithHTTPClient(httpClient),
)
```

## Безопасность

1. **Никогда не коммитьте credentials в Git**
   ```bash
   # Добавьте в .gitignore
   .env
   *.key
   credentials.json
   ```

2. **Используйте переменные окружения**
   ```go
   authKey := os.Getenv("GIGACHAT_AUTH_KEY")
   ```

3. **Ротация ключей**
   Регулярно обновляйте Client ID и Client Secret в личном кабинете

4. **Валидация пользовательского ввода**
   ```go
   func sanitizeInput(input string) string {
       // Удаление потенциально опасных символов
       input = strings.TrimSpace(input)
       if len(input) > 10000 {
           input = input[:10000]
       }
       return input
   }
   ```

5. **Rate limiting на стороне приложения**
   ```go
   import "golang.org/x/time/rate"
   
   limiter := rate.NewLimiter(rate.Every(time.Second), 5) // 5 запросов в секунду
   
   func rateLimitedChat(client *gigachat.Client, messages []gigachat.Message) (*gigachat.ChatResponse, error) {
       if err := limiter.Wait(context.Background()); err != nil {
           return nil, err
       }
       return client.Chat(messages)
   }
   ```

## Участие в разработке

Pull Request'ы приветствуются.

1. Сделайте Fork репозитория
2. Создайте ветку для новой функции (`git checkout -b feature/amazing-feature`)
3. Зафиксируйте изменения (`git commit -m 'Add some amazing feature'`)
4. Отправьте в ветку (`git push origin feature/amazing-feature`)
5. Откройте Pull Request

### Правила

- Используйте `gofmt`
- Добавляйте тесты
- Обновляйте документацию при изменении API
- Понятные commit messages
- Все тесты должны проходить

## Лицензия

Этот проект лицензирован под лицензией MIT - см. файл [LICENSE](LICENSE) для деталей.

## Автор

**Igor Sazonov**

- GitHub: [@tigusigalpa](https://github.com/tigusigalpa)
- Email: sovletig@gmail.com

## Благодарности

- [Sber GigaChat](https://developers.sber.ru/docs/ru/gigachat/api/overview)
- Go-сообщество

## Статус

Проект поддерживается. Проблемы и предложения — в [issues](https://github.com/tigusigalpa/gigachat-go/issues).
