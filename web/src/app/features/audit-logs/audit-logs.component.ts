import { Component, inject, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { DatePipe, JsonPipe } from '@angular/common';

interface AuditLog {
  id: string;
  action: string;
  resource_type: string;
  actor: string;
  ip_address: string;
  metadata: Record<string, unknown>;
  created_at: string;
}

interface AuditLogResponse {
  data: AuditLog[];
  meta: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
}

@Component({
  selector: 'app-audit-logs',
  standalone: true,
  imports: [DatePipe, JsonPipe],
  templateUrl: './audit-logs.component.html',
  styleUrl: './audit-logs.component.scss',
})
export class AuditLogsComponent {
  private readonly http = inject(HttpClient);

  readonly logs = signal<AuditLog[]>([]);
  readonly meta = signal<AuditLogResponse['meta'] | null>(null);
  readonly isLoading = signal(true);
  readonly currentPage = signal(1);
  readonly perPage = 20;
  readonly expandedId = signal<string | null>(null);

  readonly showingFrom = computed(() => {
    const m = this.meta();
    if (!m || m.total === 0) return 0;
    return (m.page - 1) * m.per_page + 1;
  });

  readonly showingTo = computed(() => {
    const m = this.meta();
    if (!m) return 0;
    return Math.min(m.page * m.per_page, m.total);
  });

  readonly hasPrev = computed(() => this.currentPage() > 1);

  readonly hasNext = computed(() => {
    const m = this.meta();
    return m ? this.currentPage() < m.total_pages : false;
  });

  constructor() {
    this.loadLogs();
  }

  loadLogs(): void {
    this.isLoading.set(true);
    this.http
      .get<AuditLogResponse>(
        `/v1/audit-logs?page=${this.currentPage()}&per_page=${this.perPage}`,
      )
      .subscribe({
        next: (res) => {
          this.logs.set(res.data);
          this.meta.set(res.meta);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  toggleExpand(id: string): void {
    this.expandedId.update((current) => (current === id ? null : id));
  }

  isExpanded(id: string): boolean {
    return this.expandedId() === id;
  }

  hasMetadata(log: AuditLog): boolean {
    return log.metadata !== null && Object.keys(log.metadata).length > 0;
  }

  prevPage(): void {
    if (this.hasPrev()) {
      this.currentPage.update((p) => p - 1);
      this.loadLogs();
    }
  }

  nextPage(): void {
    if (this.hasNext()) {
      this.currentPage.update((p) => p + 1);
      this.loadLogs();
    }
  }
}
