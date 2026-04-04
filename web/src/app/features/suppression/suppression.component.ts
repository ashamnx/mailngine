import { Component, inject, signal, computed } from '@angular/core';
import { DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { SuppressionService } from './suppression.service';
import {
  Suppression,
  SuppressionReason,
  PaginationMeta,
} from './suppression.models';

@Component({
  selector: 'app-suppression',
  standalone: true,
  imports: [DatePipe, FormsModule],
  templateUrl: './suppression.component.html',
  styleUrl: './suppression.component.scss',
})
export class SuppressionComponent {
  private readonly suppressionService = inject(SuppressionService);

  readonly suppressions = signal<Suppression[]>([]);
  readonly meta = signal<PaginationMeta | null>(null);
  readonly isLoading = signal(true);
  readonly currentPage = signal(1);
  readonly perPage = 20;

  // Add form state
  readonly showAddForm = signal(false);
  readonly newEmail = signal('');
  readonly newReason = signal<SuppressionReason>('manual');
  readonly isAdding = signal(false);

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
    this.loadSuppressions();
  }

  loadSuppressions(): void {
    this.isLoading.set(true);
    this.suppressionService
      .getSuppressions(this.currentPage(), this.perPage)
      .subscribe({
        next: (res) => {
          this.suppressions.set(res.data);
          this.meta.set(res.meta);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  toggleAddForm(): void {
    this.showAddForm.update((v) => !v);
    this.newEmail.set('');
    this.newReason.set('manual');
  }

  addSuppression(): void {
    const email = this.newEmail().trim();
    if (!email) return;

    this.isAdding.set(true);
    this.suppressionService.addSuppression(email, this.newReason()).subscribe({
      next: (res) => {
        this.suppressions.update((list) => [res.data, ...list]);
        this.showAddForm.set(false);
        this.newEmail.set('');
        this.isAdding.set(false);
        // Reload to get correct pagination
        this.loadSuppressions();
      },
      error: () => this.isAdding.set(false),
    });
  }

  removeSuppression(suppression: Suppression): void {
    this.suppressionService.removeSuppression(suppression.id).subscribe({
      next: () => {
        this.suppressions.update((list) =>
          list.filter((s) => s.id !== suppression.id),
        );
        // Reload to get correct pagination
        this.loadSuppressions();
      },
    });
  }

  prevPage(): void {
    if (this.hasPrev()) {
      this.currentPage.update((p) => p - 1);
      this.loadSuppressions();
    }
  }

  nextPage(): void {
    if (this.hasNext()) {
      this.currentPage.update((p) => p + 1);
      this.loadSuppressions();
    }
  }

  reasonBadgeClass(reason: SuppressionReason): string {
    switch (reason) {
      case 'hard_bounce':
        return 'badge badge--error';
      case 'complaint':
        return 'badge badge--warning';
      case 'manual':
        return 'badge badge--neutral';
    }
  }

  reasonLabel(reason: SuppressionReason): string {
    switch (reason) {
      case 'hard_bounce':
        return 'Hard Bounce';
      case 'complaint':
        return 'Complaint';
      case 'manual':
        return 'Manual';
    }
  }
}
