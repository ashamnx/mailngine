import { Component, inject } from '@angular/core';
import { ActivatedRoute } from '@angular/router';

@Component({
  selector: 'app-placeholder',
  standalone: true,
  template: `
    <div class="placeholder">
      <span class="material-symbols-rounded lg">construction</span>
      <h2>{{ title }}</h2>
      <p>This feature is coming soon.</p>
    </div>
  `,
  styles: [
    `
      .placeholder {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 80px 24px;
        text-align: center;
        gap: 16px;
      }

      .placeholder span.material-symbols-rounded {
        color: #9aa0a6;
      }

      .placeholder h2 {
        font-size: 1.5rem;
        color: #202124;
      }

      .placeholder p {
        font-size: 0.875rem;
        color: #5f6368;
      }
    `,
  ],
})
export class PlaceholderComponent {
  title = inject(ActivatedRoute).snapshot.data['title'] ?? 'Coming Soon';
}
