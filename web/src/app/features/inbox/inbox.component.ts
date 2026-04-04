import { Component, inject, signal, computed } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { InboxService } from './inbox.service';
import {
  InboxThread,
  InboxMessage,
  InboxLabel,
  ThreadDetail,
  PaginationMeta,
  SystemLabel,
  SystemLabelItem,
} from './inbox.models';

@Component({
  selector: 'app-inbox',
  standalone: true,
  imports: [FormsModule],
  templateUrl: './inbox.component.html',
  styleUrl: './inbox.component.scss',
})
export class InboxComponent {
  private readonly inboxService = inject(InboxService);

  // System labels
  readonly systemLabels: SystemLabelItem[] = [
    { id: 'inbox', name: 'Inbox', icon: 'inbox' },
    { id: 'starred', name: 'Starred', icon: 'star' },
    { id: 'sent', name: 'Sent', icon: 'send' },
    { id: 'archive', name: 'Archive', icon: 'archive' },
    { id: 'trash', name: 'Trash', icon: 'delete' },
  ];

  // State signals
  readonly threads = signal<InboxThread[]>([]);
  readonly meta = signal<PaginationMeta | null>(null);
  readonly customLabels = signal<InboxLabel[]>([]);
  readonly activeLabel = signal<string>('inbox');
  readonly selectedThread = signal<ThreadDetail | null>(null);
  readonly selectedThreadIds = signal<Set<string>>(new Set());
  readonly expandedMessageIds = signal<Set<string>>(new Set());

  readonly isLoadingThreads = signal(true);
  readonly isLoadingThread = signal(false);
  readonly isLoadingLabels = signal(true);

  readonly currentPage = signal(1);
  readonly perPage = 20;

  readonly searchQuery = signal('');
  readonly isSearching = signal(false);

  readonly showCreateLabel = signal(false);
  readonly newLabelName = signal('');
  readonly newLabelColor = signal('#1A73E8');

  // Computed
  readonly allSelected = computed(() => {
    const threads = this.threads();
    const selected = this.selectedThreadIds();
    return threads.length > 0 && threads.every((t) => selected.has(t.id));
  });

  readonly someSelected = computed(() => {
    const selected = this.selectedThreadIds();
    return selected.size > 0;
  });

