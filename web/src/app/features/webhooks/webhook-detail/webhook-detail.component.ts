import { Component, inject, signal, computed } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';

interface Webhook {
  id: string;
  url: string;
  events: string[];
  active: boolean;
  secret: string;
  created_at: string;
  updated_at: string;
}

interface Delivery {
  id: string;
  event_type: string;
  status: string;
  response_code: number | null;
  attempt: number;
  created_at: string;
}

interface DeliveryResponse {
  data: Delivery[];
  meta: {
    page: number;
    per_page: number;
    total: number;
    total_pages: number;
  };
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
  selector: 'app-webhook-detail',
  standalone: true,
  imports: [FormsModule, DatePipe],
  templateUrl: './webhook-detail.component.html',
  styleUrl: './webhook-detail.component.scss',
})
export class WebhookDetailComponent {
  private readonly http = inject(HttpClient);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly webhook = signal<Webhook | null>(null);
  readonly isLoading = signal(true);
  readonly isSaving = signal(false);
  readonly saveError = signal<string | null>(null);

  // Editable fields
  readonly editUrl = signal('');
  readonly editEvents = signal<Set<string>>(new Set());
  readonly editActive = signal(false);

  // Secret
  readonly secretRevealed = signal(false);
  readonly secretCopied = signal(false);

  // Deliveries
  readonly deliveries = signal<Delivery[]>([]);
  readonly deliveriesMeta = signal<DeliveryResponse['meta'] | null>(null);
  readonly deliveriesLoading = signal(false);
  readonly deliveriesPage = signal(1);
  readonly deliveriesPerPage = 10;

  // Delete
  readonly showDeleteConfirm = signal(false);
  readonly isDeleting = signal(false);

  readonly allEvents = ALL_EVENTS;

  private webhookId = '';

  readonly showingFrom = computed(() => {
    const m = this.deliveriesMeta();
    if (!m || m.total === 0) return 0;
    return (m.page - 1) * m.per_page + 1;
  });

  readonly showingTo = computed(() => {
    const m = this.deliveriesMeta();
    if (!m) return 0;
    return Math.min(m.page * m.per_page, m.total);
  });

  readonly hasPrevDeliveries = computed(() => this.deliveriesPage() > 1);

  readonly hasNextDeliveries = computed(() => {
    const m = this.deliveriesMeta();
    return m ? this.deliveriesPage() < m.total_pages : false;
  });

  constructor() {
    this.webhookId = this.route.snapshot.params['id'];
    this.loadWebhook();
    this.loadDeliveries();
  }

  loadWebhook(): void {
    this.isLoading.set(true);
    this.http
      .get<{ data: Webhook }>(`/v1/webhooks/${this.webhookId}`)
      .subscribe({
        next: (res) => {
          this.webhook.set(res.data);
          this.editUrl.set(res.data.url);
          this.editEvents.set(new Set(res.data.events));
          this.editActive.set(res.data.active);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  loadDeliveries(): void {
    this.deliveriesLoading.set(true);
    this.http
      .get<DeliveryResponse>(
        `/v1/webhooks/${this.webhookId}/deliveries?page=${this.deliveriesPage()}&per_page=${this.deliveriesPerPage}`,
      )
      .subscribe({
        next: (res) => {
          this.deliveries.set(res.data);
          this.deliveriesMeta.set(res.meta);
          this.deliveriesLoading.set(false);
        },
        error: () => this.deliveriesLoading.set(false),
      });
  }

  goBack(): void {
    this.router.navigate(['/webhooks']);
  }

  toggleEvent(event: string): void {
    this.editEvents.update((current) => {
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
    return this.editEvents().has(event);
  }

  saveWebhook(): void {
    const url = this.editUrl().trim();
    const events = Array.from(this.editEvents());
    if (!url || events.length === 0) return;

    this.isSaving.set(true);
    this.saveError.set(null);

    this.http
      .patch<{ data: Webhook }>(`/v1/webhooks/${this.webhookId}`, {
        url,
        events,
        active: this.editActive(),
      })
      .subscribe({
        next: (res) => {
          this.webhook.set(res.data);
          this.isSaving.set(false);
        },
        error: (err) => {
          this.isSaving.set(false);
          this.saveError.set(
            err.error?.message ?? 'Failed to save webhook. Please try again.',
          );
        },
      });
  }

  toggleSecretReveal(): void {
    this.secretRevealed.update((v) => !v);
  }

  copySecret(): void {
    const secret = this.webhook()?.secret;
    if (!secret) return;

    navigator.clipboard.writeText(secret).then(() => {
      this.secretCopied.set(true);
      setTimeout(() => this.secretCopied.set(false), 2000);
    });
  }

  maskedSecret(): string {
    const secret = this.webhook()?.secret;
    if (!secret) return '';
    return secret.substring(0, 8) + '\u2022'.repeat(24);
  }

  confirmDelete(): void {
    this.showDeleteConfirm.set(true);
  }

  cancelDelete(): void {
    this.showDeleteConfirm.set(false);
  }

  deleteWebhook(): void {
    this.isDeleting.set(true);
    this.http.delete(`/v1/webhooks/${this.webhookId}`).subscribe({
      next: () => {
        this.isDeleting.set(false);
        this.router.navigate(['/webhooks']);
      },
      error: () => this.isDeleting.set(false),
    });
  }

  prevDeliveries(): void {
    if (this.hasPrevDeliveries()) {
      this.deliveriesPage.update((p) => p - 1);
      this.loadDeliveries();
    }
  }

  nextDeliveries(): void {
    if (this.hasNextDeliveries()) {
      this.deliveriesPage.update((p) => p + 1);
      this.loadDeliveries();
    }
  }

  deliveryStatusClass(status: string): string {
    switch (status) {
      case 'success':
        return 'badge badge--success';
      case 'pending':
        return 'badge badge--warning';
      case 'failed':
        return 'badge badge--error';
      default:
        return 'badge badge--neutral';
    }
  }

  formatEventName(event: string): string {
    return event.replace('email.', '');
  }
}
