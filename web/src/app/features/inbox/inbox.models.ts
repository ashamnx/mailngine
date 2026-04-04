export interface InboxThread {
  id: string;
  subject: string;
  participant_addresses: string[];
  last_message_at: string;
  message_count: number;
  snippet: string;
  is_read: boolean;
  is_starred: boolean;
  label_ids: string[];
}

export interface InboxMessage {
  id: string;
  thread_id: string;
  from_address: string;
  from_name: string;
  to_addresses: string[];
  cc_addresses: string[];
  subject: string;
  snippet: string;
  body_html: string;
  body_text: string;
  is_read: boolean;
  is_starred: boolean;
  is_archived: boolean;
  is_trashed: boolean;
  received_at: string;
}

export interface InboxLabel {
  id: string;
  name: string;
  color: string;
}

export interface ThreadDetail {
  thread: InboxThread;
  messages: InboxMessage[];
}

export interface PaginationMeta {
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export type SystemLabel = 'inbox' | 'starred' | 'sent' | 'archive' | 'trash';

export interface SystemLabelItem {
  id: SystemLabel;
  name: string;
  icon: string;
}
