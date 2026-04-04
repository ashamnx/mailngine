import { Component, inject, signal } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
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

@Component({
  selector: 'app-templates',
  standalone: true,
  imports: [FormsModule, DatePipe],
  templateUrl: './templates.component.html',
  styleUrl: './templates.component.scss',
})
export class TemplatesComponent {
  private readonly http = inject(HttpClient);
  private readonly router = inject(Router);

  readonly templates = signal<Template[]>([]);
  readonly isLoading = signal(true);
  readonly showCreateForm = signal(false);
  readonly isCreating = signal(false);
  readonly createError = signal<string | null>(null);

  // Create form fields
  readonly newName = signal('');
  readonly newSubject = signal('');
  readonly newHtmlBody = signal('');
  readonly newTextBody = signal('');
  readonly newVariables = signal('');

  constructor() {
    this.loadTemplates();
  }

  loadTemplates(): void {
    this.isLoading.set(true);
    this.http
      .get<{ data: Template[] }>('/v1/templates')
      .subscribe({
        next: (res) => {
          this.templates.set(res.data);
          this.isLoading.set(false);
        },
        error: () => this.isLoading.set(false),
      });
  }

  openCreateForm(): void {
    this.showCreateForm.set(true);
    this.resetForm();
  }

  closeCreateForm(): void {
    this.showCreateForm.set(false);
    this.resetForm();
  }

  createTemplate(): void {
    const name = this.newName().trim();
    const subject = this.newSubject().trim();
    if (!name || !subject) return;

    this.isCreating.set(true);
    this.createError.set(null);

    const variables = this.newVariables()
      .split(',')
      .map((v) => v.trim())
      .filter((v) => v.length > 0);

    this.http
      .post<{ data: Template }>('/v1/templates', {
        name,
        subject,
        html_body: this.newHtmlBody(),
        text_body: this.newTextBody(),
        variables,
      })
      .subscribe({
        next: (res) => {
          this.isCreating.set(false);
          this.closeCreateForm();
          this.router.navigate(['/templates', res.data.id]);
        },
        error: (err) => {
          this.isCreating.set(false);
          this.createError.set(
            err.error?.message ?? 'Failed to create template. Please try again.',
          );
        },
      });
  }

  deleteTemplate(id: string, event: Event): void {
    event.stopPropagation();
    this.http.delete(`/v1/templates/${id}`).subscribe({
      next: () => this.loadTemplates(),
    });
  }

  navigateToTemplate(id: string): void {
    this.router.navigate(['/templates', id]);
  }

  private resetForm(): void {
    this.newName.set('');
    this.newSubject.set('');
    this.newHtmlBody.set('');
    this.newTextBody.set('');
    this.newVariables.set('');
    this.createError.set(null);
  }
}
