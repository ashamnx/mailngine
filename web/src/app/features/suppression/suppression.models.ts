export interface Suppression {
  id: string;
  email: string;
  reason: SuppressionReason;
  created_at: string;
}

export type SuppressionReason = 'hard_bounce' | 'complaint' | 'manual';

export interface PaginationMeta {
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}
