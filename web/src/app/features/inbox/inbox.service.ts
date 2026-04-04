import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import {
  InboxThread,
  InboxLabel,
  ThreadDetail,
  PaginationMeta,
} from './inbox.models';

@Injectable({ providedIn: 'root' })
export class InboxService {
  private readonly http = inject(HttpClient);

  getThreads(
    page: number,
    perPage: number,
    label?: string,
  ): Observable<{ data: InboxThread[]; meta: PaginationMeta }> {
    let params = new HttpParams()
      .set('page', page)
      .set('per_page', perPage);
    if (label) {
      params = params.set('label', label);
    }
    return this.http.get<{ data: InboxThread[]; meta: PaginationMeta }>(
      '/v1/inbox/threads',
      { params },
    );
  }

  getThread(id: string): Observable<{ data: ThreadDetail }> {
    return this.http.get<{ data: ThreadDetail }>(`/v1/inbox/threads/${id}`);
  }

  searchThreads(
    query: string,
    page: number,
    perPage: number,
  ): Observable<{ data: InboxThread[]; meta: PaginationMeta }> {
    const params = new HttpParams()
      .set('q', query)
      .set('page', page)
      .set('per_page', perPage);
    return this.http.get<{ data: InboxThread[]; meta: PaginationMeta }>(
      '/v1/inbox/search',
      { params },
    );
  }

  getLabels(): Observable<{ data: InboxLabel[] }> {
    return this.http.get<{ data: InboxLabel[] }>('/v1/inbox/labels');
  }

  createLabel(name: string, color: string): Observable<{ data: InboxLabel }> {
    return this.http.post<{ data: InboxLabel }>('/v1/inbox/labels', {
      name,
      color,
    });
  }

  starThread(threadId: string, starred: boolean): Observable<void> {
    return this.http.patch<void>(`/v1/inbox/threads/${threadId}`, {
      is_starred: starred,
    });
  }

  markRead(threadIds: string[], read: boolean): Observable<void> {
    return this.http.post<void>('/v1/inbox/threads/mark-read', {
      thread_ids: threadIds,
      is_read: read,
    });
  }

  archiveThreads(threadIds: string[]): Observable<void> {
    return this.http.post<void>('/v1/inbox/threads/archive', {
      thread_ids: threadIds,
    });
  }

  trashThreads(threadIds: string[]): Observable<void> {
    return this.http.post<void>('/v1/inbox/threads/trash', {
      thread_ids: threadIds,
    });
  }
}
