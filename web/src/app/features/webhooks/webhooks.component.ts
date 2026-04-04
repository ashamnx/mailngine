import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';

interface Webhook {
  id: string;
  url: string;
  events: string[];
  active: boolean;
  created_at: string;
}

const ALL_EVENTS = [
  'email.sent',
  'email.delivered',
  'email.bounced',
  'email.opened',
  'email.clicked',
  'email.complained',
] as const;

@Component({
  selector: 'app-webhooks',
  standalone: true,
  imports: [FormsModule, DatePipe],
  templateUrl: './webhooks.component.html',
  styleUrl: './webhooks.component.scss',
})
export class WebhooksComponent {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);

  readonly webhooks = signal<Webhook[]>([]);
  readonly isLoading = signal(true);
  readonly showCreateForm = signal(false);
  readonly isCreating = signal(false);
  readonly createError = signal<string | null>(null);

  readonly newUrl = signal('');
  readonly selectedEvents = signal<Set<string>>(new Set());

  readonly allEvents = ALL_EVENTS;

  constructor() {
    this.loadWebhooks();
  }

  loadWebhooks(): void {
    this.isLoading.set(true);
    this.http
      .get<{ data: Webhook[] }>('/v1/webhooks')
      .subscribe({
        next: (res) => {
          this.webhooks.set(res.data);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  openCreateForm(): void {
    this.showCreateForm.set(true);
    this.newUrl.set('');
    this.selectedEvents.set(new Set());
    this.createError.set(null);
  }

  closeCreateForm(): void {
    this.showCreateForm.set(false);
    this.newUrl.set('');
    this.selectedEvents.set(new Set());
    this.createError.set(null);
  }

  toggleEvent(event: string): void {
    this.selectedEvents.update((current) => {
      const next = new Set(current);
      if (next.has(event)) {
        next.delete(event);
      } else {
        next.add(event);
      }
      return next;
    });
  }

  isEventSelected(event: string): boolean {
    return this.selectedEvents().has(event);
  }

  createWebhook(): void {
    const url = this.newUrl().trim();
    const events = Array.from(this.selectedEvents());
    if (!url || events.length === 0) return;

    this.isCreating.set(true);
    this.createError.set(null);

    this.http
      .post<{ data: Webhook }>('/v1/webhooks', { url, events })
      .subscribe({
        next: (res) => {
          this.isCreating.set(false);
          this.closeCreateForm();
          this.router.navigate(['/webhooks', res.data.id]);
        },
        error: (err) => {
          this.isCreating.set(false);
          this.createError.set(
            err.error?.message ?? 'Failed to create webhook. Please try again.',
          );
        },
      });
  }

  navigateToWebhook(id: string): void {
    this.router.navigate(['/webhooks', id]);
  }

  truncateUrl(url: string): string {
    return url.length > 50 ? url.substring(0, 50) + '...' : url;
  }

  formatEventName(event: string): string {
    return event.replace('email.', '');
  }
}
