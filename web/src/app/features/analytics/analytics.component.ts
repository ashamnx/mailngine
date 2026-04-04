import { Component, inject, signal, computed } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { DatePipe, DecimalPipe } from '@angular/common';

interface OverviewData {
  emails_sent: number;
  emails_delivered: number;
  bounce_rate: number;
  emails_received: number;
}

interface TimeseriesEntry {
  date: string;
  sent: number;
  delivered: number;
  bounced: number;
  received: number;
}

interface EventBreakdown {
  event: string;
  count: number;
}

@Component({
  selector: 'app-analytics',
  standalone: true,
  imports: [FormsModule, DatePipe, DecimalPipe],
  templateUrl: './analytics.component.html',
  styleUrl: './analytics.component.scss',
})
export class AnalyticsComponent {
  private readonly http = inject(HttpClient);

  readonly overview = signal<OverviewData | null>(null);
  readonly timeseries = signal<TimeseriesEntry[]>([]);
  readonly events = signal<EventBreakdown[]>([]);
  readonly isLoading = signal(true);

  readonly dateFrom = signal('');
  readonly dateTo = signal('');

  readonly maxEventCount = computed(() => {
    const items = this.events();
    if (items.length === 0) return 1;
    return Math.max(...items.map((e) => e.count), 1);
  });

  readonly kpiCards = computed(() => {
    const data = this.overview();
    return [
      {
        icon: 'mail',
        title: 'Emails Sent',
        value: data ? data.emails_sent.toLocaleString() : '--',
      },
      {
        icon: 'check_circle',
        title: 'Delivered',
        value: data ? data.emails_delivered.toLocaleString() : '--',
      },
      {
        icon: 'error',
        title: 'Bounce Rate',
        value: data ? data.bounce_rate.toFixed(1) + '%' : '--',
      },
      {
        icon: 'inbox',
        title: 'Received',
        value: data ? data.emails_received.toLocaleString() : '--',
      },
    ];
  });

  constructor() {
    this.initDateRange();
    this.loadAll();
  }

  private initDateRange(): void {
    const today = new Date();
    const thirtyDaysAgo = new Date(today);
    thirtyDaysAgo.setDate(today.getDate() - 30);

    this.dateTo.set(this.formatDate(today));
    this.dateFrom.set(this.formatDate(thirtyDaysAgo));
  }

  private formatDate(date: Date): string {
    const y = date.getFullYear();
    const m = String(date.getMonth() + 1).padStart(2, '0');
    const d = String(date.getDate()).padStart(2, '0');
    return `${y}-${m}-${d}`;
  }

  applyDateRange(): void {
    this.loadAll();
  }

  private buildParams(): HttpParams {
    let params = new HttpParams();
    if (this.dateFrom()) {
      params = params.set('from', this.dateFrom());
    }
    if (this.dateTo()) {
      params = params.set('to', this.dateTo());
    }
    return params;
  }

  loadAll(): void {
    this.isLoading.set(true);
    const params = this.buildParams();

    let completed = 0;
    const checkDone = () => {
      completed++;
      if (completed >= 3) {
        this.isLoading.set(false);
      }
    };

    this.http
      .get<{ data: OverviewData }>('/v1/analytics/overview', { params })
      .subscribe({
        next: (res) => {
          this.overview.set(res.data);
          checkDone();
        },
        error: () => checkDone(),
      });

    this.http
      .get<{ data: TimeseriesEntry[] }>('/v1/analytics/timeseries', { params })
      .subscribe({
        next: (res) => {
          this.timeseries.set(res.data);
          checkDone();
        },
        error: () => checkDone(),
      });

    this.http
      .get<{ data: EventBreakdown[] }>('/v1/analytics/events', { params })
      .subscribe({
        next: (res) => {
          this.events.set(res.data);
          checkDone();
        },
        error: () => checkDone(),
      });
  }

  eventBarWidth(count: number): string {
    const max = this.maxEventCount();
    return ((count / max) * 100).toFixed(1) + '%';
  }

  formatEventLabel(event: string): string {
    return event
      .replace(/[._]/g, ' ')
      .replace(/\b\w/g, (c) => c.toUpperCase());
  }
}
