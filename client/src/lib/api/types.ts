/** Mirrors internal/core/models.go and handler request bodies. */

export type User = {
	id: string;
	name: string;
	email: string;
};

export type ChatDto = {
	id: string;
	name: string;
	image_url?: string;
};

export type ChatMessage = {
	id: string;
	user_id: string;
	chat_id: string;
	content: string;
	media_url?: string;
	reply_to_id?: string;
	created_at: string;
};

export type RegisterBody = {
	username: string;
	email: string;
	password: string;
};

export type RegisterResponse = {
	user_id: string;
};

export type LoginBody = {
	email: string;
	password: string;
};

export type LoginResponse = {
	token: string;
};

export type CreateChatBody = {
	name: string;
	chat_member_ids: string[];
	image_url?: string;
};

export type CreateChatResponse = {
	id: string;
};

export type SendMessageBody = {
	content: string;
	media_url?: string;
	reply_to_id?: string;
};

export type ProblemDetails = {
	type: string;
	title: string;
	status: number;
	detail?: string;
	instance?: string;
};