  readonly hasUnreadSelected = computed(() => {
    const selected = this.selectedThreadIds();
    return this.threads().some((t) => selected.has(t.id) && !t.is_read);
  });

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
    this.loadThreads();
    this.loadLabels();
  }

  // -- Data Loading --

  loadThreads(): void {
    this.isLoadingThreads.set(true);
    const query = this.searchQuery();

    if (query.trim()) {
      this.inboxService
        .searchThreads(query, this.currentPage(), this.perPage)
        .subscribe({
          next: (res) => {
            this.threads.set(res.data);
            this.meta.set(res.meta);
            this.isLoadingThreads.set(false);
            this.isSearching.set(true);
          },
          error: () => this.isLoadingThreads.set(false),
        });
    } else {
      this.isSearching.set(false);
      this.inboxService
        .getThreads(this.currentPage(), this.perPage, this.activeLabel())
        .subscribe({
          next: (res) => {
            this.threads.set(res.data);
            this.meta.set(res.meta);
            this.isLoadingThreads.set(false);
          },
          error: () => this.isLoadingThreads.set(false),
        });
    }
  }

  loadLabels(): void {
    this.isLoadingLabels.set(true);
    this.inboxService.getLabels().subscribe({
      next: (res) => {
        this.customLabels.set(res.data);
        this.isLoadingLabels.set(false);
      },
      error: () => this.isLoadingLabels.set(false),
    });
  }

  loadThread(id: string): void {
    this.isLoadingThread.set(true);
    this.inboxService.getThread(id).subscribe({
      next: (res) => {
        this.selectedThread.set(res.data);
        // Expand only the last message by default
        const messages = res.data.messages;
        if (messages.length > 0) {
          this.expandedMessageIds.set(
            new Set([messages[messages.length - 1].id]),
          );
        }
        this.isLoadingThread.set(false);
      },
      error: () => this.isLoadingThread.set(false),
    });
  }

  // -- Label Navigation --

  selectLabel(labelId: string): void {
    this.activeLabel.set(labelId);
    this.currentPage.set(1);
    this.selectedThread.set(null);
    this.selectedThreadIds.set(new Set());
    this.searchQuery.set('');
    this.isSearching.set(false);
    this.loadThreads();
  }

  // -- Thread Selection --

  selectThread(thread: InboxThread): void {
    this.loadThread(thread.id);
    // Mark as read if unread
    if (!thread.is_read) {
      this.inboxService.markRead([thread.id], true).subscribe({
        next: () => {
          this.threads.update((threads) =>
            threads.map((t) =>
              t.id === thread.id ? { ...t, is_read: true } : t,
            ),
          );
        },
      });
    }
  }

  toggleThreadSelection(event: Event, threadId: string): void {
    event.stopPropagation();
    this.selectedThreadIds.update((ids) => {
      const next = new Set(ids);
      if (next.has(threadId)) {
        next.delete(threadId);
      } else {
        next.add(threadId);
      }
      return next;
    });
  }

  toggleSelectAll(): void {
    if (this.allSelected()) {
      this.selectedThreadIds.set(new Set());
    } else {
      this.selectedThreadIds.set(
        new Set(this.threads().map((t) => t.id)),
      );
    }
  }

  isThreadSelected(threadId: string): boolean {
    return this.selectedThreadIds().has(threadId);
  }

  isActiveThread(threadId: string): boolean {
    return this.selectedThread()?.thread.id === threadId;
  }

  // -- Star --

  toggleStar(event: Event, thread: InboxThread): void {
    event.stopPropagation();
    const newState = !thread.is_starred;
    this.inboxService.starThread(thread.id, newState).subscribe({
      next: () => {
        this.threads.update((threads) =>
          threads.map((t) =>
            t.id === thread.id ? { ...t, is_starred: newState } : t,
          ),
        );
      },
    });
  }

  // -- Bulk Actions --

  archiveSelected(): void {
    const ids = Array.from(this.selectedThreadIds());
    if (ids.length === 0) return;
    this.inboxService.archiveThreads(ids).subscribe({
      next: () => {
        this.threads.update((threads) =>
          threads.filter((t) => !this.selectedThreadIds().has(t.id)),
        );
        this.selectedThreadIds.set(new Set());
        if (
          this.selectedThread() &&
          ids.includes(this.selectedThread()!.thread.id)
        ) {
          this.selectedThread.set(null);
        }
      },
    });
  }

  trashSelected(): void {
    const ids = Array.from(this.selectedThreadIds());
    if (ids.length === 0) return;
    this.inboxService.trashThreads(ids).subscribe({
      next: () => {
        this.threads.update((threads) =>
          threads.filter((t) => !this.selectedThreadIds().has(t.id)),
        );
        this.selectedThreadIds.set(new Set());
        if (
          this.selectedThread() &&
          ids.includes(this.selectedThread()!.thread.id)
        ) {
          this.selectedThread.set(null);
        }
      },
    });
  }

  toggleReadSelected(): void {
    const ids = Array.from(this.selectedThreadIds());
    if (ids.length === 0) return;
    const markAsRead = this.hasUnreadSelected();
    this.inboxService.markRead(ids, markAsRead).subscribe({
      next: () => {
        this.threads.update((threads) =>
          threads.map((t) =>
            this.selectedThreadIds().has(t.id)
              ? { ...t, is_read: markAsRead }
              : t,
          ),
        );
        this.selectedThreadIds.set(new Set());
      },
    });
  }

  // -- Search --

  onSearch(): void {
    this.currentPage.set(1);
    this.selectedThread.set(null);
    this.loadThreads();
  }

  clearSearch(): void {
    this.searchQuery.set('');
    this.isSearching.set(false);
    this.currentPage.set(1);
    this.loadThreads();
  }

  // -- Pagination --

  prevPage(): void {
    if (this.hasPrev()) {
      this.currentPage.update((p) => p - 1);
      this.loadThreads();
    }
  }

  nextPage(): void {
    if (this.hasNext()) {
      this.currentPage.update((p) => p + 1);
      this.loadThreads();
    }
  }

  // -- Message Expand/Collapse --

  toggleMessage(messageId: string): void {
    this.expandedMessageIds.update((ids) => {
      const next = new Set(ids);
      if (next.has(messageId)) {
        next.delete(messageId);
      } else {
        next.add(messageId);
      }
      return next;
    });
  }

  isMessageExpanded(messageId: string): boolean {
    return this.expandedMessageIds().has(messageId);
  }

  // -- Thread Detail Actions --

  starMessage(message: InboxMessage): void {
    const thread = this.selectedThread();
    if (!thread) return;
    this.inboxService
      .starThread(thread.thread.id, !message.is_starred)
      .subscribe();
  }

  archiveThread(threadId: string): void {
    this.inboxService.archiveThreads([threadId]).subscribe({
      next: () => {
        this.threads.update((threads) =>
          threads.filter((t) => t.id !== threadId),
        );
        this.selectedThread.set(null);
      },
    });
  }

  trashThread(threadId: string): void {
    this.inboxService.trashThreads([threadId]).subscribe({
      next: () => {
        this.threads.update((threads) =>
          threads.filter((t) => t.id !== threadId),
        );
        this.selectedThread.set(null);
      },
    });
  }

  // -- Create Label --

  toggleCreateLabel(): void {
    this.showCreateLabel.update((v) => !v);
    this.newLabelName.set('');
    this.newLabelColor.set('#1A73E8');
  }

  createLabel(): void {
    const name = this.newLabelName().trim();
    if (!name) return;
    this.inboxService.createLabel(name, this.newLabelColor()).subscribe({
      next: (res) => {
        this.customLabels.update((labels) => [...labels, res.data]);
        this.showCreateLabel.set(false);
        this.newLabelName.set('');
      },
    });
  }

  // -- Helpers --

  formatParticipants(addresses: string[]): string {
    if (addresses.length === 0) return 'Unknown';
    if (addresses.length === 1) return this.extractName(addresses[0]);
    return `${this.extractName(addresses[0])} +${addresses.length - 1}`;
  }

  private extractName(email: string): string {
    const atIndex = email.indexOf('@');
    if (atIndex === -1) return email;
    return email.substring(0, atIndex);
  }

  formatDate(dateStr: string): string {
    const date = new Date(dateStr);
    const now = new Date();
    const isToday = date.toDateString() === now.toDateString();

    if (isToday) {
      return date.toLocaleTimeString([], {
        hour: '2-digit',
        minute: '2-digit',
      });
    }

    const isThisYear = date.getFullYear() === now.getFullYear();
    if (isThisYear) {
      return date.toLocaleDateString([], { month: 'short', day: 'numeric' });
    }

    return date.toLocaleDateString([], {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  }

  closeThreadDetail(): void {
    this.selectedThread.set(null);
  }
}
