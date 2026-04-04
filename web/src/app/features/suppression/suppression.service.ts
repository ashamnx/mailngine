import { Injectable, inject } from '@angular/core';
import { HttpClient, HttpParams } from '@angular/common/http';
import { Observable } from 'rxjs';
import { Suppression, SuppressionReason, PaginationMeta } from './suppression.models';

@Injectable({ providedIn: 'root' })
export class SuppressionService {
  private readonly http = inject(HttpClient);

  getSuppressions(
    page: number,
    perPage: number,
  ): Observable<{ data: Suppression[]; meta: PaginationMeta }> {
    const params = new HttpParams()
      .set('page', page)
      .set('per_page', perPage);
    return this.http.get<{ data: Suppression[]; meta: PaginationMeta }>(
      '/v1/suppressions',
      { params },
    );
  }

  addSuppression(
    email: string,
    reason: SuppressionReason,
  ): Observable<{ data: Suppression }> {
    return this.http.post<{ data: Suppression }>('/v1/suppressions', {
      email,
      reason,
    });
  }

  removeSuppression(id: string): Observable<void> {
    return this.http.delete<void>(`/v1/suppressions/${id}`);
  }
}
