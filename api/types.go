package api

// Response 是 Telegram Bot API 的通用响应结构
type Response[T any] struct {
	OK          bool   `json:"ok"`
	Result      T      `json:"result"`
	Description string `json:"description"`
	ErrorCode   int    `json:"error_code"`
}

type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name,omitempty"`
	Username     string `json:"username,omitempty"`
	LanguageCode string `json:"language_code,omitempty"`
}

type Chat struct {
	ID       int64  `json:"id"`
	Type     string `json:"type"`
	Title    string `json:"title,omitempty"`
	Username string `json:"username,omitempty"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	From      *User  `json:"from,omitempty"`
	Chat      Chat   `json:"chat"`
	Text      string `json:"text,omitempty"`
	Caption   string `json:"caption,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

type InlineKeyboardButton struct {
	Text string `json:"text"`
	URL  string `json:"url,omitempty"`
}

// SendMessageParams 对应 sendMessage 的请求参数
type SendMessageParams struct {
	ChatID                string                `json:"chat_id"`
	Text                  string                `json:"text"`
	ParseMode             string                `json:"parse_mode,omitempty"`
	MessageThreadID       int                   `json:"message_thread_id,omitempty"`
	ReplyToMessageID      int                   `json:"reply_to_message_id,omitempty"`
	DisableWebPagePreview bool                  `json:"disable_web_page_preview,omitempty"`
	DisableNotification   bool                  `json:"disable_notification,omitempty"`
	ProtectContent        bool                  `json:"protect_content,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// SendMediaParams 用于发送媒体文件时的附加参数（multipart 表单字段）
type SendMediaParams struct {
	ChatID              string
	FilePath            string
	Caption             string
	ParseMode           string
	MessageThreadID     int
	ReplyToMessageID    int
	DisableNotification bool
	ProtectContent      bool
	ReplyMarkup         *InlineKeyboardMarkup
}

// EditMessageTextParams 对应 editMessageText 的请求参数
type EditMessageTextParams struct {
	ChatID                string                `json:"chat_id"`
	MessageID             int                   `json:"message_id"`
	Text                  string                `json:"text"`
	ParseMode             string                `json:"parse_mode,omitempty"`
	DisableWebPagePreview bool                  `json:"disable_web_page_preview,omitempty"`
	ReplyMarkup           *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// DeleteMessageParams 对应 deleteMessage 的请求参数
type DeleteMessageParams struct {
	ChatID    string `json:"chat_id"`
	MessageID int    `json:"message_id"`
}

// ForwardMessageParams 对应 forwardMessage 的请求参数
type ForwardMessageParams struct {
	ChatID              string `json:"chat_id"`
	FromChatID          string `json:"from_chat_id"`
	MessageID           int    `json:"message_id"`
	MessageThreadID     int    `json:"message_thread_id,omitempty"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
	ProtectContent      bool   `json:"protect_content,omitempty"`
}

// CopyMessageParams 对应 copyMessage 的请求参数
type CopyMessageParams struct {
	ChatID              string                `json:"chat_id"`
	FromChatID          string                `json:"from_chat_id"`
	MessageID           int                   `json:"message_id"`
	Caption             string                `json:"caption,omitempty"`
	ParseMode           string                `json:"parse_mode,omitempty"`
	MessageThreadID     int                   `json:"message_thread_id,omitempty"`
	ReplyToMessageID    int                   `json:"reply_to_message_id,omitempty"`
	DisableNotification bool                  `json:"disable_notification,omitempty"`
	ProtectContent      bool                  `json:"protect_content,omitempty"`
	ReplyMarkup         *InlineKeyboardMarkup `json:"reply_markup,omitempty"`
}

// PinChatMessageParams 对应 pinChatMessage 的请求参数
type PinChatMessageParams struct {
	ChatID              string `json:"chat_id"`
	MessageID           int    `json:"message_id"`
	DisableNotification bool   `json:"disable_notification,omitempty"`
}

// UnpinChatMessageParams 对应 unpinChatMessage 的请求参数
type UnpinChatMessageParams struct {
	ChatID    string `json:"chat_id"`
	MessageID int    `json:"message_id,omitempty"`
}
