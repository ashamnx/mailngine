import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { ActivatedRoute, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';

interface Template {
  id: string;
  name: string;
  subject: string;
  html_body: string;
  text_body: string;
  variables: string[];
  created_at: string;
  updated_at: string;
}

interface PreviewResponse {
  data: {
    subject: string;
    html: string;
  };
}

@Component({
  selector: 'app-template-detail',
  standalone: true,
  imports: [FormsModule, DatePipe],
  templateUrl: './template-detail.component.html',
  styleUrl: './template-detail.component.scss',
})
export class TemplateDetailComponent {
  private readonly http = inject(HttpClient);
  private readonly route = inject(ActivatedRoute);
  private readonly router = inject(Router);

  readonly template = signal<Template | null>(null);
  readonly isLoading = signal(true);
  readonly isSaving = signal(false);
  readonly saveError = signal<string | null>(null);
  readonly saveSuccess = signal(false);

  // Editable fields
  readonly editName = signal('');
  readonly editSubject = signal('');
  readonly editHtmlBody = signal('');
  readonly editTextBody = signal('');
  readonly editVariables = signal('');

  // Preview
  readonly showPreview = signal(false);
  readonly previewSubject = signal('');
  readonly previewHtml = signal('');
  readonly isLoadingPreview = signal(false);
  readonly previewError = signal<string | null>(null);
  readonly sampleData = signal('');

  // Delete
  readonly showDeleteConfirm = signal(false);
  readonly isDeleting = signal(false);

  private templateId = '';

  constructor() {
    this.templateId = this.route.snapshot.params['id'];
    this.loadTemplate();
  }

  loadTemplate(): void {
    this.isLoading.set(true);
    this.http
      .get<{ data: Template }>(`/v1/templates/${this.templateId}`)
      .subscribe({
        next: (res) => {
          const t = res.data;
          this.template.set(t);
          this.editName.set(t.name);
          this.editSubject.set(t.subject);
          this.editHtmlBody.set(t.html_body);
          this.editTextBody.set(t.text_body);
          this.editVariables.set(t.variables.join(', '));
          this.buildSampleData(t.variables);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  goBack(): void {
    this.router.navigate(['/templates']);
  }

  saveTemplate(): void {
    const name = this.editName().trim();
    const subject = this.editSubject().trim();
    if (!name || !subject) return;

    this.isSaving.set(true);
    this.saveError.set(null);
    this.saveSuccess.set(false);

    const variables = this.editVariables()
      .split(',')
      .map((v) => v.trim())
      .filter((v) => v.length > 0);

    this.http
      .patch<{ data: Template }>(`/v1/templates/${this.templateId}`, {
        name,
        subject,
        html_body: this.editHtmlBody(),
        text_body: this.editTextBody(),
        variables,
      })
      .subscribe({
        next: (res) => {
          this.template.set(res.data);
          this.isSaving.set(false);
          this.saveSuccess.set(true);
          setTimeout(() => this.saveSuccess.set(false), 3000);
        },
        error: (err) => {
          this.isSaving.set(false);
          this.saveError.set(
            err.error?.message ?? 'Failed to save template. Please try again.',
          );
        },
      });
  }

  loadPreview(): void {
    this.isLoadingPreview.set(true);
    this.previewError.set(null);

    let variables: Record<string, string> = {};
    try {
      const raw = this.sampleData().trim();
      if (raw) {
        variables = JSON.parse(raw);
      }
    } catch {
      this.previewError.set('Invalid JSON for sample data.');
      this.isLoadingPreview.set(false);
      return;
    }

    this.http
      .post<PreviewResponse>(`/v1/templates/${this.templateId}/preview`, {
        variables,
      })
      .subscribe({
        next: (res) => {
          this.previewSubject.set(res.data.subject);
          this.previewHtml.set(res.data.html);
          this.showPreview.set(true);
          this.isLoadingPreview.set(false);
        },
        error: (err) => {
          this.isLoadingPreview.set(false);
          this.previewError.set(
            err.error?.message ?? 'Failed to load preview.',
          );
        },
      });
  }

  closePreview(): void {
    this.showPreview.set(false);
  }

  confirmDelete(): void {
    this.showDeleteConfirm.set(true);
  }

  cancelDelete(): void {
    this.showDeleteConfirm.set(false);
  }

  deleteTemplate(): void {
    this.isDeleting.set(true);
    this.http.delete(`/v1/templates/${this.templateId}`).subscribe({
      next: () => {
        this.isDeleting.set(false);
        this.router.navigate(['/templates']);
      },
      error: () => this.isDeleting.set(false),
    });
  }

  private buildSampleData(variables: string[]): void {
    if (variables.length === 0) {
      this.sampleData.set('{}');
      return;
    }
    const sample: Record<string, string> = {};
    for (const v of variables) {
      sample[v] = `sample_${v}`;
    }
    this.sampleData.set(JSON.stringify(sample, null, 2));
  }
}
