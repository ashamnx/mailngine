import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { FormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';

interface ApiKey {
  id: string;
  name: string;
  key_prefix: string;
  permission: string;
  created_at: string;
  last_used_at: string | null;
}

interface CreateApiKeyResponse {
  data: {
    id: string;
    name: string;
    key: string;
    permission: string;
  };
}

@Component({
  selector: 'app-api-keys',
  standalone: true,
  imports: [FormsModule, DatePipe],
  templateUrl: './api-keys.component.html',
  styleUrl: './api-keys.component.scss',
})
export class ApiKeysComponent {
  private readonly http = inject(HttpClient);

  readonly keys = signal<ApiKey[]>([]);
  readonly isLoading = signal(true);
  readonly showCreateForm = signal(false);
  readonly newKeyName = signal('');
  readonly newKeyPermission = signal('full');
  readonly createdKey = signal<string | null>(null);
  readonly isCreating = signal(false);
  readonly copied = signal(false);

  constructor() {
    this.loadKeys();
  }

  loadKeys(): void {
    this.isLoading.set(true);
    this.http
      .get<{ data: ApiKey[] }>('/v1/api-keys')
      .subscribe({
        next: (res) => {
          this.keys.set(res.data);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  openCreateForm(): void {
    this.showCreateForm.set(true);
    this.createdKey.set(null);
    this.newKeyName.set('');
    this.newKeyPermission.set('full');
  }

  closeCreateForm(): void {
    this.showCreateForm.set(false);
    this.createdKey.set(null);
    this.newKeyName.set('');
  }

  createKey(): void {
    const name = this.newKeyName().trim();
    if (!name) return;

    this.isCreating.set(true);
    this.http
      .post<CreateApiKeyResponse>('/v1/api-keys', {
        name,
        permission: this.newKeyPermission(),
      })
      .subscribe({
        next: (res) => {
          this.createdKey.set(res.data.key);
          this.isCreating.set(false);
          this.loadKeys();
        },
        error: () => this.isCreating.set(false),
      });
  }

  revokeKey(id: string): void {
    this.http.delete(`/v1/api-keys/${id}`).subscribe({
      next: () => this.loadKeys(),
    });
  }

  copyKey(): void {
    const key = this.createdKey();
    if (!key) return;

    navigator.clipboard.writeText(key).then(() => {
      this.copied.set(true);
      setTimeout(() => this.copied.set(false), 2000);
    });
  }

  formatPermission(permission: string): string {
    return permission
      .split('_')
      .map((word) => word.charAt(0).toUpperCase() + word.slice(1))
      .join(' ');
  }
}
